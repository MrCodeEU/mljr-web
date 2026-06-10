package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FooterProps struct {
	Attrs []g.Node
}

func Footer(p FooterProps, children ...g.Node) g.Node {
	return h.Footer(
		g.Attr("data-component", "footer"),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
