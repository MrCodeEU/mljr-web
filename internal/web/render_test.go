package web

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// cspDirective extracts one directive (e.g. "script-src ...") from a CSP.
func cspDirective(csp, name string) string {
	for _, d := range strings.Split(csp, ";") {
		if d = strings.TrimSpace(d); strings.HasPrefix(d, name+" ") {
			return d
		}
	}
	return ""
}

func scriptHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
}

func TestRenderAllowlistsInlineScriptsViaHash(t *testing.T) {
	const script = `console.log("hi")`
	e := NewEcho()
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, h.HTML(
			h.Head(h.Script(g.Raw(script))),
			h.Body(
				h.Script(h.Src("/static/app.js")),
				h.Div(g.Text("x")),
			),
		))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), script) {
		t.Fatalf("body %q missing inline script", rec.Body.String())
	}

	scriptSrc := cspDirective(rec.Header().Get("Content-Security-Policy"), "script-src")
	if !strings.Contains(scriptSrc, scriptHash(script)) {
		t.Fatalf("script-src = %q, want hash %q for inline script", scriptSrc, scriptHash(script))
	}
	if strings.Contains(scriptSrc, "'unsafe-inline'") {
		t.Fatalf("script-src = %q, must not allow unsafe-inline", scriptSrc)
	}
	if !strings.Contains(scriptSrc, "'unsafe-eval'") {
		t.Fatalf("script-src = %q, must keep unsafe-eval for Datastar", scriptSrc)
	}
}

func TestRenderDeduplicatesIdenticalScriptHashes(t *testing.T) {
	const script = `window.x=1;`
	e := NewEcho()
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, h.HTML(
			h.Body(
				h.Script(g.Raw(script)),
				h.Script(g.Raw(script)),
			),
		))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(rec, req)

	scriptSrc := cspDirective(rec.Header().Get("Content-Security-Policy"), "script-src")
	if n := strings.Count(scriptSrc, scriptHash(script)); n != 1 {
		t.Fatalf("script-src = %q, hash appears %d times, want 1", scriptSrc, n)
	}
}

func TestRenderWithoutInlineScriptsKeepsStrictScriptSrc(t *testing.T) {
	e := NewEcho()
	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, h.HTML(h.Body(h.Div(g.Text("plain")))))
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(rec, req)

	scriptSrc := cspDirective(rec.Header().Get("Content-Security-Policy"), "script-src")
	if strings.Contains(scriptSrc, "'unsafe-inline'") || strings.Contains(scriptSrc, "sha256-") {
		t.Fatalf("script-src = %q, want strict 'self' 'unsafe-eval' only", scriptSrc)
	}
}
