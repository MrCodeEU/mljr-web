package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FlipCardProps struct {
	// Height is the CSS height of the card (default "220px").
	Height string
	// Trigger: "hover" (default) | "click"
	Trigger string
	// Signal is the Datastar signal name for click mode (default "_fc").
	Signal string
}

// FlipCard renders a CSS 3D flip card with a front and back face.
// Hover trigger uses pure CSS; click trigger uses a Datastar signal.
// Usage: FlipCard(props, frontNode, backNode)
func FlipCard(p FlipCardProps, front g.Node, back g.Node) g.Node {
	if p.Height == "" {
		p.Height = "220px"
	}
	if p.Trigger == "" {
		p.Trigger = "hover"
	}

	if p.Trigger == "click" {
		sig := p.Signal
		if sig == "" {
			sig = "_fc"
		}
		return h.Div(
			g.Attr("data-component", "flip-card"),
			g.Attr("data-trigger", "click"),
			g.Attr("data-signals", fmt.Sprintf(`{"%s":false}`, sig)),
			g.Attr("data-attr", fmt.Sprintf(`{"data-flipped":""+$%s}`, sig)),
			g.Attr("data-on:click", fmt.Sprintf("$%s=!$%s", sig, sig)),
			h.Style("height:"+p.Height),
			h.Div(g.Attr("data-slot", "inner"),
				h.Div(g.Attr("data-slot", "front"), front),
				h.Div(g.Attr("data-slot", "back"), back),
			),
		)
	}

	return h.Div(
		g.Attr("data-component", "flip-card"),
		g.Attr("data-trigger", "hover"),
		h.Style("height:"+p.Height),
		h.Div(g.Attr("data-slot", "inner"),
			h.Div(g.Attr("data-slot", "front"), front),
			h.Div(g.Attr("data-slot", "back"), back),
		),
	)
}
