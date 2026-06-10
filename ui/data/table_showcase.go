//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "table", Name: "Table", Category: "data",
		Summary: "Design-system styled data table. Supports striped rows and hover state.",
		Code: `data.Table(data.TableProps{Striped: true},
    h.THead(h.Tr(
        h.Th(g.Text("Name")),
        h.Th(g.Text("Status")),
        h.Th(g.Text("Amount")),
    )),
    h.TBody(
        h.Tr(h.Td(g.Text("Alice")), h.Td(g.Text("Active")), h.Td(g.Text("$120"))),
        h.Tr(h.Td(g.Text("Bob")),   h.Td(g.Text("Pending")), h.Td(g.Text("$85"))),
    ),
)`,
		Controls: []registry.Control{
			{Name: "striped", Type: registry.ControlBool, Default: "false"},
		},
		Render: func(p map[string]string) g.Node {
			return Table(TableProps{Striped: p["striped"] == "true"},
				h.THead(h.Tr(
					h.Th(g.Text("Name")),
					h.Th(g.Text("Role")),
					h.Th(g.Text("Status")),
					h.Th(g.Text("Revenue")),
				)),
				h.TBody(
					h.Tr(h.Td(g.Text("Alice Müller")), h.Td(g.Text("Engineer")), h.Td(g.Text("Active")), h.Td(g.Text("$4,200"))),
					h.Tr(h.Td(g.Text("Bob Chen")), h.Td(g.Text("Designer")), h.Td(g.Text("Active")), h.Td(g.Text("$3,800"))),
					h.Tr(h.Td(g.Text("Carol Singh")), h.Td(g.Text("PM")), h.Td(g.Text("Pending")), h.Td(g.Text("$5,100"))),
					h.Tr(h.Td(g.Text("Dave Park")), h.Td(g.Text("Engineer")), h.Td(g.Text("Inactive")), h.Td(g.Text("$0"))),
				),
			)
		},
	})
}
