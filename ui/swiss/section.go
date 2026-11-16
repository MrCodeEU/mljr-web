// Package swiss holds the layout components of the swiss theme (International
// Typographic Style): a 12-column sheet of section rows with an index anchor,
// ruled cells, the sticky colophon with its section ruler, the register
// dialog, the split-flap skills board and the Fahrplan timeline.
// Styling lives in ui/css/_swiss.css; components only emit data-* hooks.
package swiss

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Page wraps the whole sheet. Background layers (e.g. special.LogoGrid) go in
// bg; they fill the page behind the sheet. sheetID names the content box the
// logo grid snaps to.
func Page(sheetID string, bg g.Node, sheet ...g.Node) g.Node {
	return h.Div(
		g.Attr("data-component", "swiss-page"),
		bg,
		h.Div(append([]g.Node{g.Attr("data-slot", "sheet"), h.ID(sheetID)}, sheet...)...),
	)
}

type SectionProps struct {
	ID    string // section id, used by the ruler and register
	Num   string // index numeral, e.g. "03"
	Title string // German title, rendered uppercase
	Sub   string // English sublabel
	Note  string // short note at the foot of the anchor column
	// Wide turns the anchor into a header strip and gives the body all 12 columns.
	Wide bool
	// NoScatter keeps logo-grid bars out of this section.
	NoScatter bool
	// Field fills the anchor column with a solid colour field:
	// "block" (Klein blue), "second" (orange/blue), "third" (magenta), "ink".
	Field string
	// AnchorExtra is placed between head and note (e.g. the hero logo slot).
	AnchorExtra g.Node
	Attrs       []g.Node
}

// Section renders one row of the sheet: anchor column + 9-column body.
// Children are body items; use Cell for ruled boxes.
func Section(p SectionProps, body ...g.Node) g.Node {
	attrs := []g.Node{
		g.Attr("data-component", "swiss-section"),
		h.ID(p.ID),
		g.Attr("data-num", p.Num),
		g.Attr("data-de", p.Title),
		g.Attr("data-en", p.Sub),
		g.Attr("aria-labelledby", p.ID+"-title"),
		g.If(p.Wide, g.Attr("data-wide", "")),
		g.If(p.NoScatter, g.Attr("data-noscatter", "")),
		g.If(p.Field != "", g.Attr("data-field", p.Field)),
	}
	attrs = append(attrs, p.Attrs...)
	titleBlock := []g.Node{
		h.H2(h.ID(p.ID+"-title"), g.Attr("data-slot", "title"), g.Text(p.Title),
			g.If(p.Sub != "", h.Span(g.Attr("data-slot", "sub"), g.Attr("lang", "en"), g.Text(p.Sub))),
		),
	}
	head := h.Div(g.Attr("data-slot", "head"),
		h.Span(g.Attr("data-slot", "num"), g.Attr("aria-hidden", "true"), g.Text(p.Num)),
		g.If(p.Wide, h.Div(titleBlock...)),
		g.If(!p.Wide, g.Group(titleBlock)),
	)
	return h.Section(append(attrs,
		h.Div(g.Attr("data-slot", "anchor"),
			head,
			p.AnchorExtra,
			g.If(p.Note != "", h.P(g.Attr("data-slot", "note"), g.Text(p.Note))),
		),
		h.Div(append([]g.Node{g.Attr("data-slot", "body")}, body...)...),
	)...)
}

type CellProps struct {
	// Span is the column span inside the body (1–9, or 1–12 in a wide row).
	// 0 = full width.
	Span int
	// Variant: "" (paper), "block" (Klein blue), "second", "third", "ink",
	// "elevated", "flush" (no padding).
	Variant string
	Attrs   []g.Node
}

// Cell is a ruled box inside a section body.
func Cell(p CellProps, children ...g.Node) g.Node {
	attrs := []g.Node{g.Attr("data-component", "swiss-cell")}
	if p.Span > 0 {
		attrs = append(attrs, g.Attr("data-span", strconv.Itoa(p.Span)), h.Style("--span:"+strconv.Itoa(p.Span)))
	}
	if p.Variant != "" {
		attrs = append(attrs, g.Attr("data-variant", p.Variant))
	}
	attrs = append(attrs, p.Attrs...)
	return h.Div(append(attrs, children...)...)
}

// Kicker is the mono label line that opens a cell: label left, aside right.
func Kicker(label string, aside g.Node) g.Node {
	return h.Div(g.Attr("data-component", "swiss-kicker"),
		h.Span(g.Text(label)),
		g.If(aside != nil, h.Span(g.Attr("data-slot", "aside"), aside)),
	)
}

// Figure renders a big tabular number (pair with primitive.NumberTicker).
func Figure(children ...g.Node) g.Node {
	return h.Div(append([]g.Node{g.Attr("data-component", "swiss-figure")}, children...)...)
}

// Micro renders a small uppercase mono label.
func Micro(children ...g.Node) g.Node {
	return h.Span(append([]g.Node{g.Attr("data-component", "swiss-micro")}, children...)...)
}

// Pill is a mono spec tag with a hairline border.
func Pill(children ...g.Node) g.Node {
	return h.Span(append([]g.Node{g.Attr("data-component", "swiss-pill")}, children...)...)
}

// Pills wraps pills in a wrapping row.
func Pills(children ...g.Node) g.Node {
	return h.Div(append([]g.Node{g.Attr("data-component", "swiss-pills")}, children...)...)
}

// Btn is the square mono button.
func Btn(primary bool, attrs []g.Node, children ...g.Node) g.Node {
	a := btnAttrs(primary)
	a = append(a, h.Type("button"))
	return h.Button(append(append(a, attrs...), children...)...)
}

// BtnLink is Btn rendered as a link.
func BtnLink(href string, primary bool, attrs []g.Node, children ...g.Node) g.Node {
	a := append(btnAttrs(primary), h.Href(href))
	return h.A(append(append(a, attrs...), children...)...)
}

func btnAttrs(primary bool) []g.Node {
	a := []g.Node{g.Attr("data-component", "swiss-btn")}
	if primary {
		a = append(a, g.Attr("data-variant", "primary"))
	}
	return a
}

// SROnly hides content visually but keeps it for screen readers.
func SROnly(children ...g.Node) g.Node {
	return h.Span(append([]g.Node{g.Attr("data-component", "swiss-sr")}, children...)...)
}
