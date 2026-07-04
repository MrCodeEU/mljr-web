package pages

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// landingCSS is scoped to the landing page only (not the shared ui/css
// partials) since it's the one page that needs the postmark "stamp down"
// load animation and the postcard's tilt — both purely decorative, so they
// stay out of the component library.
const landingCSS = `
@keyframes nl-stamp-down {
  from { opacity: 0; transform: rotate(-18deg) scale(1.6) translateY(-10px); }
  to   { opacity: 1; transform: rotate(-12deg) scale(1) translateY(0); }
}
.nl-postmark {
  animation: nl-stamp-down .5s cubic-bezier(.2,1.4,.4,1) .35s both;
}
.nl-postcard {
  transform: rotate(-2deg);
  transition: transform .25s ease;
}
.nl-postcard:hover {
  transform: rotate(0deg) translateY(-4px);
}
@media (prefers-reduced-motion: reduce) {
  .nl-postmark { animation: none; }
  .nl-postcard, .nl-postcard:hover { transition: none; transform: none; }
}
`

// linkButton renders an <a> styled exactly like primitive.Button (same
// data-component contract the CSS keys off), for CTAs that must navigate
// rather than submit — primitive.Button always emits a <button>.
func linkButton(href, variant, size string, children ...g.Node) g.Node {
	nodes := []g.Node{
		g.Attr("data-component", "button"),
		g.Attr("data-variant", variant),
		g.Attr("data-size", size),
		h.Href(href),
		h.Style("text-decoration:none"),
	}
	nodes = append(nodes, children...)
	return h.A(nodes...)
}

// landingHeader is the public, unauthenticated header: brand mark, language
// toggle, sign-in link. Lighter than appHeader (shell.go) since there's no
// user session or group nav to show yet.
func landingHeader(re *core.RequestEvent, t func(string, ...any) string) g.Node {
	langToggle := special.LanguageToggle(special.LanguageToggleProps{
		Languages: []special.Language{
			{Code: "en", Label: "EN", Title: "English", Flag: "circle-flags:gb"},
			{Code: "de", Label: "DE", Title: "Deutsch", Flag: "circle-flags:de"},
		},
		Current:        currentLang(re),
		ReloadOnChange: true,
	})
	right := h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
		langToggle,
		linkButton("/login", "ghost", "sm", g.Text(t("newsletter.landing.sign_in"))),
	)
	return layout.Navbar(layout.NavbarProps{}, brand(t), g.Group(nil), right)
}

// postcard renders the signature element: a tilted vellum "edition" preview
// with a stamped postmark circle, standing in for a real mailed edition so
// the hero shows the product instead of describing it.
func postcard(t func(string, ...any) string) g.Node {
	return h.Div(
		h.Style("position:relative;max-width:340px;margin:0 auto"),
		h.Div(
			g.Attr("class", "nl-postcard"),
			g.Attr("data-component", "card"),
			h.Style("padding:var(--sp-5);background:var(--surface)"),
			h.Div(
				h.Style("font-family:var(--font-mono);font-size:var(--t-xs);letter-spacing:.08em;text-transform:uppercase;color:var(--muted)"),
				g.Text(t("newsletter.landing.postcard_kicker")),
			),
			h.Div(h.Style("height:1px;background:var(--line);opacity:.25;margin:var(--sp-2) 0 var(--sp-3)")),
			h.P(
				h.Style("font-family:var(--font-serif);font-size:var(--t-md);line-height:var(--leading-snug);margin:0 0 var(--sp-4)"),
				g.Text(t("newsletter.landing.postcard_question")),
			),
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);font-size:var(--t-sm)"),
				h.Div(h.Style("display:flex;gap:var(--sp-2);align-items:baseline"),
					h.Span(h.Style("color:var(--accent)"), g.Text("●")),
					h.Span(g.Text(t("newsletter.landing.postcard_reply_1"))),
				),
				h.Div(h.Style("display:flex;gap:var(--sp-2);align-items:baseline"),
					h.Span(h.Style("color:var(--accent)"), g.Text("●")),
					h.Span(g.Text(t("newsletter.landing.postcard_reply_2"))),
				),
				h.Div(h.Style("color:var(--muted)"), g.Text(t("newsletter.landing.postcard_more"))),
			),
		),
		h.Div(
			g.Attr("class", "nl-postmark"),
			h.Style("position:absolute;top:-18px;right:-18px;width:84px;height:84px;border-radius:50%;border:2px solid var(--accent);color:var(--accent);display:flex;flex-direction:column;align-items:center;justify-content:center;background:var(--bg);font-family:var(--font-mono);text-align:center;line-height:1.1"),
			h.Span(h.Style("font-size:.55rem;letter-spacing:.06em;text-transform:uppercase"), g.Text(t("newsletter.landing.postmark_label"))),
			icon.Icon("lucide:mail", icon.Props{Size: "1.3rem"}),
		),
	)
}

// stage is one step of the real open→reminder→grace→sent state machine
// (scheduler.go) — numbered here because the order is the actual mechanic
// a new user needs to understand, not decoration.
type stage struct {
	n, label, body string
}

func stagesSection(t func(string, ...any) string) g.Node {
	stages := []stage{
		{"01", t("newsletter.landing.stage_open_label"), t("newsletter.landing.stage_open_body")},
		{"02", t("newsletter.landing.stage_reminder_label"), t("newsletter.landing.stage_reminder_body")},
		{"03", t("newsletter.landing.stage_grace_label"), t("newsletter.landing.stage_grace_body")},
		{"04", t("newsletter.landing.stage_sent_label"), t("newsletter.landing.stage_sent_body")},
	}
	var cols []g.Node
	for _, s := range stages {
		cols = append(cols, h.Div(
			h.Style("flex:1;min-width:160px"),
			h.Div(h.Style("font-family:var(--font-mono);color:var(--accent);font-size:var(--t-sm);margin-bottom:var(--sp-1)"), g.Text(s.n)),
			h.Div(h.Style("font-weight:700;margin-bottom:var(--sp-1)"), g.Text(s.label)),
			h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin:0"), g.Text(s.body)),
		))
	}
	return h.Section(
		h.Style("max-width:920px;margin:0 auto;padding:var(--sp-12) var(--sp-4)"),
		h.Div(
			h.Style("font-family:var(--font-body);font-size:var(--t-xs);letter-spacing:.1em;text-transform:uppercase;color:var(--muted);margin-bottom:var(--sp-5);text-align:center"),
			g.Text(t("newsletter.landing.stages_eyebrow")),
		),
		h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-6);justify-content:center"), g.Group(cols)),
	)
}

// Landing is the public, unauthenticated home page — explains the product
// and gets a visitor to either sign in or create an account. Rendered by
// Dashboard (dashboard.go) when there's no session, instead of an immediate
// redirect to /login.
func Landing(re *core.RequestEvent) error {
	t := translator(re)
	return renderPage(re, 200, layout.PageShell(
		layout.PageProps{
			Title:       t("newsletter.landing.headline_line1") + " " + t("newsletter.landing.headline_line2") + " — Newsletter",
			Description: t("newsletter.landing.subhead"),
			Theme:       token.ThemeInk,
			Mode:        token.ModeLight,
			Lang:        currentLang(re),
			HeadExtra:   []g.Node{g.El("style", g.Raw(landingCSS))},
		},
		landingHeader(re, t),

		h.Main(
			h.Section(
				h.Style("max-width:1080px;margin:0 auto;padding:var(--sp-10) var(--sp-4) var(--sp-12);display:grid;grid-template-columns:1.1fr .9fr;gap:var(--sp-10);align-items:center"),
				g.Attr("class", "nl-hero"),
				h.Div(
					h.Div(
						h.Style("font-family:var(--font-body);font-size:var(--t-xs);letter-spacing:.12em;text-transform:uppercase;color:var(--accent);font-weight:700;margin-bottom:var(--sp-3)"),
						g.Text(t("newsletter.landing.eyebrow")),
					),
					h.H1(
						h.Style("font-family:var(--font-serif);font-size:var(--t-display);line-height:var(--leading-tight);margin:0 0 var(--sp-5);color:var(--ink)"),
						g.Text(t("newsletter.landing.headline_line1")), h.Br(),
						g.Text(t("newsletter.landing.headline_line2")),
					),
					h.P(
						h.Style("font-size:var(--t-md);color:var(--muted);max-width:46ch;margin:0 0 var(--sp-6);line-height:var(--leading-body)"),
						g.Text(t("newsletter.landing.subhead")),
					),
					h.Div(h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
						linkButton("/login", "primary", "lg", g.Text(t("newsletter.landing.sign_in"))),
						linkButton("/signup", "outline", "lg", g.Text(t("newsletter.landing.create_account"))),
					),
				),
				postcard(t),
			),

			h.Div(h.Style("height:1px;background:var(--line);opacity:.15;max-width:1080px;margin:0 auto")),
			stagesSection(t),
			h.Div(h.Style("height:1px;background:var(--line);opacity:.15;max-width:1080px;margin:0 auto")),

			h.Section(
				h.Style("max-width:680px;margin:0 auto;padding:var(--sp-12) var(--sp-4);text-align:center"),
				h.P(
					h.Style("font-family:var(--font-serif);font-size:var(--t-xl);line-height:var(--leading-snug);color:var(--ink);margin:0"),
					g.Text("“"+t("newsletter.landing.quote")+"”"),
				),
			),

			h.Section(
				h.Style("max-width:680px;margin:0 auto;padding:0 var(--sp-4) var(--sp-12);text-align:center"),
				h.H2(h.Style("font-family:var(--font-serif);font-size:var(--t-xl);margin:0 0 var(--sp-3)"), g.Text(t("newsletter.landing.cta_heading"))),
				h.P(h.Style("color:var(--muted);margin:0 0 var(--sp-5)"), g.Text(t("newsletter.landing.cta_body"))),
				h.Div(h.Style("display:flex;gap:var(--sp-3);justify-content:center;flex-wrap:wrap;margin-bottom:var(--sp-3)"),
					linkButton("/login", "primary", "lg", g.Text(t("newsletter.landing.sign_in"))),
					linkButton("/signup", "outline", "lg", g.Text(t("newsletter.landing.create_account"))),
				),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin:0"), g.Text(t("newsletter.landing.cta_invite_hint"))),
			),
		),

		h.Style("@media (max-width: 720px) { .nl-hero { grid-template-columns: 1fr !important; } }"),
	))
}
