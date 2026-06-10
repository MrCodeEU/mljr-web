package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ScrollAreaProps struct {
	Height    string // CSS max-height (default "320px")
	MaxWidth  string // optional CSS max-width
	Direction string // "vertical" (default) | "horizontal" | "both"
}

// ScrollArea wraps content in a styled scroll container with thin themed scrollbars.
// Uses CSS scrollbar-width + scrollbar-color for consistent cross-browser styling.
func ScrollArea(p ScrollAreaProps, children ...g.Node) g.Node {
	if p.Height == "" {
		p.Height = "320px"
	}
	if p.Direction == "" {
		p.Direction = "vertical"
	}

	overflow := "overflow-y:auto;overflow-x:hidden"
	switch p.Direction {
	case "horizontal":
		overflow = "overflow-x:auto;overflow-y:hidden"
	case "both":
		overflow = "overflow:auto"
	}

	style := overflow + ";max-height:" + p.Height
	if p.MaxWidth != "" {
		style += ";max-width:" + p.MaxWidth
	}

	return h.Div(
		g.Attr("data-component", "scroll-area"),
		h.Style(style),
		g.Group(children),
	)
}
