package web

import "github.com/labstack/echo/v4"

// SecurityHeaders applies the default mljr-web security header set. All JS is
// self-hosted under /static/, so script-src 'self' is the strict baseline.
// SRI on <script> tags is unnecessary because nothing is loaded cross-origin.
func SecurityHeaders() echo.MiddlewareFunc {
	// 'unsafe-eval' is required: Datastar v1.x evaluates data-* expressions
	// via new Function() at runtime. There is no precompile mode.
	// 'unsafe-inline' covers the pre-paint FOUC-prevention inline <script>.
	const csp = "default-src 'self'; " +
		"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
		"style-src 'self' 'unsafe-inline'; " +
		"img-src 'self' data: https://*.tile.openstreetmap.org https://picsum.photos https://fastly.picsum.photos; " +
		"font-src 'self'; " +
		"connect-src 'self'; " +
		"frame-ancestors 'self'; " +
		"base-uri 'self'; " +
		"form-action 'self';"
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			h.Set("X-Frame-Options", "SAMEORIGIN")
			return next(c)
		}
	}
}
