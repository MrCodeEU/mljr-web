// Package form provides labelled form field wrappers and input primitives.
// All components follow the data-* convention; styling hooks are data attributes.
package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FieldProps struct {
	Label    string
	Hint     string
	ErrorFor string // signal name — shows error text when $<signal>Error is non-empty
	Attrs    []g.Node
}

// Field wraps a control in a labelled group with optional hint and reactive error.
func Field(p FieldProps, control g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "field"),
		g.Group(p.Attrs),
		g.If(p.Label != "", h.Label(g.Attr("data-slot", "label"), g.Text(p.Label))),
		control,
		g.If(p.Hint != "", h.Span(g.Attr("data-slot", "hint"), g.Text(p.Hint))),
		g.If(p.ErrorFor != "",
			h.Span(
				g.Attr("data-slot", "error"),
				g.Attr("data-show", "$"+p.ErrorFor+"Error !== ''"),
				g.Attr("data-text", "$"+p.ErrorFor+"Error"),
			),
		),
	)
}
