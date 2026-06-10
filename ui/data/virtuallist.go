package data

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type VirtualListProps struct {
	// EstimatedItemHeight is the CSS contain-intrinsic-size hint (default "56px").
	// Must approximate your row height for scroll-position accuracy.
	EstimatedItemHeight string
	// Height is the container CSS height (default "400px").
	Height string
	// ID is unique per page (default "vl").
	ID string
}

// VirtualList renders a scrollable container where items use CSS content-visibility:auto.
// The browser skips layout and paint for off-screen items — zero JS, near-native performance.
// Render all items server-side; the browser handles viewport culling.
func VirtualList(p VirtualListProps, items ...g.Node) g.Node {
	if p.Height == "" {
		p.Height = "400px"
	}
	if p.EstimatedItemHeight == "" {
		p.EstimatedItemHeight = "56px"
	}
	if p.ID == "" {
		p.ID = "vl"
	}

	// Wrap each item in a culled div
	wrapped := make([]g.Node, len(items))
	for i, item := range items {
		wrapped[i] = h.Div(
			g.Attr("data-slot", "item"),
			h.Style(fmt.Sprintf("content-visibility:auto;contain-intrinsic-size:0 %s", p.EstimatedItemHeight)),
			item,
		)
	}

	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "virtual-list"),
		h.Style(fmt.Sprintf("height:%s;overflow-y:auto;contain:strict", p.Height)),
		g.Group(wrapped),
	)
}

// VirtualListItem is a convenience builder for a standard list row.
type VirtualListItemProps struct {
	Title    string
	Subtitle string
	Meta     string // right-aligned metadata
}

func VirtualListRow(p VirtualListItemProps) g.Node {
	return h.Div(
		g.Attr("data-slot", "row"),
		h.Style("display:flex;align-items:center;justify-content:space-between;padding:var(--sp-3) var(--sp-4);border-bottom:var(--bw-1) solid var(--line);gap:var(--sp-3)"),
		h.Div(
			h.Strong(h.Style("display:block;font-size:var(--t-sm);font-weight:600"), g.Text(p.Title)),
			g.If(p.Subtitle != "", h.Span(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text(p.Subtitle))),
		),
		g.If(p.Meta != "", h.Span(h.Style("font-size:var(--t-xs);color:var(--muted);white-space:nowrap;flex-shrink:0"), g.Text(p.Meta))),
	)
}
