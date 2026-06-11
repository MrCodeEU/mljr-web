package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MasonryProps struct {
	// Cols is the number of columns (default 3).
	Cols int
	// Gap between items (default "var(--sp-4)").
	Gap string
	// MinColWidth for auto-fill (if set, overrides Cols with auto-fill).
	MinColWidth string
}

// Masonry renders a CSS-columns masonry layout. Zero JS.
// Items flow top-to-bottom within each column (CSS column-fill:auto).
// Use MasonryItem to prevent items from breaking across columns.
func Masonry(p MasonryProps, items ...g.Node) g.Node {
	if p.Cols == 0 {
		p.Cols = 3
	}
	if p.Gap == "" {
		p.Gap = "var(--sp-4)"
	}

	var colStyle string
	if p.MinColWidth != "" {
		colStyle = fmt.Sprintf("columns:%s;column-gap:%s", p.MinColWidth, p.Gap)
	} else {
		colStyle = fmt.Sprintf("columns:%d;column-gap:%s", p.Cols, p.Gap)
	}

	return h.Div(
		g.Attr("data-component", "masonry"),
		h.Style(colStyle),
		g.Group(items),
	)
}

// MasonryItem prevents an item from breaking across columns.
func MasonryItem(children ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "masonry-item"),
		g.Group(children),
	)
}
