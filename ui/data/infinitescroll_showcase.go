//go:build showcase

package data

import (
	"fmt"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "infinite-scroll", Name: "Infinite Scroll", Category: "data",
		Summary: "Sentinel div triggers Datastar @get when it enters the viewport. Server appends HTML. One Datastar attribute — no custom JS.",
		Code: `data.InfiniteScroll(data.InfiniteScrollProps{
    FetchURL:    "/api/items",
    PageSignal:  "_isPage",
    ContainerID: "items",
    LoadingText: "Loading more…",
}, initialItems...)`,
		Render: func(p map[string]string) g.Node {
			initial := make([]g.Node, 5)
			for i := range initial {
				initial[i] = primitive.Card(primitive.CardProps{},
					h.P(h.Style("margin:0;font-weight:700"), g.Text(fmt.Sprintf("Item %d", i+1))),
					h.P(h.Style("margin:var(--sp-1) 0 0;color:var(--muted);font-size:var(--t-sm)"), g.Text("Pre-loaded item from server.")),
				)
			}
			return h.Div(
				h.Style("max-height:400px;overflow-y:auto;padding:var(--sp-1)"),
				InfiniteScroll(InfiniteScrollProps{
					FetchURL:    "/api/showcase/infinite-items",
					PageSignal:  "_isPage",
					ContainerID: "is-demo",
					LoadingText: "Loading more items…",
				}, initial...),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin:var(--sp-3) 0 0"),
					g.Text("Scroll down to trigger @get. Page signal auto-increments.")),
			)
		},
	})
}
