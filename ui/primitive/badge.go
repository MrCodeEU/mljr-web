package primitive

import (
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BadgeVariant string

const (
	BadgeDefault BadgeVariant = ""
	BadgeSuccess BadgeVariant = "success"
	BadgeDanger  BadgeVariant = "danger"
	BadgeWarning BadgeVariant = "warning"
	BadgeInfo    BadgeVariant = "info"
	BadgeOutline BadgeVariant = "outline"
	BadgeTone    BadgeVariant = "tone"
)

type BadgeProps struct {
	Variant BadgeVariant
	Tone    token.Tone
	Size    token.Size // sm | md (default) | lg
	Attrs   []g.Node
}

// Badge renders a small pill label for counts, statuses, or labels.
func Badge(p BadgeProps, children ...g.Node) g.Node {
	nodes := []g.Node{
		g.Attr("data-component", "badge"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.If(p.Tone != "", g.Attr("data-tone", string(p.Tone))),
		g.If(p.Size != "", g.Attr("data-size", string(p.Size))),
		g.Group(p.Attrs),
	}
	nodes = append(nodes, children...)
	return h.Span(nodes...)
}
