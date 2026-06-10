package data

import (
	"fmt"
	"math"

	g "maragu.dev/gomponents"
)

type SparklineProps struct {
	Points []float64
	Width  int    // px (default 80)
	Height int    // px (default 28)
	Color  string // CSS color (default var(--primary))
	Fill   bool   // render area fill
	Thick  bool   // stroke-width 2 instead of 1.5
}

// Sparkline renders a tiny inline SVG trend line — ideal inside stat cards or tables.
func Sparkline(p SparklineProps) g.Node {
	if p.Width == 0 {
		p.Width = 80
	}
	if p.Height == 0 {
		p.Height = 28
	}
	if p.Color == "" {
		p.Color = "var(--primary)"
	}
	if len(p.Points) < 2 {
		return g.El("svg",
			g.Attr("style", fmt.Sprintf("width:%dpx;height:%dpx", p.Width, p.Height)),
		)
	}

	minV, maxV := math.MaxFloat64, -math.MaxFloat64
	for _, v := range p.Points {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if maxV == minV {
		maxV = minV + 1
	}
	rang := maxV - minV

	n := len(p.Points)
	xStep := float64(p.Width) / float64(n-1)
	yOf := func(v float64) float64 {
		pad := 2.0
		return pad + (1-(v-minV)/rang)*float64(p.Height-4)
	}

	pts := make([][2]float64, n)
	for i, v := range p.Points {
		pts[i] = [2]float64{float64(i) * xStep, yOf(v)}
	}

	d := fmt.Sprintf("M %.1f %.1f", pts[0][0], pts[0][1])
	for i := 1; i < len(pts); i++ {
		cpx := (pts[i-1][0] + pts[i][0]) / 2
		d += fmt.Sprintf(" C %.1f %.1f %.1f %.1f %.1f %.1f",
			cpx, pts[i-1][1], cpx, pts[i][1], pts[i][0], pts[i][1])
	}

	sw := "1.5"
	if p.Thick {
		sw = "2"
	}

	nodes := []g.Node{}
	if p.Fill {
		fillD := d + fmt.Sprintf(" L %.1f %d L 0 %d Z", pts[n-1][0], p.Height, p.Height)
		nodes = append(nodes, g.El("path",
			g.Attr("d", fillD),
			g.Attr("fill", p.Color),
			g.Attr("fill-opacity", "0.15"),
			g.Attr("stroke", "none"),
		))
	}
	nodes = append(nodes, g.El("path",
		g.Attr("d", d),
		g.Attr("fill", "none"),
		g.Attr("stroke", p.Color),
		g.Attr("stroke-width", sw),
		g.Attr("stroke-linecap", "round"),
	))

	return g.El("svg",
		g.Attr("viewBox", fmt.Sprintf("0 0 %d %d", p.Width, p.Height)),
		g.Attr("style", fmt.Sprintf("width:%dpx;height:%dpx;overflow:visible", p.Width, p.Height)),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(nodes),
	)
}
