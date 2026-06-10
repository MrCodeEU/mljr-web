package web

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo/v4"
	g "maragu.dev/gomponents"
)

// Render writes a gomponents node as a UTF-8 HTML response.
func Render(c echo.Context, status int, n g.Node) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/html; charset=utf-8")
	c.Response().WriteHeader(status)
	return n.Render(c.Response().Writer)
}

// RenderToString materializes a node into an HTML string, e.g. for
// sse.PatchElements(html string).
func RenderToString(n g.Node) string {
	var buf bytes.Buffer
	_ = n.Render(&buf)
	return buf.String()
}

// Ensure http import is used in tests; silences unused-import lint if shaved.
var _ = http.StatusOK
