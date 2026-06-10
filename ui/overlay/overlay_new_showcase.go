//go:build showcase

package overlay

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "drawer", Name: "Drawer", Category: "overlay",
		Summary: "Datastar-gated slide-in side panel. Right or left placement, three sizes.",
		Code: `// trigger
h.Button(g.Attr("data-on:click", "$drawerOpen=true"), g.Text("Open drawer"))

// panel (rendered once per page, outside triggers)
overlay.Drawer(overlay.DrawerProps{Title: "Settings", OpenExpr: "$drawerOpen"},
    h.P(g.Text("Drawer body content.")),
)`,
		Controls: []registry.Control{
			{Name: "placement", Type: registry.ControlEnum, Options: []string{"right", "left"}, Default: "right"},
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg"}, Default: "md"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(
				g.Attr("data-signals", `{"drawerOpen2":false}`),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary},
					g.Attr("data-on:click", "$drawerOpen2=true"),
					g.Text("Open drawer"),
				),
				Drawer(DrawerProps{
					ID:        "drawer2",
					Title:     "Drawer panel",
					OpenExpr:  "$drawerOpen2",
					Placement: DrawerPlacement(p["placement"]),
					Size:      p["size"],
				},
					h.P(g.Text("Drawer body. Click the scrim or × to close.")),
					h.P(g.Text("Place forms, settings, or detail views here.")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "dropdown", Name: "Dropdown", Category: "overlay",
		Summary: "Datastar-driven contextual menu. Closes on outside-click via window listener.",
		Code: `overlay.Dropdown(
    overlay.DropdownProps{Signal: "menu"},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Actions ▾")),
    overlay.DropdownItem{Label: "Edit",   Icon: "lucide:edit"},
    overlay.DropdownItem{Label: "Copy",   Icon: "lucide:copy"},
    overlay.DropdownItem{Divider: true, Label: "Delete", Icon: "lucide:trash-2", Variant: "danger"},
)`,
		Render: func(p map[string]string) g.Node {
			return Dropdown(
				DropdownProps{Signal: "ddemo"},
				primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Actions ▾")),
				DropdownItem{Label: "Edit profile", Icon: "lucide:edit"},
				DropdownItem{Label: "Copy link", Icon: "lucide:copy"},
				DropdownItem{Label: "Share", Icon: "lucide:share"},
				DropdownItem{Divider: true, Label: "Delete", Icon: "lucide:trash-2", Variant: "danger"},
			)
		},
	})
}
