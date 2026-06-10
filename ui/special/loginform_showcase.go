//go:build showcase

package special

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "login-form", Name: "Login Form", Category: "special",
		Summary: "Composite sign-in form: email + password (show/hide) + remember me + submit. Composed from form.Field, form.Input, form.PasswordInput, primitive.Button.",
		Code: `special.LoginForm(special.LoginFormProps{
    Title:       "Welcome back",
    Description: "Sign in to your account to continue.",
    ForgotHref:  "/forgot-password",
    SignupHref:  "/register",
    SubmitLabel: "Sign in",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(auto-fit,minmax(320px,1fr));gap:var(--sp-8)"),
				// Normal
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-3);color:var(--muted)"), g.Text("Default:")),
					LoginForm(LoginFormProps{
						Title:       "Welcome back",
						Description: "Sign in to your account.",
						ForgotHref:  "#",
						SignupHref:  "#",
					}),
				),
				// With error
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-3);color:var(--muted)"), g.Text("With error:")),
					LoginForm(LoginFormProps{
						Title:       "Welcome back",
						ForgotHref:  "#",
						SignupHref:  "#",
						SubmitLabel: "Sign in",
						Error:       "Invalid email or password. Please try again.",
					}),
				),
			)
		},
	})
}
