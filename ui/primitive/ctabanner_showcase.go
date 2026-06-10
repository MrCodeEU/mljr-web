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
		Slug: "cta-banner", Name: "CTA Banner", Category: "primitive",
		Summary: "Full-width call-to-action strip with title, description, and one or two action buttons.",
		Code: `primitive.CTABanner(primitive.CTABannerProps{
    Title:       "Ready to ship faster?",
    Description: "Join 5,000+ teams using mljr-ui.",
    CTAText:     "Get started free",
    CTAHref:     "/signup",
    SecondCTA:   "View docs",
    SecondHref:  "/docs",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				CTABanner(CTABannerProps{
					Title:       "Ready to ship faster?",
					Description: "Join 5,000+ teams already using mljr-ui in production.",
					CTAText:     "Get started free",
					CTAHref:     "#",
					SecondCTA:   "View docs",
					SecondHref:  "#",
				}),
				CTABanner(CTABannerProps{
					Title:   "New: Datastar 1.0.2 support",
					CTAText: "Read the changelog",
					CTAHref: "#",
					Variant: token.Secondary,
				}),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "announcement-bar", Name: "Announcement Bar", Category: "primitive",
		Summary: "Dismissable top-of-page strip. Datastar signal controls visibility.",
		Code: `primitive.AnnouncementBar(primitive.AnnouncementBarProps{
    CTAText: "Read the changelog",
    Href:    "/changelog",
},
    g.Text("✨ mljr-ui v2.0 is now available"),
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
				AnnouncementBar(AnnouncementBarProps{SignalName: "_a1", CTAText: "Learn more", Href: "#"},
					g.Text("🚀 mljr-ui v2.0 is now available with 90+ components"),
				),
				AnnouncementBar(AnnouncementBarProps{SignalName: "_a2", Variant: token.Secondary, CTAText: "Upgrade", Href: "#"},
					g.Text("Datastar 1.0.2 brings signal-patch, computed signals, and more"),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "countdown", Name: "Countdown", Category: "primitive",
		Summary: "Live countdown timer to a target datetime. Updates every second via JS setInterval.",
		Code: `primitive.Countdown(primitive.CountdownProps{
    Target:   "2026-12-31T23:59:59",
    OnExpire: "alert('Time is up!')",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-3)"), g.Text("Full (days : hours : min : sec)")),
					Countdown(CountdownProps{Target: "2026-12-31T23:59:59", ID: "cd1"}),
				),
				h.Div(
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-3)"), g.Text("Compact (seconds)")),
					Countdown(CountdownProps{Target: "2026-12-31T23:59:59", ID: "cd2", Compact: true}),
				),
			)
		},
	})
}
