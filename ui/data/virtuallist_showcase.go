//go:build showcase

package data

import (
	"fmt"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "virtual-list", Name: "Virtual List", Category: "data",
		Summary: "Scrollable list using CSS content-visibility:auto for zero-JS viewport culling. Browser skips layout/paint for off-screen rows.",
		Code: `// import "mljr-web/ui/data"
items := make([]g.Node, 500)
for i := range items {
    items[i] = data.VirtualListRow(data.VirtualListItemProps{
        Title:    fmt.Sprintf("Row %d", i+1),
        Subtitle: "Subtitle text here",
        Meta:     "12:34",
    })
}
data.VirtualList(data.VirtualListProps{
    Height: "400px",
    EstimatedItemHeight: "56px",
}, items...)`,
		Render: func(p map[string]string) g.Node {
			rows := make([]g.Node, 200)
			for i := range rows {
				rows[i] = VirtualListRow(VirtualListItemProps{
					Title:    fmt.Sprintf("Item #%d — Server-rendered, zero JS culling", i+1),
					Subtitle: fmt.Sprintf("Subtitle for row %d", i+1),
					Meta:     fmt.Sprintf("%d KB", (i+1)*3),
				})
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.P(h.Style("font-size:var(--t-sm);color:var(--muted)"), g.Text("200 rows rendered server-side. Browser culls off-screen rows via CSS — no JS windowing library needed.")),
				VirtualList(VirtualListProps{Height: "360px"}, rows...),
			)
		},
	})
}
