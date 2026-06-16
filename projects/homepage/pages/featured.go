package pages

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

// featuredSection renders a spotlight grid of featured projects. The first
// project gets a double-height hero cell; the rest stack beside it — an
// asymmetric editorial layout.
func featuredSection(featured []hpdata.Project, lang string) g.Node {
	if len(featured) == 0 {
		return nil
	}

	tones := []token.Tone{token.ToneYellow, token.ToneCyan, token.ToneViolet, token.TonePink, token.ToneLime}

	cells := make([]g.Node, len(featured))
	for i, p := range featured {
		big := i == 0
		cells[i] = h.Div(
			h.Style("display:flex;min-width:0;margin-bottom:var(--sp-4);break-inside:avoid"),
			featuredCard(p, tones[i%len(tones)], big, lang),
		)
	}

	return h.Section(
		h.ID("featured"),
		h.Style("padding:var(--sp-12) 0"),
		layout.Container(layout.ContainerProps{},
			sectionHeader("01", i18n.T(lang, "sections.featured.title"), i18n.T(lang, "sections.featured.sub"), token.ToneYellow),
			h.Div(
				h.Class("featured-grid"),
				h.Style("column-count:2;column-gap:var(--sp-4)"),
				g.Group(cells),
			),
		),
	)
}

func featuredCard(p hpdata.Project, tone token.Tone, big bool, lang string) g.Node {
	imgs := p.LocalImages()

	aspect := "16/9"
	if big {
		aspect = "16/8"
	}

	var media g.Node
	switch {
	case len(imgs) > 1:
		media = h.Div(
			h.Style("border-bottom:var(--bw-2) solid var(--ink);aspect-ratio:"+aspect),
			uidata.Carousel(uidata.CarouselProps{
				ID:     "fc" + slugify(p.Name),
				Images: imgs,
				Alt:    p.Name,
			}),
		)
	case len(imgs) == 1:
		media = h.Img(
			h.Src(imgs[0]),
			h.Alt(p.Name),
			g.Attr("loading", "lazy"),
			h.Style("width:100%;aspect-ratio:"+aspect+";object-fit:contain;background:var(--surface-2);border-bottom:var(--bw-2) solid var(--ink);display:block"),
		)
	}

	titleSize := "var(--t-xl)"
	if big {
		titleSize = "clamp(1.8rem,3vw,2.6rem)"
	}

	// Deduplicate: language first, then unique topics (case-insensitive)
	seen := map[string]bool{}
	allTags := make([]string, 0, len(p.Topics)+1)
	if p.Language != "" {
		allTags = append(allTags, p.Language)
		seen[strings.ToLower(p.Language)] = true
	}
	for _, t := range p.Topics {
		if !seen[strings.ToLower(t)] {
			allTags = append(allTags, t)
			seen[strings.ToLower(t)] = true
		}
	}
	topicNodes := make([]g.Node, 0, 5)
	for _, t := range allTags {
		if len(topicNodes) >= 5 {
			break
		}
		topicNodes = append(topicNodes, primitive.Tag(
			primitive.TagProps{Icon: hpdata.TechIcon(t), Tone: hpdata.TagTone(t)}, g.Text(t),
		))
	}

	var linkNodes []g.Node
	if p.URL != "" {
		linkNodes = append(linkNodes,
			h.A(h.Href(p.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
					icon.Icon("simple-icons:github"),
					g.Text(i18n.T(lang, "projects.source")),
					icon.Icon("lucide:arrow-up-right"),
				),
			),
		)
	}
	for _, lnk := range p.Links {
		if lnk.URL == "" {
			continue
		}
		linkNodes = append(linkNodes,
			h.A(h.Href(lnk.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
					g.Text(lnk.Name),
					icon.Icon("lucide:arrow-up-right"),
				),
			),
		)
	}
	ghLink := h.Div(h.Style("margin-top:auto;display:flex;flex-wrap:wrap;gap:var(--sp-2)"), g.Group(linkNodes))

	return primitive.Card(primitive.CardProps{Tone: tone, Attrs: []g.Node{h.Style("width:100%;display:flex;flex-direction:column")}},
		media,
		h.Div(h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-2)"),
			h.H3(h.Style("font-weight:900;font-size:"+titleSize+";line-height:1.05;margin:0"),
				g.Text(friendlyName(p.Name))),
			g.If(p.Stars > 0,
				primitive.Tag(primitive.TagProps{Tone: token.ToneYellow}, g.Text(fmt.Sprintf("★ %d", p.Stars))),
			),
		),
		h.P(h.Style("margin:0;font-size:var(--t-sm);line-height:1.55;opacity:.85"), g.Text(truncate(p.DescFor(lang), 350))),
		h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
			g.Group(topicNodes),
		),
		ghLink,
	)
}
