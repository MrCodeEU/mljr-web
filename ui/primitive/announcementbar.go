package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AnnouncementBarProps struct {
	// SignalName is the Datastar signal (default "_announceDismissed").
	SignalName string
	Variant    token.Variant
	Href       string // optional CTA link
	CTAText    string // optional CTA label
}

// AnnouncementBar renders a dismissable top-of-page strip.
// Place at the very top of <body>, before Navbar.
// The Datastar signal controls visibility — persisting dismiss state requires
// localStorage via data-effect or a cookie handler.
func AnnouncementBar(p AnnouncementBarProps, children ...g.Node) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_announceDismissed"
	}
	if p.Variant == "" {
		p.Variant = token.Primary
	}
	sig := p.SignalName

	inner := []g.Node{
		h.Span(g.Attr("data-slot", "message"), g.Group(children)),
	}
	if p.CTAText != "" && p.Href != "" {
		inner = append(inner, h.A(
			h.Href(p.Href),
			g.Attr("data-slot", "cta"),
			g.Text(p.CTAText+" →"),
		))
	}
	inner = append(inner, h.Button(
		g.Attr("data-slot", "close"),
		h.Type("button"),
		g.Attr("aria-label", "Dismiss"),
		g.Attr("data-on:click", "$"+sig+"=true"),
		icon.Icon("lucide:x", icon.Props{Size: "1rem"}),
	))

	return h.Div(
		g.Attr("data-component", "announcement-bar"),
		g.Attr("data-variant", string(p.Variant)),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		g.Attr("data-show", "!$"+sig),
		h.Div(g.Attr("data-slot", "inner"), g.Group(inner)),
	)
}
