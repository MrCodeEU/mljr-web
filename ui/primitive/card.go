package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CardProps struct {
	Tone        token.Tone
	Interactive bool
	Attrs       []g.Node
}

func Card(p CardProps, children ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "card"),
		g.If(p.Tone != token.ToneNone, g.Attr("data-tone", string(p.Tone))),
		g.If(p.Interactive, g.Attr("data-interactive")),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
