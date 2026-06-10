package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ButtonGroupProps struct {
	Attached bool   // true = connected border (shared edge), false = gap between buttons
	Label    string // aria-label for the group
}

// ButtonGroup wraps a set of buttons in a role="group" container.
// Attached mode removes border radius on inner edges and collapses margins.
// Pass primitive.Button nodes as children.
func ButtonGroup(p ButtonGroupProps, children ...g.Node) g.Node {
	attrs := []g.Node{
		g.Attr("data-component", "button-group"),
		h.Role("group"),
	}
	if p.Label != "" {
		attrs = append(attrs, g.Attr("aria-label", p.Label))
	}
	if p.Attached {
		attrs = append(attrs, g.Attr("data-attached", ""))
	}
	return h.Div(append(attrs, g.Group(children))...)
}
