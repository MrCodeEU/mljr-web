// Package feedback provides Alert, Spinner, and Skeleton components.
package feedback

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AlertVariant string

const (
	AlertInfo    AlertVariant = "info"
	AlertSuccess AlertVariant = "success"
	AlertWarning AlertVariant = "warning"
	AlertDanger  AlertVariant = "danger"
)

var alertIcons = map[AlertVariant]string{
	AlertInfo:    "lucide:info",
	AlertSuccess: "lucide:circle-check",
	AlertWarning: "lucide:alert-triangle",
	AlertDanger:  "lucide:circle-x",
}

type AlertProps struct {
	Variant AlertVariant
	Title   string
	Dismiss bool // show × dismiss button
	Attrs   []g.Node
}

// Alert renders a contextual banner. Dismiss button removes it from the DOM.
func Alert(p AlertProps, children ...g.Node) g.Node {
	if p.Variant == "" {
		p.Variant = AlertInfo
	}
	iconName := alertIcons[p.Variant]

	var dismissBtn g.Node
	if p.Dismiss {
		dismissBtn = h.Button(
			g.Attr("data-slot", "dismiss"),
			g.Attr("aria-label", "Dismiss"),
			g.Attr("data-on:click", "evt.target.closest('[data-component=alert]').remove()"),
			g.Text("×"),
		)
	}

	return h.Div(
		g.Attr("data-component", "alert"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("role", "alert"),
		g.Group(p.Attrs),
		h.Div(g.Attr("data-slot", "icon"), icon.Icon(iconName)),
		h.Div(
			g.Attr("data-slot", "content"),
			g.If(p.Title != "", h.Div(g.Attr("data-slot", "title"), g.Text(p.Title))),
			g.Group(children),
		),
		dismissBtn,
	)
}
