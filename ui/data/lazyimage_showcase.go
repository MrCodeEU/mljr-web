//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "lazy-image", Name: "Lazy Image", Category: "data",
		Summary: "IntersectionObserver lazy-loads image on scroll. Skeleton placeholder fades to image on load.",
		Code: `data.LazyImage(data.LazyImageProps{
    Src:    "/static/img/hero.jpg",
    Alt:    "Hero image",
    Width:  "100%",
    Height: "240px",
    Rounded: true,
})`,
		Render: func(p map[string]string) g.Node {
			// Using picsum photos — CSP blocked in some environments; use relative path in prod
			images := []struct{ src, alt string }{
				{"https://picsum.photos/seed/a/400/240", "Nature landscape"},
				{"https://picsum.photos/seed/b/400/240", "City skyline"},
				{"https://picsum.photos/seed/c/400/240", "Abstract art"},
				{"https://picsum.photos/seed/d/400/240", "Architecture"},
			}
			cells := make([]g.Node, len(images))
			for i, img := range images {
				cells[i] = h.Div(
					h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"),
					LazyImage(LazyImageProps{
						Src:     img.src,
						Alt:     img.alt,
						Width:   "100%",
						Height:  "160px",
						Rounded: true,
					}),
					h.Span(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text(img.alt)),
				)
			}
			return h.Div(
				h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-4)"),
				g.Group(cells),
			)
		},
	})
}
