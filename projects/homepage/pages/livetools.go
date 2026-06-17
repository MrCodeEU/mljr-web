package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func liveToolsSection(num string, tools []hpdata.Project, lang string) g.Node {
	if len(tools) == 0 {
		return nil
	}

	tones := []token.Tone{token.ToneViolet, token.ToneCyan, token.ToneLime, token.ToneYellow, token.TonePink, token.ToneMint}
	cards := make([]g.Node, 0, len(tools))
	for i, t := range tools {
		cards = append(cards, liveToolCard(t, lang, tones[i%len(tones)]))
	}

	subtitle := fmt.Sprintf(i18n.T(lang, "sections.livetools.count"), len(tools))

	return h.Section(
		h.ID("tools"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader(num, i18n.T(lang, "sections.livetools.title"), subtitle, token.ToneAccent),
			h.Div(
				h.Class("tools-grid"),
				g.Group(cards),
			),
		),
	)
}

func liveToolCard(p hpdata.Project, lang string, tone token.Tone) g.Node {
	liveURL := p.LiveURL()

	tagNodes := make([]g.Node, 0, len(p.Topics))
	for _, t := range p.Topics {
		if t == "live-tool" {
			continue
		}
		tagNodes = append(tagNodes, primitive.Tag(
			primitive.TagProps{Tone: hpdata.TagTone(t), Icon: hpdata.TechIcon(t)},
			g.Text(t),
		))
	}

	return primitive.Card(primitive.CardProps{
		Tone:  tone,
		Attrs: []g.Node{h.Style("display:flex;flex-direction:column;gap:var(--sp-4);height:100%")},
	},
		// Header row: name + live badge
		h.Div(
			h.Style("display:flex;align-items:flex-start;gap:var(--sp-3)"),
			h.H3(
				h.Style("font-size:var(--t-xl);font-weight:900;line-height:1.1;margin:0;flex:1"),
				g.Text(p.Name),
			),
			primitive.Tag(primitive.TagProps{Tone: token.ToneLime},
				icon.Icon("lucide:radio", icon.Props{Size: ".8rem"}),
				g.Text("LIVE"),
			),
		),
		// Description
		h.P(
			h.Style("margin:0;font-size:var(--t-sm);line-height:1.6;opacity:.85;flex:1"),
			g.Text(p.DescFor(lang)),
		),
		// Tags row
		g.If(len(tagNodes) > 0,
			h.Div(
				h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
				g.Group(tagNodes),
			),
		),
		// CTA
		g.If(liveURL != "",
			h.A(
				h.Href(liveURL),
				g.Attr("target", "_blank"),
				g.Attr("rel", "noopener"),
				h.Style("margin-top:auto"),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeMD},
					icon.Icon("lucide:external-link", icon.Props{Size: "1rem"}),
					g.Text(i18n.T(lang, "sections.livetools.open")),
				),
			),
		),
	)
}
