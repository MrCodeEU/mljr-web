package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TooltipPlacement string

const (
	TooltipTop    TooltipPlacement = "top"
	TooltipBottom TooltipPlacement = "bottom"
	TooltipLeft   TooltipPlacement = "left"
	TooltipRight  TooltipPlacement = "right"
)

type TooltipProps struct {
	Text      string
	Placement TooltipPlacement // default: top
	Attrs     []g.Node
}

// Tooltip wraps a trigger element with a CSS-only hover tooltip. No JS required.
func Tooltip(p TooltipProps, trigger g.Node) g.Node {
	placement := string(p.Placement)
	if placement == "" {
		placement = "top"
	}
	return h.Div(
		g.Attr("data-component", "tooltip"),
		g.Attr("data-placement", placement),
		g.Group(p.Attrs),
		trigger,
		h.Span(g.Attr("data-slot", "tip"), g.Text(p.Text)),
	)
}
