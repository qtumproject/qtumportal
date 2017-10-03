package portal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var testRPCURL *url.URL

func init() {
	qtumRPC, found := os.LookupEnv("QTUM_RPC")
	if !found {
		fmt.Println("Please specify RPC url with QTUM_RPC environment variable")
		os.Exit(1)
	}

	qtumRPCURL, err := url.Parse(qtumRPC)
	if err != nil {
		log.Println("Invalid QTUM_RPC", qtumRPC)
	}

	testRPCURL = qtumRPCURL
}

func testProxy(jsonreq *jsonRPCRequest) (*http.Response, error) {
	jsonreqBodyBytes, err := json.Marshal(jsonreq)
	if err != nil {
		return nil, err
	}

	jsonreqBody := bytes.NewReader(jsonreqBodyBytes)

	req := httptest.NewRequest("POST", "/qtumrpc", jsonreqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)

	opts := ServerOption{
		// Port:        9999,
		QtumdRPCURL: testRPCURL,
	}
	s := &Server{
		Options: opts,
	}

	err = s.proxyRPC(c)
	if err != nil {
		return nil, err
	}

	res := rec.Result()

	return res, nil
}

func TestProxyMethodNotFound(t *testing.T) {
	is := assert.New(t)

	jsonreq := &jsonRPCRequest{
		Method: "no-such-method",
	}

	res, err := testProxy(jsonreq)
	is.NoError(err)
	is.Equal(404, res.StatusCode)
}

func TestProxyMethodFound(t *testing.T) {
	is := assert.New(t)

	jsonreq := &jsonRPCRequest{
		Method: "getinfo",
	}

	res, err := testProxy(jsonreq)
	is.NoError(err)
	is.Equal(200, res.StatusCode)
}

func TestProxyUserAuthorization(t *testing.T) {
	is := assert.New(t)

	jsonreq := &jsonRPCRequest{
		Method: "getnewaddress",
	}

	res, err := testProxy(jsonreq)
	is.NoError(err)
	is.Equal(402, res.StatusCode)
}
