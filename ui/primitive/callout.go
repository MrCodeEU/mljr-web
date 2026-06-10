package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CalloutVariant string

const (
	CalloutDefault CalloutVariant = ""
	CalloutInfo    CalloutVariant = "info"
	CalloutSuccess CalloutVariant = "success"
	CalloutWarning CalloutVariant = "warning"
	CalloutDanger  CalloutVariant = "danger"
)

type CalloutProps struct {
	Variant CalloutVariant
	Title   string
	Attrs   []g.Node
}

// Callout renders a left-bordered highlight block for notes, tips, or warnings.
func Callout(p CalloutProps, children ...g.Node) g.Node {
	nodes := []g.Node{
		g.Attr("data-component", "callout"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Group(p.Attrs),
	}
	if p.Title != "" {
		nodes = append(nodes, h.Div(g.Attr("data-slot", "title"), g.Text(p.Title)))
	}
	nodes = append(nodes, children...)
	return h.Div(nodes...)
}
