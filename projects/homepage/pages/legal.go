package pages

import (
	"strconv"
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

type AnalyticsConfig struct {
	UmamiScriptSrc string
	UmamiWebsiteID string
	UmamiHostURL   string
	UmamiDomains   string
}

func Impressum(lang string, a AnalyticsConfig) g.Node {
	t := func(key string) string { return i18n.T(lang, key) }
	return legalPage(lang, a, t("legal.impressum.title"), t("legal.impressum.description"),
		legalSection(t("legal.impressum.section_ecg.title"),
			legalDefinition(t("legal.impressum.section_ecg.name"), "Michael Reinegger"),
			legalDefinition(t("legal.impressum.section_ecg.address"), "Ahornstrasse 8, 4484 Kronstorf, Österreich"),
			legalDefinition(t("legal.impressum.section_ecg.email"), "hello@mljr.eu"),
			legalDefinition(t("legal.impressum.section_ecg.other_contact"), "reinemic2.0@gmail.com, michael-reinegger@tuta.io"),
			legalDefinition(t("legal.impressum.section_ecg.website"), "mljr.eu"),
			legalDefinition(t("legal.impressum.section_ecg.purpose"), t("legal.impressum.section_ecg.purpose_value")),
			legalDefinition(t("legal.impressum.section_ecg.owner"), "Michael Reinegger"),
		),
		legalSection(t("legal.impressum.section_notice.title"),
			h.P(g.Text(t("legal.impressum.section_notice.body"))),
		),
	)
}

func Datenschutz(lang string, a AnalyticsConfig) g.Node {
	t := func(key string) string { return i18n.T(lang, key) }
	return legalPage(lang, a, t("legal.datenschutz.title"), t("legal.datenschutz.description"),
		legalSection(t("legal.datenschutz.section_controller.title"),
			h.P(g.Text(t("legal.datenschutz.section_controller.address"))),
			h.P(g.Text(t("legal.datenschutz.section_controller.email"))),
		),
		legalSection(t("legal.datenschutz.section_purpose.title"),
			h.P(g.Text(t("legal.datenschutz.section_purpose.body"))),
		),
		legalSection(t("legal.datenschutz.section_hosting.title"),
			h.P(g.Text(t("legal.datenschutz.section_hosting.body1"))),
			h.P(g.Text(t("legal.datenschutz.section_hosting.body2"))),
		),
		legalSection(t("legal.datenschutz.section_form.title"),
			h.P(g.Text(t("legal.datenschutz.section_form.body"))),
		),
		legalSection(t("legal.datenschutz.section_umami.title"),
			h.P(g.Text(t("legal.datenschutz.section_umami.body1"))),
			h.P(g.Text(t("legal.datenschutz.section_umami.body2"))),
			h.P(g.Text(t("legal.datenschutz.section_umami.body3"))),
		),
		legalSection(t("legal.datenschutz.section_cookies.title"),
			h.P(g.Text(t("legal.datenschutz.section_cookies.body"))),
		),
		legalSection(t("legal.datenschutz.section_external.title"),
			h.P(g.Text(t("legal.datenschutz.section_external.body1"))),
			h.P(g.Text(t("legal.datenschutz.section_external.body2"))),
		),
		legalSection(t("legal.datenschutz.section_legal_basis.title"),
			h.P(g.Text(t("legal.datenschutz.section_legal_basis.body"))),
		),
		legalSection(t("legal.datenschutz.section_retention.title"),
			h.P(g.Text(t("legal.datenschutz.section_retention.body"))),
		),
		legalSection(t("legal.datenschutz.section_rights.title"),
			h.P(g.Text(t("legal.datenschutz.section_rights.body1"))),
			h.P(g.Text(t("legal.datenschutz.section_rights.body2"))),
		),
	)
}

func legalPage(lang string, a AnalyticsConfig, title, description string, content ...g.Node) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       title + " - Michael Reinegger",
			Description: description,
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra: append([]g.Node{
				g.El("style", g.Raw(homepageCSS+legalCSS)),
			}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("legal-page"),
			h.Div(h.Class("legal-outer"),
				// Sticky TOC sidebar (desktop)
				h.Aside(h.Class("legal-toc"),
					layout.TableOfContents(layout.TOCProps{
						Title:           i18n.T(lang, "legal.toc_title"),
						Sticky:          true,
						ContentSelector: ".legal-shell",
					}),
				),
				h.Div(h.Class("legal-shell"),
					h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
					h.H1(g.Text(title)),
					g.Group(content),
				),
			),
		),
		siteFooter(lang),
	)
}

func legalSection(title string, children ...g.Node) g.Node {
	return h.Section(h.Class("legal-section"),
		h.H2(g.Text(title)),
		g.Group(children),
	)
}

func legalDefinition(term, value string) g.Node {
	return h.Div(h.Class("legal-definition"),
		h.Dt(g.Text(term)),
		h.Dd(g.Text(value)),
	)
}

func AnalyticsHead(a AnalyticsConfig) []g.Node {
	if strings.TrimSpace(a.UmamiScriptSrc) == "" || strings.TrimSpace(a.UmamiWebsiteID) == "" {
		return nil
	}
	attrs := []g.Node{
		h.Defer(),
		h.Src(a.UmamiScriptSrc),
		g.Attr("data-website-id", a.UmamiWebsiteID),
		g.Attr("data-do-not-track", "true"),
		g.Attr("data-exclude-search", "true"),
	}
	if strings.TrimSpace(a.UmamiHostURL) != "" {
		attrs = append(attrs, g.Attr("data-host-url", a.UmamiHostURL))
	}
	if strings.TrimSpace(a.UmamiDomains) != "" {
		attrs = append(attrs, g.Attr("data-domains", a.UmamiDomains))
	}
	return []g.Node{h.Script(attrs...)}
}

func siteNavbar(lang string) g.Node {
	return layout.Navbar(layout.NavbarProps{},
		h.A(h.Href("/"),
			h.Img(
				h.Src("/static/img/logo/Logo-h.png"),
				h.Alt("mljr.eu"),
				h.Width("172"),
				h.Height("32"),
				h.Style("height:32px;width:auto"),
			),
		),
		g.Group{
			h.A(h.Href("/#experience"), h.Class("nav-link-hide"), g.Text(i18n.T(lang, "nav.experience"))),
			h.A(h.Href("/#projects"), g.Text(i18n.T(lang, "nav.projects"))),
			h.A(h.Href("/#opensource"), h.Class("nav-link-hide"), g.Text(i18n.T(lang, "nav.opensource"))),
			h.A(h.Href("/#activity"), h.Class("nav-link-hide"), g.Text(i18n.T(lang, "nav.activity"))),
			h.A(h.Href("/#skills"), h.Class("nav-link-hide"), g.Text(i18n.T(lang, "nav.skills"))),
			h.A(h.Href("/#contact"), g.Text(i18n.T(lang, "nav.contact"))),
		},
		g.Group{
			special.LanguageToggle(special.LanguageToggleProps{
				Languages: []special.Language{
					{Code: "en", Label: "EN", Title: "English"},
					{Code: "de", Label: "DE", Title: "Deutsch"},
				},
				Current:        lang,
				ReloadOnChange: true,
			}),
			special.ThemeToggle(),
			special.ModeToggle(),
			h.A(
				h.Href("https://github.com/MrCodeEU"),
				g.Attr("target", "_blank"),
				g.Attr("rel", "noopener noreferrer"),
				g.Attr("aria-label", "GitHub"),
				primitive.Button(
					primitive.ButtonProps{
						Variant: token.Outline,
						Size:    token.SizeIcon,
						Attrs:   []g.Node{g.Attr("aria-hidden", "true"), g.Attr("tabindex", "-1")},
					},
					icon.Icon("lucide:github"),
				),
			),
		},
	)
}

func siteFooter(lang string) g.Node {
	t := func(key string) string { return i18n.T(lang, key) }
	year, _ := strconv.Atoi(time.Now().Format("2006"))
	return layout.Footer(layout.FooterProps{
		Attrs: []g.Node{g.Attr("data-homepage-footer", "true")},
		Brand: h.Div(
			h.Img(h.Src("/static/img/logo/Logo-h.png"), h.Alt("mljr.eu"), h.Width("172"), h.Height("32"), h.Style("height:32px;width:auto")),
		),
		Tagline: t("footer.tagline"),
		Columns: []layout.FooterColumn{
			{Title: t("footer.site_title"), Links: []layout.FooterLink{
				{Label: t("nav.experience"), Href: "/#experience"},
				{Label: t("nav.projects"), Href: "/#projects"},
				{Label: t("nav.opensource"), Href: "/#opensource"},
				{Label: t("footer.homelab"), Href: "/#homelab"},
				{Label: t("nav.activity"), Href: "/#activity"},
				{Label: t("nav.skills"), Href: "/#skills"},
			}},
			{Title: t("footer.elsewhere_title"), Links: []layout.FooterLink{
				{Label: t("footer.github"), Href: "https://github.com/MrCodeEU", External: true},
				{Label: t("footer.linkedin"), Href: "https://www.linkedin.com/in/mrcodeeu/", External: true},
				{Label: t("footer.strava"), Href: "https://www.strava.com/athletes/123496455", External: true},
				{Label: t("footer.status_page"), Href: "https://uptime.mljr.eu/status/all", External: true},
			}},
			{Title: t("footer.legal_title"), Links: []layout.FooterLink{
				{Label: t("footer.impressum"), Href: "/impressum"},
				{Label: t("footer.datenschutz"), Href: "/datenschutz"},
				{Label: t("footer.contact"), Href: "/#contact"},
			}},
		},
		Bottom: h.Div(
			h.Style("display:flex;flex-wrap:wrap;justify-content:space-between;align-items:center;gap:var(--sp-3);width:100%"),
			h.Span(g.Text(i18n.T(lang, "footer.copyright", year))),
			h.Span(g.Text(t("footer.built_with"))),
		),
	})
}

const legalCSS = `
.legal-page {
  min-height: 70vh;
  padding: var(--sp-12) var(--sp-4);
}
.legal-outer {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: var(--sp-8);
  max-width: 1140px;
  margin: 0 auto;
  align-items: start;
}
.legal-toc {
  position: sticky;
  top: var(--sp-6);
}
.legal-shell {
  min-width: 0;
  background: var(--surface);
  border: var(--bw-2) solid var(--line);
  box-shadow: var(--shadow);
  padding: clamp(var(--sp-5), 4vw, var(--sp-10));
}
@media (max-width: 768px) {
  .legal-outer { grid-template-columns: 1fr; }
  .legal-toc { display: none; }
}
.legal-shell h1 {
  font-size: clamp(2.4rem, 8vw, 5rem);
  line-height: .95;
  margin: var(--sp-4) 0 var(--sp-8);
  font-weight: 950;
}
.legal-section {
  border-top: var(--bw-1) solid var(--line);
  padding: var(--sp-6) 0 0;
  margin-top: var(--sp-6);
}
.legal-section h2 {
  font-size: var(--t-xl);
  font-weight: 900;
  margin: 0 0 var(--sp-3);
}
.legal-section p {
  max-width: 74ch;
  margin: 0 0 var(--sp-3);
  color: var(--muted);
}
.legal-definition {
  display: grid;
  grid-template-columns: minmax(160px, 220px) 1fr;
  gap: var(--sp-3);
  padding: var(--sp-2) 0;
  border-bottom: var(--bw-1) solid color-mix(in srgb, var(--line) 20%, transparent);
}
.legal-definition dt {
  font-weight: 900;
}
.legal-definition dd {
  margin: 0;
  color: var(--muted);
  overflow-wrap: anywhere;
}
.legal-back {
  display: inline-flex;
  align-items: center;
  gap: var(--sp-2);
  font-weight: 900;
  text-decoration: none;
}
@media (max-width: 640px) {
  .legal-page { padding: var(--sp-6) var(--sp-3); }
  .legal-definition { grid-template-columns: 1fr; gap: var(--sp-1); }
}
`
