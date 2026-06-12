//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "pagination", Name: "Pagination", Category: "data",
		Summary: "Client-side pagination controls driven by a Datastar signal.",
		Code: `// init signal alongside the component
data.PaginationSignals("pg", 6)
data.Pagination(data.PaginationProps{
    ID:      "pg",
    Total:   48,
    PerPage: 6,
})`,
		Controls: []registry.Control{
			{Name: "pages", Type: registry.ControlEnum, Options: []string{"4", "8", "12"}, Default: "8"},
		},
		Render: func(p map[string]string) g.Node {
			total := 48
			switch p["pages"] {
			case "4":
				total = 24
			case "12":
				total = 72
			}
			perPage := 6
			return h.Div(
				// PaginationSignals must be rendered alongside the component to init the signal
				PaginationSignals("show", perPage),
				Pagination(PaginationProps{
					ID:      "show",
					Total:   total,
					PerPage: perPage,
				}),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "paginated-pages", Name: "Paginated Pages", Category: "data",
		Summary: "Animated page container driven by the Pagination signal — replays an entrance animation once per page switch.",
		Code: `data.PaginationSignals("pg", 1)
data.Pagination(data.PaginationProps{ID: "pg", Total: 4, PerPage: 1})
data.PaginatedPages(data.PaginatedPagesProps{
    ID:        "pg",
    Animation: data.PageAnimSlideUp, // slide-up | slide-left | fade | scale | flip | none
},
    page1, page2, page3, page4,
)`,
		Controls: []registry.Control{
			{Name: "animation", Type: registry.ControlEnum, Options: []string{"slide-up", "slide-left", "fade", "scale", "flip"}, Default: "slide-up"},
		},
		Render: func(p map[string]string) g.Node {
			anim := PageAnimation(p["animation"])
			if anim == "" {
				anim = PageAnimSlideUp
			}
			pageBox := func(label, bg string) g.Node {
				return h.Div(
					h.Style("border:var(--bw-2,2px) solid var(--ink);box-shadow:var(--shadow);padding:var(--sp-8);text-align:center;font-weight:900;font-size:var(--t-xl);background:"+bg),
					g.Text(label),
				)
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);align-items:center"),
				PaginationSignals("ppdemo", 1),
				Pagination(PaginationProps{ID: "ppdemo", Total: 4, PerPage: 1}),
				h.Div(h.Style("width:100%"),
					PaginatedPages(PaginatedPagesProps{ID: "ppdemo", Animation: anim},
						pageBox("Page one", "var(--yellow-bg,#fef08a)"),
						pageBox("Page two", "var(--cyan-bg,#a5f3fc)"),
						pageBox("Page three", "var(--violet-bg,#ddd6fe)"),
						pageBox("Page four", "var(--lime-bg,#d9f99d)"),
					),
				),
			)
		},
	})
}
