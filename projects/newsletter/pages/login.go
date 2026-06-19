package pages

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/special"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type LoginProps struct {
	Error       string
	InviteToken string
}

func Login(re *core.RequestEvent, p LoginProps) g.Node {
	if p.InviteToken == "" {
		p.InviteToken = re.Request.URL.Query().Get("invite")
	}

	var hidden []g.Node
	if p.InviteToken != "" {
		hidden = []g.Node{h.Input(h.Type("hidden"), h.Name("invite"), h.Value(p.InviteToken))}
	}

	return authPage(re, "Sign in",
		layout.AuthLayout(layout.AuthLayoutProps{}, special.LoginForm(special.LoginFormProps{
			Action:       "/login",
			SignupHref:   "/signup",
			Error:        p.Error,
			SubmitLabel:  "Sign in",
			HiddenFields: hidden,
		})),
	)
}

func HandleLogin(re *core.RequestEvent) error {
	email := re.Request.FormValue("email")
	password := re.Request.FormValue("password")
	inviteToken := re.Request.FormValue("invite")

	user, err := re.App.FindAuthRecordByEmail("users", email)
	if err != nil || !user.ValidatePassword(password) {
		return renderPage(re, 401, Login(re, LoginProps{Error: "Invalid email or password", InviteToken: inviteToken}))
	}

	if err := setSession(re, user); err != nil {
		return err
	}
	if inviteToken != "" {
		return redirect(re, "/invites/"+inviteToken)
	}
	return redirect(re, "/")
}
