package feedback

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SkeletonVariant string

const (
	SkeletonText   SkeletonVariant = "text"
	SkeletonRect   SkeletonVariant = "rect"
	SkeletonCircle SkeletonVariant = "circle"
)

type SkeletonProps struct {
	Variant SkeletonVariant
	Width   string // CSS width (default "100%")
	Height  string // CSS height (default "")
	Attrs   []g.Node
}

// Skeleton renders a shimmering placeholder for loading content.
func Skeleton(p SkeletonProps, attrs ...g.Node) g.Node {
	if p.Variant == "" {
		p.Variant = SkeletonRect
	}
	style := ""
	if p.Width != "" {
		style += "width:" + p.Width + ";"
	}
	if p.Height != "" {
		style += "height:" + p.Height + ";"
	}
	return h.Div(
		g.Attr("data-component", "skeleton"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("aria-hidden", "true"),
		g.If(style != "", h.Style(style)),
		g.Group(p.Attrs),
		g.Group(attrs),
	)
}
