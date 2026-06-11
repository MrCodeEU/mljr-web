package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"mljr-web/internal/config"

	"github.com/labstack/echo/v4"
)

func registerAnalyticsProxy(e *echo.Echo, cfg config.AnalyticsConfig) error {
	return registerAnalyticsProxyWithTransport(e, cfg, nil)
}

func registerAnalyticsProxyWithTransport(e *echo.Echo, cfg config.AnalyticsConfig, transport http.RoundTripper) error {
	if strings.TrimSpace(cfg.UmamiProxyTarget) == "" {
		return nil
	}

	target, err := url.Parse(cfg.UmamiProxyTarget)
	if err != nil {
		return err
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	if transport != nil {
		proxy.Transport = transport
	}
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = joinURLPath(target.Path, strings.TrimPrefix(req.URL.Path, "/umami"))
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		req.Header.Del("Accept-Encoding")
	}
	proxy.ModifyResponse = func(res *http.Response) error {
		if strings.HasSuffix(res.Request.URL.Path, "/script.js") {
			res.Header.Set("Content-Type", "application/javascript; charset=utf-8")
		}
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if strings.HasSuffix(r.URL.Path, "/script.js") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("console.warn('Umami proxy unavailable');"))
			return
		}
		http.Error(w, "Umami proxy unavailable", http.StatusBadGateway)
	}

	handler := echo.WrapHandler(http.StripPrefix("", proxy))
	e.Any("/umami", handler)
	e.Any("/umami/*", handler)
	return nil
}

func joinURLPath(base, requestPath string) string {
	if base == "" || base == "/" {
		if strings.HasPrefix(requestPath, "/") {
			return requestPath
		}
		return "/" + requestPath
	}
	if requestPath == "" || requestPath == "/" {
		return base
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(requestPath, "/")
}
