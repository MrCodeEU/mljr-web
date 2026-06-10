package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type StackProps struct {
	Axis    string // "v" (default) | "h"
	Gap     string // "" | "sm" | "md" | "lg"
	Align   string // "" | "center"
	Justify string // "" | "between"
	Attrs   []g.Node
}

func Stack(p StackProps, children ...g.Node) g.Node {
	if p.Axis == "" {
		p.Axis = "v"
	}
	return h.Div(
		g.Attr("data-component", "stack"),
		g.Attr("data-axis", p.Axis),
		g.If(p.Gap != "", g.Attr("data-gap", p.Gap)),
		g.If(p.Align != "", g.Attr("data-align", p.Align)),
		g.If(p.Justify != "", g.Attr("data-justify", p.Justify)),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
