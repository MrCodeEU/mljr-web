package pages

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// Sub pages (games, legal, project details) use the homepage's sheet: a
// 3-column anchor with back link, title and lede beside a 9-column body.
// Markup: main.sub-page > subHead(...) + .sub-body.

// subHead renders the anchor column. back is where the back link leads,
// crumb names it (e.g. "Labor"); mark picks the Bauhaus composition; head
// holds the h1 and any lede or links.
func subHead(back, crumb string, mark int, head ...g.Node) g.Node {
	return h.Div(h.Class("sub-head"),
		h.Div(h.Class("sub-head-in"),
			h.A(h.Href(back), h.Class("sub-back"), g.Text("← mljr.eu / "+crumb)),
			g.Group(head),
		),
		bauhausMark(mark),
	)
}

// subRow is a further row on a sub page: a small anchor title and a body.
// The heading names the section, which makes it a region landmark.
func subRow(class, title, sub string, body ...g.Node) g.Node {
	id := "row-" + strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		return '-'
	}, strings.ToLower(title))
	return h.Section(h.Class(strings.TrimSpace("sub-page sub-row "+class)), g.Attr("aria-labelledby", id),
		h.Div(h.Class("sub-head"),
			h.Div(h.Class("sub-head-in"),
				h.H2(h.ID(id), h.Class("sub-row-title"), g.Text(title)),
				g.If(sub != "", swiss.Micro(g.Text(sub))),
			),
		),
		h.Div(append([]g.Node{h.Class("sub-body")}, body...)...),
	)
}

const subCSS = `
.sub-page { display: grid; grid-template-columns: repeat(12, minmax(0, 1fr)); margin: 0 var(--swiss-gutter); border-bottom: 3px solid var(--ink); min-height: 0; padding: 0; }
.sub-head { grid-column: span 3; border-right: 1px solid var(--ink); padding: var(--swiss-pad) var(--swiss-pad) var(--swiss-pad) 0; display: flex; flex-direction: column; gap: 1.5rem; min-width: 0; }
.sub-head-in { position: sticky; top: calc(var(--swiss-hdr) + 1rem); z-index: 1; display: flex; flex-direction: column; gap: .9rem; }
.sub-head .bh { margin-top: auto; }
.sub-back { font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .05em; color: var(--ink); text-decoration: none; align-self: flex-start; border-bottom: 2px solid var(--act); }
.sub-back:hover { background: var(--act); color: var(--act-fg); }
.sub-head h1 { margin: 0; font-size: clamp(1.9rem, 1rem + 1.9vw, 3.2rem); font-weight: 900; line-height: .95; letter-spacing: -.035em; overflow-wrap: break-word; }
.sub-head .game-lede, .sub-lede { margin: 0; color: var(--muted); font-size: .95rem; line-height: 1.5; }
.sub-row-title { margin: 0; font-size: 1.05rem; font-weight: 800; text-transform: uppercase; letter-spacing: .01em; }
.sub-body { grid-column: span 9; min-width: 0; padding: var(--swiss-pad); }
.sub-page .game-shell, .sub-page .game-shell-wide { max-width: none; margin: 0; background: none; border: 0; box-shadow: none; }
@media (max-width: 1024px) {
  .sub-page { grid-template-columns: minmax(0, 1fr); }
  .sub-head { grid-column: auto; border-right: 0; border-bottom: 1px solid var(--ink); padding-right: 0; }
  .sub-head-in { position: static; }
  .sub-body { grid-column: auto; padding-inline: 0; }
}
html.embedded .sub-page { display: block; margin: 0; border: 0; }
html.embedded .sub-head { display: none; }
html.embedded .sub-body { padding: .75rem; }
`

// NotFound is the 404 page: what happened and where to go instead.
func NotFound(lang string, a AnalyticsConfig) g.Node {
	links := []struct{ href, de, en string }{
		{"/#projects", "Werkkatalog", "Catalogue"},
		{"/#playground", "Kinetisches Labor", "Kinetic Lab"},
		{"/#experience", "Laufbahn", "Career"},
		{"/#contact", "Kontakt", "Contact"},
	}
	nodes := make([]g.Node, len(links))
	for i, l := range links {
		nodes[i] = h.Li(swiss.BtnLink(l.href, i == 0, nil, g.Text(tr(lang, l.de, l.en)+" →")))
	}
	return layout.PageShell(
		layout.PageProps{
			Title:     tr(lang, "Seite nicht gefunden", "Page not found") + " - Michael Reinegger",
			Theme:     token.ThemeSwiss,
			AutoMode:  true,
			LockTheme: true,
			Mode:      token.ModeLight,
			Lang:      lang,
			HeadExtra: append([]g.Node{g.El("style", g.Raw(homepageCSS+notFoundCSS)), h.Meta(h.Name("robots"), h.Content("noindex"))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"), h.Class("sub-page nf"),
			subHead("/", tr(lang, "Start", "Home"), 4,
				h.P(h.Class("nf-code"), g.Attr("aria-hidden", "true"), g.Text("404")),
				h.H1(g.Text(tr(lang, "Seite nicht gefunden", "Page not found"))),
			),
			h.Div(h.Class("sub-body"),
				h.P(h.Class("sub-lede"), g.Text(tr(lang,
					"Diese Adresse gibt es (nicht mehr). Vielleicht hilft einer dieser Wege weiter:",
					"This address does not exist (any more). One of these may help:"))),
				h.Ul(h.Class("nf-links"), g.Group(nodes)),
			),
		),
		siteFooter(lang),
	)
}

const notFoundCSS = `
.nf { min-height: 60vh; }
.nf-code { margin: 0; font-size: clamp(4rem, 3rem + 5vw, 9rem); font-weight: 900; line-height: .85; letter-spacing: -.05em; color: var(--accent); }
.nf-links { list-style: none; margin: 1.25rem 0 0; padding: 0; display: flex; flex-wrap: wrap; gap: .5rem; }
`
