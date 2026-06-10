package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CheckboxProps struct {
	Label   string
	Signal  string
	Name    string
	Checked bool
	Attrs   []g.Node
}

// Checkbox renders a styled checkbox with a custom visual box and optional label.
// Signal binds to a boolean Datastar signal.
func Checkbox(p CheckboxProps, attrs ...g.Node) g.Node {
	return h.Label(
		g.Attr("data-component", "checkbox"),
		h.Input(
			h.Type("checkbox"),
			g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
			g.If(p.Name != "", h.Name(p.Name)),
			g.If(p.Checked, g.Attr("checked")),
			g.Group(p.Attrs),
			g.Group(attrs),
		),
		h.Span(g.Attr("data-slot", "box")),
		g.If(p.Label != "", h.Span(g.Attr("data-slot", "label"), g.Text(p.Label))),
	)
}

type RadioOption struct {
	Value string
	Label string
}

type RadioProps struct {
	Label  string
	Value  string
	Signal string
	Name   string
	Attrs  []g.Node
}

// Radio renders a single styled radio button. Group multiple with RadioGroup.
func Radio(p RadioProps, attrs ...g.Node) g.Node {
	return h.Label(
		g.Attr("data-component", "radio"),
		h.Input(
			h.Type("radio"),
			g.If(p.Value != "", h.Value(p.Value)),
			g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
			g.If(p.Name != "", h.Name(p.Name)),
			g.Group(p.Attrs),
			g.Group(attrs),
		),
		h.Span(g.Attr("data-slot", "box")),
		g.If(p.Label != "", h.Span(g.Attr("data-slot", "label"), g.Text(p.Label))),
	)
}

type RadioGroupProps struct {
	Signal  string
	Name    string
	Options []RadioOption
	Attrs   []g.Node
}

// RadioGroup renders a vertical group of Radio buttons bound to a shared signal.
func RadioGroup(p RadioGroupProps, attrs ...g.Node) g.Node {
	opts := make([]g.Node, len(p.Options))
	for i, o := range p.Options {
		opts[i] = Radio(RadioProps{
			Label:  o.Label,
			Value:  o.Value,
			Signal: p.Signal,
			Name:   p.Name,
		})
	}
	return h.Div(
		g.Attr("data-component", "radio-group"),
		g.Group(p.Attrs),
		g.Group(attrs),
		g.Group(opts),
	)
}
