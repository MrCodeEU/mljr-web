package pages

import (
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/layout"
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
			h.Dl(h.Class("legal-dl"),
				legalDefinition(t("legal.impressum.section_ecg.name"), "Michael Reinegger"),
				legalDefinition(t("legal.impressum.section_ecg.address"), "Ahornstrasse 8, 4484 Kronstorf, Österreich"),
				legalDefinition(t("legal.impressum.section_ecg.email"), "hello@mljr.eu"),
				legalDefinition(t("legal.impressum.section_ecg.other_contact"), "reinemic2.0@gmail.com, michael-reinegger@tuta.io"),
				legalDefinition(t("legal.impressum.section_ecg.website"), "mljr.eu"),
				legalDefinition(t("legal.impressum.section_ecg.purpose"), t("legal.impressum.section_ecg.purpose_value")),
				legalDefinition(t("legal.impressum.section_ecg.owner"), "Michael Reinegger"),
			),
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
			h.P(g.Text(t("legal.datenschutz.section_external.body3"))),
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
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra: append([]g.Node{
				g.El("style", g.Raw(homepageCSS+legalCSS)),
			}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"), h.Class("legal-page sub-page"),
			subHead("/", tr(lang, "Start", "Home"), len(title),
				h.H1(g.Text(title)),
				// generated from the headings of .legal-shell
				h.Div(h.Class("legal-toc"),
					layout.TableOfContents(layout.TOCProps{
						Title:           i18n.T(lang, "legal.toc_title"),
						ContentSelector: ".legal-shell",
					}),
				),
			),
			h.Div(h.Class("legal-shell sub-body"), g.Group(content)),
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

const legalCSS = `
.legal-toc { font-size: .85rem; }
@media (max-width: 1024px) { .legal-toc { display: none; } }
.legal-section {
  border-top: 1px solid var(--ink);
  padding: var(--sp-5) 0 var(--sp-2);
}
.legal-section:first-child { border-top: 0; padding-top: 0; }
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
@media (max-width: 640px) {
  .legal-definition { grid-template-columns: 1fr; gap: var(--sp-1); }
}
`
