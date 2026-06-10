package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TagProps struct {
	Tone  token.Tone
	Icon  string // optional icon name e.g. "simple-icons:go" — rendered before text
	Attrs []g.Node
}

func Tag(p TagProps, children ...g.Node) g.Node {
	content := make([]g.Node, 0, len(children)+1)
	if p.Icon != "" && icon.Has(p.Icon) {
		content = append(content, icon.Icon(p.Icon, icon.Props{Size: "1em"}))
	}
	content = append(content, g.Group(children))

	return h.Span(
		g.Attr("data-component", "tag"),
		g.If(p.Tone != token.ToneNone, g.Attr("data-tone", string(p.Tone))),
		g.Group(p.Attrs),
		g.Group(content),
	)
}
