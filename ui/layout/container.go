package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ContainerProps struct {
	Attrs []g.Node
}

func Container(p ContainerProps, children ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "container"),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
