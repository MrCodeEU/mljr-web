package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type IconButtonProps struct {
	Icon    string       // lucide or simple-icons name
	Label   string       // aria-label (required for accessibility)
	Variant token.Variant // default Outline
	Size    token.Size   // default SizeIcon
	Href    string       // renders as <a> when set
	Attrs   []g.Node     // pass-through: data-on:click, id, etc.
}

// IconButton renders a square button containing only an icon.
// Thin convenience wrapper over Button with Size=icon.
func IconButton(p IconButtonProps) g.Node {
	if p.Variant == "" {
		p.Variant = token.Outline
	}
	if p.Size == "" {
		p.Size = token.SizeIcon
	}
	if p.Label == "" {
		p.Label = p.Icon
	}

	attrs := append([]g.Node{g.Attr("aria-label", p.Label)}, p.Attrs...)

	if p.Href != "" {
		return h.A(
			h.Href(p.Href),
			g.Attr("data-component", "button"),
			g.Attr("data-variant", string(p.Variant)),
			g.Attr("data-size", string(p.Size)),
			g.Attr("aria-label", p.Label),
			g.Group(p.Attrs),
			icon.Icon(p.Icon),
		)
	}

	return Button(ButtonProps{Variant: p.Variant, Size: p.Size, Attrs: attrs},
		icon.Icon(p.Icon),
	)
}
