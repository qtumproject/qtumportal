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
}

type ServerOption struct {
	Port          int
	AuthPort      int
	StaticBaseDir string
	QtumdRPCURL   *url.URL
}

func NewServer(opts ServerOption) *Server {
	authStore := newAuthorizationStore()

	return &Server{
		Options:   opts,
		authStore: authStore,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.POST("/qtumd", s.proxyRPC)

	addr := fmt.Sprintf(":%d", s.Options.Port)
	fmt.Sprintln("Server listening", addr)
	return e.Start(addr)
}

func (s *Server) StartAuthService() error {
	e := echo.New()

	e.GET("/authorizations", s.listAuthorizations)

	addr := fmt.Sprintf(":%d", s.Options.AuthPort)
	fmt.Sprintln("Authorization service listening", addr)
	return e.Start(addr)
}

func (s *Server) listAuthorizations(c echo.Context) error {
	auths := s.authStore.pendingAuthorizations()
	return c.JSON(http.StatusOK, auths)
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

	if !method.NoAuth {
		auth, err := s.authStore.create(&jsonRPCReq)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusPaymentRequired, auth)
	}

	rpcURL := s.Options.QtumdRPCURL

	rpcBodyBytes, err := json.Marshal(&jsonRPCReq)
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
