//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "word-rotate", Name: "Word Rotate", Category: "primitive",
		Summary: "Cycles through a list of words with a fade+slide transition. No dependencies — pure CSS transitions + setInterval.",
		Code: `// import "mljr-web/ui/primitive"
h.H2(
    g.Text("Build "),
    primitive.WordRotate(primitive.WordRotateProps{
        Words:    []string{"faster", "smarter", "cleaner", "better"},
        Interval: 2000,
        ID:       "hero-rotate",
    }),
    g.Text(" today."),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-8)"),
				h.Div(
					h.H2(
						h.Style("font-size:var(--t-2xl);font-weight:800;line-height:1.2"),
						g.Text("Build "),
						WordRotate(WordRotateProps{
							Words:    []string{"faster", "smarter", "cleaner", "better"},
							Interval: 2000,
							ID:       "wr-hero",
							Class:    "color:var(--accent)",
						}),
						g.Text(" with mljr-web."),
					),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin-bottom:var(--sp-2)"), g.Text("Fast interval (800ms):")),
					h.Div(
						h.Style("font-size:var(--t-lg);font-weight:700"),
						g.Text("Status: "),
						WordRotate(WordRotateProps{
							Words:    []string{"Online", "Ready", "Active", "Live"},
							Interval: 800,
							ID:       "wr-status",
							Class:    "color:var(--success)",
						}),
					),
				),
			)
		},
	})
}
