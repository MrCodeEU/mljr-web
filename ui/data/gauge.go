package data

import (
	"fmt"
	stdhtml "html"
	"math"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GaugeProps struct {
	// Value is the current value (0–Max).
	Value float64
	// Max is the maximum value (default 100).
	Max float64
	// Min is the minimum value (default 0).
	Min float64
	// Label shown in center below value (e.g. "CPU", "Score").
	Label string
	// Unit appended to value display (e.g. "%", "°C").
	Unit string
	// Size is SVG width/height in px (default 200).
	Size int
	// Ticks renders tick marks around the arc.
	Ticks bool
	// Color overrides the arc color (default "var(--accent)").
	Color string
	// TrackColor is the background arc color (default "var(--surface-2)").
	TrackColor string
}

// Gauge renders a circular gauge/meter as inline SVG. Zero JS, zero dependencies.
func Gauge(p GaugeProps) g.Node {
	if p.Max == 0 {
		p.Max = 100
	}
	if p.Size == 0 {
		p.Size = 200
	}
	if p.Color == "" {
		p.Color = "var(--accent)"
	}
	if p.TrackColor == "" {
		p.TrackColor = "var(--surface-2)"
	}

	cx, cy := float64(p.Size)/2, float64(p.Size)/2
	radius := float64(p.Size) * 0.38
	strokeW := float64(p.Size) * 0.1

	// Arc spans 270° (from 135° to 135°+270° = 405°=45°)
	startAngle := 135.0 // degrees
	totalAngle := 270.0
	pct := (p.Value - p.Min) / (p.Max - p.Min)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}

	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }
	arcPt := func(deg float64) (float64, float64) {
		a := toRad(deg)
		return cx + radius*math.Cos(a), cy + radius*math.Sin(a)
	}
	arcPath := func(startDeg, endDeg float64, color string) string {
		x1, y1 := arcPt(startDeg)
		x2, y2 := arcPt(endDeg)
		large := 0
		if math.Abs(endDeg-startDeg) > 180 {
			large = 1
		}
		return fmt.Sprintf(`<path d="M %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f" fill="none" stroke="%s" stroke-width="%.1f" stroke-linecap="round"/>`,
			x1, y1, radius, radius, large, x2, y2, color, strokeW)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" data-component="gauge">`,
		p.Size, p.Size, p.Size, p.Size))

	// Track (full arc)
	sb.WriteString(arcPath(startAngle, startAngle+totalAngle, p.TrackColor))
	// Value arc
	if pct > 0.001 {
		sb.WriteString(arcPath(startAngle, startAngle+totalAngle*pct, p.Color))
	}

	// Ticks
	if p.Ticks {
		nTicks := 10
		for i := 0; i <= nTicks; i++ {
			deg := startAngle + totalAngle*float64(i)/float64(nTicks)
			a := toRad(deg)
			inner := radius - strokeW*0.7
			outer := radius + strokeW*0.1
			x1 := cx + inner*math.Cos(a)
			y1 := cy + inner*math.Sin(a)
			x2 := cx + outer*math.Cos(a)
			y2 := cy + outer*math.Sin(a)
			sb.WriteString(fmt.Sprintf(`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="var(--bg)" stroke-width="2"/>`, x1, y1, x2, y2))
		}
	}

	// Center value text
	valueStr := fmt.Sprintf("%.0f%s", p.Value, stdhtml.EscapeString(p.Unit))
	sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" dominant-baseline="middle" fill="var(--ink)" font-size="%d" font-family="var(--font-display)" font-weight="900">%s</text>`,
		cx, cy-float64(p.Size)*0.04, p.Size/6, valueStr))
	if p.Label != "" {
		sb.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" text-anchor="middle" fill="var(--muted)" font-size="%d" font-family="var(--font-display)" font-weight="700" letter-spacing="1">%s</text>`,
			cx, cy+float64(p.Size)*0.12, p.Size/14, stdhtml.EscapeString(strings.ToUpper(p.Label))))
	}

	sb.WriteString(`</svg>`)
	return h.Div(g.Attr("data-component", "gauge-wrap"), g.Raw(sb.String()))
}
