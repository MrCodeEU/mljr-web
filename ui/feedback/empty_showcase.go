//go:build showcase

package feedback

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "empty-state", Name: "Empty State", Category: "feedback",
		Summary: "Centred placeholder for zero-data views. Icon, title, message, action buttons.",
		Code: `feedback.EmptyState(feedback.EmptyStateProps{
    Icon:    "lucide:inbox",
    Title:   "No results",
    Message: "Try adjusting your filters or search query.",
},
    primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Clear filters")),
)`,
		Controls: []registry.Control{
			{Name: "icon", Type: registry.ControlEnum, Options: []string{"lucide:inbox", "lucide:search", "lucide:folder", "lucide:file"}, Default: "lucide:inbox"},
		},
		Render: func(p map[string]string) g.Node {
			return EmptyState(EmptyStateProps{
				Icon:    p["icon"],
				Title:   "Nothing here yet",
				Message: "Add your first item to get started.",
			},
				primitive.Button(primitive.ButtonProps{Variant: token.Primary}, g.Text("Add item")),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline}, g.Text("Learn more")),
			)
		},
	})
}
