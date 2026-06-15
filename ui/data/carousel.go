package data

import (
	"fmt"

	"mljr-web/ui"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CarouselProps struct {
	ID           string
	Images       []string
	Alt          string
	AutoInterval string // Datastar interval modifier, e.g. "3s". Default "4s".
	Attrs        []g.Node
}

func Carousel(p CarouselProps) g.Node {
	id := p.ID
	if id == "" {
		id = "carousel"
	}
	interval := p.AutoInterval
	if interval == "" {
		interval = "4s"
	}
	sig := id + "Idx"
	n := len(p.Images)
	if n == 0 {
		return g.Raw("")
	}

	nextExpr := fmt.Sprintf("$%s = ($%s+1)%%%d", sig, sig, n)
	prevExpr := fmt.Sprintf("$%s = ($%s-1+%d)%%%d", sig, sig, n, n)

	trackChildren := make([]g.Node, 0, n)
	for i, src := range p.Images {
		trackChildren = append(trackChildren,
			h.Div(
				ui.Show(fmt.Sprintf("$%s === %d", sig, i)),
				h.Img(
					h.Src(src),
					h.Alt(fmt.Sprintf("%s %d", p.Alt, i+1)),
					g.Attr("loading", "lazy"),
				),
			),
		)
	}

	dotChildren := make([]g.Node, 0, n)
	for i := range p.Images {
		dotChildren = append(dotChildren,
			h.Button(
				h.Type("button"),
				g.Attr("data-slot", "dots-btn"),
				g.Attr("aria-label", fmt.Sprintf("Go to slide %d", i+1)),
				// active state via data-attr
				g.Attr("data-attr", fmt.Sprintf(`{"data-state": $%s === %d ? "active" : ""}`, sig, i)),
				ui.On("click", fmt.Sprintf("$%s = %d", sig, i)),
			),
		)
	}

	attrs := []g.Node{
		g.Attr("data-component", "carousel"),
		ui.Signals(fmt.Sprintf(`{%s:0}`, sig)),
		// auto-advance via Datastar interval
		ui.On(fmt.Sprintf("interval__%s", interval), nextExpr),
		g.Group(p.Attrs),
	}
	attrs = append(attrs,
		h.Div(g.Attr("data-slot", "track"), g.Group(trackChildren)),
		h.Button(g.Attr("data-slot", "prev"), ui.On("click", prevExpr), g.Text("‹")),
		h.Button(g.Attr("data-slot", "next"), ui.On("click", nextExpr), g.Text("›")),
		h.Div(g.Attr("data-slot", "dots"), g.Group(dotChildren)),
	)
	return h.Div(attrs...)
}
