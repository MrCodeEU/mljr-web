//go:build showcase

package data

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "data-grid", Name: "Data Grid", Category: "data",
		PreviewHeight: "520px",
		Summary: "Sortable, filterable, paginated data table. All logic client-side — no server roundtrip for sort/filter.",
		Code: `data.DataGrid(data.DataGridProps{
    ID:       "users",
    Search:   true,
    PageSize: 5,
    Columns: []data.DataGridColumn{
        {Key: "name",  Label: "Name",       Sortable: true},
        {Key: "role",  Label: "Role",       Sortable: true},
        {Key: "email", Label: "Email"},
        {Key: "status",Label: "Status",     Sortable: true},
    },
    Rows: []map[string]string{
        {"name":"Alex Chen","role":"Engineer","email":"alex@...","status":"Active"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			rows := []map[string]string{
				{"name": "Alex Chen", "role": "Senior Engineer", "dept": "Backend", "joined": "2022-03", "status": "Active"},
				{"name": "Jordan Lee", "role": "Product Designer", "dept": "Design", "joined": "2021-07", "status": "Active"},
				{"name": "Sam Park", "role": "DevOps Lead", "dept": "Infrastructure", "joined": "2020-11", "status": "Active"},
				{"name": "Morgan Wu", "role": "Data Scientist", "dept": "Analytics", "joined": "2023-01", "status": "Active"},
				{"name": "Riley Kim", "role": "Frontend Engineer", "dept": "Frontend", "joined": "2022-08", "status": "On leave"},
				{"name": "Casey Brown", "role": "Product Manager", "dept": "Product", "joined": "2021-02", "status": "Active"},
				{"name": "Dana Miller", "role": "QA Engineer", "dept": "QA", "joined": "2023-06", "status": "Active"},
				{"name": "Avery Davis", "role": "Backend Engineer", "dept": "Backend", "joined": "2022-05", "status": "Active"},
				{"name": "Quinn Wilson", "role": "UX Researcher", "dept": "Design", "joined": "2021-09", "status": "Inactive"},
				{"name": "Drew Martinez", "role": "Site Reliability", "dept": "Infrastructure", "joined": "2020-04", "status": "Active"},
				{"name": "Sage Thompson", "role": "Mobile Engineer", "dept": "Mobile", "joined": "2023-02", "status": "Active"},
				{"name": "Blake Anderson", "role": "Security Engineer", "dept": "Security", "joined": "2022-01", "status": "Active"},
			}
			return DataGrid(DataGridProps{
				ID:       "demo-grid",
				Search:   true,
				PageSize: 5,
				Columns: []DataGridColumn{
					{Key: "name", Label: "Name", Sortable: true},
					{Key: "role", Label: "Role", Sortable: true},
					{Key: "dept", Label: "Dept", Sortable: true},
					{Key: "joined", Label: "Joined", Sortable: true, Width: "90px"},
					{Key: "status", Label: "Status", Sortable: true, Width: "80px"},
				},
				Rows: rows,
			})
		},
	})
}
