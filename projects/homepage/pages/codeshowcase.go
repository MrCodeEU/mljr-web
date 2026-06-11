package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// codeShowcaseSection explains how this site is built and shows a real
// component from the UI library, highlighted server-side with chroma.
func codeShowcaseSection() g.Node {
	return h.Section(
		h.ID("under-the-hood"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("08", "Under the hood", "this site, in code", token.ToneSky),
			h.Div(
				h.Class("hood-grid"),
				h.Style("display:grid;grid-template-columns:1fr 1.2fr;gap:var(--sp-6);align-items:start"),
				// Left: narrative + principle chips
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-4);position:sticky;top:var(--sp-8)"),
					h.P(h.Style("font-size:var(--t-lg);font-weight:700;line-height:1.5;margin:0"),
						g.Text("Every component on this page is a Go function. No React, no npm, no build pipeline for the markup — the server renders HTML, Datastar adds reactivity, and SVG charts are computed at request time."),
					),
					h.P(h.Style("font-size:var(--t-base);color:var(--muted);line-height:1.6;margin:0"),
						g.Text("On the right: the actual arc math behind the gauges in the open-source section above. What you see is what runs."),
					),
					h.Div(
						h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-3)"),
						hoodFact("lucide:boxes", "175", "UI components"),
						hoodFact("lucide:package-x", "0", "npm runtime deps"),
						hoodFact("lucide:palette", "4×2", "themes × modes"),
						hoodFact("lucide:server", "1", "binary to deploy"),
					),
					h.A(h.Href("https://github.com/MrCodeEU"), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary},
							icon.Icon("simple-icons:go"),
							g.Text("Explore the source"),
							icon.Icon("lucide:arrow-up-right"),
						),
					),
				),
				// Right: real component source, highlighted server-side
				uidata.SyntaxHighlighter(uidata.SyntaxHighlighterProps{
					Language: "go",
					Theme:    "monokai",
					Filename: "ui/data/gauge.go",
				}, gaugeExcerpt),
			),
		),
	)
}

func hoodFact(ic, value, label string) g.Node {
	return h.Div(
		h.Style("border:var(--bw-2) solid var(--ink);background:var(--surface);box-shadow:var(--shadow-sm);padding:var(--sp-3) var(--sp-4);display:flex;align-items:center;gap:var(--sp-3)"),
		icon.Icon(ic, icon.Props{Size: "1.4rem"}),
		h.Div(
			h.Div(h.Style("font-weight:900;font-size:var(--t-xl);line-height:1"), g.Text(value)),
			h.Div(h.Style("font-size:var(--t-xs);font-weight:700;color:var(--muted);text-transform:uppercase;letter-spacing:.08em;margin-top:2px"), g.Text(label)),
		),
	)
}

// gaugeExcerpt is the real arc-path math from ui/data/gauge.go, shown in the
// "Under the hood" section. Keep in sync if the component changes materially.
const gaugeExcerpt = `// Gauge renders a circular gauge as inline SVG.
// Zero JS, zero dependencies — pure server-side Go.
func Gauge(p GaugeProps) g.Node {
	cx, cy := float64(p.Size)/2, float64(p.Size)/2
	radius := float64(p.Size) * 0.38
	strokeW := float64(p.Size) * 0.1

	// Arc spans 270° (from 135° to 45°)
	startAngle := 135.0
	totalAngle := 270.0
	pct := (p.Value - p.Min) / (p.Max - p.Min)

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}
	arcPt := func(deg float64) (float64, float64) {
		a := toRad(deg)
		return cx + radius*math.Cos(a),
			cy + radius*math.Sin(a)
	}
	arcPath := func(start, end float64, color string) string {
		x1, y1 := arcPt(start)
		x2, y2 := arcPt(end)
		large := 0
		if math.Abs(end-start) > 180 {
			large = 1
		}
		return fmt.Sprintf(
			'<path d="M %.1f %.1f A %.1f %.1f 0 %d 1 %.1f %.1f"
			   fill="none" stroke="%s" stroke-width="%.1f"
			   stroke-linecap="round"/>',
			x1, y1, radius, radius, large, x2, y2,
			color, strokeW)
	}

	// Track (full arc), then value arc on top
	sb.WriteString(arcPath(startAngle, startAngle+totalAngle, p.TrackColor))
	sb.WriteString(arcPath(startAngle, startAngle+totalAngle*pct, p.Color))
	// …
}`
