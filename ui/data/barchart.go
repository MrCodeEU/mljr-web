package data

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BarDatum struct {
	Label string
	Value float64
	Color string // CSS color; defaults to var(--primary)
}

type BarChartProps struct {
	Data        []BarDatum
	Height      int    // px height of chart area (default 160)
	ShowValues  bool   // render value labels above bars
	ShowGrid    bool   // render horizontal grid lines
	Caption     string // optional <figcaption>
}

// BarChart renders a pure-SVG bar chart. No JS required.
func BarChart(p BarChartProps) g.Node {
	if p.Height == 0 {
		p.Height = 160
	}
	if len(p.Data) == 0 {
		return h.Div(g.Attr("data-component", "bar-chart"), g.Text("no data"))
	}

	const labelH = 24  // px reserved for x-axis labels
	const padL = 8
	const padR = 8
	const barGapRatio = 0.3 // gap as fraction of bar+gap width

	n := len(p.Data)
	maxVal := 0.0
	for _, d := range p.Data {
		if d.Value > maxVal {
			maxVal = d.Value
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	totalW := 400
	slotW := float64(totalW-padL-padR) / float64(n)
	barW := slotW * (1 - barGapRatio)
	gap := slotW * barGapRatio

	svgH := p.Height + labelH
	chartH := float64(p.Height)

	nodes := []g.Node{}

	// Grid lines
	if p.ShowGrid {
		for _, frac := range []float64{0.25, 0.5, 0.75, 1.0} {
			y := chartH * (1 - frac)
			nodes = append(nodes, g.El("line",
				g.Attr("x1", "0"), g.Attr("y1", fmt.Sprintf("%.1f", y)),
				g.Attr("x2", fmt.Sprintf("%d", totalW)), g.Attr("y2", fmt.Sprintf("%.1f", y)),
				g.Attr("stroke", "var(--line)"), g.Attr("stroke-dasharray", "4 4"),
			))
		}
	}

	// Bars and labels
	for i, d := range p.Data {
		x := float64(padL) + float64(i)*slotW + gap/2
		barH := math.Max(2, (d.Value/maxVal)*chartH)
		y := chartH - barH
		color := d.Color
		if color == "" {
			color = "var(--primary)"
		}

		nodes = append(nodes,
			// Bar rect
			g.El("rect",
				g.Attr("data-component", "bar"),
				g.Attr("x", fmt.Sprintf("%.1f", x)),
				g.Attr("y", fmt.Sprintf("%.1f", y)),
				g.Attr("width", fmt.Sprintf("%.1f", barW)),
				g.Attr("height", fmt.Sprintf("%.1f", barH)),
				g.Attr("fill", color),
				g.Attr("rx", "2"),
			),
		)

		// Value label above bar
		if p.ShowValues {
			nodes = append(nodes, g.El("text",
				g.Attr("x", fmt.Sprintf("%.1f", x+barW/2)),
				g.Attr("y", fmt.Sprintf("%.1f", y-4)),
				g.Attr("text-anchor", "middle"),
				g.Attr("font-size", "10"),
				g.Attr("fill", "var(--fg)"),
				g.Attr("opacity", "0.7"),
				g.Raw(formatVal(d.Value)),
			))
		}

		// X-axis label
		nodes = append(nodes, g.El("text",
			g.Attr("x", fmt.Sprintf("%.1f", x+barW/2)),
			g.Attr("y", fmt.Sprintf("%d", p.Height+labelH-4)),
			g.Attr("text-anchor", "middle"),
			g.Attr("font-size", "10"),
			g.Attr("fill", "var(--muted)"),
			g.Raw(d.Label),
		))
	}

	svgNode := g.El("svg",
		g.Attr("data-component", "bar-chart"),
		g.Attr("viewBox", fmt.Sprintf("0 0 %d %d", totalW, svgH)),
		g.Attr("style", fmt.Sprintf("width:100%%;height:%dpx;overflow:visible", p.Height+labelH)),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(nodes),
	)

	if p.Caption != "" {
		return g.El("figure",
			g.Attr("style", "margin:0"),
			svgNode,
			g.El("figcaption",
				g.Attr("style", "font-size:var(--t-xs);color:var(--muted);text-align:center;margin-top:var(--sp-2)"),
				g.Text(p.Caption),
			),
		)
	}
	return svgNode
}

func formatVal(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%.0f", v)
	}
	return fmt.Sprintf("%.1f", v)
}
