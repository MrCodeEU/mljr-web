//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "table-of-contents", Name: "Table of Contents", Category: "layout",
		Summary: "Scroll-spy navigation for long-form content. Highlights active section via IntersectionObserver. Can auto-detect headings.",
		Code: `layout.TableOfContents(layout.TOCProps{Sticky: true},
    layout.TOCItem{ID: "intro", Label: "Introduction", Level: 2},
    layout.TOCItem{ID: "usage", Label: "Usage", Level: 2},
    layout.TOCItem{ID: "api", Label: "API Reference", Level: 3},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:200px 1fr;gap:var(--sp-6);align-items:start"),
				TableOfContents(TOCProps{Title: "On this page", Sticky: true},
					TOCItem{ID: "overview", Label: "Overview", Level: 2},
					TOCItem{ID: "install", Label: "Installation", Level: 2},
					TOCItem{ID: "usage", Label: "Usage", Level: 2},
					TOCItem{ID: "variants", Label: "Variants", Level: 3},
					TOCItem{ID: "props", Label: "Props", Level: 3},
					TOCItem{ID: "advanced", Label: "Advanced", Level: 2},
				),
				h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
					h.Div(h.ID("overview"), h.H2(g.Text("Overview")), h.P(g.Text("This component provides automatic scroll-spy tracking for long documents."))),
					h.Div(h.ID("install"), h.H2(g.Text("Installation")), h.P(g.Text("Import the layout package and call TableOfContents with optional items."))),
					h.Div(h.ID("usage"), h.H2(g.Text("Usage")), h.P(g.Text("Pass TOCItem slices or leave empty for auto-detection from headings."))),
					h.Div(h.ID("variants"), h.H3(g.Text("Variants")), h.P(g.Text("Sticky mode keeps the TOC visible while scrolling the content."))),
					h.Div(h.ID("props"), h.H3(g.Text("Props")), h.P(g.Text("ContentSelector targets a CSS selector for the heading observation root."))),
					h.Div(h.ID("advanced"), h.H2(g.Text("Advanced")), h.P(g.Text("Pass no items to auto-populate from h2/h3/h4 elements in the target container."))),
				),
			)
		},
	})
}
