package web

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/labstack/echo/v4"
)

// cspTail is the part of the Content-Security-Policy shared by all responses.
// script-src is prepended separately because it varies per response (hashes).
const cspTail = "style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data: https://*.tile.openstreetmap.org https://picsum.photos https://fastly.picsum.photos; " +
	"font-src 'self'; " +
	"connect-src 'self'; " +
	"media-src 'self'; " +
	"object-src 'none'; " +
	"frame-ancestors 'self'; " +
	"base-uri 'self'; " +
	"form-action 'self';"

// SecurityHeaders applies the default mljr-web security header set. All JS is
// self-hosted under /static/, so script-src 'self' is the strict baseline.
// SRI on <script> tags is unnecessary because nothing is loaded cross-origin.
//
// script-src has NO 'unsafe-inline': inline <script> blocks are allowlisted
// per response via sha256 hashes that web.Render computes from the final HTML
// (SetCSPWithScriptHashes). The policy set here is only the fallback for
// non-HTML responses. 'unsafe-eval' is still required: Datastar v1.x
// evaluates all data-* expressions via new Function() at runtime and has no
// precompile mode. 'unsafe-inline' in style-src covers the inline style=
// attributes components use for dynamic values.
func SecurityHeaders() echo.MiddlewareFunc {
	fallback := "default-src 'none'; script-src 'self' 'unsafe-eval'; " + cspTail
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("Content-Security-Policy", fallback)
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			h.Set("X-Frame-Options", "SAMEORIGIN")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-origin")
			return next(c)
		}
	}
}

// SetCSPWithScriptHashes replaces the fallback CSP with one that allowlists
// each given inline <script> content via its sha256 hash. Hashes are computed
// over the exact bytes between <script> and </script>; duplicates are emitted
// once. Called by web.Render for every HTML response.
func SetCSPWithScriptHashes(c echo.Context, scripts [][]byte) {
	seen := make(map[string]struct{}, len(scripts))
	var b strings.Builder
	b.WriteString("default-src 'none'; script-src 'self' 'unsafe-eval'")
	for _, s := range scripts {
		sum := sha256.Sum256(s)
		hash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
		if _, dup := seen[hash]; dup {
			continue
		}
		seen[hash] = struct{}{}
		b.WriteString(" ")
		b.WriteString(hash)
	}
	b.WriteString("; ")
	b.WriteString(cspTail)
	c.Response().Header().Set("Content-Security-Policy", b.String())
}
