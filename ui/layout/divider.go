package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DividerProps struct {
	Attrs []g.Node
}

// Divider renders a horizontal rule with an optional centered text label.
func Divider(p DividerProps, label ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "divider"),
		g.Group(p.Attrs),
		g.Group(label),
	)
}
