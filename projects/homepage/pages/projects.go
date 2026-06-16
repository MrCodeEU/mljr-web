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

func projectsSection(projects []hpdata.Project, lang string) g.Node {
	total := len(projects)
	pages := (total + perPage - 1) / perPage

	// Build all project cards — one Grid per page, animated by PaginatedPages
	var pageNodes []g.Node
	for p := 0; p < pages; p++ {
		start := p * perPage
		end := start + perPage
		if end > total {
			end = total
		}
		slice := projects[start:end]

		var cols []g.Node
		tones := []token.Tone{token.ToneNone, token.ToneSky, token.ToneLime, token.ToneViolet, token.TonePink, token.ToneMint}
		for i, proj := range slice {
			cols = append(cols, projectCard(proj, tones[i%len(tones)], lang))
		}

		pageNodes = append(pageNodes, layout.Grid(layout.GridProps{}, g.Group(cols)))
	}

	return h.Section(
		h.ID("projects"),
		h.Style("padding:var(--sp-12) 0"),
		uidata.PaginationSignals("pg", perPage),
		layout.Container(layout.ContainerProps{},
			sectionHeader("02", i18n.T(lang, "sections.projects.title"), fmt.Sprintf("%d projects", total), token.ToneLime),
			// top pagination
			h.Div(h.Style("margin-bottom:var(--sp-5);display:flex;justify-content:center"),
				uidata.Pagination(uidata.PaginationProps{ID: "pg", Total: total, PerPage: perPage}),
			),
			uidata.PaginatedPages(uidata.PaginatedPagesProps{ID: "pg", Animation: uidata.PageAnimSlideUp}, pageNodes...),
			// bottom pagination
			h.Div(h.Style("margin-top:var(--sp-6);display:flex;justify-content:center"),
				uidata.Pagination(uidata.PaginationProps{ID: "pg", Total: total, PerPage: perPage}),
			),
		),
	)
}

func projectCard(p hpdata.Project, tone token.Tone, lang string) g.Node {
	imgs := p.LocalImages()

	var carouselNode g.Node
	if len(imgs) > 1 {
		carouselNode = uidata.Carousel(uidata.CarouselProps{
			ID:     "c" + slugify(p.Name),
			Images: imgs,
			Alt:    p.Name,
		})
	} else if len(imgs) == 1 {
		carouselNode = h.Img(
			h.Src(imgs[0]),
			h.Alt(p.Name),
			g.Attr("loading", "lazy"),
			h.Style("width:100%;aspect-ratio:16/9;object-fit:contain;background:var(--surface-2);border-bottom:var(--border-w) solid var(--line);display:block"),
		)
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
	topicNodes := make([]g.Node, 0, 4)
	for _, t := range allTags {
		if len(topicNodes) >= 4 {
			break
		}
		topicNodes = append(topicNodes, primitive.Tag(
			primitive.TagProps{Icon: hpdata.TechIcon(t), Tone: hpdata.TagTone(t)},
			g.Text(t),
		))
	}

	return layout.Col(layout.ColProps{Span: 4},
		primitive.Card(primitive.CardProps{Tone: tone},
			g.If(carouselNode != nil, carouselNode),
			h.Div(h.Style("display:flex;justify-content:space-between;align-items:flex-start;gap:var(--sp-2)"),
				primitive.Heading(primitive.HeadingProps{Level: 3},
					g.Text(friendlyName(p.Name))),
				g.If(p.Stars > 0,
					primitive.Tag(primitive.TagProps{Tone: token.ToneYellow},
						g.Text(fmt.Sprintf("★ %d", p.Stars))),
				),
			),
			h.P(h.Style("margin:0;font-size:var(--t-sm)"), g.Text(truncate(p.DescFor(lang), 200))),
			h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
				g.Group(topicNodes),
			),
			g.If(len(linkNodes) > 0,
				h.Div(h.Style("display:flex;gap:var(--sp-1);flex-wrap:wrap;margin-top:auto"), g.Group(linkNodes)),
			),
		),
	)
}
