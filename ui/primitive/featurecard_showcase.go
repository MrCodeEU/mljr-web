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
		Slug: "feature-card", Name: "Feature Card", Category: "primitive",
		Summary: "Marketing card with colored icon box, title, and description. Optionally a link.",
		Code: `primitive.FeatureCard(primitive.FeatureCardProps{
    Icon:  "lucide:zap",
    Title: "Blazing fast",
    Tone:  token.ToneYellow,
}, h.P(g.Text("Sub-millisecond server renders.")))`,
		Render: func(p map[string]string) g.Node {
			features := []struct {
				icon, title, desc string
				tone              token.Tone
			}{
				{"lucide:zap", "Blazing fast", "Sub-millisecond server renders. No hydration overhead.", token.ToneYellow},
				{"lucide:shield", "Type safe", "Go generics and gomponents give compile-time guarantees.", token.ToneCyan},
				{"lucide:puzzle", "Composable", "Every component is a plain Go function — easy to extend.", token.ToneViolet},
				{"lucide:globe", "Zero CDN", "Self-hosted static assets. Full control of your bundle.", token.ToneLime},
				{"lucide:server", "SSR first", "HTML rendered on the server. Works without JS enabled.", token.TonePink},
				{"lucide:layers", "Themeable", "4 themes × 2 modes via CSS custom properties.", token.ToneMint},
			}
			nodes := make([]g.Node, len(features))
			for i, f := range features {
				nodes[i] = FeatureCard(FeatureCardProps{Icon: f.icon, Title: f.title, Tone: f.tone},
					h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text(f.desc)),
				)
			}
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:var(--sp-3)"),
				g.Group(nodes),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "pricing-card", Name: "Pricing Card", Category: "primitive",
		PreviewHeight: "560px",
		Summary: "Pricing tier card with price display, feature checklist, and CTA button.",
		Code: `primitive.PricingCard(primitive.PricingCardProps{
    Name:  "Pro",
    Price: "€29",
    Period: "/month",
    Features: []string{"Unlimited projects", "Custom domain", "Priority support"},
    CTA: "Get started",
    CTAHref: "/signup",
    Highlighted: true,
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:var(--sp-4)"),
				PricingCard(PricingCardProps{
					Name:        "Free",
					Price:       "€0",
					Period:      "/month",
					Description: "Perfect for side projects",
					Features:    []string{"3 projects", "Community support", "Basic analytics", "1 GB storage"},
					CTA:         "Start free",
					CTAHref:     "#",
				}),
				PricingCard(PricingCardProps{
					Name:        "Pro",
					Price:       "€19",
					Period:      "/month",
					Description: "For indie developers",
					Features:    []string{"Unlimited projects", "Email support", "Advanced analytics", "20 GB storage", "Custom domain"},
					CTA:         "Start trial",
					CTAHref:     "#",
					Highlighted: true,
				}),
				PricingCard(PricingCardProps{
					Name:        "Team",
					Price:       "€49",
					Period:      "/month",
					Description: "For small teams",
					Features:    []string{"Everything in Pro", "10 seats", "Priority support", "100 GB storage", "SSO / SAML", "SLA"},
					CTA:         "Contact sales",
					CTAHref:     "#",
				}),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "color-swatch", Name: "Color Swatch", Category: "primitive",
		Summary: "Single color preview tile. Use ColorSwatchGroup for palettes. Supports square and circle shapes.",
		Code: `primitive.ColorSwatch(primitive.ColorSwatchProps{
    Color: "#f28d1d",
    Label: "#f28d1d",
})
primitive.ColorSwatchGroup(
    primitive.ColorSwatch(primitive.ColorSwatchProps{Color: "var(--primary)"}),
    primitive.ColorSwatch(primitive.ColorSwatchProps{Color: "var(--accent)"}),
)`,
		Render: func(p map[string]string) g.Node {
			themeColors := []struct{ color, label string }{
				{"var(--primary)", "--primary"},
				{"var(--accent)", "--accent"},
				{"var(--success)", "--success"},
				{"var(--warning)", "--warning"},
				{"var(--danger)", "--danger"},
				{"var(--muted)", "--muted"},
				{"var(--surface-2)", "--surface-2"},
				{"var(--line)", "--line"},
			}
			swatches := make([]g.Node, len(themeColors))
			for i, c := range themeColors {
				swatches[i] = ColorSwatch(ColorSwatchProps{Color: c.color, Label: c.label, Outline: true})
			}
			brandColors := []struct{ color, label string }{
				{"#f28d1d", "#f28d1d"},
				{"#a40054", "#a40054"},
				{"#00ADD8", "Go blue"},
				{"#3178C6", "TS blue"},
				{"#CE422B", "Rust"},
			}
			brand := make([]g.Node, len(brandColors))
			for i, c := range brandColors {
				brand[i] = ColorSwatch(ColorSwatchProps{Color: c.color, Label: c.label, Shape: "circle"})
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Theme tokens")),
					ColorSwatchGroup(swatches...),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Brand (circle)")),
					ColorSwatchGroup(brand...),
				),
			)
		},
	})
}
