package overlay

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PopoverPlacement string

const (
	PopoverBottom PopoverPlacement = "bottom"
	PopoverTop    PopoverPlacement = "top"
	PopoverLeft   PopoverPlacement = "left"
	PopoverRight  PopoverPlacement = "right"
)

type PopoverProps struct {
	Signal    string           // Datastar signal name (default "pop")
	Placement PopoverPlacement // default: bottom
	Attrs     []g.Node
}

// Popover wraps a trigger and floating content panel gated by a Datastar signal.
// Click outside closes via window listener.
func Popover(p PopoverProps, trigger g.Node, content g.Node) g.Node {
	sig := p.Signal
	if sig == "" {
		sig = "pop"
	}
	placement := string(p.Placement)
	if placement == "" {
		placement = "bottom"
	}
	openExpr := "$" + sig
	toggleExpr := fmt.Sprintf("$%s=!$%s", sig, sig)
	closeExpr := fmt.Sprintf("$%s=false", sig)

	return h.Div(
		g.Attr("data-component", "popover"),
		g.Attr("data-signals", fmt.Sprintf(`{%s:false}`, sig)),
		g.Attr("data-on:click__window", closeExpr),
		g.Group(p.Attrs),
		h.Div(
			g.Attr("data-on:click__stop", toggleExpr),
			trigger,
		),
		h.Div(
			g.Attr("data-component", "popover-content"),
			g.Attr("data-placement", placement),
			g.Attr("data-show", openExpr),
			g.Attr("data-on:click__stop", "void 0"),
			h.Style("display:none"),
			content,
		),
	)
}
