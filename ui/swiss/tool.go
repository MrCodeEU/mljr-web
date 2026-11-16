package swiss

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Building blocks for the interactive tools (regex, cron, codec): a page
// head, labelled form fields, result bars and rows, a cheat-sheet list and
// an error box. Styling lives in ui/css/_swiss.css (section "Tools").

// ToolHead is the opening cell of a tool: big name, one-line lede and spec
// pills (e.g. "PCRE", "Go").
func ToolHead(name, lede string, pills ...string) g.Node {
	ps := make([]g.Node, len(pills))
	for i, p := range pills {
		ps[i] = Pill(g.Text(p))
	}
	return h.Div(g.Attr("data-component", "swiss-tool-head"),
		h.H1(g.Attr("data-slot", "title"), g.Text(name)),
		g.If(lede != "", h.P(g.Attr("data-slot", "lede"), g.Text(lede))),
		g.If(len(ps) > 0, Pills(ps...)),
	)
}

// Field is a labelled form control. id must match the control's id, so the
// label names it; children are the control and anything beside it (flags,
// toggles). Inputs, textareas and selects inside get the Swiss form look.
func Field(id, label string, children ...g.Node) g.Node {
	return h.Div(g.Attr("data-component", "swiss-field"),
		h.Label(g.Attr("data-slot", "label"), h.For(id), g.Text(label)),
		h.Div(append([]g.Node{g.Attr("data-slot", "control")}, children...)...),
	)
}

// Affix renders fixed text glued to an input inside a Field, e.g. the
// slashes around a regex pattern.
func Affix(text string) g.Node {
	return h.Span(g.Attr("data-slot", "affix"), g.Attr("aria-hidden", "true"), g.Text(text))
}

// Choice is a group of pressed/unpressed buttons (segmented control).
func Choice(label string, buttons ...g.Node) g.Node {
	return h.Div(append([]g.Node{g.Attr("data-component", "swiss-choice"), h.Role("group"), g.Attr("aria-label", label)}, buttons...)...)
}

// ResultBar is the ink strip that heads a result: what is shown, plus an
// aside (count, mode).
func ResultBar(label string, aside g.Node) g.Node {
	return h.Div(g.Attr("data-component", "swiss-result-bar"),
		h.Span(g.Text(label)),
		g.If(aside != nil, h.Span(g.Attr("data-slot", "aside"), aside)),
	)
}

// ResultBlock is one titled part of a result below the bar.
func ResultBlock(title string, children ...g.Node) g.Node {
	return h.Div(append([]g.Node{g.Attr("data-component", "swiss-result-block"), Micro(g.Text(title))}, children...)...)
}

// Row is one line of a numbered result list: index, main value, aside on
// the right, optional detail lines underneath.
func Row(num string, main, aside g.Node, detail ...g.Node) g.Node {
	return h.Div(g.Attr("data-component", "swiss-row"),
		h.Span(g.Attr("data-slot", "num"), g.Text(num)),
		h.Div(g.Attr("data-slot", "main"), main),
		h.Span(g.Attr("data-slot", "aside"), aside),
		g.If(len(detail) > 0, h.Div(append([]g.Node{g.Attr("data-slot", "detail")}, detail...)...)),
	)
}

// RefList is a cheat-sheet block: syntax on the left, meaning on the right.
func RefList(title string, rows [][2]string) g.Node {
	items := make([]g.Node, 0, 2*len(rows))
	for _, r := range rows {
		items = append(items, h.Dt(h.Code(g.Text(r[0]))), h.Dd(g.Text(r[1])))
	}
	return g.Group{
		Kicker(title, nil),
		h.Dl(append([]g.Node{g.Attr("data-component", "swiss-ref")}, items...)...),
	}
}

// ErrorBox reports invalid input; role=alert announces it.
func ErrorBox(msg string) g.Node {
	return h.Div(g.Attr("data-component", "swiss-error"), h.Role("alert"),
		h.B(g.Text("!")), h.Span(g.Text(msg)),
	)
}

// FooterLink is one link in Footer.
type FooterLink struct {
	Label, Href string
}

// Footer is the slim footer of standalone pages: wordmark, link row and a
// bottom line.
func Footer(note string, links []FooterLink, bottom ...string) g.Node {
	items := make([]g.Node, len(links))
	for i, l := range links {
		items[i] = h.Li(h.A(h.Href(l.Href), g.Text(l.Label)))
	}
	bot := make([]g.Node, len(bottom))
	for i, b := range bottom {
		bot[i] = h.Span(g.Text(b))
	}
	return h.Footer(g.Attr("data-component", "swiss-footer"),
		h.Div(g.Attr("data-slot", "poster"),
			h.Div(g.Attr("data-slot", "word"), g.Attr("aria-hidden", "true"), h.Span(g.Text("MLJR")), h.Span(g.Text(".EU"))),
			h.I(g.Attr("data-slot", "circle"), g.Attr("aria-hidden", "true")),
			h.I(g.Attr("data-slot", "tri"), g.Attr("aria-hidden", "true")),
			h.Div(g.Attr("data-slot", "note"), h.P(g.Text(note)), h.Ul(items...)),
		),
		h.Div(g.Attr("data-slot", "bottom"), g.Group(bot)),
	)
}
