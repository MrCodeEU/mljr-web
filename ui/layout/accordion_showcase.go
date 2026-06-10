//go:build showcase

package layout

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "accordion", Name: "Accordion", Category: "layout",
		Summary: "Collapsible content sections using native <details>/<summary> — no JS required.",
		Code: `layout.Accordion(layout.AccordionProps{},
    layout.AccordionItem(layout.AccordionItemProps{Title: "What is this?", Open: true},
        h.P(g.Text("Accordion content goes here.")),
    ),
    layout.AccordionItem(layout.AccordionItemProps{Title: "How does it work?"},
        h.P(g.Text("Pure HTML details/summary — no JavaScript.")),
    ),
)`,
		Controls: []registry.Control{
			{Name: "open", Type: registry.ControlBool, Default: "true"},
		},
		Render: func(p map[string]string) g.Node {
			return Accordion(AccordionProps{},
				AccordionItem(AccordionItemProps{Title: "What is mljr-web?", Open: p["open"] == "true"},
					h.P(g.Text("A self-hosted Go web stack: gomponents, Datastar, Tailwind v4. No JS framework, no CDN, no trackers.")),
				),
				AccordionItem(AccordionItemProps{Title: "How does the accordion work?"},
					h.P(g.Text("Native HTML <details>/<summary> elements — the browser handles open/close with zero JavaScript.")),
				),
				AccordionItem(AccordionItemProps{Title: "Can I nest content?"},
					h.P(g.Text("Yes — any gomponents nodes work inside an AccordionItem. Stack, Grid, lists, images, all valid.")),
				),
			)
		},
	})
}
