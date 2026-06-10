package data

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type StatCardProps struct {
	Label   string
	Value   string
	Delta   string // e.g. "+12%" — empty hides the slot
	DeltaUp bool   // true=up/green, false=down/red (only relevant when Delta != "")
	Attrs   []g.Node
}

// StatCard displays a single key metric with label, value, and optional delta.
func StatCard(p StatCardProps) g.Node {
	var delta g.Node
	if p.Delta != "" {
		state := "down"
		if p.DeltaUp {
			state = "up"
		}
		delta = h.Span(g.Attr("data-slot", "delta"), g.Attr("data-state", state), g.Text(p.Delta))
	}
	return h.Div(
		g.Attr("data-component", "stat-card"),
		g.Group(p.Attrs),
		h.Div(g.Attr("data-slot", "label"), g.Text(p.Label)),
		h.Div(g.Attr("data-slot", "value"), g.Text(p.Value)),
		delta,
	)
}
