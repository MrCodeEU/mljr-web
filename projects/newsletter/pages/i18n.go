package pages

import (
	"mljr-web/internal/i18n"

	"github.com/pocketbase/pocketbase/core"
)

// currentLang resolves the request's language from the "lang" cookie (set
// by the LanguageToggle component), falling back to i18n.DefaultLang. There
// is no PocketBase-side auth middleware like the homepage's Echo stack, so
// every handler reads it directly off the cookie, mirroring how
// currentUser (session.go) reads the session cookie.
func currentLang(re *core.RequestEvent) string {
	cookie, err := re.Request.Cookie("lang")
	if err != nil || !i18n.IsSupported(cookie.Value) {
		return i18n.DefaultLang
	}
	return cookie.Value
}

// translator returns a t(key, args...) closure bound to the request's
// language, so page-building code reads naturally without repeating
// i18n.T(lang, ...) at every call site.
func translator(re *core.RequestEvent) func(key string, args ...any) string {
	lang := currentLang(re)
	return func(key string, args ...any) string {
		return i18n.T(lang, key, args...)
	}
}
