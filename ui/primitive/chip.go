package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ChipProps struct {
	Tone      token.Tone
	Dismiss   bool   // show × button
	OnDismiss string // Datastar expression on dismiss click
	Attrs     []g.Node
}

// Chip is a dismissible tag-shaped label.
// Default dismiss removes the chip from the DOM; provide OnDismiss to override.
func Chip(p ChipProps, children ...g.Node) g.Node {
	dismissExpr := p.OnDismiss
	if dismissExpr == "" {
		dismissExpr = "evt.target.closest('[data-component=chip]').remove()"
	}

	var dismissBtn g.Node
	if p.Dismiss {
		dismissBtn = h.Button(
			g.Attr("data-slot", "dismiss"),
			h.Type("button"),
			g.Attr("aria-label", "Remove"),
			g.Attr("data-on:click", dismissExpr),
			g.Text("×"),
		)
	}

	nodes := []g.Node{
		g.Attr("data-component", "chip"),
		g.If(p.Tone != "", g.Attr("data-tone", string(p.Tone))),
		g.Group(p.Attrs),
	}
	nodes = append(nodes, children...)
	if dismissBtn != nil {
		nodes = append(nodes, dismissBtn)
	}
	return h.Span(nodes...)
}
