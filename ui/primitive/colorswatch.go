package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ColorSwatchProps struct {
	Color   string // CSS color value e.g. "#f28d1d" or "var(--primary)"
	Label   string // optional label below the swatch
	Size    string // CSS size default "2.5rem"
	Shape   string // "circle" or "square" (default "square")
	Outline bool   // show border
}

// ColorSwatch renders a single color preview square/circle.
func ColorSwatch(p ColorSwatchProps) g.Node {
	if p.Size == "" {
		p.Size = "2.5rem"
	}
	radius := "var(--radius)"
	if p.Shape == "circle" {
		radius = "50%"
	}
	borderStyle := ""
	if p.Outline {
		borderStyle = "border:var(--bw-1) solid var(--line);"
	}
	style := "width:" + p.Size + ";height:" + p.Size + ";border-radius:" + radius +
		";background:" + p.Color + ";" + borderStyle + "flex-shrink:0"

	swatch := h.Span(
		h.Style(style),
		h.Title(p.Color),
	)
	if p.Label == "" {
		return swatch
	}
	return h.Div(
		g.Attr("data-component", "color-swatch"),
		h.Style("display:inline-flex;flex-direction:column;align-items:center;gap:var(--sp-1)"),
		swatch,
		h.Span(
			h.Style("font-size:var(--t-xs);font-family:var(--font-mono);color:var(--muted)"),
			g.Text(p.Label),
		),
	)
}

// ColorSwatchGroup renders a row of swatches.
func ColorSwatchGroup(swatches ...g.Node) g.Node {
	return h.Div(
		h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
		g.Group(swatches),
	)
}
