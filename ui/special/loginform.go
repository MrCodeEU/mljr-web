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
	// LoginHref links back to the sign-in page (for signup forms).
	LoginHref string
	// LoginLabel text (default "Sign in instead").
	LoginLabel string
	// SubmitLabel is the button text (default "Sign in").
	SubmitLabel string
	// Error message to show (e.g. "Invalid credentials").
	Error string
	// PasswordAutocomplete hints the browser whether this is a new password
	// (signup, default "new-password") or an existing one — set explicitly
	// to "current-password" for sign-in forms.
	PasswordAutocomplete string
	// PasswordMinLength enables native browser validation (e.g. 8 to match
	// the backend's minimum). 0 disables the check.
	PasswordMinLength int
	// EmailValue pre-fills the email field (e.g. from an invite link).
	EmailValue string
	// EmailReadOnly locks the email field when EmailValue is set, so an
	// invited address can't be swapped on submit.
	EmailReadOnly bool
	// HiddenFields are rendered as hidden inputs inside the form (e.g. to
	// carry an invite token through the POST).
	HiddenFields []g.Node
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
	if p.LoginLabel == "" {
		p.LoginLabel = "Sign in instead"
	}
	if p.PasswordAutocomplete == "" {
		p.PasswordAutocomplete = "current-password"
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
			g.Group(p.HiddenFields),
			// Email
			form.Field(form.FieldProps{Label: "Email"},
				form.Input(form.InputProps{
					Type:        "email",
					Name:        "email",
					Placeholder: "you@example.com",
					Required:    true,
					Attrs: []g.Node{
						h.ID("lf-email"),
						g.If(p.EmailValue != "", h.Value(p.EmailValue)),
						g.If(p.EmailReadOnly, h.ReadOnly()),
					},
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
				Name:         "password",
				ID:           "lf-pass",
				Placeholder:  "••••••••",
				Autocomplete: p.PasswordAutocomplete,
				Required:     true,
				MinLength:    p.PasswordMinLength,
			}),
			g.If(p.PasswordMinLength > 0, h.Span(
				h.Style("font-size:var(--t-xs);color:var(--muted)"),
				g.Textf("At least %d characters", p.PasswordMinLength),
			)),
			// Remember me
			h.Div(
				h.Style("margin-top:var(--sp-3)"),
				form.Checkbox(form.CheckboxProps{Name: "remember", Label: "Remember me"}),
			),
			// Submit
			primitive.Button(primitive.ButtonProps{
				Variant: token.Primary,
				Type:    "submit",
				Attrs: []g.Node{
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
		g.If(p.LoginHref != "", h.Div(
			g.Attr("data-slot", "footer"),
			g.Text("Already have an account? "),
			h.A(h.Href(p.LoginHref), h.Style("font-weight:700;color:var(--accent)"), g.Text(p.LoginLabel)),
		)),
	)
}
