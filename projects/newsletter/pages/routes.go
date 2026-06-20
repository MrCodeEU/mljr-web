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
	e.Router.POST("/login", wrapErrors(HandleLogin))
	e.Router.GET("/signup", func(re *core.RequestEvent) error {
		return web.RenderPB(re, 200, Signup(re, SignupProps{}))
	})
	e.Router.POST("/signup", wrapErrors(HandleSignup))
	e.Router.POST("/logout", wrapErrors(HandleLogout))
	e.Router.GET("/profile", wrapErrors(Profile))
	e.Router.POST("/profile", wrapErrors(HandleProfile))
	e.Router.POST("/profile/avatar", wrapErrors(HandleProfileAvatar))
	e.Router.GET("/", wrapErrors(Dashboard))
	e.Router.GET("/g/{slug}", wrapErrors(GroupHome))
	e.Router.GET("/g/{slug}/settings", wrapErrors(GroupSettings))
	e.Router.POST("/g/{slug}/settings", wrapErrors(HandleGroupSettings))
	e.Router.POST("/groups", wrapErrors(HandleCreateGroup))
	e.Router.POST("/g/{slug}/leave", wrapErrors(HandleLeaveGroup))
	e.Router.POST("/g/{slug}/reactivate", wrapErrors(HandleReactivateGroup))
	e.Router.GET("/g/{slug}/editions/{id}/questions", wrapErrors(EditionQuestions))
	e.Router.POST("/g/{slug}/editions/{id}/questions", wrapErrors(HandleEditionQuestions))
	e.Router.GET("/g/{slug}/invites", wrapErrors(ListInvites))
	e.Router.POST("/g/{slug}/invites", wrapErrors(HandleCreateInvite))
	e.Router.POST("/g/{slug}/invites/{id}/revoke", wrapErrors(HandleRevokeInvite))
	e.Router.GET("/invites/{token}", wrapErrors(InviteAccept))
	e.Router.POST("/invites/{token}/accept", wrapErrors(HandleAcceptInvite))
	e.Router.POST("/notifications/read-all", wrapErrors(HandleMarkAllNotificationsRead))
	e.Router.GET("/g/{slug}/questions", wrapErrors(ListQuestions))
	e.Router.POST("/g/{slug}/questions", wrapErrors(HandleCreateQuestion))
	e.Router.POST("/g/{slug}/questions/{id}/toggle", wrapErrors(HandleToggleQuestion))
	e.Router.GET("/g/{slug}/editions", wrapErrors(ListEditions))
	e.Router.POST("/g/{slug}/editions", wrapErrors(HandleCreateEdition))
	e.Router.GET("/g/{slug}/editions/{id}", wrapErrors(EditionAnswer))
	e.Router.POST("/g/{slug}/editions/{id}", wrapErrors(HandleSubmitAnswers))
	e.Router.POST("/g/{slug}/editions/{id}/close", wrapErrors(HandleCloseEdition))
	e.Router.GET("/g/{slug}/editions/{id}/view", wrapErrors(EditionView))
	e.Router.POST("/g/{slug}/editions/{id}/answers/{answerID}/react", wrapErrors(HandleToggleReaction))
	e.Router.POST("/g/{slug}/editions/{id}/answers/{answerID}/comments", wrapErrors(HandleCreateComment))
	e.Router.GET("/g/{slug}/suggestions", wrapErrors(ListSuggestions))
	e.Router.POST("/g/{slug}/suggestions", wrapErrors(HandleCreateSuggestion))
	e.Router.POST("/g/{slug}/suggestions/{id}/vote", wrapErrors(HandleToggleVote))
	e.Router.POST("/g/{slug}/suggestions/{id}/approve", wrapErrors(HandleApproveSuggestion))
	e.Router.POST("/g/{slug}/suggestions/{id}/reject", wrapErrors(HandleRejectSuggestion))
	e.Router.GET("/g/{slug}/recap", wrapErrors(Recap))
	e.Router.GET("/healthz", func(re *core.RequestEvent) error {
		return re.String(200, "ok")
	})
	return nil
}
