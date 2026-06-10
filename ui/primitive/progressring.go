package primitive

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ProgressRingProps struct {
	Value     int    // 0–100
	Size      int    // px, default 64
	Thickness int    // stroke px, default 6
	Label     string // center label (e.g. "72%")
	Variant   ProgressVariant
	Attrs     []g.Node
}

// ProgressRing renders an SVG circular progress indicator.
func ProgressRing(p ProgressRingProps) g.Node {
	size := p.Size
	if size <= 0 {
		size = 64
	}
	thickness := p.Thickness
	if thickness <= 0 {
		thickness = 6
	}
	value := p.Value
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}

	cx := float64(size) / 2
	r := cx - float64(thickness)/2
	circumference := 2 * math.Pi * r
	offset := circumference * (1.0 - float64(value)/100.0)

	sz := fmt.Sprintf("%d", size)
	cxs := fmt.Sprintf("%.1f", cx)
	rs := fmt.Sprintf("%.1f", r)
	tw := fmt.Sprintf("%d", thickness)

	var label g.Node
	if p.Label != "" {
		label = g.El("text",
			g.Attr("x", cxs),
			g.Attr("y", cxs),
			g.Attr("text-anchor", "middle"),
			g.Attr("dominant-baseline", "middle"),
			g.Attr("data-slot", "label"),
			g.Text(p.Label),
		)
	}

	return h.Div(
		g.Attr("data-component", "progress-ring"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Group(p.Attrs),
		g.El("svg",
			g.Attr("width", sz),
			g.Attr("height", sz),
			g.Attr("viewBox", fmt.Sprintf("0 0 %s %s", sz, sz)),
			g.Attr("aria-valuenow", fmt.Sprintf("%d", value)),
			g.Attr("role", "progressbar"),
			g.El("circle",
				g.Attr("cx", cxs), g.Attr("cy", cxs), g.Attr("r", rs),
				g.Attr("fill", "none"),
				g.Attr("stroke-width", tw),
				g.Attr("data-slot", "track"),
			),
			g.El("circle",
				g.Attr("cx", cxs), g.Attr("cy", cxs), g.Attr("r", rs),
				g.Attr("fill", "none"),
				g.Attr("stroke-width", tw),
				g.Attr("stroke-dasharray", fmt.Sprintf("%.2f", circumference)),
				g.Attr("stroke-dashoffset", fmt.Sprintf("%.2f", offset)),
				g.Attr("stroke-linecap", "round"),
				g.Attr("transform", fmt.Sprintf("rotate(-90 %s %s)", cxs, cxs)),
				g.Attr("data-slot", "fill"),
			),
			label,
		),
	)
}
