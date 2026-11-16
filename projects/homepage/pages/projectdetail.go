package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// HasProjectDetail reports whether a project has a dedicated detail page.
func HasProjectDetail(p hpdata.Project) bool {
	return p.HasDetailPage()
}

// ProjectDetail renders the full architecture deep-dive page for one project.
// All content (long description, Mermaid diagram, code snippets) is
// authored in mljr-data/projects.json and flows through hpdata.Project —
// this page has no hardcoded per-project content of its own.
func ProjectDetail(d hpdata.SiteData, p hpdata.Project, lang string, a AnalyticsConfig) g.Node {
	longDesc := p.LongDescFor(lang)

	imgs := p.LocalImages()
	var hero g.Node
	if len(imgs) > 1 {
		hero = h.Div(h.Class("pd-hero-carousel"),
			uidata.Carousel(uidata.CarouselProps{ID: "pd-" + p.ID, Images: imgs, Alt: p.Name}),
		)
	} else if len(imgs) == 1 {
		hero = h.Img(h.Src(imgs[0]), h.Alt(p.Name), g.Attr("loading", "lazy"),
			h.Style("width:100%;aspect-ratio:16/9;object-fit:contain;background:var(--surface-2);border:1px solid var(--ink)"))
	}

	var linkNodes []g.Node
	if p.URL != "" {
		linkNodes = append(linkNodes, swiss.BtnLink(p.URL, true, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")},
			icon.Icon("simple-icons:github", icon.Props{Size: "1em"}), g.Text(i18n.T(lang, "projects.source")+" ↗")))
	}
	for _, lnk := range p.Links {
		if lnk.URL == "" {
			continue
		}
		linkNodes = append(linkNodes, swiss.BtnLink(lnk.URL, false, []g.Node{g.Attr("target", "_blank"), g.Attr("rel", "noopener")},
			g.Text(lnk.Name+" ↗")))
	}

	var snippetNodes []g.Node
	for _, s := range p.Snippets {
		snippetNodes = append(snippetNodes,
			h.Div(h.Style("margin-bottom:var(--sp-6)"),
				h.P(h.Style("margin:0 0 var(--sp-3);font-size:var(--t-sm);font-weight:700;color:var(--muted);line-height:1.5"),
					g.Text(s.Caption)),
				uidata.SyntaxHighlighter(uidata.SyntaxHighlighterProps{
					Language: s.Language,
					Theme:    "monokai",
					Filename: s.Filename,
				}, s.Code),
			),
		)
	}

	headExtra := append([]g.Node{
		g.El("style", g.Raw(homepageCSS+legalCSS+projectDetailCSS)),
	}, AnalyticsHead(a)...)
	if p.Diagram != "" {
		// Self-hosted (CSP is script-src 'self', no CDN imports allowed).
		// mermaid.min.js is the UMD bundle; it sets globalThis.mermaid.
		headExtra = append(headExtra,
			h.Script(h.Src("/static/mermaid.min.js")),
			h.Script(g.Raw(`mermaid.initialize({startOnLoad:true,theme:'neutral'});`)),
		)
	}

	return layout.PageShell(
		layout.PageProps{
			Title:       friendlyName(p.Name) + " - Michael Reinegger",
			Description: p.DescFor(lang),
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   headExtra,
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"), h.Class("pd-page sub-page"),
			subHead("/#projects", tr(lang, "Werkkatalog", "Catalogue"), len(p.ID),
				h.H1(g.Text(friendlyName(p.Name))),
				h.P(h.Class("sub-lede"), g.Text(p.DescFor(lang))),
				g.If(len(linkNodes) > 0, h.Div(h.Class("cat-links"), g.Group(linkNodes))),
			),
			h.Div(h.Class("sub-body pd-body"),
				g.If(hero != nil, h.Div(h.Class("pd-hero"), hero)),
				h.P(h.Class("pd-long"), g.Text(longDesc)),
			),
		),
		g.If(p.Diagram != "", subRow("", tr(lang, "Architektur", "Architecture"), "",
			h.Pre(h.Class("mermaid"), g.Raw(p.Diagram)),
		)),
		g.If(len(snippetNodes) > 0, subRow("", tr(lang, "Unter der Haube", "Under the hood"), "",
			g.Group(snippetNodes),
		)),
		siteFooter(lang),
	)
}

// The shared carousel component fixes its track at 220px tall, sized for
// compact project-grid cards. The detail-page hero is much wider (up to
// 920px), so the same height letterboxes screenshots badly — override it
// here with a taller, more cinematic track just for this context.
const projectDetailCSS = `
.pd-body { display: flex; flex-direction: column; gap: var(--swiss-pad); }
.pd-long { margin: 0; font-size: 1.05rem; line-height: 1.7; max-width: 72ch; }
.sub-head .cat-links { display: flex; flex-wrap: wrap; gap: .4rem; }
.sub-row .mermaid { margin: 0; background: none; }
.pd-hero-carousel [data-component="carousel"] [data-slot="track"] {
  height: clamp(280px, 45vw, 480px);
}
`
