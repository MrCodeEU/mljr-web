package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SwitchProps struct {
	Label   string
	Signal  string
	Name    string
	Checked bool
	Attrs   []g.Node
}

// Switch renders a styled toggle switch bound to a boolean Datastar signal.
func Switch(p SwitchProps, attrs ...g.Node) g.Node {
	return h.Label(
		g.Attr("data-component", "switch"),
		h.Input(
			h.Type("checkbox"),
			g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
			g.If(p.Name != "", h.Name(p.Name)),
			g.If(p.Checked, g.Attr("checked")),
			g.Group(p.Attrs),
			g.Group(attrs),
		),
		h.Div(
			g.Attr("data-slot", "track"),
			h.Div(g.Attr("data-slot", "thumb")),
		),
		g.If(p.Label != "", h.Span(g.Attr("data-slot", "label"), g.Text(p.Label))),
	)
}
