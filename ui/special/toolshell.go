package special

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui/layout"
	"mljr-web/ui/swiss"
	"mljr-web/ui/token"
)

// ToolShellProps configures a standalone live tool page (regex, cron, codec).
type ToolShellProps struct {
	ID          string // tool id from liveTools; its own link is left out
	Name        string // e.g. "Regex Lab", shown beside the brand
	Title       string // <title>
	Description string // meta description
	CSS         string // page-specific styles
	// Sections feed the colophon ruler and the mobile register; they must
	// match the IDs of the swiss.Section rows passed as body.
	Sections []swiss.NavItem
}

// ToolShell renders a live tool in the mljr.eu Swiss layout: colophon with
// section ruler, sibling tools and light/dark toggle, the 12-column sheet
// with the given section rows, and the slim footer linking back to mljr.eu.
func ToolShell(p ToolShellProps, sections ...g.Node) g.Node {
	actions := []g.Node{}
	for _, t := range liveTools {
		if t.ID != p.ID {
			actions = append(actions, swiss.BtnLink(t.URL, false, nil, g.Text(t.Label)))
		}
	}
	actions = append(actions, swiss.Btn(false, []g.Node{
		g.Attr("data-on:click", `$mode = $mode === 'light' ? 'dark' : 'light'`),
		g.Attr("data-attr:aria-pressed", `$mode === 'dark' ? 'true' : 'false'`),
		g.Attr("aria-label", "Dark mode"),
	}, h.Span(g.Attr("data-text", `$mode === 'dark' ? 'Light' : 'Dark'`), g.Text("Dark"))))

	nav := swiss.ColophonProps{
		Brand:      "mljr.eu",
		BrandHref:  "https://mljr.eu",
		BrandSub:   "— " + p.Name,
		Items:      p.Sections,
		Actions:    actions,
		PageID:     "page",
		NavLabel:   "Sections",
		CloseLabel: "Close",
		SkipLabel:  "Skip to content",
	}
	links := []swiss.FooterLink{{Label: "mljr.eu ↗", Href: "https://mljr.eu"}}
	for _, t := range liveTools {
		links = append(links, swiss.FooterLink{Label: t.Label, Href: t.URL})
	}
	links = append(links,
		swiss.FooterLink{Label: "GitHub ↗", Href: "https://github.com/MrCodeEU"},
		swiss.FooterLink{Label: "Impressum", Href: "https://mljr.eu/impressum"},
		swiss.FooterLink{Label: "Datenschutz", Href: "https://mljr.eu/datenschutz"},
	)

	return layout.PageShell(
		layout.PageProps{
			Title:       p.Title,
			Description: p.Description,
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			HeadExtra:   []g.Node{g.If(p.CSS != "", g.El("style", g.Raw(p.CSS)))},
		},
		ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		swiss.Colophon(nav),
		h.Div(h.ID("page"),
			swiss.Page("sheet", nil, h.Main(append([]g.Node{h.ID("main")}, sections...)...)),
			swiss.Footer("A live tool by Michael Reinegger. Runs server-side in Go; input is not stored or logged.", links,
				"© Michael Reinegger · mljr.eu", "Go · gomponents · Datastar"),
		),
		swiss.Register(nav),
	)
}
