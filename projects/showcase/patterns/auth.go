//go:build showcase

package patterns

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.RegisterPattern(&registry.Pattern{
		Slug:        "auth-login",
		Name:        "Login Page",
		Category:    "auth",
		Description: "Centered auth card with email + password + remember me. Composes LoginForm, AuthLayout, and theming.",
		Render: func(theme, mode string) g.Node {
			th := token.Theme(theme)
			mo := token.Mode(mode)
			if th == "" {
				th = token.ThemeSwissBrut
			}
			if mo == "" {
				mo = token.ModeLight
			}
			return fullPage(th, mo,
				layout.AuthLayout(layout.AuthLayoutProps{},
					special.LoginForm(special.LoginFormProps{
						Title:       "Welcome back",
						Description: "Enter your credentials to continue.",
						ForgotHref:  "#",
						SignupHref:  "#",
						SignupLabel: "Create account",
					}),
				),
			)
		},
	})

	registry.RegisterPattern(&registry.Pattern{
		Slug:        "auth-register",
		Name:        "Register Page",
		Category:    "auth",
		Description: "Sign-up form with name, email, password and terms checkbox.",
		Render: func(theme, mode string) g.Node {
			th := token.Theme(theme)
			mo := token.Mode(mode)
			if th == "" {
				th = token.ThemeSwissBrut
			}
			if mo == "" {
				mo = token.ModeLight
			}
			return fullPage(th, mo,
				layout.AuthLayout(layout.AuthLayoutProps{},
					h.Div(
						g.Attr("data-component", "login-form"),
						h.Div(
							g.Attr("data-slot", "header"),
							h.H2(g.Attr("data-slot", "title"), g.Text("Create account")),
							h.P(g.Attr("data-slot", "description"), g.Text("Start building with mljr-ui today.")),
						),
						h.Form(h.Action("#"), h.Method("post"), g.Attr("data-slot", "form"),
							h.Div(
								h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-3)"),
								formField("First name", h.Input(h.Type("text"), h.Name("fname"), h.Placeholder("Jane"),
									g.Attr("data-component", "input"), h.Style("width:100%"))),
								formField("Last name", h.Input(h.Type("text"), h.Name("lname"), h.Placeholder("Smith"),
									g.Attr("data-component", "input"), h.Style("width:100%"))),
							),
							formField("Email", h.Input(h.Type("email"), h.Name("email"), h.Placeholder("you@example.com"),
								g.Attr("data-component", "input"), h.Style("width:100%"))),
							formField("Password", h.Input(h.Type("password"), h.Name("password"), h.Placeholder("••••••••"),
								g.Attr("data-component", "input"), h.Style("width:100%"))),
							h.Div(h.Style("margin-top:var(--sp-2)"),
								h.Label(h.Style("display:flex;align-items:flex-start;gap:var(--sp-2);font-size:var(--t-sm)"),
									h.Input(h.Type("checkbox"), h.Name("terms"), h.Required()),
									g.Text("I agree to the Terms of Service and Privacy Policy"),
								),
							),
							primitive.Button(primitive.ButtonProps{
								Variant: token.Primary,
								Type:    "submit",
								Attrs:   []g.Node{h.Style("width:100%;margin-top:var(--sp-4)")},
							}, g.Text("Create account")),
						),
						h.Div(g.Attr("data-slot", "footer"),
							g.Text("Already have an account? "),
							h.A(h.Href("#"), h.Style("font-weight:700;color:var(--accent)"), g.Text("Sign in")),
						),
					),
				),
			)
		},
	})
}

func formField(label string, control g.Node) g.Node {
	return h.Div(
		h.Style("margin-bottom:var(--sp-3)"),
		h.Label(h.Style("display:block;font-weight:700;font-size:var(--t-sm);margin-bottom:var(--sp-1)"), g.Text(label)),
		control,
	)
}
