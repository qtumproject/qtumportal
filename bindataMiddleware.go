package portal

import (
	"bytes"
	"encoding/json"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func newBindataMiddleware(config bindataConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		s := bindataMiddleware{
			bindataConfig: config,
			next:          next,
		}

		return s.Handle
	}
}

type bindataConfig struct {
	// the prefix path to serve bindata resources
	prefix      string
	getter      bindataGetter
	jsConstants map[string]interface{}
}

type bindataGetter = func(key string) ([]byte, error)

type bindataMiddleware struct {
	bindataConfig
	next echo.HandlerFunc
}

func (s *bindataMiddleware) injectJSConstantsIntoHTML(html []byte) ([]byte, error) {
	var buf bytes.Buffer

	buf.Write([]byte(`<body>
		<script type="text/javascript">
		//<!CDATA[[
		Object.assign(window,
		`))

	enc := json.NewEncoder(&buf)
	err := enc.Encode(s.jsConstants)
	if err != nil {
		return nil, errors.Wrap(err, "index.html JS constants inject")
	}

	buf.Write([]byte(`);
		//]]>
		</script>`))

	return bytes.Replace(html, []byte("<body>"), buf.Bytes(), 1), nil
}

func (s *bindataMiddleware) Handle(c echo.Context) error {
	next := s.next

	if c.Request().Method != "GET" {
		return next(c)
	}

	p := c.Request().URL.Path

	if s.prefix != "" {
		if !strings.HasPrefix(p, s.prefix) {
			return next(c)
		}

		p = strings.TrimPrefix(p, s.prefix)
	}

	// enforce "/" at the end of url
	if p == "" {
		c.Redirect(http.StatusPermanentRedirect, c.Request().URL.String()+"/")
		return nil
	}

	// TODO support prefix paths
	isIndex := p == "/"

	if isIndex {
		p = "/index.html"
	}

	// strip off leading /
	assetName := strings.TrimPrefix(p, "/")
	data, err := s.getter(assetName)

	if isIndex && s.jsConstants != nil {
		data, err = s.injectJSConstantsIntoHTML(data)
		if err != nil {
			return err
		}
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

// func newAssetHandler() echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 	}
// }
