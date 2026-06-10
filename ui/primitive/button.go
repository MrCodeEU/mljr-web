package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ButtonProps struct {
	Variant  token.Variant
	Size     token.Size
	Tone     token.Tone
	Type     string // "button" (default) | "submit" | "reset"
	Disabled bool
	Attrs    []g.Node
}

func Button(p ButtonProps, children ...g.Node) g.Node {
	if p.Variant == "" {
		p.Variant = token.Primary
	}
	if p.Size == "" {
		p.Size = token.SizeMD
	}
	if p.Type == "" {
		p.Type = "button"
	}
	return h.Button(
		g.Attr("data-component", "button"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("data-size", string(p.Size)),
		g.If(p.Tone != token.ToneNone, g.Attr("data-tone", string(p.Tone))),
		h.Type(p.Type),
		g.If(p.Disabled, g.Attr("disabled")),
		g.Group(p.Attrs),
		g.Group(children),
	)
}
