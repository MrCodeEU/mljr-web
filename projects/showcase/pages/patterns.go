//go:build showcase

package pages

import (
	"fmt"

	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// PatternsListing renders the /patterns overview page.
func PatternsListing() g.Node {
	patterns := registry.AllPatterns()

	cards := make([]g.Node, len(patterns))
	for i, p := range patterns {
		slug := p.Slug
		cards[i] = h.Div(
			g.Attr("data-component", "card"),
			h.Style("overflow:hidden;cursor:pointer"),
			g.Attr("data-on:click", fmt.Sprintf("window.location='/patterns/%s'", slug)),
			// Iframe preview
			h.Div(
				h.Style("height:240px;overflow:hidden;pointer-events:none;border-bottom:var(--bw-2) solid var(--ink)"),
				h.IFrame(
					h.Style("width:200%;height:200%;border:none;transform:scale(0.5);transform-origin:top left"),
					g.Attr("loading", "lazy"),
					g.Attr("data-attr", fmt.Sprintf(`{"src":"/patterns/%s/preview?theme="+$theme+"&mode="+$mode}`, slug)),
					h.Src(fmt.Sprintf("/patterns/%s/preview", slug)),
				),
			),
			// Info
			h.Div(
				h.Style("padding:var(--sp-4)"),
				h.Div(
					h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-1)"),
					h.Strong(h.Style("font-size:var(--t-base)"), g.Text(p.Name)),
					h.Span(
						g.Attr("data-component", "badge"),
						h.Style("font-size:var(--t-xs)"),
						g.Text(p.Category),
					),
				),
				h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin:0"), g.Text(p.Description)),
			),
		)
	}

	var emptyState g.Node
	if len(patterns) == 0 {
		emptyState = h.Div(
			h.Style("padding:var(--sp-12) 0;text-align:center;color:var(--muted)"),
			g.Text("No patterns registered yet."),
		)
	}

	return layout.PageShell(
		layout.PageProps{Title: "Patterns — mljr-ui", Theme: token.ThemeSwissBrut, Mode: token.ModeLight},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		layout.Navbar(layout.NavbarProps{},
			g.Text("mljr-ui · patterns"),
			g.Group{
				h.A(h.Href("/"), g.Text("Catalogue")),
				h.A(h.Href("/patterns"), h.Style("font-weight:800"), g.Text("Patterns")),
			},
			g.Group{special.ThemeToggle(), special.ModeToggle()},
		),
		h.Main(
			layout.Container(layout.ContainerProps{},
				primitive.Display(primitive.DisplayProps{}, g.Text("UI "), h.Em(g.Text("patterns"))),
				h.P(
					h.Style("color:var(--muted);margin-bottom:var(--sp-8);max-width:56ch"),
					g.Text("Full-page compositions showing multiple components working together. Real-world layouts ready to copy."),
				),
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(min(340px,100%),1fr));gap:var(--sp-6)"),
					g.Group(cards),
					emptyState,
				),
			),
		),
	)
}

// PatternDetail renders /patterns/{slug} — full detail with open-in-new-tab link.
func PatternDetail(p *registry.Pattern, theme token.Theme, mode token.Mode) g.Node {
	if theme == "" {
		theme = token.ThemeSwissBrut
	}
	if mode == "" {
		mode = token.ModeLight
	}

	return layout.PageShell(
		layout.PageProps{Title: p.Name + " — Patterns", Theme: theme, Mode: mode},
		special.ThemeToggleRoot(theme, mode),
		layout.Navbar(layout.NavbarProps{},
			h.A(h.Href("/patterns"), g.Text("← Patterns")),
			g.Group{
				h.Span(h.Style("font-weight:800"), g.Text(p.Name)),
				h.Span(g.Attr("data-component", "badge"), g.Text(p.Category)),
			},
			g.Group{
				h.A(
					h.Href(fmt.Sprintf("/patterns/%s/preview", p.Slug)),
					h.Target("_blank"),
					g.Attr("rel", "noopener"),
					h.Style("font-size:var(--t-sm);font-weight:600"),
					g.Text("Open full page ↗"),
				),
				special.ThemeToggle(), special.ModeToggle(),
			},
		),
		h.Main(
			h.Style("display:flex;flex-direction:column;gap:var(--sp-4);padding:var(--sp-6) 0"),
			layout.Container(layout.ContainerProps{},
				h.P(h.Style("color:var(--muted);margin:0"), g.Text(p.Description)),
			),
			h.Div(
				h.Style("border:var(--bw-2) solid var(--ink);border-radius:var(--radius);overflow:hidden;margin:0 var(--sp-4)"),
				h.IFrame(
					h.Style("width:100%;height:80vh;border:none"),
					g.Attr("data-attr", fmt.Sprintf(`{"src":"/patterns/%s/preview?theme="+$theme+"&mode="+$mode}`, p.Slug)),
					h.Src(fmt.Sprintf("/patterns/%s/preview?theme=%s&mode=%s", p.Slug, theme, mode)),
				),
			),
		),
	)
}

// PatternPreview renders just the pattern content (used in iframes).
func PatternPreview(p *registry.Pattern, theme token.Theme, mode token.Mode) g.Node {
	if theme == "" {
		theme = token.ThemeSwissBrut
	}
	if mode == "" {
		mode = token.ModeLight
	}
	return p.Render(string(theme), string(mode))
}
