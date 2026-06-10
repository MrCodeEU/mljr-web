//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "typewriter", Name: "Typewriter", Category: "primitive",
		Summary: "Animated type-and-delete loop across multiple phrases. No dependencies — pure JS setInterval with configurable speed/pause.",
		Code: `// import "mljr-web/ui/primitive"
primitive.Typewriter(primitive.TypewriterProps{
    Lines: []string{
        "go build ./...",
        "go test ./...",
        "git push origin main",
    },
    Speed:  50,
    Pause:  2000,
    ID:     "tw-demo",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-8)"),
				// Hero typewriter
				h.Div(
					h.H2(
						h.Style("font-size:var(--t-2xl);font-weight:800"),
						g.Text("We "),
						Typewriter(TypewriterProps{
							Lines:    []string{"build fast.", "ship clean.", "iterate quickly.", "care about DX."},
							Speed:    65,
							Pause:    2000,
							ID:       "tw-hero",
						}),
					),
				),
				// Terminal style
				h.Div(
					h.Style("background:#1a1a1a;border-radius:var(--radius);padding:var(--sp-4) var(--sp-5);font-family:var(--font-mono)"),
					h.Div(h.Style("display:flex;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
						h.Span(h.Style("width:10px;height:10px;border-radius:50%;background:#ff5f57")),
						h.Span(h.Style("width:10px;height:10px;border-radius:50%;background:#febc2e")),
						h.Span(h.Style("width:10px;height:10px;border-radius:50%;background:#28c840")),
					),
					h.Div(
						h.Style("color:#e0e0e0;font-size:var(--t-sm)"),
						h.Span(h.Style("color:#7eff7e"), g.Text("$ ")),
						Typewriter(TypewriterProps{
							Lines:    []string{"go build ./...", "go test -v ./...", "bin/tailwindcss --minify", "git push origin main"},
							Speed:    55,
							DeleteSpeed: 25,
							Pause:    1500,
							ID:       "tw-term",
						}),
					),
				),
			)
		},
	})
}
