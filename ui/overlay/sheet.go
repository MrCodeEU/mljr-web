package overlay

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SheetProps struct {
	SignalName string // default "_sheetOpen"
	Title      string
	// Placement: "bottom" (default) | "right" | "left" | "top"
	Placement string
}

// Sheet renders a full-viewport-edge slide-in panel.
// Unlike Drawer, Sheet always slides from bottom (mobile-style) or right (desktop).
// Open with: data-on:click="$_sheetOpen=true"
func Sheet(p SheetProps, children ...g.Node) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_sheetOpen"
	}
	if p.Placement == "" {
		p.Placement = "bottom"
	}
	sig := p.SignalName

	return h.Div(
		g.Attr("data-component", "sheet"),
		g.Attr("data-placement", p.Placement),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		g.Attr("data-show", "$"+sig),
		h.Style("display:none"),

		h.Div(g.Attr("data-slot", "backdrop"), g.Attr("data-on:click", "$"+sig+"=false")),

		h.Div(
			g.Attr("data-slot", "panel"),
			g.Attr("role", "dialog"),
			g.Attr("aria-modal", "true"),
			h.Div(
				g.Attr("data-slot", "header"),
				g.If(p.Title != "", h.Strong(g.Text(p.Title))),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeIcon},
					g.Attr("data-on:click", "$"+sig+"=false"),
					g.Attr("aria-label", "Close"),
					icon.Icon("lucide:x"),
				),
			),
			h.Div(g.Attr("data-slot", "content"), g.Group(children)),
		),
	)
}
