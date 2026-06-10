package form

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type InputProps struct {
	Type        string // "text" (default) | "email" | "password" | "tel" | "hidden"
	Placeholder string
	Signal      string // Datastar signal name for data-bind
	Name        string // HTML name attribute (for native form fallback)
	Required    bool
	Attrs       []g.Node
}

func Input(p InputProps, attrs ...g.Node) g.Node {
	if p.Type == "" {
		p.Type = "text"
	}
	nodes := []g.Node{
		g.Attr("data-component", "input"),
		h.Type(p.Type),
		g.If(p.Placeholder != "", h.Placeholder(p.Placeholder)),
		g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
		g.If(p.Name != "", h.Name(p.Name)),
		g.If(p.Required, g.Attr("required")),
		g.Group(p.Attrs),
		g.Group(attrs),
	}
	return h.Input(nodes...)
}

type TextareaProps struct {
	Placeholder string
	Signal      string
	Name        string
	Rows        int
	Required    bool
	Attrs       []g.Node
}

func Textarea(p TextareaProps, attrs ...g.Node) g.Node {
	if p.Rows <= 0 {
		p.Rows = 5
	}
	return h.Textarea(
		g.Attr("data-component", "textarea"),
		h.Placeholder(p.Placeholder),
		g.If(p.Signal != "", g.Attr("data-bind:"+p.Signal)),
		g.If(p.Name != "", h.Name(p.Name)),
		g.If(p.Required, g.Attr("required")),
		g.Attr("rows", strconv.Itoa(p.Rows)),
		g.Group(p.Attrs),
		g.Group(attrs),
	)
}
