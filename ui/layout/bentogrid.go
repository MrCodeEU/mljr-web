package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BentoGridProps struct {
	// Cols is the number of grid columns (default 3).
	Cols int
	// Gap is the CSS gap value (default "var(--sp-4)").
	Gap string
}

type BentoItemProps struct {
	// ColSpan: how many columns this item spans (1–4, default 1).
	ColSpan int
	// RowSpan: how many rows this item spans (1–2, default 1).
	RowSpan int
}

// BentoGrid renders a CSS grid mosaic layout.
// Use BentoItem to control how much space each cell occupies.
func BentoGrid(p BentoGridProps, items ...g.Node) g.Node {
	if p.Cols == 0 {
		p.Cols = 3
	}
	if p.Gap == "" {
		p.Gap = "var(--sp-4)"
	}
	style := fmt.Sprintf("display:grid;grid-template-columns:repeat(%d,1fr);gap:%s;", p.Cols, p.Gap)
	return h.Div(
		g.Attr("data-component", "bento-grid"),
		g.Attr("data-cols", fmt.Sprintf("%d", p.Cols)),
		h.Style(style),
		g.Group(items),
	)
}

// BentoItem wraps content with grid-column and grid-row span styles.
func BentoItem(p BentoItemProps, children ...g.Node) g.Node {
	if p.ColSpan < 1 {
		p.ColSpan = 1
	}
	if p.RowSpan < 1 {
		p.RowSpan = 1
	}
	style := fmt.Sprintf("grid-column:span %d;grid-row:span %d", p.ColSpan, p.RowSpan)
	return h.Div(
		g.Attr("data-component", "bento-item"),
		g.Attr("data-col-span", fmt.Sprintf("%d", p.ColSpan)),
		g.Attr("data-row-span", fmt.Sprintf("%d", p.RowSpan)),
		h.Style(style),
		g.Group(children),
	)
}
