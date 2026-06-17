package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

type projectSnippet struct {
	Caption  string
	Filename string
	Language string
	Code     string
}

type projectDetailContent struct {
	LongDesc   string
	LongDescDE string
	Diagram    string // Mermaid graph definition
	Snippets   []projectSnippet
}

// projectDetails holds deep-dive content for flagship projects, keyed by
// curated project id. Cards for ids present here link to /projects/<id>;
// everything else stays a card-only entry.
var projectDetails = map[string]projectDetailContent{
	"godrive":            godriveDetail,
	"homelab-automation": homelabAutomationDetail,
	"mljr-web":           mljrWebDetail,
	"nightscout-tray":    nightscoutTrayDetail,
}

// HasProjectDetail reports whether a project has a dedicated detail page.
func HasProjectDetail(id string) bool {
	_, ok := projectDetails[id]
	return ok
}

// ProjectDetail renders the full architecture deep-dive page for one project.
func ProjectDetail(d hpdata.SiteData, p hpdata.Project, lang string, a AnalyticsConfig) g.Node {
	det, ok := projectDetails[p.ID]
	if !ok {
		det = projectDetailContent{LongDesc: p.DescFor(lang)}
	}
	longDesc := det.LongDesc
	if lang == "de" && det.LongDescDE != "" {
		longDesc = det.LongDescDE
	}

	imgs := p.LocalImages()
	var hero g.Node
	if len(imgs) > 1 {
		hero = h.Div(h.Class("pd-hero-carousel"),
			uidata.Carousel(uidata.CarouselProps{ID: "pd-" + p.ID, Images: imgs, Alt: p.Name}),
		)
	} else if len(imgs) == 1 {
		hero = h.Img(h.Src(imgs[0]), h.Alt(p.Name), g.Attr("loading", "lazy"),
			h.Style("width:100%;aspect-ratio:16/9;object-fit:contain;background:var(--surface-2);border:var(--bw-2) solid var(--ink)"))
	}

	var linkNodes []g.Node
	if p.URL != "" {
		linkNodes = append(linkNodes, h.A(h.Href(p.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
			primitive.Button(primitive.ButtonProps{Variant: token.Outline},
				icon.Icon("simple-icons:github"), g.Text(i18n.T(lang, "projects.source")), icon.Icon("lucide:arrow-up-right")),
		))
	}
	for _, lnk := range p.Links {
		if lnk.URL == "" {
			continue
		}
		linkNodes = append(linkNodes, h.A(h.Href(lnk.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
			primitive.Button(primitive.ButtonProps{Variant: token.Outline},
				g.Text(lnk.Name), icon.Icon("lucide:arrow-up-right")),
		))
	}

	var snippetNodes []g.Node
	for _, s := range det.Snippets {
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
		g.El("style", g.Raw(homepageCSS + legalCSS + projectDetailCSS)),
	}, AnalyticsHead(a)...)
	if det.Diagram != "" {
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
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   headExtra,
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("legal-page"),
			h.Div(h.Style("max-width:920px;margin:0 auto"),
				layout.Breadcrumb(layout.BreadcrumbProps{},
					layout.BreadcrumbItem{Label: i18n.T(lang, "nav.projects"), Href: "/#projects"},
					layout.BreadcrumbItem{Label: friendlyName(p.Name)},
				),
				h.H1(h.Style("font-size:clamp(2.4rem,8vw,4.5rem);line-height:.95;margin:var(--sp-4) 0 var(--sp-3);font-weight:950"),
					g.Text(friendlyName(p.Name))),
				h.P(h.Style("font-size:var(--t-lg);color:var(--muted);font-weight:700;margin:0 0 var(--sp-5)"),
					g.Text(p.DescFor(lang))),
				g.If(len(linkNodes) > 0,
					h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap;margin-bottom:var(--sp-6)"), g.Group(linkNodes)),
				),
				g.If(hero != nil, h.Div(h.Style("margin-bottom:var(--sp-8)"), hero)),
				h.Div(h.Class("legal-shell"),
					h.P(h.Style("font-size:var(--t-md);line-height:1.7;max-width:none;margin:0 0 var(--sp-8)"), g.Text(longDesc)),
					g.If(det.Diagram != "",
						h.Div(h.Style("margin-bottom:var(--sp-8)"),
							h.H2(g.Text("Architecture")),
							h.Pre(h.Class("mermaid"), g.Raw(det.Diagram)),
						),
					),
					g.If(len(snippetNodes) > 0,
						h.Div(
							h.H2(g.Text("Under the hood")),
							g.Group(snippetNodes),
						),
					),
				),
			),
		),
		siteFooter(lang),
	)
}

// The shared carousel component fixes its track at 220px tall, sized for
// compact project-grid cards. The detail-page hero is much wider (up to
// 920px), so the same height letterboxes screenshots badly — override it
// here with a taller, more cinematic track just for this context.
const projectDetailCSS = `
.pd-hero-carousel [data-component="carousel"] [data-slot="track"] {
  height: clamp(280px, 45vw, 480px);
}
`
