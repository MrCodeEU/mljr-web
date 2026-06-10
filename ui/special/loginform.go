package special

import (
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type LoginFormProps struct {
	// Action is the form POST URL (default "#").
	Action string
	// Title shown at top (default "Welcome back").
	Title string
	// Description below title.
	Description string
	// ForgotHref links to forgot-password page.
	ForgotHref string
	// SignupHref links to the register page.
	SignupHref string
	// SignupLabel text (default "Create account").
	SignupLabel string
	// SubmitLabel is the button text (default "Sign in").
	SubmitLabel string
	// Error message to show (e.g. "Invalid credentials").
	Error string
}

// LoginForm renders a complete sign-in form.
// Composite of form.Field, form.Input, form.PasswordInput, primitive.Button.
func LoginForm(p LoginFormProps) g.Node {
	if p.Action == "" {
		p.Action = "#"
	}
	if p.Title == "" {
		p.Title = "Welcome back"
	}
	if p.SubmitLabel == "" {
		p.SubmitLabel = "Sign in"
	}
	if p.SignupLabel == "" {
		p.SignupLabel = "Create account"
	}

	return h.Div(
		g.Attr("data-component", "login-form"),
		// Header
		h.Div(
			g.Attr("data-slot", "header"),
			h.H2(g.Attr("data-slot", "title"), g.Text(p.Title)),
			g.If(p.Description != "", h.P(g.Attr("data-slot", "description"), g.Text(p.Description))),
		),
		// Error
		g.If(p.Error != "", h.Div(
			g.Attr("data-component", "alert"),
			g.Attr("data-variant", "danger"),
			h.Style("margin-bottom:var(--sp-4)"),
			g.Text(p.Error),
		)),
		// Form
		h.Form(
			h.Action(p.Action),
			h.Method("post"),
			g.Attr("data-slot", "form"),
			// Email
			form.Field(form.FieldProps{Label: "Email"},
				form.Input(form.InputProps{
					Type:        "email",
					Name:        "email",
					Placeholder: "you@example.com",
					Required:    true,
					Attrs:       []g.Node{h.ID("lf-email")},
				}),
			),
			// Password row: label + forgot link
			h.Div(
				h.Style("display:flex;align-items:baseline;justify-content:space-between;margin-bottom:var(--sp-1)"),
				h.Label(h.For("lf-pass"), g.Text("Password")),
				g.If(p.ForgotHref != "", h.A(
					h.Href(p.ForgotHref),
					h.Style("font-size:var(--t-sm);color:var(--accent);text-decoration:none;font-weight:600"),
					g.Text("Forgot password?"),
				)),
			),
			form.PasswordInput(form.PasswordInputProps{
				Name:        "password",
				ID:          "lf-pass",
				Placeholder: "••••••••",
			}),
			// Remember me
			h.Div(
				h.Style("margin-top:var(--sp-3)"),
				form.Checkbox(form.CheckboxProps{Name: "remember", Label: "Remember me"}),
			),
			// Submit
			primitive.Button(primitive.ButtonProps{
				Variant: token.Primary,
				Attrs: []g.Node{
					h.Type("submit"),
					h.Style("width:100%;margin-top:var(--sp-4)"),
				},
			}, g.Text(p.SubmitLabel)),
		),
		// Footer
		g.If(p.SignupHref != "", h.Div(
			g.Attr("data-slot", "footer"),
			g.Text("Don't have an account? "),
			h.A(h.Href(p.SignupHref), h.Style("font-weight:700;color:var(--accent)"), g.Text(p.SignupLabel)),
		)),
	)
}
