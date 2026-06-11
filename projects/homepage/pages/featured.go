package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// featuredSection renders a spotlight grid of featured projects using 3D tilt
// cards. The first project gets a double-height hero cell; the rest stack
// beside it — an asymmetric editorial layout.
func featuredSection(featured []hpdata.Project) g.Node {
	if len(featured) == 0 {
		return nil
	}
	if len(featured) > 3 {
		featured = featured[:3]
	}

	tones := []token.Tone{token.ToneYellow, token.ToneCyan, token.ToneViolet}

	cells := make([]g.Node, len(featured))
	for i, p := range featured {
		big := i == 0
		span := "grid-row:span 1"
		if big && len(featured) > 1 {
			span = "grid-row:span 2"
		}
		cells[i] = h.Div(
			h.Style(span+";display:flex;min-width:0"),
			primitive.TiltCard(primitive.TiltCardProps{MaxTilt: 6, Scale: 1.015, Shine: true},
				featuredCard(p, tones[i%len(tones)], big),
			),
		)
	}

	return h.Section(
		h.ID("featured"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("02", "Featured work", "hand-picked", token.ToneYellow),
			h.Div(
				h.Class("featured-grid"),
				h.Style("display:grid;grid-template-columns:1.25fr 1fr;grid-auto-rows:minmax(220px,auto);gap:var(--sp-4)"),
				g.Group(cells),
			),
		),
	)
}

func featuredCard(p hpdata.Project, tone token.Tone, big bool) g.Node {
	imgs := p.LocalImages()

	var media g.Node
	if big && len(imgs) > 0 {
		media = h.Img(
			h.Src(imgs[0]),
			h.Alt(p.Name),
			g.Attr("loading", "lazy"),
			h.Style("width:100%;aspect-ratio:16/8;object-fit:cover;border-bottom:var(--bw-2) solid var(--ink);display:block"),
		)
	}

	titleSize := "var(--t-xl)"
	if big {
		titleSize = "clamp(1.8rem,3vw,2.6rem)"
	}

	topicNodes := make([]g.Node, 0, 4)
	for _, t := range p.Topics {
		if len(topicNodes) >= 4 {
			break
		}
		topicNodes = append(topicNodes, primitive.Tag(
			primitive.TagProps{Icon: hpdata.TechIcon(t)}, g.Text(t),
		))
	}

	var ghLink g.Node
	if p.URL != "" {
		ghLink = h.A(h.Href(p.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
			h.Style("margin-top:auto"),
			primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
				icon.Icon("simple-icons:github"),
				g.Text("Source"),
				icon.Icon("lucide:arrow-up-right"),
			),
		)
	}

	return primitive.Card(primitive.CardProps{Tone: tone, Attrs: []g.Node{h.Style("width:100%;display:flex;flex-direction:column")}},
		media,
		h.Div(h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-2)"),
			h.H3(h.Style("font-weight:900;font-size:"+titleSize+";line-height:1.05;margin:0"),
				g.Text(friendlyName(p.Name))),
			g.If(p.Stars > 0,
				primitive.Tag(primitive.TagProps{Tone: token.ToneYellow}, g.Text(fmt.Sprintf("★ %d", p.Stars))),
			),
		),
		h.P(h.Style("margin:0;font-size:var(--t-sm);line-height:1.55;opacity:.85"), g.Text(truncate(p.Desc, 160))),
		h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
			g.If(p.Language != "", primitive.Tag(primitive.TagProps{Tone: token.ToneAccent}, g.Text(p.Language))),
			g.Group(topicNodes),
		),
		ghLink,
	)
}
