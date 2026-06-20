package pages

import (
	"net/http"
	"time"

	"mljr-web/internal/web"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
)

func renderPage(re *core.RequestEvent, status int, n g.Node) error {
	return web.RenderPB(re, status, n)
}

const sessionCookieName = "nl_session"

// setSession issues a long-lived auth-token cookie for the given user record.
func setSession(re *core.RequestEvent, user *core.Record) error {
	token, err := user.NewAuthToken()
	if err != nil {
		return err
	}
	http.SetCookie(re.Response, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})
	return nil
}

func clearSession(re *core.RequestEvent) {
	http.SetCookie(re.Response, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// currentUser resolves the logged-in user from the session cookie, if any.
func currentUser(re *core.RequestEvent) *core.Record {
	cookie, err := re.Request.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}
	user, err := re.App.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth)
	if err != nil {
		return nil
	}
	return user
}

// HandleLogout clears the session cookie and sends the user back to /login.
func HandleLogout(re *core.RequestEvent) error {
	clearSession(re)
	return redirect(re, "/login")
}

func redirect(re *core.RequestEvent, url string) error {
	re.Response.Header().Set("Location", url)
	re.Response.WriteHeader(http.StatusSeeOther)
	return nil
}
