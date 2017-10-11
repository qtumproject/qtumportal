package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/labstack/echo/middleware"

	"github.com/hayeah/qtum-portal/ui"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/olebedev/emitter"
)

type Server struct {
	Options ServerOption

	authStore *authorizationStore

	authApp  *echo.Echo
	proxyApp *echo.Echo

	wsUpgrader *websocket.Upgrader
	emitter    *emitter.Emitter
}

type ServerOption struct {
	DAppPort      int
	AuthPort      int
	StaticBaseDir string
	QtumdRPCURL   *url.URL
	DebugMode     bool
}

type qtumPortalUIConfig struct {
	AuthBaseURL string `json:"AUTH_BASEURL"`
}

func NewServer(opts ServerOption) *Server {
	authStore := newAuthorizationStore()

	var wsCheckOrigin func(req *http.Request) bool

	if opts.QtumdRPCURL.User == nil {
		panic("must specify user and password in QTUM_RPC URL")
	}

	if opts.DebugMode {
		wsCheckOrigin = func(req *http.Request) bool {
			return true
		}
	}

	s := &Server{
		Options:   opts,
		authStore: authStore,
		emitter:   &emitter.Emitter{},
		wsUpgrader: &websocket.Upgrader{
			CheckOrigin: wsCheckOrigin,
		},
	}

	staticDir, err := filepath.Abs(opts.StaticBaseDir)
	if err != nil {
		log.Errorf("Invalid DApp dir %s: %s", staticDir, err)
		os.Exit(1)
	}

	e := echo.New()
	if opts.DebugMode {
		e.Use(middleware.CORS())
	}

	e.Logger.SetOutput(ioutil.Discard)
	e.HideBanner = true
	e.HTTPErrorHandler = errorHandler
	s.proxyApp = e
	e.POST("/", s.proxyRPC)
	e.GET("/api/authorizations/:id", s.getAuthorization)
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		// support for SPA
		HTML5: true,
		Root:  opts.StaticBaseDir,
		Skipper: func(c echo.Context) bool {
			return c.Request().Method != "GET"
		},
		Index: "index.html",
	}))
	log.Println("Serving DApp from", staticDir)

	e = echo.New()
	e.Logger.SetOutput(ioutil.Discard)
	e.HideBanner = true
	e.HTTPErrorHandler = errorHandler
	s.authApp = e

	if opts.DebugMode {
		e.Use(middleware.CORS())
	}
	// e.Use(middleware.Static("ui/build"))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != "GET" {
				return next(c)
			}

			p := c.Request().URL.Path

			isIndex := p == "/"

			if isIndex {
				p += "index.html"
			}

			// strip off leading /
			assetName := p[1:]
			data, err := ui.Asset(assetName)

			if isIndex {
				config := qtumPortalUIConfig{
					AuthBaseURL: fmt.Sprintf("http://localhost:%d", opts.AuthPort),
				}

				var buf bytes.Buffer

				buf.Write([]byte(`<body>
<script type="text/javascript">
//<!CDATA[[
QTUMPORTAL_CONFIG =
`))

				enc := json.NewEncoder(&buf)
				err := enc.Encode(config)
				if err != nil {
					return errors.Wrap(err, "index.html JS inject")
				}
				buf.Write([]byte(`
//]]>
</script>`))

				data = bytes.Replace(data, []byte("<body>"), buf.Bytes(), 1)
			}

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
	e.Any("/events", s.subscribeToEvents)
	e.GET("/authorizations", s.listAuthorizations)
	e.GET("/authorizations/:id", s.getAuthorization)
	e.POST("/authorizations/:id/accept", s.acceptAuthorization)
	e.POST("/authorizations/:id/deny", s.denyAuthorization)

	return s
}

func (s *Server) Start() error {
	errC := make(chan error)
	go func() {
		errC <- s.startDAppService()
	}()

	go func() {
		errC <- s.startAuthService()
	}()

	return <-errC
}

func (s *Server) startDAppService() error {
	addr := fmt.Sprintf(":%d", s.Options.DAppPort)
	log.Println("DApp service listening", addr)
	return s.proxyApp.Start(addr)
}

func (s *Server) startAuthService() error {
	addr := fmt.Sprintf(":%d", s.Options.AuthPort)
	log.Println("Auth service listening", addr)
	return s.authApp.Start(addr)
}

func (s *Server) subscribeToEvents(c echo.Context) error {
	conn, err := s.wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for e := range s.emitter.On(eventRefresh) {
		conn.WriteMessage(websocket.TextMessage, []byte(e.OriginalTopic))
	}

	return nil
}

func (s *Server) listAuthorizations(c echo.Context) error {
	auths := s.authStore.allAuthorizations()
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

	s.emitter.Emit(eventRefresh)

	auth, _ := s.authStore.get(id)
	return c.JSON(http.StatusOK, auth)
}

func (s *Server) denyAuthorization(c echo.Context) error {
	id := c.Param("id")
	err := s.authStore.deny(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	s.emitter.Emit(eventRefresh)

	auth, _ := s.authStore.get(id)
	return c.JSON(http.StatusOK, auth)
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

	log.Println("RPC Authorization requested:", jsonRPCReq.Method)
	return s.doRPCCallAuth(c, &jsonRPCReq)
}

func (s *Server) doRPCCallAuth(c echo.Context, jsonRPCReq *jsonRPCRequest) error {
	// If no auth token is provided, create authorization object and return 402
	if jsonRPCReq.Auth == "" {
		auth, err := s.authStore.create(jsonRPCReq)
		if err != nil {
			return errors.Wrap(err, "auth")
		}

		jsonRPCReq.Auth = auth.ID

		s.emitter.Emit(eventRefresh)
		return c.JSON(http.StatusPaymentRequired, auth)
	}

	// If auth token is provided, verify then proxy
	if s.authStore.verify(jsonRPCReq.Auth, jsonRPCReq) {
		s.emitter.Emit(eventRefresh)
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
	if auth == "" {
		rpcReq.Header.Set("Authorization", auth)
	} else {
		user := rpcURL.User.Username()
		pass, hasPass := rpcURL.User.Password()

		if !hasPass || user == "" {
			return errors.New("Need to specify Authorization header for RPC call")
		}

		rpcReq.SetBasicAuth(user, pass)
	}

	rpcRes, err := http.DefaultClient.Do(rpcReq)
	if err != nil {
		return errors.Wrap(err, "proxy RPC")
	}
	defer rpcRes.Body.Close()

	return c.Stream(http.StatusOK, rpcRes.Header.Get("Content-Type"), rpcRes.Body)
}

func errorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		// msg = err.Error()
		msg = err
	}

	log.Println("code, msg", code, msg)
	log.Println("committed", c.Response().Committed)

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			if err := c.NoContent(code); err != nil {
				goto ERROR
			}
		} else {
			type errorMsg struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}

			errrmsg := &errorMsg{
				Code:    code,
				Message: fmt.Sprintf("%v", msg),
			}

			if err := c.JSON(code, errrmsg); err != nil {
				log.Println("error handle json", err)
				goto ERROR
			}
		}
	}
ERROR:
	log.Errorln(err)
}
