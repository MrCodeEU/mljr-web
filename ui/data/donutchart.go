package data

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DonutSlice struct {
	Label string
	Value float64
	Color string // CSS color
}

type DonutChartProps struct {
	Slices    []DonutSlice
	Size      int    // px diameter (default 200)
	Thickness int    // stroke thickness in px (default 36)
	Label     string // center label (e.g. total or %)
	Sublabel  string // center sublabel
	Caption   string
}

// DonutChart renders a pure-SVG donut/pie chart using stroke-dasharray. No JS.
func DonutChart(p DonutChartProps) g.Node {
	if p.Size == 0 {
		p.Size = 180
	}
	if p.Thickness == 0 {
		p.Thickness = 36
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
	r := cx - float64(p.Thickness)/2 - 2
	circ := 2 * math.Pi * r

	defaultColors := []string{"var(--primary)", "var(--accent)", "var(--success)", "var(--warning)", "var(--danger)"}

	nodes := []g.Node{}

	// Background ring
	nodes = append(nodes, g.El("circle",
		g.Attr("cx", fmt.Sprintf("%.1f", cx)),
		g.Attr("cy", fmt.Sprintf("%.1f", cy)),
		g.Attr("r", fmt.Sprintf("%.1f", r)),
		g.Attr("fill", "none"),
		g.Attr("stroke", "var(--line)"),
		g.Attr("stroke-width", fmt.Sprintf("%d", p.Thickness)),
	))

	// Slices — drawn by rotating each circle
	offset := 0.0
	for i, s := range p.Slices {
		color := s.Color
		if color == "" {
			color = defaultColors[i%len(defaultColors)]
		}
		pct := s.Value / total
		dash := pct * circ
		gap := circ - dash
		// SVG starts at 3 o'clock; rotate -90deg (start at 12 o'clock)
		// Each slice starts at 'offset' fraction of the circumference
		rotate := -90 + offset/circ*360

		nodes = append(nodes, g.El("circle",
			g.Attr("cx", fmt.Sprintf("%.1f", cx)),
			g.Attr("cy", fmt.Sprintf("%.1f", cy)),
			g.Attr("r", fmt.Sprintf("%.1f", r)),
			g.Attr("fill", "none"),
			g.Attr("stroke", color),
			g.Attr("stroke-width", fmt.Sprintf("%d", p.Thickness)),
			g.Attr("stroke-dasharray", fmt.Sprintf("%.2f %.2f", dash, gap)),
			g.Attr("transform", fmt.Sprintf("rotate(%.2f %.1f %.1f)", rotate, cx, cy)),
			g.El("title", g.Raw(fmt.Sprintf("%s: %.1f", s.Label, s.Value))),
		))
		offset += dash
	}

	// Center text
	if p.Label != "" {
		nodes = append(nodes, g.El("text",
			g.Attr("x", fmt.Sprintf("%.1f", cx)),
			g.Attr("y", fmt.Sprintf("%.1f", cy+5)),
			g.Attr("text-anchor", "middle"),
			g.Attr("dominant-baseline", "middle"),
			g.Attr("font-size", "22"),
			g.Attr("font-weight", "900"),
			g.Attr("fill", "var(--fg)"),
			g.Attr("font-family", "var(--font-display)"),
			g.Raw(p.Label),
		))
	}
	if p.Sublabel != "" {
		nodes = append(nodes, g.El("text",
			g.Attr("x", fmt.Sprintf("%.1f", cx)),
			g.Attr("y", fmt.Sprintf("%.1f", cy+24)),
			g.Attr("text-anchor", "middle"),
			g.Attr("font-size", "11"),
			g.Attr("fill", "var(--muted)"),
			g.Raw(p.Sublabel),
		))
	}

	svgNode := g.El("svg",
		g.Attr("viewBox", fmt.Sprintf("0 0 %d %d", p.Size, p.Size)),
		g.Attr("style", fmt.Sprintf("width:%dpx;height:%dpx", p.Size, p.Size)),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(nodes),
	)

	// Legend
	legendItems := make([]g.Node, len(p.Slices))
	for i, s := range p.Slices {
		color := s.Color
		if color == "" {
			color = defaultColors[i%len(defaultColors)]
		}
		pct := s.Value / total * 100
		legendItems[i] = h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-2)"),
			h.Span(h.Style("width:10px;height:10px;border-radius:2px;background:"+color+";flex-shrink:0")),
			h.Span(h.Style("font-size:var(--t-xs);flex:1"), g.Text(s.Label)),
			h.Span(h.Style("font-size:var(--t-xs);font-weight:700;opacity:.7"), g.Text(fmt.Sprintf("%.0f%%", pct))),
		)
	}

	chart := h.Div(
		h.Style("display:flex;align-items:center;gap:var(--sp-5);flex-wrap:wrap"),
		svgNode,
		h.Div(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-2);min-width:120px"),
			g.Group(legendItems),
		),
	)

	if p.Caption != "" {
		return g.El("figure",
			g.Attr("style", "margin:0"),
			chart,
			g.El("figcaption",
				g.Attr("style", "font-size:var(--t-xs);color:var(--muted);text-align:center;margin-top:var(--sp-2)"),
				g.Raw(p.Caption),
			),
		)
	}
	return chart
}
