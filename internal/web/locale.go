package web

import (
	"github.com/labstack/echo/v4"

	"mljr-web/internal/i18n"
)

// LangContextKey is the echo.Context key the active locale is stored under.
const LangContextKey = "lang"

// LangCookieName is the cookie used to persist the user's language choice.
const LangCookieName = "lang"

// LocaleMiddleware reads the "lang" cookie, validates it against the
// supported locales, and stores the resolved language code in the request
// context under LangContextKey (defaulting to i18n.DefaultLang).
func LocaleMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			lang := i18n.DefaultLang
			if cookie, err := c.Cookie(LangCookieName); err == nil && i18n.IsSupported(cookie.Value) {
				lang = cookie.Value
			}
			c.Set(LangContextKey, lang)
			return next(c)
		}
	}
}

// Lang returns the active locale for the request, defaulting to
// i18n.DefaultLang if unset.
func Lang(c echo.Context) string {
	if lang, ok := c.Get(LangContextKey).(string); ok && lang != "" {
		return lang
	}
	return i18n.DefaultLang
}
