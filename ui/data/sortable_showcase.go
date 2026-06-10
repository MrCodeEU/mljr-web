//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "sortable", Name: "Sortable List", Category: "data",
		Summary: "Drag-to-reorder list using HTML5 drag-and-drop API. Optional grab handle icon. OnChange callback with new value order.",
		Code: `data.Sortable(data.SortableProps{ID: "my-list"},
    data.SortableRow("a", "First item", false),
    data.SortableRow("b", "Second item", false),
    data.SortableRow("c", "Third item", false),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:1fr 1fr;gap:var(--sp-6)"),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-3)"), g.Text("Full row draggable:")),
					Sortable(SortableProps{ID: "sort1"},
						SortableRow("task1", "Design system tokens", false),
						SortableRow("task2", "Button component", false),
						SortableRow("task3", "Form inputs", false),
						SortableRow("task4", "Data table", false),
						SortableRow("task5", "Deploy to production", false),
					),
				),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-3)"), g.Text("Handle only:")),
					Sortable(SortableProps{ID: "sort2", Handle: true},
						SortableRow("item1", "Priority 1", true),
						SortableRow("item2", "Priority 2", true),
						SortableRow("item3", "Priority 3", true),
						SortableRow("item4", "Priority 4", true),
					),
				),
			)
		},
	})
}
