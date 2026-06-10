package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Tab struct {
	Slug  string
	Label string
	Body  g.Node
}

type TabsProps struct {
	Signal  string // Datastar signal name (default "tab")
	Default string // default active slug (defaults to first tab)
	Attrs   []g.Node
}

// Tabs renders a Datastar-driven tabbed panel. Signal switches the active tab.
func Tabs(p TabsProps, tabs []Tab) g.Node {
	if p.Signal == "" {
		p.Signal = "tab"
	}
	defaultSlug := p.Default
	if defaultSlug == "" && len(tabs) > 0 {
		defaultSlug = tabs[0].Slug
	}

	sig := p.Signal
	buttons := make([]g.Node, len(tabs))
	panels := make([]g.Node, len(tabs))

	for i, t := range tabs {
		slug := t.Slug
		isDefault := slug == defaultSlug
		activeExpr := fmt.Sprintf(`$%s===%q`, sig, slug)

		buttons[i] = h.Button(
			g.Attr("data-component", "tab-btn"),
			g.Attr("data-attr", fmt.Sprintf(`{"data-state":%s?"active":""}`, activeExpr)),
			g.Attr("data-on:click", fmt.Sprintf(`$%s=%q`, sig, slug)),
			h.Type("button"),
			g.Text(t.Label),
		)
		panels[i] = h.Div(
			g.Attr("data-component", "tab-panel"),
			g.Attr("data-show", activeExpr),
			g.If(!isDefault, h.Style("display:none")),
			t.Body,
		)
	}

	return h.Div(
		g.Attr("data-component", "tabs"),
		g.Attr("data-signals", fmt.Sprintf(`{%s:%q}`, sig, defaultSlug)),
		g.Group(p.Attrs),
		h.Div(g.Attr("data-component", "tab-list"), g.Group(buttons)),
		h.Div(g.Attr("data-component", "tab-panels"), g.Group(panels)),
	)
}
