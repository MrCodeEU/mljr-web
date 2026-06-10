//go:build showcase

package primitive

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "split-button", Name: "Split Button", Category: "primitive",
		Summary: "Primary action button with chevron dropdown for alternate actions. Datastar signal controls dropdown open state.",
		Code: `// import "mljr-web/ui/primitive"
primitive.SplitButton(primitive.SplitButtonProps{
    Label:   "Deploy",
    Variant: token.Primary,
    OnClick: "alert('Deploying to production')",
    Items: []primitive.SplitButtonItem{
        {Label: "Deploy to Staging", OnClick: "alert('staging')"},
        {Label: "Deploy to Preview", OnClick: "alert('preview')"},
        {Label: "View Logs", Href: "/logs"},
    },
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.H4(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-3)"), g.Text("Variants")),
					h.Div(
						h.Style("display:flex;gap:var(--sp-4);flex-wrap:wrap;align-items:flex-start"),
						SplitButton(SplitButtonProps{
							Label:      "Deploy",
							Variant:    token.Primary,
							OnClick:    "alert('Deploying to production')",
							SignalName: "_sb1",
							Items: []SplitButtonItem{
								{Label: "Deploy to Staging", OnClick: "alert('staging')"},
								{Label: "Deploy to Preview", OnClick: "alert('preview')"},
								{Label: "View Logs"},
							},
						}),
						SplitButton(SplitButtonProps{
							Label:      "Save",
							Variant:    token.Outline,
							OnClick:    "alert('Saved')",
							SignalName: "_sb2",
							Items: []SplitButtonItem{
								{Label: "Save as Draft", OnClick: "alert('draft')"},
								{Label: "Save & Publish", OnClick: "alert('published')"},
							},
						}),
						SplitButton(SplitButtonProps{
							Label:      "Download",
							Variant:    token.Ghost,
							OnClick:    "alert('Downloading PDF')",
							SignalName: "_sb3",
							Items: []SplitButtonItem{
								{Label: "Export as CSV", OnClick: "alert('csv')"},
								{Label: "Export as JSON", OnClick: "alert('json')"},
								{Label: "Export as XML", OnClick: "alert('xml')"},
							},
						}),
					),
				),
			)
		},
	})
}
