package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path"

	"github.com/labstack/echo/middleware"

	"github.com/hayeah/qtum-portal/ui"

	"github.com/pkg/errors"

	"github.com/labstack/echo"
)

type Server struct {
	Options ServerOption

	authStore *authorizationStore

	authApp  *echo.Echo
	proxyApp *echo.Echo
}

type ServerOption struct {
	Port          int
	AuthPort      int
	StaticBaseDir string
	QtumdRPCURL   *url.URL
	DebugMode     bool
}

func NewServer(opts ServerOption) *Server {
	authStore := newAuthorizationStore()

	s := &Server{
		Options:   opts,
		authStore: authStore,
	}

	e := echo.New()
	e.Logger.SetOutput(ioutil.Discard)
	e.HideBanner = true
	s.proxyApp = e
	e.POST("/qtumd", s.proxyRPC)

	e = echo.New()
	e.Logger.SetOutput(ioutil.Discard)
	e.HideBanner = true
	s.authApp = e

	if opts.DebugMode {
		e.Use(middleware.CORS())
	}
	// e.Use(middleware.Static("ui/build"))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := c.Request().URL.Path

			if p[len(p)-1] == '/' {
				p += "index.html"
			}

			// strip off leading /
			assetName := p[1:]
			data, err := ui.Asset(assetName)
			if err == nil {
				ext := path.Ext(p)
				contentType := mime.TypeByExtension(ext)
				if contentType == "" {
					contentType = http.DetectContentType(data)
				}
				return c.Blob(http.StatusOK, contentType, data)
			}

			return next(c)
		}
	})
	e.GET("/authorizations", s.listAuthorizations)
	e.GET("/authorizations/:id", s.getAuthorization)
	e.POST("/authorizations/:id/accept", s.acceptAuthorization)
	e.POST("/authorizations/:id/deny", s.acceptAuthorization)

	return s
}

func (s *Server) Start() error {
	errC := make(chan error)
	go func() {
		errC <- s.startProxyService()
	}()

	go func() {
		errC <- s.startAuthService()
	}()

	return <-errC
}

func (s *Server) startProxyService() error {
	addr := fmt.Sprintf(":%d", s.Options.Port)
	log.Println("RPC service listening", addr)
	return s.proxyApp.Start(addr)
}

func (s *Server) startAuthService() error {
	addr := fmt.Sprintf(":%d", s.Options.AuthPort)
	log.Println("Auth service listening", addr)
	return s.authApp.Start(addr)
}

func (s *Server) listAuthorizations(c echo.Context) error {
	auths := s.authStore.pendingAuthorizations()
	return c.JSON(http.StatusOK, auths)
}

func (s *Server) getAuthorization(c echo.Context) error {
	id := c.Param("id")
	auth, found := s.authStore.get(id)

	if !found {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, auth)
}

func (s *Server) acceptAuthorization(c echo.Context) error {
	id := c.Param("id")
	err := s.authStore.accept(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) denyAuthorization(c echo.Context) error {
	id := c.Param("id")
	err := s.authStore.deny(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) proxyRPC(c echo.Context) error {
	// {
	// 	"jsonrpc": "1.0",
	// 	"id":"1",
	// 	"method": "getinfo",
	// 	"params": []
	// }

	var jsonRPCReq jsonRPCRequest

	err := c.Bind(&jsonRPCReq)
	if err != nil {
		return errors.Wrap(err, "rpc request")
	}

	methodName := jsonRPCReq.Method

	method, found := rpcMethods[methodName]

	if !found {
		return c.JSON(http.StatusNotFound, &jsonRPCError{
			Code:    0,
			Message: fmt.Sprintf("RPC method not supported: %s", methodName),
		})
	}

	if method.NoAuth {
		return s.doProxyRPCCall(c, &jsonRPCReq)
	}

	return s.doRPCCallAuth(c, &jsonRPCReq)
}

func (s *Server) doRPCCallAuth(c echo.Context, jsonRPCReq *jsonRPCRequest) error {
	// If no auth token is provided, create authorization object and return 402
	if jsonRPCReq.Auth == "" {
		auth, err := s.authStore.create(jsonRPCReq)
		if err != nil {
			return errors.Wrap(err, "auth")
		}

		return c.JSON(http.StatusPaymentRequired, auth)
	}

	// If auth token is provided, verify then proxy
	if s.authStore.verify(jsonRPCReq.Auth, jsonRPCReq) {
		return s.doProxyRPCCall(c, jsonRPCReq)
	}

	return echo.NewHTTPError(http.StatusForbidden, "Cannot verify RPC request")
}

func (s *Server) doProxyRPCCall(c echo.Context, jsonRPCReq *jsonRPCRequest) error {
	rpcURL := s.Options.QtumdRPCURL

	rpcBodyBytes, err := json.Marshal(jsonRPCReq)
	if err != nil {
		return errors.Wrap(err, "proxy rpc")
	}
	rpcBody := bytes.NewReader(rpcBodyBytes)

	rpcReq, err := http.NewRequest(http.MethodPost, rpcURL.String(), rpcBody)
	if err != nil {
		return errors.Wrap(err, "proxy RPC")
	}

	auth := c.Request().Header.Get("Authorization")
	rpcReq.Header.Set("Authorization", auth)

	rpcRes, err := http.DefaultClient.Do(rpcReq)
	if err != nil {
		return errors.Wrap(err, "proxy RPC")
	}
	defer rpcRes.Body.Close()

	return c.Stream(http.StatusOK, rpcRes.Header.Get("Content-Type"), rpcRes.Body)
}
