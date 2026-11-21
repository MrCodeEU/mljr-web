package web

import (
	"html"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"

	"mljr-web/internal/i18n"
)

// LangContextKey is the echo.Context key the active locale is stored under.
const LangContextKey = "lang"

// LangCookieName is the cookie used to persist the user's language choice.
const LangCookieName = "lang"

// SiteURL is the public origin used for canonical and hreflang links.
const SiteURL = "https://mljr.eu"

// Context keys set by LangPrefix.
const (
	prefixModeKey = "langPrefixMode" // LangPrefix is active for this server
	urlLangKey    = "urlLang"        // language taken from the URL prefix
	pagePathKey   = "pagePath"       // request path without the language prefix
)

// LangPrefix (use with e.Pre) gives German pages their own URLs: /de/...
// is served by the same routes as /... with the language set to "de". The
// prefix is stripped before routing; LocaleMiddleware then prefers the URL
// language over the cookie, and Render adds canonical and hreflang links.
func LangPrefix() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(prefixModeKey, true)
			r := c.Request()
			if r.URL.Path == "/de" || strings.HasPrefix(r.URL.Path, "/de/") {
				rest := strings.TrimPrefix(r.URL.Path, "/de")
				if rest == "" {
					rest = "/"
				}
				r.URL.Path, r.URL.RawPath = rest, ""
				c.Set(urlLangKey, "de")
			}
			c.Set(pagePathKey, r.URL.Path)
			return next(c)
		}
	}
}

// LocaleMiddleware resolves the request language and stores it in the
// request context under LangContextKey (defaulting to i18n.DefaultLang).
//
// Without LangPrefix the "lang" cookie decides. With it, the URL decides
// for pages: /de/... is German, everything else English. A page visit under
// /de/ remembers the choice in the cookie, and a visitor whose cookie says
// "de" who opens an unprefixed page is redirected to its /de/ twin, so
// plain internal links keep the chosen language. Non-page requests (API,
// SSE) still follow the cookie.
func LocaleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookieLang := ""
			if cookie, err := c.Cookie(LangCookieName); err == nil && i18n.IsSupported(cookie.Value) {
				cookieLang = cookie.Value
			}
			lang := i18n.DefaultLang
			if cookieLang != "" {
				lang = cookieLang
			}
			urlLang, _ := c.Get(urlLangKey).(string)
			if urlLang != "" {
				// /de/... is German for every client, also crawlers and
				// tools that do not send Accept: text/html.
				lang = urlLang
			}
			if on, _ := c.Get(prefixModeKey).(bool); on && isPageRequest(c) {
				if u := urlLang; u != "" {
					if cookieLang != u {
						c.SetCookie(&http.Cookie{Name: LangCookieName, Value: u, Path: "/", MaxAge: 31536000, SameSite: http.SameSiteLaxMode})
					}
				} else if cookieLang == "de" {
					target := "/de" + c.Request().URL.Path
					if q := c.Request().URL.RawQuery; q != "" {
						target += "?" + q
					}
					return c.Redirect(http.StatusFound, target)
				} else {
					lang = i18n.DefaultLang
				}
			}
			c.Set(LangContextKey, lang)
			return next(c)
		}
	}
}

// isPageRequest reports whether c is a browser navigation to an HTML page
// (a document or an iframe), as opposed to an asset, API or SSE request.
func isPageRequest(c echo.Context) bool {
	r := c.Request()
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}
	p := r.URL.Path
	for _, pre := range []string{"/api/", "/sse", "/static/", "/assets/", "/.well-known/"} {
		if strings.HasPrefix(p, pre) {
			return false
		}
	}
	if strings.Contains(p[strings.LastIndex(p, "/")+1:], ".") {
		return false // robots.txt, sitemap.xml, files
	}
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

// localizedHead returns the canonical and hreflang links for the page
// being rendered, or nil when LangPrefix is not active.
func localizedHead(c echo.Context) []byte {
	raw, ok := c.Get(pagePathKey).(string)
	if !ok {
		return nil
	}
	// The path comes from the request: re-encode it as a URL path and
	// HTML-escape it so it can never break out of the href attribute.
	path := html.EscapeString((&url.URL{Path: raw}).EscapedPath())
	de := "/de" + path
	if path == "/" {
		de = "/de/"
	}
	own := path
	if u, _ := c.Get(urlLangKey).(string); u == "de" {
		own = de
	}
	return []byte(`<link rel="canonical" href="` + SiteURL + own + `">` +
		`<link rel="alternate" hreflang="en" href="` + SiteURL + path + `">` +
		`<link rel="alternate" hreflang="de" href="` + SiteURL + de + `">` +
		`<link rel="alternate" hreflang="x-default" href="` + SiteURL + path + `">`)
}

// Lang returns the active locale for the request, defaulting to
// i18n.DefaultLang if unset.
func Lang(c echo.Context) string {
	if lang, ok := c.Get(LangContextKey).(string); ok && lang != "" {
		return lang
	}
	return i18n.DefaultLang
}
