package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type InputGroupProps struct {
	Prefix g.Node // text or icon rendered before the input
	Suffix g.Node // text or icon rendered after the input
	Attrs  []g.Node
}

// InputGroup wraps an Input (or other form control) with prefix/suffix addons.
// Place the actual Input/Select as the child.
func InputGroup(p InputGroupProps, control g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "input-group"),
		g.Group(p.Attrs),
		g.If(p.Prefix != nil, h.Span(g.Attr("data-slot", "prefix"), p.Prefix)),
		control,
		g.If(p.Suffix != nil, h.Span(g.Attr("data-slot", "suffix"), p.Suffix)),
	)
}
