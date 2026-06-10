package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
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

	return h.Section(
		h.ID("skills"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("Skills", "My tech stack", token.ToneViolet),
		),
		// Marquee rows (outside container for full width)
		h.Div(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-3);padding:var(--sp-4) 0"),
			row(langWeb, "left", "28s"),
			row(infra, "right", "22s"),
			row(secOther, "left", "25s"),
		),
	)
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

// sectionHeader renders a section heading with a tone tag badge.
func sectionHeader(heading, sub string, tone token.Tone) g.Node {
	return h.Div(
		h.Style("display:flex;align-items:baseline;justify-content:space-between;flex-wrap:wrap;gap:var(--sp-3);margin-bottom:var(--sp-8)"),
		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(heading)),
		g.If(sub != "", primitive.Tag(primitive.TagProps{Tone: tone}, g.Text(sub))),
	)
}
