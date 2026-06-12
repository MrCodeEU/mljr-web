package layout

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// FooterLink is one link inside a footer column.
type FooterLink struct {
	Label    string
	Href     string
	External bool // adds target=_blank rel=noopener
}

// FooterColumn is a titled group of links.
type FooterColumn struct {
	Title string
	Links []FooterLink
}

type FooterProps struct {
	// Brand is rendered in the wide leading cell (logo, name, …).
	Brand g.Node
	// Tagline is rendered under the brand.
	Tagline string
	// Columns are titled link groups laid out in a grid next to the brand.
	Columns []FooterColumn
	// Bottom is the right side of the bottom bar (legal links, ©, …).
	Bottom g.Node
	Attrs  []g.Node
}

// Footer renders the site footer. With Columns set it uses the structured
// neo-brutalist layout (brand cell + link columns + bottom bar); without,
// children are laid out in the legacy flex row.
func Footer(p FooterProps, children ...g.Node) g.Node {
	if len(p.Columns) == 0 && p.Brand == nil {
		return h.Footer(
			g.Attr("data-component", "footer"),
			g.Group(p.Attrs),
			g.Group(children),
		)
	}

	cols := make([]g.Node, 0, len(p.Columns))
	for _, col := range p.Columns {
		links := make([]g.Node, 0, len(col.Links))
		for _, lnk := range col.Links {
			attrs := []g.Node{h.Href(lnk.Href)}
			if lnk.External {
				attrs = append(attrs, g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer"))
			}
			links = append(links, h.Li(h.A(append(attrs, g.Text(lnk.Label))...)))
		}
		cols = append(cols, h.Div(
			g.Attr("data-slot", "col"),
			h.Div(g.Attr("data-slot", "col-title"), g.Text(col.Title)),
			h.Ul(g.Attr("data-slot", "col-links"), g.Group(links)),
		))
	}

	return h.Footer(
		g.Attr("data-component", "footer"),
		g.Attr("data-variant", "structured"),
		g.Group(p.Attrs),
		h.Div(
			g.Attr("data-slot", "grid"),
			h.Div(
				g.Attr("data-slot", "brand"),
				p.Brand,
				g.If(p.Tagline != "", h.P(g.Attr("data-slot", "tagline"), g.Text(p.Tagline))),
			),
			g.Group(cols),
		),
		g.Group(children),
		g.If(p.Bottom != nil,
			h.Div(g.Attr("data-slot", "bottom"), p.Bottom),
		),
	)
}
