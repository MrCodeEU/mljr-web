package form

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PasswordInputProps struct {
	Name        string
	ID          string
	Placeholder string
	Value       string
	Autocomplete string // e.g. "current-password", "new-password"
}

// PasswordInput renders a password field with a show/hide toggle button.
// Uses Datastar signal $_pwvN (unique per Name) to track visibility state.
func PasswordInput(p PasswordInputProps) g.Node {
	if p.Placeholder == "" {
		p.Placeholder = "Password"
	}
	if p.Autocomplete == "" {
		p.Autocomplete = "current-password"
	}
	sig := "_pwv"
	if p.Name != "" {
		for _, c := range p.Name {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
				sig += string(c)
			}
		}
	}

	return h.Div(
		g.Attr("data-component", "password-input"),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		h.Input(
			g.Attr("data-component", "input"),
			g.Attr("data-attr", `{"type":$`+sig+`?"text":"password"}`),
			h.Type("password"),
			h.Name(p.Name),
			g.If(p.ID != "", h.ID(p.ID)),
			h.Placeholder(p.Placeholder),
			h.Value(p.Value),
			g.Attr("autocomplete", p.Autocomplete),
		),
		h.Button(
			g.Attr("data-slot", "toggle"),
			h.Type("button"),
			g.Attr("aria-label", "Toggle password visibility"),
			g.Attr("data-on:click", `$`+sig+`=!$`+sig),
			h.Span(g.Attr("data-show", "!$"+sig), icon.Icon("lucide:eye")),
			h.Span(g.Attr("data-show", "$"+sig), h.Style("display:none"), icon.Icon("lucide:eye-off")),
		),
	)
}
