package primitive

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type HeadingProps struct {
	Level int // 1..5; default 2
	Attrs []g.Node
}

func Heading(p HeadingProps, children ...g.Node) g.Node {
	if p.Level < 1 || p.Level > 5 {
		p.Level = 2
	}
	tag := map[int]func(...g.Node) g.Node{
		1: h.H1, 2: h.H2, 3: h.H3, 4: h.H4, 5: h.H5,
	}[p.Level]
	return tag(
		g.Attr("data-component", "heading"),
		g.Attr("data-level", strconv.Itoa(p.Level)),
		g.Group(p.Attrs),
		g.Group(children),
	)
}

type DisplayProps struct {
	Attrs []g.Node
}

func Display(p DisplayProps, children ...g.Node) g.Node {
	return h.H1(
		g.Attr("data-component", "display"),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
