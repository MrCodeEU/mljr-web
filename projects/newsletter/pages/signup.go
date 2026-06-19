package pages

import (
	"errors"
	"strings"

	"mljr-web/ui/feedback"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SignupProps struct {
	Error       string
	InviteToken string
}

func Signup(re *core.RequestEvent, p SignupProps) g.Node {
	if p.InviteToken == "" {
		p.InviteToken = re.Request.URL.Query().Get("invite")
	}

	var hidden []g.Node
	var emailValue string
	var emailReadOnly bool
	var banner g.Node
	if p.InviteToken != "" {
		hidden = []g.Node{h.Input(h.Type("hidden"), h.Name("invite"), h.Value(p.InviteToken))}
		if invite, err := findInviteByToken(re, p.InviteToken); err == nil && invite.GetString("status") == "pending" {
			emailValue = invite.GetString("email")
			emailReadOnly = true
			if group, err := re.App.FindRecordById("groups", invite.GetString("group")); err == nil {
				banner = feedback.Alert(
					feedback.AlertProps{Variant: feedback.AlertInfo, Attrs: []g.Node{h.Style("margin-bottom:var(--sp-4)")}},
					g.Text("You're creating an account to join "+group.GetString("name")+"."),
				)
			}
		}
	}

	return authPage(re, "Sign up",
		layout.AuthLayout(layout.AuthLayoutProps{}, h.Div(
			banner,
			special.LoginForm(special.LoginFormProps{
				Action:               "/signup",
				Title:                "Create your account",
				LoginHref:            "/login",
				Error:                p.Error,
				SubmitLabel:          "Sign up",
				PasswordAutocomplete: "new-password",
				PasswordMinLength:    8,
				EmailValue:           emailValue,
				EmailReadOnly:        emailReadOnly,
				HiddenFields:         hidden,
			}),
		)),
	)
}

var fieldLabels = map[string]string{
	"email":    "Email",
	"password": "Password",
	"name":     "Name",
}

// signupErrorMessage turns a record-save error into field-specific,
// human-readable text instead of a generic "could not create account".
func signupErrorMessage(err error) string {
	var verrs validation.Errors
	if !errors.As(err, &verrs) {
		return "Could not create account — please try again"
	}

	var parts []string
	for _, field := range []string{"email", "password", "name"} {
		if fieldErr, ok := verrs[field]; ok {
			label := fieldLabels[field]
			if label == "" {
				label = field
			}
			parts = append(parts, label+": "+fieldErr.Error())
		}
	}
	if len(parts) == 0 {
		return "Could not create account — please try again"
	}
	return strings.Join(parts, ". ")
}

func HandleSignup(re *core.RequestEvent) error {
	email := re.Request.FormValue("email")
	password := re.Request.FormValue("password")
	name := re.Request.FormValue("name")
	inviteToken := re.Request.FormValue("invite")

	if email == "" || password == "" {
		return renderPage(re, 400, Signup(re, SignupProps{Error: "Email and password are required", InviteToken: inviteToken}))
	}

	users, err := re.App.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	user := core.NewRecord(users)
	user.SetEmail(email)
	user.SetPassword(password)
	user.SetVerified(true)
	user.Set("name", name)

	if err := re.App.Save(user); err != nil {
		return renderPage(re, 400, Signup(re, SignupProps{Error: signupErrorMessage(err), InviteToken: inviteToken}))
	}

	if err := setSession(re, user); err != nil {
		return err
	}

	if inviteToken != "" {
		if invite, err := findInviteByToken(re, inviteToken); err == nil && invite.GetString("status") == "pending" && inviteTargetsUser(invite, user) {
			if group, err := acceptInvite(re, invite, user); err == nil {
				return redirect(re, "/g/"+group.GetString("slug"))
			}
		}
	}

	return redirect(re, "/")
}
