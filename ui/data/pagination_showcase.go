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
}
