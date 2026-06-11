//go:build showcase

package layout

import (
	"fmt"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "masonry", Name: "Masonry Grid", Category: "layout",
		Summary: "CSS-columns masonry layout. Items flow top-to-bottom. Zero JS — pure CSS column-count.",
		Code: `layout.Masonry(layout.MasonryProps{Cols: 3},
    layout.MasonryItem(card1),
    layout.MasonryItem(card2),
)`,
		Render: func(p map[string]string) g.Node {
			heights := []string{"120px", "200px", "160px", "100px", "180px", "140px", "220px", "90px", "150px"}
			colors := []string{"var(--accent)", "var(--surface-2)", "var(--ink)", "var(--surface-2)", "var(--accent)", "var(--surface-2)", "var(--ink)", "var(--accent)", "var(--surface-2)"}
			textColors := []string{"var(--accent-ink)", "var(--fg)", "var(--bg)", "var(--fg)", "var(--accent-ink)", "var(--fg)", "var(--bg)", "var(--accent-ink)", "var(--fg)"}

			items := make([]g.Node, len(heights))
			for i, h2 := range heights {
				idx := i
				items[idx] = MasonryItem(
					h.Div(
						h.Style("height:"+h2+";background:"+colors[idx]+";color:"+textColors[idx]+";border:var(--bw-2) solid var(--ink);border-radius:var(--radius);display:flex;align-items:center;justify-content:center;font-weight:800;font-size:var(--t-sm)"),
						g.Text(fmt.Sprintf("Item %d", idx+1)),
					),
				)
			}
			return Masonry(MasonryProps{Cols: 3}, items...)
		},
	})
}
