package layout

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GridProps struct {
	Attrs []g.Node
}

func Grid(p GridProps, children ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "grid"),
		g.Group(p.Attrs),
		g.Group(children),
	)
}

type ColProps struct {
	Span  int // 1..12, default 12
	Attrs []g.Node
}

func Col(p ColProps, children ...g.Node) g.Node {
	if p.Span < 1 || p.Span > 12 {
		p.Span = 12
	}
	return h.Div(
		g.Attr("data-component", "col"),
		g.Attr("data-span", strconv.Itoa(p.Span)),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
