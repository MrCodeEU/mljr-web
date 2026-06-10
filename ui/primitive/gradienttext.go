package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GradientTextProps struct {
	// From is the start color (default "var(--accent)").
	From string
	// To is the end color (default "var(--ink)").
	To string
	// Via is an optional middle color.
	Via string
	// Angle is the CSS gradient angle (default "135deg").
	Angle string
	// Tag is the HTML element: "span" (default) | "h1" | "h2" | "h3" | "p"
	Tag string
}

// GradientText renders text with a CSS gradient fill via background-clip:text.
func GradientText(p GradientTextProps, children ...g.Node) g.Node {
	if p.From == "" {
		p.From = "var(--accent)"
	}
	if p.To == "" {
		p.To = "var(--ink)"
	}
	if p.Angle == "" {
		p.Angle = "135deg"
	}

	grad := "linear-gradient(" + p.Angle + "," + p.From
	if p.Via != "" {
		grad += "," + p.Via
	}
	grad += "," + p.To + ")"

	style := "background:" + grad + ";-webkit-background-clip:text;background-clip:text;-webkit-text-fill-color:transparent;text-fill-color:transparent;display:inline-block"

	attrs := append([]g.Node{
		g.Attr("data-component", "gradient-text"),
		h.Style(style),
	}, children...)

	switch p.Tag {
	case "h1":
		return h.H1(attrs...)
	case "h2":
		return h.H2(attrs...)
	case "h3":
		return h.H3(attrs...)
	case "p":
		return h.P(attrs...)
	default:
		return h.Span(attrs...)
	}
}
