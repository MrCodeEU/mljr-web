package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AccordionProps struct {
	Attrs []g.Node
}

type AccordionItemProps struct {
	Title string
	Open  bool
	Attrs []g.Node
}

// Accordion groups AccordionItem elements into a bordered container.
func Accordion(p AccordionProps, items ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "accordion"),
		g.Group(p.Attrs),
		g.Group(items),
	)
}

// AccordionItem renders a <details>/<summary> pair — no JS required.
func AccordionItem(p AccordionItemProps, children ...g.Node) g.Node {
	return g.El("details",
		g.Attr("data-component", "accordion-item"),
		g.If(p.Open, g.Attr("open")),
		g.Group(p.Attrs),
		g.El("summary", g.Text(p.Title)),
		h.Div(g.Attr("data-slot", "content"), g.Group(children)),
	)
}
