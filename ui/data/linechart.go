package data

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
)

type LineChartSeries struct {
	Label  string
	Points []float64
	Color  string // CSS color; defaults to var(--primary)
	Fill   bool   // render area fill under the line
}

type LineChartProps struct {
	Series   []LineChartSeries
	Labels   []string // x-axis labels; len should match Points
	Height   int      // px (default 160)
	ShowDots bool
	ShowGrid bool
	Caption  string
}

// LineChart renders a pure-SVG multi-series line/area chart. No JS required.
func LineChart(p LineChartProps) g.Node {
	if p.Height == 0 {
		p.Height = 160
	}
	if len(p.Series) == 0 {
		return g.El("svg")
	}

	const padL = 8
	const padR = 8
	const labelH = 20
	totalW := 400
	chartH := float64(p.Height)

	// Global min/max across all series
	minV, maxV := math.MaxFloat64, -math.MaxFloat64
	for _, s := range p.Series {
		for _, v := range s.Points {
			if v < minV {
				minV = v
			}
			if v > maxV {
				maxV = v
			}
		}
	}
	if maxV == minV {
		maxV = minV + 1
	}
	rang := maxV - minV

	// Determine number of x points
	nPts := 0
	for _, s := range p.Series {
		if len(s.Points) > nPts {
			nPts = len(s.Points)
		}
	}
	if nPts < 2 {
		nPts = 2
	}
	xStep := float64(totalW-padL-padR) / float64(nPts-1)

	yOf := func(v float64) float64 {
		return chartH - (v-minV)/rang*chartH
	}
	xOf := func(i int) float64 {
		return float64(padL) + float64(i)*xStep
	}

	nodes := []g.Node{}

	// Grid lines
	if p.ShowGrid {
		for _, frac := range []float64{0.25, 0.5, 0.75, 1.0} {
			y := chartH * (1 - frac)
			nodes = append(nodes, g.El("line",
				g.Attr("x1", fmt.Sprintf("%d", padL)),
				g.Attr("y1", fmt.Sprintf("%.1f", y)),
				g.Attr("x2", fmt.Sprintf("%d", totalW-padR)),
				g.Attr("y2", fmt.Sprintf("%.1f", y)),
				g.Attr("stroke", "var(--line)"),
				g.Attr("stroke-dasharray", "4 4"),
			))
		}
	}

	// Series
	defaultColors := []string{"var(--primary)", "var(--accent)", "var(--success)", "var(--warning)"}
	for si, s := range p.Series {
		color := s.Color
		if color == "" {
			color = defaultColors[si%len(defaultColors)]
		}

		pts := make([][2]float64, len(s.Points))
		for i, v := range s.Points {
			pts[i] = [2]float64{xOf(i), yOf(v)}
		}

		// Build SVG path
		d := ""
		for i, pt := range pts {
			if i == 0 {
				d += fmt.Sprintf("M %.1f %.1f", pt[0], pt[1])
			} else {
				// Smooth curve via cubic bezier
				prev := pts[i-1]
				cpx := (prev[0] + pt[0]) / 2
				d += fmt.Sprintf(" C %.1f %.1f %.1f %.1f %.1f %.1f", cpx, prev[1], cpx, pt[1], pt[0], pt[1])
			}
		}

		// Area fill
		if s.Fill && len(pts) > 0 {
			fillD := d +
				fmt.Sprintf(" L %.1f %.1f L %.1f %.1f Z", pts[len(pts)-1][0], chartH, pts[0][0], chartH)
			nodes = append(nodes, g.El("path",
				g.Attr("d", fillD),
				g.Attr("fill", color),
				g.Attr("fill-opacity", "0.12"),
				g.Attr("stroke", "none"),
			))
		}

		// Line
		nodes = append(nodes, g.El("path",
			g.Attr("d", d),
			g.Attr("fill", "none"),
			g.Attr("stroke", color),
			g.Attr("stroke-width", "2"),
			g.Attr("stroke-linecap", "round"),
			g.Attr("stroke-linejoin", "round"),
		))

		// Dots
		if p.ShowDots {
			for _, pt := range pts {
				nodes = append(nodes, g.El("circle",
					g.Attr("cx", fmt.Sprintf("%.1f", pt[0])),
					g.Attr("cy", fmt.Sprintf("%.1f", pt[1])),
					g.Attr("r", "3"),
					g.Attr("fill", color),
					g.Attr("stroke", "var(--bg)"),
					g.Attr("stroke-width", "1.5"),
				))
			}
		}
	}

	// X labels
	if len(p.Labels) > 0 {
		for i, lbl := range p.Labels {
			if i >= nPts {
				break
			}
			nodes = append(nodes, g.El("text",
				g.Attr("x", fmt.Sprintf("%.1f", xOf(i))),
				g.Attr("y", fmt.Sprintf("%d", p.Height+labelH-4)),
				g.Attr("text-anchor", "middle"),
				g.Attr("font-size", "10"),
				g.Attr("fill", "var(--muted)"),
				g.Raw(lbl),
			))
		}
	}

	svgH := p.Height
	if len(p.Labels) > 0 {
		svgH += labelH
	}

	svgNode := g.El("svg",
		g.Attr("viewBox", fmt.Sprintf("0 0 %d %d", totalW, svgH)),
		g.Attr("style", fmt.Sprintf("width:100%%;height:%dpx;overflow:visible", svgH)),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(nodes),
	)

	if p.Caption != "" {
		return g.El("figure",
			g.Attr("style", "margin:0"),
			svgNode,
			g.El("figcaption",
				g.Attr("style", "font-size:var(--t-xs);color:var(--muted);text-align:center;margin-top:var(--sp-2)"),
				g.Raw(p.Caption),
			),
		)
	}
	return svgNode
}
