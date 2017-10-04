package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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
}

func NewServer(opts ServerOption) *Server {
	authStore := newAuthorizationStore()

	s := &Server{
		Options:   opts,
		authStore: authStore,
	}

	e := echo.New()
	s.proxyApp = e
	e.POST("/qtumd", s.proxyRPC)

	e = echo.New()
	s.authApp = e
	e.GET("/authorizations", s.listAuthorizations)
	e.GET("/authorizations/:id", s.getAuthorization)
	e.POST("/authorizations/:id/accept", s.acceptAuthorization)
	e.POST("/authorizations/:id/deny", s.acceptAuthorization)

	return s
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Options.Port)
	fmt.Sprintln("Server listening", addr)
	return s.proxyApp.Start(addr)
}

func (s *Server) StartAuthService() error {
	addr := fmt.Sprintf(":%d", s.Options.AuthPort)
	fmt.Sprintln("Authorization service listening", addr)
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
