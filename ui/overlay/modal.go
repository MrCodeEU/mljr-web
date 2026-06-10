package overlay

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ModalProps struct {
	ID       string     // root element id (default "modal")
	Size     token.Size // sm | md (default) | lg
	OpenExpr string     // Datastar expression toggling visibility (default "$modalOpen")
	Title    string
	Footer   g.Node
	Attrs    []g.Node
}

// Modal renders an always-present scrim+dialog gated by data-show.
// Pair with a trigger button that flips the open signal.
func Modal(p ModalProps, body ...g.Node) g.Node {
	if p.ID == "" {
		p.ID = "modal"
	}
	if p.OpenExpr == "" {
		p.OpenExpr = "$modalOpen"
	}
	scrimAttrs := []g.Node{
		h.ID(p.ID),
		g.Attr("data-component", "modal-scrim"),
		g.Attr("data-show", p.OpenExpr),
		g.Attr("role", "dialog"),
		g.Attr("aria-modal", "true"),
		// Initial state is hidden so the modal does not flash before Datastar
		// resolves the data-show expression. Datastar removes the inline
		// display:none when OpenExpr becomes truthy.
		h.Style("display:none"),
		// click on scrim closes; clicks inside modal stop propagation
		g.Attr("data-on:click", p.OpenExpr+" = false"),
	}
	modalAttrs := []g.Node{
		g.Attr("data-component", "modal"),
		g.If(p.Size != "", g.Attr("data-size", string(p.Size))),
		g.Attr("data-on:click__stop", "void 0"),
	}
	modalAttrs = append(modalAttrs, p.Attrs...)

	header := g.If(p.Title != "",
		h.Div(
			g.Attr("data-slot", "header"),
			h.Span(g.Text(p.Title)),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-variant", "ghost"),
				g.Attr("data-size", "icon"),
				g.Attr("aria-label", "Close"),
				g.Attr("data-on:click", p.OpenExpr+" = false"),
				g.Text("×"),
			),
		),
	)

	footer := g.If(p.Footer != nil,
		h.Div(g.Attr("data-slot", "footer"), p.Footer),
	)

	modalAttrs = append(modalAttrs, header)
	modalAttrs = append(modalAttrs, g.Group(body))
	modalAttrs = append(modalAttrs, footer)

	scrimAttrs = append(scrimAttrs, h.Div(modalAttrs...))
	return h.Div(scrimAttrs...)
}
