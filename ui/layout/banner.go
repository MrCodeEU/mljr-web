package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type BannerVariant string

const (
	BannerDefault BannerVariant = ""
	BannerInfo    BannerVariant = "info"
	BannerSuccess BannerVariant = "success"
	BannerWarning BannerVariant = "warning"
	BannerDanger  BannerVariant = "danger"
)

type BannerProps struct {
	Variant     BannerVariant
	Dismiss     bool   // show × dismiss button
	DismissExpr string // Datastar expression to run on dismiss (e.g. "$bannerOpen=false")
	Attrs       []g.Node
}

// Banner renders a full-width announcement strip, typically placed above the navbar.
func Banner(p BannerProps, children ...g.Node) g.Node {
	var dismissBtn g.Node
	if p.Dismiss {
		expr := p.DismissExpr
		if expr == "" {
			expr = "evt.target.closest('[data-component=banner]').remove()"
		}
		dismissBtn = h.Button(
			g.Attr("data-slot", "dismiss"),
			g.Attr("aria-label", "Dismiss"),
			g.Attr("data-on:click", expr),
			g.Text("×"),
		)
	}
	nodes := []g.Node{
		g.Attr("data-component", "banner"),
		g.Attr("role", "banner"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Group(p.Attrs),
	}
	nodes = append(nodes, children...)
	if dismissBtn != nil {
		nodes = append(nodes, dismissBtn)
	}
	return h.Div(nodes...)
}
