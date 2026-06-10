package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BreadcrumbItem struct {
	Label string
	Href  string // empty = current page (last item)
}

type BreadcrumbProps struct {
	Attrs []g.Node
}

// Breadcrumb renders a navigation trail from BreadcrumbItems.
// Last item is treated as current page (no link, bold).
func Breadcrumb(p BreadcrumbProps, items ...BreadcrumbItem) g.Node {
	nodes := make([]g.Node, len(items))
	for i, item := range items {
		var inner g.Node
		if item.Href != "" {
			inner = h.A(h.Href(item.Href), g.Text(item.Label))
		} else {
			inner = h.Span(g.Text(item.Label))
		}
		nodes[i] = h.Li(g.Attr("data-component", "breadcrumb-item"), inner)
	}
	return h.Nav(
		g.Attr("data-component", "breadcrumb"),
		g.Attr("aria-label", "Breadcrumb"),
		g.Group(p.Attrs),
		h.Ol(g.Group(nodes)),
	)
}
