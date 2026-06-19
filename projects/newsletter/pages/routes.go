package pages

import (
	"mljr-web/internal/web"

	"github.com/pocketbase/pocketbase/core"
)

// RegisterRoutes wires every page/handler route onto e.Router. Split out from
// main.go so pages/*_test.go can exercise full HTTP request/response cycles
// via tests.ApiScenario without booting a real server or registering static
// assets/cron (those stay project-specific in main.go).
func RegisterRoutes(e *core.ServeEvent) error {
	e.Router.GET("/login", func(re *core.RequestEvent) error {
		return web.RenderPB(re, 200, Login(re, LoginProps{}))
	})
	e.Router.POST("/login", func(re *core.RequestEvent) error {
		return HandleLogin(re)
	})
	e.Router.GET("/signup", func(re *core.RequestEvent) error {
		return web.RenderPB(re, 200, Signup(re, SignupProps{}))
	})
	e.Router.POST("/signup", func(re *core.RequestEvent) error {
		return HandleSignup(re)
	})
	e.Router.POST("/logout", func(re *core.RequestEvent) error {
		return HandleLogout(re)
	})
	e.Router.GET("/profile", func(re *core.RequestEvent) error {
		return Profile(re)
	})
	e.Router.POST("/profile", func(re *core.RequestEvent) error {
		return HandleProfile(re)
	})
	e.Router.GET("/", func(re *core.RequestEvent) error {
		return Dashboard(re)
	})
	e.Router.GET("/g/{slug}", func(re *core.RequestEvent) error {
		return GroupHome(re)
	})
	e.Router.GET("/g/{slug}/settings", func(re *core.RequestEvent) error {
		return GroupSettings(re)
	})
	e.Router.POST("/g/{slug}/settings", func(re *core.RequestEvent) error {
		return HandleGroupSettings(re)
	})
	e.Router.POST("/groups", func(re *core.RequestEvent) error {
		return HandleCreateGroup(re)
	})
	e.Router.GET("/g/{slug}/invites", func(re *core.RequestEvent) error {
		return ListInvites(re)
	})
	e.Router.POST("/g/{slug}/invites", func(re *core.RequestEvent) error {
		return HandleCreateInvite(re)
	})
	e.Router.POST("/g/{slug}/invites/{id}/revoke", func(re *core.RequestEvent) error {
		return HandleRevokeInvite(re)
	})
	e.Router.GET("/invites/{token}", func(re *core.RequestEvent) error {
		return InviteAccept(re)
	})
	e.Router.POST("/invites/{token}/accept", func(re *core.RequestEvent) error {
		return HandleAcceptInvite(re)
	})
	e.Router.POST("/notifications/read-all", func(re *core.RequestEvent) error {
		return HandleMarkAllNotificationsRead(re)
	})
	e.Router.GET("/g/{slug}/questions", func(re *core.RequestEvent) error {
		return ListQuestions(re)
	})
	e.Router.POST("/g/{slug}/questions", func(re *core.RequestEvent) error {
		return HandleCreateQuestion(re)
	})
	e.Router.POST("/g/{slug}/questions/{id}/toggle", func(re *core.RequestEvent) error {
		return HandleToggleQuestion(re)
	})
	e.Router.GET("/g/{slug}/editions", func(re *core.RequestEvent) error {
		return ListEditions(re)
	})
	e.Router.POST("/g/{slug}/editions", func(re *core.RequestEvent) error {
		return HandleCreateEdition(re)
	})
	e.Router.GET("/g/{slug}/editions/{id}", func(re *core.RequestEvent) error {
		return EditionAnswer(re)
	})
	e.Router.POST("/g/{slug}/editions/{id}", func(re *core.RequestEvent) error {
		return HandleSubmitAnswers(re)
	})
	e.Router.POST("/g/{slug}/editions/{id}/close", func(re *core.RequestEvent) error {
		return HandleCloseEdition(re)
	})
	e.Router.GET("/g/{slug}/editions/{id}/view", func(re *core.RequestEvent) error {
		return EditionView(re)
	})
	e.Router.GET("/healthz", func(re *core.RequestEvent) error {
		return re.String(200, "ok")
	})
	return nil
}
