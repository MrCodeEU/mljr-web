package data

import (
	"fmt"
	stdhtml "html"
	"math"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type RadarSeries struct {
	Label  string
	Values []float64 // one per axis, 0–100
	Color  string    // CSS color (default cycles through palette)
}

type RadarChartProps struct {
	// Axes are the dimension labels.
	Axes []string
	// Size is the SVG width/height in px (default 280).
	Size int
	// Max value for scaling (default 100).
	Max float64
	// ShowGrid renders concentric polygon grid lines.
	ShowGrid bool
	// GridLevels is number of concentric rings (default 5).
	GridLevels int
}

var radarColors = []string{"var(--accent)", "var(--success)", "var(--warning)", "var(--danger)", "var(--info)"}

// RadarChart renders a multi-series radar/spider chart as inline SVG.
// All rendering is server-side Go — zero JS, zero dependencies.
func RadarChart(p RadarChartProps, series ...RadarSeries) g.Node {
	if p.Size == 0 {
		p.Size = 280
	}
	if p.Max == 0 {
		p.Max = 100
	}
	if p.GridLevels == 0 {
		p.GridLevels = 5
	}

	n := len(p.Axes)
	if n < 3 {
		return g.Text("RadarChart: need at least 3 axes")
	}

	cx, cy := float64(p.Size)/2, float64(p.Size)/2
	r := float64(p.Size) * 0.38 // chart radius
	labelR := r + 18

	// angle for axis i (start at top = -90°)
	angle := func(i int) float64 {
		return float64(i)*2*math.Pi/float64(n) - math.Pi/2
	}
	pt := func(radius float64, i int) (float64, float64) {
		a := angle(i)
		return cx + radius*math.Cos(a), cy + radius*math.Sin(a)
	}
	poly := func(radius float64) string {
		pts := make([]string, n)
		for i := range p.Axes {
			x, y := pt(radius, i)
			pts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
		}
		return strings.Join(pts, " ")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" data-component="radar-chart">`,
		p.Size, p.Size, p.Size, p.Size))

	// Grid polygons
	if p.ShowGrid {
		for level := 1; level <= p.GridLevels; level++ {
			lvlR := r * float64(level) / float64(p.GridLevels)
			sb.WriteString(fmt.Sprintf(`<polygon points="%s" fill="none" stroke="var(--line)" stroke-width="1" opacity="0.5"/>`, poly(lvlR)))
		}
	}

	// Axis lines + labels
	for i, axis := range p.Axes {
		x, y := pt(r, i)
		sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="var(--line)" stroke-width="1"/>`,
			cx, cy, x, y))
		lx, ly := pt(labelR, i)
		anchor := "middle"
		if lx < cx-1 {
			anchor = "end"
		} else if lx > cx+1 {
			anchor = "start"
		}
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="%s" dominant-baseline="middle" fill="var(--muted)" font-size="11" font-family="var(--font-display)" font-weight="700">%s</text>`,
			lx, ly, anchor, stdhtml.EscapeString(axis)))
	}

	// Data series
	for si, s := range series {
		color := s.Color
		if color == "" {
			color = radarColors[si%len(radarColors)]
		}
		pts := make([]string, len(p.Axes))
		for i := range p.Axes {
			val := 0.0
			if i < len(s.Values) {
				val = s.Values[i]
			}
			scaled := (val / p.Max) * r
			x, y := pt(scaled, i)
			pts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
		}
		sb.WriteString(fmt.Sprintf(`<polygon points="%s" fill="%s" fill-opacity="0.18" stroke="%s" stroke-width="2"/>`,
			strings.Join(pts, " "), color, color))
		// Dots at each vertex
		for i := range p.Axes {
			val := 0.0
			if i < len(s.Values) {
				val = s.Values[i]
			}
			scaled := (val / p.Max) * r
			x, y := pt(scaled, i)
			sb.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" fill="%s" stroke="var(--bg)" stroke-width="2"/>`, x, y, color))
		}
	}

	// Legend
	if len(series) > 1 {
		for si, s := range series {
			color := s.Color
			if color == "" {
				color = radarColors[si%len(radarColors)]
			}
			lx := 8.0
			ly := float64(p.Size) - float64((len(series)-si)*18) - 4
			sb.WriteString(fmt.Sprintf(`<rect x="%.0f" y="%.0f" width="10" height="10" fill="%s" rx="2"/>`, lx, ly, color))
			sb.WriteString(fmt.Sprintf(`<text x="%.0f" y="%.0f" fill="var(--muted)" font-size="10" font-family="var(--font-display)">%s</text>`,
				lx+14, ly+9, stdhtml.EscapeString(s.Label)))
		}
	}

	sb.WriteString(`</svg>`)
	return h.Div(g.Attr("data-component", "radar-chart-wrap"), g.Raw(sb.String()))
}
