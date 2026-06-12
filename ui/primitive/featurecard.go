package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FeatureCardProps struct {
	Icon  string // lucide or simple-icons name
	Title string
	Tone  token.Tone // accent color for icon box (default Primary)
	Href  string     // if set, whole card is a link
}

// FeatureCard renders a card with colored icon box, title, and description slot.
func FeatureCard(p FeatureCardProps, description ...g.Node) g.Node {
	iconBoxStyle := "width:2.5rem;height:2.5rem;border-radius:var(--radius);display:flex;align-items:center;justify-content:center;background:var(--accent);color:var(--accent-ink);flex-shrink:0"
	if p.Tone != token.ToneNone && p.Tone != "" {
		iconBoxStyle = "width:2.5rem;height:2.5rem;border-radius:var(--radius);display:flex;align-items:center;justify-content:center;background:var(--tone-bg);color:var(--tone-ink);flex-shrink:0"
	}

	inner := h.Div(
		g.Attr("data-component", "feature-card"),
		g.If(p.Tone != token.ToneNone && p.Tone != "", g.Attr("data-tone", string(p.Tone))),
		h.Div(h.Style(iconBoxStyle),
			g.If(p.Icon != "", icon.Icon(p.Icon, icon.Props{Size: "1.25rem"})),
		),
		h.Div(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-1)"),
			h.Strong(h.Style("font-weight:700;font-size:var(--t-base)"), g.Text(p.Title)),
			g.Group(description),
		),
	)

	if p.Href != "" {
		return h.A(
			h.Href(p.Href),
			g.Attr("data-component", "card"),
			g.Attr("data-interactive", ""),
			inner,
		)
	}
	return Card(CardProps{}, inner)
}

// PricingCard renders a pricing tier card.
type PricingCardProps struct {
	Name        string
	Price       string // e.g. "€29"
	Period      string // e.g. "/month"
	Description string
	Features    []string
	CTA         string
	CTAHref     string
	Highlighted bool // primary accent styling
}

func PricingCard(p PricingCardProps) g.Node {
	featureNodes := make([]g.Node, len(p.Features))
	for i, f := range p.Features {
		featureNodes[i] = h.Li(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);font-size:var(--t-sm)"),
			icon.Icon("lucide:check", icon.Props{Size: "1rem"}),
			g.Text(f),
		)
	}

	variant := token.Outline
	if p.Highlighted {
		variant = token.Primary
	}

	attrs := []g.Node{}
	if p.Highlighted {
		attrs = append(attrs, g.Attr("data-tone", "sky"))
	}

	return Card(CardProps{Attrs: attrs},
		h.Div(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-4);height:100%"),
			h.Div(
				h.Strong(h.Style("font-size:var(--t-base);font-weight:700"), g.Text(p.Name)),
				g.If(p.Description != "", h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin-top:var(--sp-1)"), g.Text(p.Description))),
			),
			h.Div(
				h.Style("display:flex;align-items:baseline;gap:var(--sp-1)"),
				h.Span(h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display);letter-spacing:-.02em"), g.Text(p.Price)),
				g.If(p.Period != "", h.Span(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text(p.Period))),
			),
			h.Ul(h.Style("list-style:none;padding:0;margin:0;display:flex;flex-direction:column;gap:var(--sp-2);flex:1"),
				g.Group(featureNodes),
			),
			g.If(p.CTA != "", h.A(
				h.Href(p.CTAHref),
				g.Attr("data-component", "button"),
				g.Attr("data-variant", string(variant)),
				h.Style("width:100%;justify-content:center"),
				g.Text(p.CTA),
			)),
		),
	)
}
