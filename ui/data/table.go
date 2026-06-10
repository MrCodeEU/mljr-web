package data

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TableProps struct {
	Striped bool
	Hover   bool
	Attrs   []g.Node
}

// Table wraps a standard HTML table with design-system styling.
// Use h.THead / h.TBody / h.Tr / h.Th / h.Td for structure.
func Table(p TableProps, children ...g.Node) g.Node {
	var variant string
	if p.Striped {
		variant = "striped"
	}
	return h.Div(
		g.Attr("data-component", "table-wrap"),
		h.Table(
			g.Attr("data-component", "table"),
			g.If(variant != "", g.Attr("data-variant", variant)),
			g.If(p.Hover, g.Attr("data-state", "hover")),
			g.Group(p.Attrs),
			g.Group(children),
		),
	)
}
