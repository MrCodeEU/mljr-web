package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func TestNewEchoAppliesSecurityHeaders(t *testing.T) {
	e := NewEcho()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	headers := rec.Header()
	if got := headers.Get("Content-Security-Policy"); !strings.Contains(got, "default-src 'self'") {
		t.Fatalf("Content-Security-Policy = %q, want self default-src", got)
	}
	if got := headers.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
	if got := headers.Get("X-Frame-Options"); got != "SAMEORIGIN" {
		t.Fatalf("X-Frame-Options = %q, want SAMEORIGIN", got)
	}
}

func TestRenderWritesHTMLResponse(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := Render(c, http.StatusCreated, h.Div(g.Text("hello"))); err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get(echo.HeaderContentType); got != "text/html; charset=utf-8" {
		t.Fatalf("content type = %q, want text/html; charset=utf-8", got)
	}
	if got := rec.Body.String(); got != "<div>hello</div>" {
		t.Fatalf("body = %q, want rendered div", got)
	}
}

func TestRenderToString(t *testing.T) {
	if got := RenderToString(h.Span(g.Text("ready"))); got != "<span>ready</span>" {
		t.Fatalf("RenderToString() = %q", got)
	}
}
