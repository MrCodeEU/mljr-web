package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CTABannerProps struct {
	Title       string
	Description string
	CTAText     string
	CTAHref     string
	Variant     token.Variant // default Primary
	SecondCTA   string        // optional second CTA text (ghost variant)
	SecondHref  string
}

// CTABanner renders a full-width call-to-action strip.
func CTABanner(p CTABannerProps) g.Node {
	if p.Variant == "" {
		p.Variant = token.Primary
	}

	var ctas []g.Node
	if p.CTAText != "" {
		ctas = append(ctas, h.A(
			h.Href(p.CTAHref),
			g.Attr("data-component", "button"),
			g.Attr("data-variant", string(p.Variant)),
			g.Text(p.CTAText),
		))
	}
	if p.SecondCTA != "" {
		ctas = append(ctas, h.A(
			h.Href(p.SecondHref),
			g.Attr("data-component", "button"),
			g.Attr("data-variant", "outline"),
			g.Text(p.SecondCTA),
		))
	}

	return h.Div(
		g.Attr("data-component", "cta-banner"),
		h.Div(
			g.Attr("data-slot", "content"),
			h.Div(
				g.Attr("data-slot", "text"),
				h.Strong(h.Style("font-size:var(--t-xl);font-weight:900;font-family:var(--font-display)"), g.Text(p.Title)),
				g.If(p.Description != "", h.P(h.Style("color:var(--muted);margin:var(--sp-1) 0 0"), g.Text(p.Description))),
			),
			g.If(len(ctas) > 0, h.Div(
				g.Attr("data-slot", "actions"),
				g.Group(ctas),
			)),
		),
	)
}
