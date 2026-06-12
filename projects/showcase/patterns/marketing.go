//go:build showcase

package patterns

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.RegisterPattern(&registry.Pattern{
		Slug:        "marketing-pricing",
		Name:        "Pricing Page",
		Category:    "marketing",
		Description: "Hero + three-tier pricing cards + FAQ accordion. Common SaaS marketing pattern.",
		Render: func(theme, mode string) g.Node {
			th := token.Theme(theme)
			mo := token.Mode(mode)
			if th == "" {
				th = token.ThemeSwissBrut
			}
			if mo == "" {
				mo = token.ModeLight
			}

			plans := []pricingPlan{
				{
					Name:      "Free",
					Price:     "$0",
					Period:    "forever",
					Highlight: false,
					Features:  []string{"5 projects", "1 team member", "100MB storage", "Community support"},
					CTA:       "Get started",
				},
				{
					Name:      "Pro",
					Price:     "$19",
					Period:    "per month",
					Highlight: true,
					Badge:     "Most popular",
					Features:  []string{"Unlimited projects", "10 team members", "50GB storage", "Priority email support", "Custom domains"},
					CTA:       "Start free trial",
				},
				{
					Name:      "Enterprise",
					Price:     "$79",
					Period:    "per month",
					Highlight: false,
					Features:  []string{"Everything in Pro", "Unlimited team members", "500GB storage", "24/7 phone support", "SSO + SAML", "SLA guarantee"},
					CTA:       "Contact sales",
				},
			}

			planCards := make([]g.Node, len(plans))
			for i, plan := range plans {
				planCards[i] = pricingCard(plan)
			}

			faqs := []struct{ q, a string }{
				{"Can I change plans later?", "Yes, you can upgrade or downgrade at any time. Changes take effect immediately."},
				{"What payment methods do you accept?", "We accept all major credit cards, PayPal, and wire transfer for annual plans."},
				{"Is there a free trial?", "The Pro plan includes a 14-day free trial with no credit card required."},
				{"Can I cancel anytime?", "Yes, you can cancel your subscription at any time. No lock-in contracts."},
			}
			faqNodes := make([]g.Node, len(faqs))
			for i, faq := range faqs {
				faqNodes[i] = h.Details(
					h.Style("border-bottom:var(--bw-1) solid var(--line);padding:var(--sp-4) 0"),
					h.Summary(h.Style("font-weight:700;cursor:pointer;list-style:none;display:flex;justify-content:space-between"),
						g.Text(faq.q),
						icon.Icon("lucide:plus", icon.Props{Size: "1rem"}),
					),
					h.P(h.Style("margin-top:var(--sp-3);color:var(--muted);font-size:var(--t-sm)"), g.Text(faq.a)),
				)
			}

			return fullPage(th, mo,
				h.Div(
					layout.Navbar(layout.NavbarProps{},
						g.Text("mljr-ui"),
						g.Group{
							h.A(h.Href("#"), g.Text("Features")),
							h.A(h.Href("#"), h.Style("font-weight:800"), g.Text("Pricing")),
							h.A(h.Href("#"), g.Text("Docs")),
						},
						primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM}, g.Text("Sign up free")),
					),
					h.Main(
						layout.Container(layout.ContainerProps{},
							// Hero
							h.Div(
								h.Style("text-align:center;padding:var(--sp-12) 0 var(--sp-8)"),
								primitive.Display(primitive.DisplayProps{},
									primitive.GradientText(primitive.GradientTextProps{
										From: "var(--accent)", To: "var(--ink)",
									}, g.Text("Simple pricing")),
								),
								h.P(h.Style("font-size:var(--t-lg);color:var(--muted);max-width:48ch;margin:var(--sp-4) auto 0"),
									g.Text("Start free. Scale as you grow. No hidden fees."),
								),
							),
							// Plans
							h.Div(
								h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-6);margin-bottom:var(--sp-12)"),
								g.Group(planCards),
							),
							// Stats
							h.Div(
								h.Style("display:grid;grid-template-columns:repeat(4,1fr);gap:var(--sp-4);margin-bottom:var(--sp-12);text-align:center"),
								metricCell("10,000+", "developers"),
								metricCell("99.9%", "uptime SLA"),
								metricCell("< 50ms", "p99 latency"),
								metricCell("4.9/5", "rating"),
							),
							// FAQ
							h.Div(
								h.Style("max-width:640px;margin:0 auto var(--sp-12)"),
								h.H2(h.Style("font-size:var(--t-xl);font-weight:800;margin-bottom:var(--sp-6);text-align:center"), g.Text("Frequently asked questions")),
								g.Group(faqNodes),
							),
						),
					),
				),
			)
		},
	})
}

type pricingPlan struct {
	Name, Price, Period, Badge, CTA string
	Features                        []string
	Highlight                       bool
}

func pricingCard(p pricingPlan) g.Node {
	border := "var(--bw-1) solid var(--line)"
	bg := "var(--surface)"
	if p.Highlight {
		border = "var(--bw-2) solid var(--ink)"
		bg = "var(--accent)"
	}

	featureNodes := make([]g.Node, len(p.Features))
	for i, f := range p.Features {
		featureNodes[i] = h.Li(
			h.Style("display:flex;align-items:center;gap:var(--sp-2);font-size:var(--t-sm);margin-bottom:var(--sp-2)"),
			icon.Icon("lucide:check", icon.Props{Size: "0.9rem"}),
			g.Text(f),
		)
	}

	variant := token.Outline
	if p.Highlight {
		variant = token.Secondary
	}

	return h.Div(
		h.Style("border:"+border+";background:"+bg+";border-radius:var(--radius);padding:var(--sp-6);position:relative;display:flex;flex-direction:column;gap:var(--sp-4)"),
		g.If(p.Badge != "", h.Div(
			g.Attr("data-component", "badge"),
			h.Style("position:absolute;top:var(--sp-3);right:var(--sp-3)"),
			g.Text(p.Badge),
		)),
		h.Div(
			h.Div(h.Style("font-size:var(--t-sm);font-weight:700;margin-bottom:var(--sp-2)"), g.Text(p.Name)),
			h.Div(
				h.Span(h.Style("font-size:var(--t-3xl);font-weight:900"), g.Text(p.Price)),
				h.Span(h.Style("font-size:var(--t-sm);color:var(--muted);margin-left:var(--sp-1)"), g.Text(p.Period)),
			),
		),
		h.Ul(h.Style("list-style:none;margin:0;padding:0;flex:1"), g.Group(featureNodes)),
		primitive.Button(primitive.ButtonProps{Variant: variant, Attrs: []g.Node{h.Style("width:100%")}}, g.Text(p.CTA)),
	)
}

func metricCell(value, label string) g.Node {
	return h.Div(
		g.Attr("data-component", "card"),
		h.Style("padding:var(--sp-5)"),
		h.Div(h.Style("font-size:var(--t-2xl);font-weight:900"), g.Text(value)),
		h.Div(h.Style("font-size:var(--t-sm);color:var(--muted)"), g.Text(label)),
	)
}
