//go:build showcase

package data

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "carousel", Name: "Carousel", Category: "data",
		Summary: "Pure Datastar image carousel with prev/next and dot navigation.",
		Code: `// wrap in a Card for proper full-bleed negative-margin bleed
primitive.Card(primitive.CardProps{},
    data.Carousel(data.CarouselProps{
        ID:     "hero",
        Images: []string{"/img/a.jpg", "/img/b.jpg"},
        Alt:    "Slide",
    }),
    h.P(g.Text("Caption")),
)`,
		Render: func(p map[string]string) g.Node {
			// SVG data URIs — no external requests, CSP-safe
			svgSlide := func(fill, label string) string {
				return "data:image/svg+xml," +
					"%3Csvg xmlns='http://www.w3.org/2000/svg' width='800' height='300'%3E" +
					"%3Crect width='800' height='300' fill='" + fill + "'/%3E" +
					"%3Ctext x='400' y='150' text-anchor='middle' dominant-baseline='middle' " +
					"font-family='sans-serif' font-size='48' font-weight='bold' fill='white'%3E" + label + "%3C/text%3E" +
					"%3C/svg%3E"
			}
			return primitive.Card(primitive.CardProps{},
				Carousel(CarouselProps{
					ID: "demo",
					Images: []string{
						svgSlide("%23e2483d", "Slide+1"),
						svgSlide("%232f5cff", "Slide+2"),
						svgSlide("%231fab57", "Slide+3"),
					},
					Alt: "Demo slide",
				}),
				h.P(g.Text("Prev / next arrows · dot navigation · Datastar signals")),
			)
		},
	})
}
