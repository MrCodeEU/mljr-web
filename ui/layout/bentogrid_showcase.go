//go:build showcase

package layout

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "bento-grid", Name: "Bento Grid", Category: "layout",
		Summary: "Mosaic CSS grid layout. BentoItem controls column and row spans. Great for dashboards, feature showcases, landing pages.",
		Code: `// import "mljr-web/ui/layout"
layout.BentoGrid(layout.BentoGridProps{Cols: 3},
    layout.BentoItem(layout.BentoItemProps{ColSpan: 2, RowSpan: 1}, content1),
    layout.BentoItem(layout.BentoItemProps{ColSpan: 1, RowSpan: 2}, content2),
    layout.BentoItem(layout.BentoItemProps{ColSpan: 1, RowSpan: 1}, content3),
    layout.BentoItem(layout.BentoItemProps{ColSpan: 1, RowSpan: 1}, content4),
)`,
		Render: func(p map[string]string) g.Node {
			cell := func(label, iconName, bg string) g.Node {
				return h.Div(
					h.Style("height:100%;min-height:120px;padding:var(--sp-5);background:"+bg+";border-radius:var(--radius);border:var(--bw-2) solid var(--ink);display:flex;flex-direction:column;justify-content:space-between"),
					icon.Icon(iconName, icon.Props{Size: "1.5rem"}),
					h.Span(h.Style("font-size:var(--t-sm);font-weight:700"), g.Text(label)),
				)
			}

			return BentoGrid(BentoGridProps{Cols: 3, Gap: "var(--sp-3)"},
				BentoItem(BentoItemProps{ColSpan: 2, RowSpan: 1}, cell("Fast Server-Side Rendering", "lucide:zap", "var(--accent)")),
				BentoItem(BentoItemProps{ColSpan: 1, RowSpan: 2}, cell("137 Components", "lucide:layout-grid", "var(--surface-2)")),
				BentoItem(BentoItemProps{ColSpan: 1, RowSpan: 1}, cell("4 Themes", "lucide:palette", "var(--surface-2)")),
				BentoItem(BentoItemProps{ColSpan: 1, RowSpan: 1}, cell("Datastar Reactive", "lucide:activity", "var(--surface-2)")),
				BentoItem(BentoItemProps{ColSpan: 2, RowSpan: 1}, cell("Motion v10 Animations", "lucide:sparkles", "var(--surface-2)")),
				BentoItem(BentoItemProps{ColSpan: 1, RowSpan: 1}, cell("Tailwind v4", "lucide:wind", "var(--surface-2)")),
				BentoItem(BentoItemProps{ColSpan: 2, RowSpan: 1}, cell("Zero Runtime Dependencies", "lucide:package", "var(--surface-2)")),
			)
		},
	})
}
