package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NumberInputProps struct {
	Signal string
	Name   string
	Min    int
	Max    int
	Step   int
	Value  int
	Attrs  []g.Node
}

// NumberInput renders a stepper-style number input with − and + buttons.
// Signal drives the value via Datastar two-way binding.
func NumberInput(p NumberInputProps, attrs ...g.Node) g.Node {
	if p.Step == 0 {
		p.Step = 1
	}
	sig := p.Signal

	decExpr := fmt.Sprintf("$%s=Math.max(%d,$%s-%d)", sig, p.Min, sig, p.Step)
	incExpr := fmt.Sprintf("$%s=Math.min(%d,$%s+%d)", sig, p.Max, sig, p.Step)

	return h.Div(
		g.Attr("data-component", "number-input"),
		g.Group(p.Attrs),
		h.Button(
			g.Attr("data-slot", "dec"),
			h.Type("button"),
			g.Attr("data-on:click", decExpr),
			g.Text("−"),
		),
		h.Input(
			h.Type("number"),
			g.Attr("data-bind:"+sig),
			h.Value(fmt.Sprintf("%d", p.Value)),
			g.If(p.Name != "", h.Name(p.Name)),
			h.Min(fmt.Sprintf("%d", p.Min)),
			h.Max(fmt.Sprintf("%d", p.Max)),
			h.Step(fmt.Sprintf("%d", p.Step)),
			g.Group(attrs),
		),
		h.Button(
			g.Attr("data-slot", "inc"),
			h.Type("button"),
			g.Attr("data-on:click", incExpr),
			g.Text("+"),
		),
	)
}
