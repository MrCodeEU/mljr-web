package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SegmentedOption struct {
	Value string
	Label string
}

type SegmentedProps struct {
	Name    string
	Options []SegmentedOption
	Default string // default selected value
}

// Segmented renders a radio-group styled as a connected button bar.
// Use with a <form> or read the checked radio's value via JS/Datastar.
func Segmented(p SegmentedProps, extra ...g.Node) g.Node {
	if p.Name == "" {
		p.Name = "segmented"
	}
	items := make([]g.Node, 0, len(p.Options)*2)
	for i, opt := range p.Options {
		id := fmt.Sprintf("%s-%d", p.Name, i)
		checked := opt.Value == p.Default
		items = append(items,
			h.Input(
				h.ID(id),
				h.Name(p.Name),
				h.Type("radio"),
				h.Value(opt.Value),
				g.If(checked, g.Attr("checked", "")),
			),
			h.Label(
				h.For(id),
				g.Text(opt.Label),
			),
		)
	}
	return h.Div(
		g.Attr("data-component", "segmented"),
		h.Role("group"),
		g.Group(extra),
		g.Group(items),
	)
}
