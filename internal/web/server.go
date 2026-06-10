package web

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewEcho returns an Echo instance pre-configured with logging, recovery,
// security headers, and gzip compression excluding SSE endpoints.
func NewEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogError:     true,
		LogLatency:   true,
		HandleError:  true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("request remote_ip=%s host=%s method=%s uri=%s status=%d latency=%s error=%v user_agent=%q",
				v.RemoteIP, v.Host, v.Method, v.URI, v.Status, v.Latency, v.Error, v.UserAgent)
			return nil
		},
	}))
	e.Use(SecurityHeaders())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			p := c.Path()
			// SSE endpoints must not be buffered/compressed.
			return p == "/sse" || (len(p) >= 5 && p[:5] == "/sse/") || (len(p) >= 5 && p[:5] == "/api/")
		},
	}))
	return e
}
