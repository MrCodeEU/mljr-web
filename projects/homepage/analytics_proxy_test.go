package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mljr-web/internal/config"

	"github.com/labstack/echo/v4"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestRegisterAnalyticsProxyForwardsUmamiScript(t *testing.T) {
	var upstreamAcceptEncoding string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		upstreamAcceptEncoding = r.Header.Get("Accept-Encoding")
		if r.URL.Path != "/script.js" {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("not found")),
				Request:    r,
			}, nil
		}
		header := make(http.Header)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(strings.NewReader("window.umami={track:function(){}};")),
			Request:    r,
		}, nil
	})

	e := echo.New()
	if err := registerAnalyticsProxyWithTransport(e, config.AnalyticsConfig{UmamiProxyTarget: "https://stats.example.test"}, transport); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/umami/script.js", nil)
	req.Header.Set("Accept-Encoding", "br, gzip, deflate")
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/javascript; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if upstreamAcceptEncoding != "" {
		t.Fatalf("upstream Accept-Encoding = %q, want empty", upstreamAcceptEncoding)
	}
}

func TestRegisterAnalyticsProxyScriptErrorUsesJavaScriptMime(t *testing.T) {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return nil, errors.New("upstream unavailable")
	})

	e := echo.New()
	if err := registerAnalyticsProxyWithTransport(e, config.AnalyticsConfig{UmamiProxyTarget: "https://stats.example.test"}, transport); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/umami/script.js", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/javascript; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
}

func TestRegisterAnalyticsProxyIgnoresEmptyTarget(t *testing.T) {
	e := echo.New()
	if err := registerAnalyticsProxy(e, config.AnalyticsConfig{}); err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/umami/script.js", nil)
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestJoinURLPath(t *testing.T) {
	tests := map[string]struct {
		base string
		path string
		want string
	}{
		"root base":          {"", "/script.js", "/script.js"},
		"nested base":        {"/analytics", "/script.js", "/analytics/script.js"},
		"empty path":         {"/analytics", "", "/analytics"},
		"missing slashes":    {"/analytics/", "api/send", "/analytics/api/send"},
		"root request":       {"/analytics", "/", "/analytics"},
		"root base no slash": {"", "script.js", "/script.js"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := joinURLPath(tt.base, tt.path); got != tt.want {
				t.Fatalf("joinURLPath(%q, %q) = %q, want %q", tt.base, tt.path, got, tt.want)
			}
		})
	}
}
