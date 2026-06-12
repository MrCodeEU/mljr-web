//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "resizable-panels", Name: "Resizable Panels", Category: "layout",
		PreviewHeight: "400px",
		Summary:       "Two panels with a drag handle. Pointer events, touch support, configurable split and min-size.",
		Code: `layout.ResizablePanels(
    layout.ResizablePanelsProps{
        Direction:    "horizontal",
        InitialSplit: 40,
        Min:          20,
    },
    firstPanelContent,
    secondPanelContent,
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("height:360px;border:var(--bw-2) solid var(--line);border-radius:var(--radius);overflow:hidden"),
				ResizablePanels(
					ResizablePanelsProps{Direction: "horizontal", InitialSplit: 35, Min: 15},
					h.Div(
						h.Style("padding:var(--sp-4);height:100%;background:var(--surface)"),
						h.P(h.Style("font-weight:700;margin:0 0 var(--sp-2)"), g.Text("File tree")),
						h.Ul(h.Style("list-style:none;padding:0;margin:0;font-size:var(--t-sm);display:flex;flex-direction:column;gap:var(--sp-1)"),
							h.Li(g.Text("📁 src/")),
							h.Li(h.Style("padding-left:var(--sp-3)"), g.Text("📄 main.go")),
							h.Li(h.Style("padding-left:var(--sp-3)"), g.Text("📁 ui/")),
							h.Li(h.Style("padding-left:var(--sp-5)"), g.Text("📄 button.go")),
							h.Li(h.Style("padding-left:var(--sp-5)"), g.Text("📄 card.go")),
							h.Li(g.Text("📄 go.mod")),
						),
					),
					h.Div(
						h.Style("padding:var(--sp-4);height:100%;background:var(--bg)"),
						h.P(h.Style("font-weight:700;margin:0 0 var(--sp-2)"), g.Text("Editor")),
						h.Pre(h.Style("font-family:var(--font-mono);font-size:var(--t-xs);opacity:.8;margin:0"),
							h.Code(g.Text("package main\n\nimport (\n    g \"maragu.dev/gomponents\"\n    h \"maragu.dev/gomponents/html\"\n)\n\nfunc main() {\n    // ...\n}")),
						),
					),
				),
			)
		},
	})
}
