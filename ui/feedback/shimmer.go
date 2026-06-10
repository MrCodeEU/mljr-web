package feedback

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ShimmerProps struct {
	// Width: CSS width (default "100%").
	Width string
	// Height: CSS height (default "1em").
	Height string
	// Radius: CSS border-radius (default "var(--radius)").
	Radius string
	// Lines renders N stacked shimmer bars (0 = single bar).
	Lines int
	// Circle renders a circular shimmer (overrides Radius).
	Circle bool
}

// Shimmer renders an animated gradient loading placeholder.
// Lighter-weight alternative to Skeleton — just the shimmer gradient, no layout.
func Shimmer(p ShimmerProps) g.Node {
	if p.Width == "" {
		p.Width = "100%"
	}
	if p.Height == "" {
		p.Height = "1em"
	}
	if p.Radius == "" {
		p.Radius = "var(--radius)"
	}
	if p.Circle {
		p.Radius = "50%"
	}

	single := func(w, h2 string) g.Node {
		return h.Div(
			g.Attr("data-component", "shimmer"),
			h.Style(fmt.Sprintf("width:%s;height:%s;border-radius:%s", w, h2, p.Radius)),
		)
	}

	if p.Lines <= 1 {
		return single(p.Width, p.Height)
	}

	bars := make([]g.Node, p.Lines)
	for i := range bars {
		w := p.Width
		if i == p.Lines-1 && p.Lines > 1 {
			w = "60%" // last line shorter
		}
		bars[i] = single(w, p.Height)
	}
	return h.Div(
		h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"),
		g.Group(bars),
	)
}
