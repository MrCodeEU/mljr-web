package form

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FileInputProps struct {
	Signal   string // signal for the filename label (updates on change via data-on:change)
	Name     string
	Accept   string // e.g. "image/*,.pdf"
	Multiple bool
	Label    string // default "Choose file…"
	Attrs    []g.Node
}

// FileInput renders a styled file upload button with a visible filename label.
func FileInput(p FileInputProps, attrs ...g.Node) g.Node {
	label := p.Label
	if label == "" {
		label = "Choose file…"
	}

	var filenameNode g.Node
	if p.Signal != "" {
		filenameNode = h.Span(
			g.Attr("data-slot", "filename"),
			g.Attr("data-show", "$"+p.Signal+"!==''"),
			g.Attr("data-text", "$"+p.Signal),
			h.Style("display:none"),
		)
	}

	changeExpr := ""
	if p.Signal != "" {
		changeExpr = "$" + p.Signal + " = evt.target.files.length ? evt.target.files[0].name : ''"
	}

	return h.Label(
		g.Attr("data-component", "file-input-wrap"),
		h.Input(
			h.Type("file"),
			g.If(p.Name != "", h.Name(p.Name)),
			g.If(p.Accept != "", h.Accept(p.Accept)),
			g.If(p.Multiple, g.Attr("multiple")),
			g.If(changeExpr != "", g.Attr("data-on:change", changeExpr)),
			g.Group(p.Attrs),
			g.Group(attrs),
		),
		h.Div(
			g.Attr("data-component", "file-input"),
			h.Div(g.Attr("data-slot", "icon"), icon.Icon("lucide:upload")),
			h.Span(g.Attr("data-slot", "label"), g.Text(label)),
			filenameNode,
		),
	)
}
