package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SelectOption struct {
	Value string
	Label string
}

type SelectProps struct {
	Options  []SelectOption
	Signal   string
	Name     string
	Required bool
	Attrs    []g.Node
}

// Select renders a styled <select> inside a wrapper div that provides the custom arrow.
func Select(p SelectProps, attrs ...g.Node) g.Node {
	opts := make([]g.Node, len(p.Options))
	for i, o := range p.Options {
		opts[i] = h.Option(h.Value(o.Value), g.Text(o.Label))
	}
	return h.Div(
		g.Attr("data-component", "select-wrap"),
		h.Select(
			g.Attr("data-component", "select"),
			g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
			g.If(p.Name != "", h.Name(p.Name)),
			g.If(p.Required, g.Attr("required")),
			g.Group(p.Attrs),
			g.Group(attrs),
			g.Group(opts),
		),
	)
}
