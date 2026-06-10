package data

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ListVariant string

const (
	ListDefault ListVariant = ""
	ListDivided ListVariant = "divided"
	ListOrdered ListVariant = "ordered"
)

type ListProps struct {
	Variant ListVariant
	Attrs   []g.Node
}

// List renders a styled list. Use ListDivided for bordered rows.
func List(p ListProps, items ...g.Node) g.Node {
	if p.Variant == ListOrdered {
		return h.Ol(
			g.Attr("data-component", "list"),
			g.Attr("data-variant", "ordered"),
			g.Group(p.Attrs),
			g.Group(items),
		)
	}
	return h.Ul(
		g.Attr("data-component", "list"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Group(p.Attrs),
		g.Group(items),
	)
}

// ListItem renders a single list entry.
func ListItem(children ...g.Node) g.Node {
	return h.Li(g.Attr("data-component", "list-item"), g.Group(children))
}

// DescriptionItem is a term/definition pair for DescriptionList.
type DescriptionItem struct {
	Term string
	Desc g.Node
}

type DescriptionListProps struct {
	Attrs []g.Node
}

// DescriptionList renders a dl with aligned term/description pairs.
func DescriptionList(p DescriptionListProps, items ...DescriptionItem) g.Node {
	rows := make([]g.Node, 0, len(items)*2)
	for _, item := range items {
		rows = append(rows,
			h.Dt(g.Attr("data-component", "dt"), g.Text(item.Term)),
			h.Dd(g.Attr("data-component", "dd"), item.Desc),
		)
	}
	return h.Dl(
		g.Attr("data-component", "dl"),
		g.Group(p.Attrs),
		g.Group(rows),
	)
}
