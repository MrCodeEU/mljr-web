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
		Slug: "command", Name: "Command Palette", Category: "overlay",
		PreviewHeight: "500px",
		Summary:       "⌘K command palette overlay. Filterable list of actions, grouped with keyboard navigation.",
		Code: `// 1. Place Command anywhere on the page
overlay.Command(overlay.CommandProps{
    Items: []overlay.CommandItem{
        {Label: "Go to Dashboard", Icon: "lucide:home", Href: "/", Group: "Navigate"},
        {Label: "New Project",     Icon: "lucide:plus", OnClick: "openNewProject()", Group: "Actions"},
    },
})

// 2. Open via Datastar signal or button
primitive.Button(..., g.Attr("data-on:click", "$_cmdOpen=true"), g.Text("⌘K"))

// Keyboard: Ctrl/⌘+K opens, Escape closes.`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4);align-items:center;padding-top:var(--sp-6)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm);text-align:center"), g.Text("Press the button (or Ctrl/⌘+K) to open the command palette.")),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary},
					g.Attr("data-on:click", "$_cmdOpen=true"),
					g.Text("Open command palette  ⌘K"),
				),
				Command(CommandProps{
					Items: []CommandItem{
						{Label: "Dashboard", Icon: "lucide:home", Href: "#", Group: "Navigate"},
						{Label: "Projects", Icon: "lucide:folder", Href: "#", Group: "Navigate"},
						{Label: "Analytics", Icon: "lucide:bar-chart-2", Href: "#", Group: "Navigate"},
						{Label: "Settings", Icon: "lucide:settings", Href: "#", Group: "Navigate"},
						{Label: "New Project", Icon: "lucide:plus", OnClick: "alert('New project!')", Group: "Actions"},
						{Label: "Export Data", Icon: "lucide:download", OnClick: "alert('Exporting...')", Group: "Actions"},
						{Label: "Toggle Dark Mode", Icon: "lucide:moon", OnClick: "$mode=$mode==='light'?'dark':'light'", Group: "Actions", Shortcut: "⌘D"},
						{Label: "GitHub", Icon: "simple-icons:github", Href: "#", Group: "Links"},
						{Label: "Documentation", Icon: "lucide:file-text", Href: "#", Group: "Links"},
					},
				}),
			)
		},
	})
}
