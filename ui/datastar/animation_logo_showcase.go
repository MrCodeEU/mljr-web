//go:build showcase

package datastar

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/special"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-logo", Name: "Logo Scatter", Category: "animation",
		PreviewHeight: "680px",
		Summary:       "Pill-shaped logo pieces scatter across the viewport on a loop, then sway gently before reassembling. Built with Motion v10 WAAPI transforms.",
		Code: `// Motion v10 scatter pattern — see ui/special/logo_scatter.go
// special.LogoScatter(special.LogoScatterProps{
//   ID:   "my-logo",
//   Mode: "loop",    // or "scroll" for IntersectionObserver trigger
// })
//
// Key constraint: CSS transform px on SVG children = SVG user units.
// recalc() multiplies targets by svgScale = viewBox.width / svgRect.width.`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("position:relative;min-height:300px;width:100%;overflow:hidden"),
				special.LogoScatter(special.LogoScatterProps{
					ID:       "logo-svg-demo",
					Size:     "280px",
					SVGStyle: "position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);overflow:visible;width:280px;height:280px",
					Mode:     "loop",
				}),
			)
		},
	})
}
