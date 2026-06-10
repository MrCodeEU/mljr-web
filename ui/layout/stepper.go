package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type StepState string

const (
	StepComplete StepState = "complete"
	StepActive   StepState = "active"
	StepUpcoming StepState = ""
)

type Step struct {
	Label string
	State StepState
}

type StepperProps struct {
	Attrs []g.Node
}

// Stepper renders a horizontal multi-step progress indicator.
func Stepper(p StepperProps, steps ...Step) g.Node {
	nodes := make([]g.Node, len(steps))
	for i, s := range steps {
		num := fmt.Sprintf("%d", i+1)
		dot := g.Text(num)
		if s.State == StepComplete {
			dot = g.Text("✓")
		}
		nodes[i] = h.Div(
			g.Attr("data-component", "step"),
			g.If(s.State != "", g.Attr("data-state", string(s.State))),
			h.Div(g.Attr("data-slot", "dot"), dot),
			h.Div(g.Attr("data-slot", "label"), g.Text(s.Label)),
		)
	}
	return h.Div(
		g.Attr("data-component", "stepper"),
		g.Group(p.Attrs),
		g.Group(nodes),
	)
}
