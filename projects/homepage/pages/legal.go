package pages

import (
	"strings"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

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

func Impressum(a AnalyticsConfig) g.Node {
	return legalPage(a, "Impressum", "Legal notice for mljr.eu.",
		legalSection("Angaben gemaess ECG und MedienG",
			legalDefinition("Name", "Michael Reinegger"),
			legalDefinition("Anschrift", "Ahornstrasse 8, 4484 Kronstorf, Oesterreich"),
			legalDefinition("E-Mail", "hello@mljr.eu"),
			legalDefinition("Weitere Kontaktadressen", "reinemic2.0@gmail.com, michael-reinegger@tuta.io"),
			legalDefinition("Website", "mljr.eu"),
			legalDefinition("Zweck der Website", "Persoenliche Portfolio-Website mit Informationen zu Projekten, Ausbildung, Berufserfahrung und Kontaktmoeglichkeit."),
			legalDefinition("Medieninhaber und Herausgeber", "Michael Reinegger"),
		),
		legalSection("Hinweis",
			h.P(g.Text("Diese Website ist eine rein persoenliche Portfolio-Website. Es werden keine Waren oder Dienstleistungen ueber diese Website verkauft.")),
		),
	)
}

func Datenschutz(a AnalyticsConfig) g.Node {
	return legalPage(a, "Datenschutz", "Privacy notice for mljr.eu.",
		legalSection("Verantwortlicher",
			h.P(g.Text("Michael Reinegger, Ahornstrasse 8, 4484 Kronstorf, Oesterreich")),
			h.P(g.Text("E-Mail: hello@mljr.eu")),
		),
		legalSection("Zweck der Website",
			h.P(g.Text("Diese Website dient als persoenliche Portfolio-Website. Sie stellt Projekte, Erfahrung, Ausbildung, Aktivitaetsdaten und Kontaktmoeglichkeiten dar.")),
		),
		legalSection("Hosting und Server-Logs",
			h.P(g.Text("Die Website wird auf einem VPS bei Contabo betrieben. Beim Aufruf der Website verarbeitet der Server technisch notwendige Zugriffsdaten, um die Website auszuliefern, Fehler zu erkennen und Missbrauch zu untersuchen.")),
			h.P(g.Text("Dabei koennen insbesondere Zeitpunkt, IP-Adresse, Host, HTTP-Methode, aufgerufene URL, HTTP-Status, Antwortzeit und User-Agent verarbeitet werden. Diese Daten werden nicht zur Erstellung persoenlicher Profile verwendet.")),
		),
		legalSection("Kontaktformular",
			h.P(g.Text("Wenn du das Kontaktformular verwendest, werden die von dir eingegebenen Daten verarbeitet, um deine Nachricht zu beantworten. Zur Spam-Abwehr wird eine ALTCHA-Challenge eingesetzt.")),
		),
		legalSection("Reichweitenmessung mit Umami",
			h.P(g.Text("Diese Website kann eine selbst gehostete Umami-Instanz zur einfachen Reichweitenmessung verwenden. Umami wird ohne Cookies eingebunden.")),
			h.P(g.Text("Erfasst werden aggregierte technische Nutzungsdaten wie Seitenaufrufe, Referrer, Browser, Betriebssystem, Geraetetyp und ungefaehre Herkunft. IP-Adressen werden nicht fuer persoenliche Profile verwendet. Die Messung dient dazu, die Nutzung der Website grob zu verstehen und technische oder inhaltliche Verbesserungen abzuleiten.")),
			h.P(g.Text("Die Einbindung ist nur aktiv, wenn die Website entsprechend konfiguriert ist. In der lokalen Entwicklungsumgebung ist sie standardmaessig deaktiviert.")),
		),
		legalSection("Cookies und lokale Speicherung",
			h.P(g.Text("Diese Website setzt keine Tracking-Cookies. Fuer die Darstellung koennen technisch notwendige lokale Einstellungen wie Theme oder Farbmodus im Browser gespeichert werden. Diese Einstellungen bleiben auf deinem Geraet und werden nicht zu Werbe- oder Profilingzwecken verwendet.")),
		),
		legalSection("Externe Inhalte",
			h.P(g.Text("Die Website laedt Schriftarten, Skripte und Stylesheets nach aktuellem Stand selbst gehostet aus. Es werden keine Google Fonts oder CDN-Schriften eingebunden.")),
			h.P(g.Text("Externe Links, zum Beispiel zu GitHub, LinkedIn oder Strava, fuehren zu Angeboten anderer Anbieter. Fuer deren Inhalte und Datenverarbeitung gelten die jeweiligen Datenschutzhinweise dieser Anbieter.")),
		),
		legalSection("Rechtsgrundlagen",
			h.P(g.Text("Die technische Auslieferung der Website, Server-Logs, Sicherheitsmassnahmen und einfache Reichweitenmessung erfolgen auf Grundlage berechtigter Interessen an Betrieb, Sicherheit und Verbesserung der Website. Die Verarbeitung von Kontaktanfragen erfolgt zur Bearbeitung deiner Anfrage.")),
		),
		legalSection("Speicherdauer",
			h.P(g.Text("Server-Logs werden nur so lange gespeichert, wie dies fuer Betrieb, Fehleranalyse und Sicherheit erforderlich ist. Kontaktanfragen werden so lange gespeichert, wie es fuer die Bearbeitung und eine nachvollziehbare Kommunikation erforderlich ist.")),
		),
		legalSection("Deine Rechte",
			h.P(g.Text("Du kannst Auskunft, Berichtigung, Loeschung, Einschraenkung der Verarbeitung und Widerspruch gegen bestimmte Verarbeitungen verlangen, soweit die gesetzlichen Voraussetzungen vorliegen. Du kannst dich dazu per E-Mail an hello@mljr.eu wenden.")),
			h.P(g.Text("Ausserdem besteht ein Beschwerderecht bei einer Datenschutzaufsichtsbehoerde, in Oesterreich insbesondere bei der Oesterreichischen Datenschutzbehoerde.")),
		),
	)
}

func legalPage(a AnalyticsConfig, title, description string, content ...g.Node) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       title + " - Michael Reinegger",
			Description: description,
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			HeadExtra: append([]g.Node{
				g.El("style", g.Raw(homepageCSS+legalCSS)),
			}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(),
		h.Main(h.Class("legal-page"),
			h.Div(h.Class("legal-outer"),
				// Sticky TOC sidebar (desktop)
				h.Aside(h.Class("legal-toc"),
					layout.TableOfContents(layout.TOCProps{
						Title:           "Contents",
						Sticky:          true,
						ContentSelector: ".legal-shell",
					}),
				),
				h.Div(h.Class("legal-shell"),
					h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text("Back")),
					h.H1(g.Text(title)),
					g.Group(content),
				),
			),
		),
		siteFooter(),
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

func siteNavbar() g.Node {
	return layout.Navbar(layout.NavbarProps{},
		h.A(h.Href("/"),
			h.Img(
				h.Src("/static/img/logo/Logo-h.png"),
				h.Alt("mljr.eu"),
				h.Style("height:32px;width:auto"),
			),
		),
		g.Group{
			h.A(h.Href("/#experience"), h.Class("nav-link-hide"), g.Text("Experience")),
			h.A(h.Href("/#projects"), g.Text("Projects")),
			h.A(h.Href("/#opensource"), h.Class("nav-link-hide"), g.Text("Open Source")),
			h.A(h.Href("/#activity"), h.Class("nav-link-hide"), g.Text("Activity")),
			h.A(h.Href("/#skills"), h.Class("nav-link-hide"), g.Text("Skills")),
			h.A(h.Href("/#contact"), g.Text("Contact")),
		},
		g.Group{
			special.ThemeToggle(),
			special.ModeToggle(),
			h.A(
				h.Href("https://github.com/MrCodeEU"),
				g.Attr("target", "_blank"),
				g.Attr("rel", "noopener noreferrer"),
				g.Attr("aria-label", "GitHub"),
				primitive.Button(
					primitive.ButtonProps{Variant: token.Outline, Size: token.SizeIcon},
					icon.Icon("lucide:github"),
				),
			),
		},
	)
}

func siteFooter() g.Node {
	return layout.Footer(layout.FooterProps{
		Attrs: []g.Node{g.Attr("data-homepage-footer", "true")},
		Brand: h.Div(
			h.Img(h.Src("/static/img/logo/Logo-h.png"), h.Alt("mljr.eu"), h.Style("height:32px;width:auto")),
		),
		Tagline: "Go, security and self-hosted infrastructure. Every component on this site is a Go function — no JS framework, no CDN, no adtech tracking.",
		Columns: []layout.FooterColumn{
			{Title: "Site", Links: []layout.FooterLink{
				{Label: "Experience", Href: "/#experience"},
				{Label: "Projects", Href: "/#projects"},
				{Label: "Open source", Href: "/#opensource"},
				{Label: "Homelab", Href: "/#homelab"},
				{Label: "Activity", Href: "/#activity"},
				{Label: "Skills", Href: "/#skills"},
			}},
			{Title: "Elsewhere", Links: []layout.FooterLink{
				{Label: "GitHub", Href: "https://github.com/MrCodeEU", External: true},
				{Label: "LinkedIn", Href: "https://www.linkedin.com/in/michael-reinegger", External: true},
				{Label: "Strava", Href: "https://www.strava.com/athletes/mrcode", External: true},
				{Label: "Status page", Href: "https://uptime.mljr.eu/status/all", External: true},
			}},
			{Title: "Legal", Links: []layout.FooterLink{
				{Label: "Impressum", Href: "/impressum"},
				{Label: "Datenschutz", Href: "/datenschutz"},
				{Label: "Contact", Href: "/#contact"},
			}},
		},
		Bottom: h.Div(
			h.Style("display:flex;flex-wrap:wrap;justify-content:space-between;align-items:center;gap:var(--sp-3);width:100%"),
			h.Span(g.Text("© "+time.Now().Format("2006")+" Michael Reinegger · mljr.eu")),
			h.Span(g.Text("built with Go · gomponents · Datastar · Tailwind v4")),
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
