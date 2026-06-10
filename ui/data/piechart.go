package data

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PieChartProps struct {
	Slices  []DonutSlice // reuse DonutSlice type
	Size    int          // px diameter (default 180)
	Caption string
}

// PieChart renders a solid SVG pie chart (no center hole). Uses path arcs.
func PieChart(p PieChartProps) g.Node {
	if p.Size == 0 {
		p.Size = 180
	}
	if len(p.Slices) == 0 {
		return g.El("svg")
	}

	total := 0.0
	for _, s := range p.Slices {
		total += s.Value
	}
	if total == 0 {
		total = 1
	}

	cx := float64(p.Size) / 2
	cy := float64(p.Size) / 2
	r := cx - 2
	defaultColors := []string{"var(--primary)", "var(--accent)", "var(--success)", "var(--warning)", "var(--danger)"}

	nodes := []g.Node{}
	startAngle := -math.Pi / 2 // 12 o'clock

	for i, s := range p.Slices {
		color := s.Color
		if color == "" {
			color = defaultColors[i%len(defaultColors)]
		}
		angle := (s.Value / total) * 2 * math.Pi

		x1 := cx + r*math.Cos(startAngle)
		y1 := cy + r*math.Sin(startAngle)
		endAngle := startAngle + angle
		x2 := cx + r*math.Cos(endAngle)
		y2 := cy + r*math.Sin(endAngle)

		largeArc := 0
		if angle > math.Pi {
			largeArc = 1
		}

		d := fmt.Sprintf("M %.2f %.2f L %.2f %.2f A %.2f %.2f 0 %d 1 %.2f %.2f Z",
			cx, cy, x1, y1, r, r, largeArc, x2, y2)

		nodes = append(nodes, g.El("path",
			g.Attr("d", d),
			g.Attr("fill", color),
			g.El("title", g.Raw(fmt.Sprintf("%s: %.1f%%", s.Label, s.Value/total*100))),
		))
		startAngle = endAngle
	}

	// Legend
	legendItems := make([]g.Node, len(p.Slices))
	for i, s := range p.Slices {
		color := s.Color
		if color == "" {
			color = defaultColors[i%len(defaultColors)]
		}
		legendItems[i] = h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2)"),
			h.Span(h.Style("width:10px;height:10px;border-radius:2px;background:"+color+";flex-shrink:0")),
			h.Span(h.Style("font-size:var(--t-xs);flex:1"), g.Text(s.Label)),
			h.Span(h.Style("font-size:var(--t-xs);font-weight:700;opacity:.7"),
				g.Text(fmt.Sprintf("%.0f%%", s.Value/total*100))),
		)
	}

	svgNode := g.El("svg",
		g.Attr("viewBox", fmt.Sprintf("0 0 %d %d", p.Size, p.Size)),
		g.Attr("style", fmt.Sprintf("width:%dpx;height:%dpx", p.Size, p.Size)),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(nodes),
	)

	chart := h.Div(
		h.Style("display:flex;align-items:center;gap:var(--sp-5);flex-wrap:wrap"),
		svgNode,
		h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"), g.Group(legendItems)),
	)

	if p.Caption != "" {
		return g.El("figure", g.Attr("style", "margin:0"),
			chart,
			g.El("figcaption",
				g.Attr("style", "font-size:var(--t-xs);color:var(--muted);text-align:center;margin-top:var(--sp-2)"),
				g.Raw(p.Caption)),
		)
	}
	return chart
}
