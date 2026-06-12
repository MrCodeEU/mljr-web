package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func skillsSection() g.Node {
	groups := hpdata.SkillGroups()

	// 3 Marquee rows alternating direction, one per "layer" of the skill stack
	// Group skills into 3 rows: languages, web+infra, security+other
	row := func(skills []hpdata.SkillGroup, dir string, speed string) g.Node {
		var items []g.Node
		idx := 0
		for _, sg := range skills {
			for _, s := range sg.Skills {
				ic := hpdata.TechIcon(s)
				items = append(items, skillPill(s, ic, idx))
				idx++
			}
		}
		return primitive.Marquee(primitive.MarqueeProps{
			Speed:        speed,
			Direction:    dir,
			PauseOnHover: true,
			Gap:          "var(--sp-2)",
		}, items...)
	}

	// Split groups into 3 rows
	var langWeb, infra, secOther []hpdata.SkillGroup
	for _, sg := range groups {
		switch sg.Label {
		case "Languages", "Web":
			langWeb = append(langWeb, sg)
		case "Infra / Homelab", "Ops / Data":
			infra = append(infra, sg)
		default:
			secOther = append(secOther, sg)
		}
	}

	// Radar: self-assessed depth per area, vertex dots color-coded to the
	// group cards on the right.
	radarAxes := make([]string, len(groups))
	radarValues := make([]float64, len(groups))
	axisColors := make([]string, len(groups))
	for i, sg := range groups {
		radarAxes[i] = sg.Short
		radarValues[i] = float64(sg.Level)
		axisColors[i] = toneBG(sg.Tone)
	}

	groupCards := make([]g.Node, len(groups))
	for i, sg := range groups {
		chips := make([]g.Node, 0, len(sg.Skills))
		for _, s := range sg.Skills {
			ic := hpdata.TechIcon(s)
			chips = append(chips, h.Span(
				h.Style("display:inline-flex;align-items:center;gap:4px;border:var(--bw-1) solid var(--ink);background:var(--bg);padding:2px var(--sp-2);font-size:var(--t-xs);font-weight:700;white-space:nowrap"),
				g.If(ic != "", icon.Icon(ic, icon.Props{Size: ".85rem"})),
				g.Text(s),
			))
		}
		groupCards[i] = primitive.Card(primitive.CardProps{Tone: token.Tone(sg.Tone)},
			h.Div(
				h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
				icon.Icon(sg.Icon, icon.Props{Size: "1.2rem"}),
				h.H3(h.Style("font-size:var(--t-base);font-weight:900;margin:0;flex:1;min-width:0"), g.Text(sg.Label)),
				h.Span(h.Style("font-size:var(--t-xs);font-weight:900;font-family:var(--font-mono,monospace);border:var(--bw-1) solid var(--ink);background:var(--bg);padding:1px var(--sp-2)"), g.Textf("%d", len(sg.Skills))),
			),
			h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-1)"), g.Group(chips)),
		)
	}

	return h.Section(
		h.ID("skills"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("07", "Skills", "depth × breadth", token.ToneViolet),
			h.Div(
				h.Class("skills-grid"),
				h.Style("display:grid;grid-template-columns:minmax(280px,360px) 1fr;gap:var(--sp-5);align-items:stretch"),
				// Left: radar card
				primitive.Card(primitive.CardProps{Tone: token.ToneNone},
					h.Div(h.Style("display:flex;flex-direction:column;align-items:center;justify-content:center;height:100%;gap:var(--sp-3)"),
						uidata.RadarChart(uidata.RadarChartProps{
							Axes:       radarAxes,
							ShowGrid:   true,
							GridLevels: 4,
							Size:       300,
							Max:        100,
							AxisColors: axisColors,
						},
							uidata.RadarSeries{Label: "Depth", Values: radarValues, Color: "var(--accent)"},
						),
						h.Div(h.Style("font-size:var(--t-xs);font-weight:800;color:var(--muted);text-transform:uppercase;letter-spacing:.08em;text-align:center"),
							g.Text("self-assessed depth per area · dots match the cards"),
						),
					),
				),
				// Right: one toned card per skill group
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:var(--sp-4);align-content:stretch"),
					g.Group(groupCards),
				),
			),
		),
		// Marquee rows (outside container for full width)
		h.Div(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-3);padding:var(--sp-6) 0 0"),
			row(langWeb, "left", "28s"),
			row(infra, "right", "22s"),
			row(secOther, "left", "25s"),
		),
	)
}

// toneBG maps a tone name to its pastel background CSS variable (with
// fallback), used to color radar vertex dots to match the group cards.
func toneBG(tone string) string {
	m := map[string]string{
		"yellow": "var(--yellow-bg,#fef08a)",
		"cyan":   "var(--cyan-bg,#a5f3fc)",
		"violet": "var(--violet-bg,#ddd6fe)",
		"lime":   "var(--lime-bg,#d9f99d)",
		"pink":   "var(--pink-bg,#fbcfe8)",
		"sky":    "var(--sky-bg,#bae6fd)",
		"mint":   "var(--mint-bg,#d1fae5)",
		"blush":  "var(--blush-bg,#fde8e8)",
	}
	if v, ok := m[tone]; ok {
		return v
	}
	return "var(--accent)"
}

// toneVars maps tone names to their CSS background + text variables for inline use.
var toneVars = []struct{ bg, ink string }{
	{"var(--yellow-bg,#fef08a)", "var(--yellow-ink,#713f12)"},
	{"var(--cyan-bg,#a5f3fc)", "var(--cyan-ink,#164e63)"},
	{"var(--violet-bg,#ddd6fe)", "var(--violet-ink,#3b0764)"},
	{"var(--lime-bg,#d9f99d)", "var(--lime-ink,#365314)"},
	{"var(--pink-bg,#fbcfe8)", "var(--pink-ink,#831843)"},
	{"var(--sky-bg,#bae6fd)", "var(--sky-ink,#0c4a6e)"},
	{"var(--mint-bg,#d1fae5)", "var(--mint-ink,#064e3b)"},
	{"var(--blush-bg,#fde8e8)", "var(--blush-ink,#7f1d1d)"},
}

func skillPill(label, ic string, idx int) g.Node {
	tv := toneVars[idx%len(toneVars)]
	style := "display:inline-flex;align-items:center;gap:var(--sp-2);padding:var(--sp-2) var(--sp-4);border:var(--bw-2) solid var(--ink);border-radius:calc(var(--radius)*2);font-size:var(--t-sm);font-weight:700;white-space:nowrap;flex-shrink:0;background:" + tv.bg + ";color:" + tv.ink
	return h.Div(
		h.Style(style),
		g.If(ic != "", icon.Icon(ic, icon.Props{Size: "1rem"})),
		g.Text(label),
	)
}

// sectionHeader renders a Swiss-editorial section heading: a large outlined
// index number, the heading, and a tone tag badge on the right.
func sectionHeader(num, heading, sub string, tone token.Tone) g.Node {
	return h.Div(
		h.Style("display:flex;align-items:baseline;justify-content:space-between;flex-wrap:wrap;gap:var(--sp-3);margin-bottom:var(--sp-8)"),
		h.Div(
			h.Style("display:flex;align-items:baseline;gap:var(--sp-4)"),
			g.If(num != "", h.Span(h.Class("section-num"), g.Text(num))),
			primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(heading)),
		),
		g.If(sub != "", primitive.Tag(primitive.TagProps{Tone: tone}, g.Text(sub))),
	)
}
