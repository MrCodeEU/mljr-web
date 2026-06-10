package overlay

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type HoverCardProps struct {
	// Placement: "top" (default) | "bottom" | "left" | "right"
	Placement string
	// Width is the CSS width of the card (default "280px").
	Width string
	// Delay is the CSS transition-delay for showing (default "0.1s").
	Delay string
}

// HoverCard shows a rich content card when the trigger element is hovered.
// Pure CSS — no JS, no Datastar required.
// Usage: HoverCard(props, trigger, content)
func HoverCard(p HoverCardProps, trigger g.Node, content g.Node) g.Node {
	if p.Placement == "" {
		p.Placement = "top"
	}
	if p.Width == "" {
		p.Width = "280px"
	}
	if p.Delay == "" {
		p.Delay = "0.1s"
	}

	cardStyle := "width:" + p.Width + ";transition-delay:" + p.Delay
	return h.Div(
		g.Attr("data-component", "hover-card"),
		g.Attr("data-placement", p.Placement),
		trigger,
		h.Div(
			g.Attr("data-slot", "content"),
			h.Style(cardStyle),
			content,
		),
	)
}
