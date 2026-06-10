package feedback

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SpinnerVariant string

const (
	SpinnerRing  SpinnerVariant = ""      // default: spinning ring
	SpinnerDots  SpinnerVariant = "dots"  // three bouncing dots
	SpinnerPulse SpinnerVariant = "pulse" // pulsing circle
	SpinnerBars  SpinnerVariant = "bars"  // four scaling bars
	SpinnerSwiss SpinnerVariant = "swiss" // stepped square (Swiss Brut)
	SpinnerInk   SpinnerVariant = "ink"   // organic blob (Ink/Paper)
)

type SpinnerProps struct {
	Variant SpinnerVariant
	Size    string // sm | md (default) | lg
	Label   string // accessible label (default "Loading…")
	Attrs   []g.Node
}

// Spinner renders a CSS-animated loading indicator.
func Spinner(p SpinnerProps, attrs ...g.Node) g.Node {
	if p.Size == "" {
		p.Size = "md"
	}
	label := p.Label
	if label == "" {
		label = "Loading…"
	}

	nodes := []g.Node{
		g.Attr("data-component", "spinner"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Attr("data-size", p.Size),
		g.Attr("role", "status"),
		g.Attr("aria-label", label),
		g.Group(p.Attrs),
		g.Group(attrs),
	}
	switch p.Variant {
	case SpinnerDots:
		nodes = append(nodes, h.Span(), h.Span(), h.Span())
	case SpinnerBars:
		nodes = append(nodes, h.Span(), h.Span(), h.Span(), h.Span())
	}
	return h.Div(nodes...)
}
