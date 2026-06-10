package primitive

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CodeBlockProps struct {
	Language string // display label (e.g. "Go", "bash")
	ID       string // unique id for copy button (auto-generated if blank)
	Attrs    []g.Node
}

// CodeBlock renders a styled pre/code block with an optional language label and copy button.
func CodeBlock(p CodeBlockProps, code string) g.Node {
	id := p.ID
	if id == "" {
		id = "cb-" + p.Language
	}

	return h.Div(
		g.Attr("data-component", "code-block"),
		g.Group(p.Attrs),
		h.Div(
			g.Attr("data-slot", "header"),
			h.Span(g.Text(p.Language)),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-variant", "ghost"),
				g.Attr("data-size", "sm"),
				g.Attr("data-on:click", "navigator.clipboard.writeText(document.getElementById('"+id+"').innerText)"),
				g.Text("Copy"),
			),
		),
		h.Pre(
			h.ID(id),
			h.Code(g.Text(code)),
		),
	)
}
