package pages

import (
	"strconv"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/swiss"
)

// tr picks the German or English string. Section titles stay German by
// design (English only as sublabel); tr is for running UI copy.
func tr(lang, de, en string) string {
	if lang == "de" {
		return de
	}
	return en
}

func cvPath(lang string) string {
	if lang == "de" {
		return "/static/cv/lebenslauf.pdf"
	}
	return "/static/cv/lebenslauf_eng.pdf"
}

// sitePrimary links the CV and the contact form from every page; contact is
// "#contact" on the homepage and the homepage's anchor elsewhere.
func sitePrimary(lang string, onHome bool) []g.Node {
	contact := "#contact"
	if !onHome {
		contact = tr(lang, "/de/", "/") + "#contact"
	}
	return []g.Node{
		swiss.BtnLink(cvPath(lang), true, []g.Node{g.Attr("type", "application/pdf")}, g.Text("CV ↓")),
		swiss.BtnLink(contact, false, nil, g.Text(tr(lang, "Kontakt", "Contact"))),
	}
}

// siteActions: language switch (cookie + reload) and light/dark toggle.
// The mode button reads the $mode signal seeded by special.ThemeToggleRoot.
func siteActions(lang string) []g.Node {
	langBtn := func(code, label string) g.Node {
		return swiss.Btn(false, []g.Node{
			g.Attr("lang", code),
			g.If(lang == code, g.Attr("aria-current", "true")),
			// same page in the other language: German lives under /de/
			g.Attr("data-on:click", `document.cookie='lang=`+code+`;path=/;max-age=31536000;samesite=lax';var p=location.pathname.replace(/^\/de(?=\/|$)/,'')||'/';location.href=('`+code+`'==='de'?'/de'+p:p)+location.search+location.hash`),
		}, g.Text(label))
	}
	return []g.Node{
		h.Div(g.Attr("data-component", "swiss-segmented"), h.Role("group"), g.Attr("aria-label", tr(lang, "Sprache", "Language")),
			langBtn("de", "DE"), langBtn("en", "EN"),
		),
		swiss.Btn(false, []g.Node{
			g.Attr("data-on:click", `$mode = $mode === 'light' ? 'dark' : 'light'`),
			g.Attr("data-attr:aria-pressed", `$mode === 'dark' ? 'true' : 'false'`),
			g.Attr("aria-label", tr(lang, "Dunkelmodus", "Dark mode")),
		},
			h.Span(g.Attr("data-text", `$mode === 'dark' ? '`+tr(lang, "Hell", "Light")+`' : '`+tr(lang, "Dunkel", "Dark")+`'`), g.Text(tr(lang, "Dunkel", "Dark"))),
		),
	}
}

// siteNavbar is the colophon for sub pages (no section ruler).
func siteNavbar(lang string) g.Node {
	return g.Group{
		// Inside the homepage's lab iframe, game pages drop their chrome.
		h.Script(g.Raw(`if(window.self!==window.top)document.documentElement.classList.add('embedded');`)),
		swiss.Colophon(swiss.ColophonProps{
			Brand:     "mljr.eu",
			BrandSub:  "— Michael Reinegger · Linz (AT)",
			Primary:   sitePrimary(lang, false),
			Actions:   siteActions(lang),
			SkipLabel: tr(lang, "Zum Inhalt springen", "Skip to content"),
		}),
	}
}

func siteFooter(lang string) g.Node {
	t := func(key string) string { return i18n.T(lang, key) }
	year, _ := strconv.Atoi(time.Now().Format("2006"))
	type link struct {
		label, href string
		ext         bool
	}
	col := func(title string, links ...link) g.Node {
		items := make([]g.Node, len(links))
		for i, l := range links {
			items[i] = h.Li(h.A(h.Href(l.href),
				g.If(l.ext, g.Group{g.Attr("target", "_blank"), g.Attr("rel", "noopener noreferrer")}),
				g.Text(l.label), g.If(l.ext, g.Text(" ↗")),
			))
		}
		return h.Div(h.Class("ft-col"), h.H2(h.Class("ft-title"), g.Text(title)), h.Ul(items...))
	}
	return h.Footer(h.Class("site-footer"),
		// Bauhaus poster band: wordmark, circle, triangle, back to top.
		h.Div(h.Class("ft-poster"),
			h.Div(h.Class("ftp-word"), g.Attr("aria-hidden", "true"), h.Span(g.Text("MLJR")), h.Span(g.Text(".EU")),
				h.Small(g.Text(t("footer.tagline")))),
			h.Div(h.Class("ftp-circle"), g.Attr("aria-hidden", "true"), h.I()),
			h.Div(h.Class("ftp-tri"), g.Attr("aria-hidden", "true"), h.I()),
			h.A(h.Class("ftp-top"), h.Href("#hero"), h.Span(h.Class("ftp-arrow"), g.Text("↑")), h.Span(g.Text(tr(lang, "Nach oben", "Back to top")))),
		),
		h.Div(h.Class("ft-grid"),
			h.Div(h.Class("ft-col ft-brand"),
				h.B(g.Text("Michael Reinegger")),
				h.P(g.Text(tr(lang, "Linz, Österreich", "Linz, Austria"))),
			),
			col(t("footer.site_title"),
				link{t("nav.experience"), "/#experience", false},
				link{t("nav.projects"), "/#projects", false},
				link{t("nav.opensource"), "/#opensource", false},
				link{t("footer.homelab"), "/#homelab", false},
				link{t("nav.activity"), "/#activity", false},
				link{t("nav.skills"), "/#skills", false},
			),
			col(t("footer.elsewhere_title"),
				link{t("footer.github"), "https://github.com/MrCodeEU", true},
				link{t("footer.linkedin"), "https://www.linkedin.com/in/mrcodeeu/", true},
				link{t("footer.strava"), "https://www.strava.com/athletes/123496455", true},
				link{t("footer.status_page"), "https://uptime.mljr.eu/status/all", true},
			),
			col(t("footer.legal_title"),
				link{t("footer.impressum"), "/impressum", false},
				link{t("footer.datenschutz"), "/datenschutz", false},
				link{t("footer.contact"), "/#contact", false},
			),
		),
		h.Div(h.Class("ft-bottom"),
			h.Span(g.Text(i18n.T(lang, "footer.copyright", year))),
			h.Span(g.Text(t("footer.built_with"))),
		),
	)
}

// homepageCSS is shared by every homepage-project page (home, games, legal,
// project detail): footer, page gutters and the homepage's section layouts.
const homepageCSS = subCSS + bauhausCSS + `
.site-footer { margin: 0 var(--swiss-gutter); border-top: 3px solid var(--ink); margin-top: 4rem; font-size: .85rem; position: relative; z-index: 1; background: var(--bg); }
.ft-poster { display: grid; grid-template-columns: repeat(12, minmax(0, 1fr)); height: clamp(200px, 22vw, 340px); margin-bottom: 0; border-bottom: 3px solid var(--ink); }
.ftp-word { grid-column: span 6; background: var(--ink); color: var(--bg); display: flex; flex-wrap: wrap; align-content: space-between; padding: var(--swiss-pad); font-weight: 900; letter-spacing: -0.07em; line-height: .8; overflow: hidden; container-type: inline-size; }
.ftp-word span { font-size: 21cqi; }
.ftp-word small { flex-basis: 100%; order: -1; font-family: var(--font-mono); font-size: .7rem; font-weight: 700; letter-spacing: .04em; text-transform: uppercase; line-height: 1.5; opacity: .75; max-width: 60ch; }
.ftp-word span:nth-of-type(2) { color: var(--second); }
.ftp-circle { grid-column: span 2; background: var(--block); display: grid; place-items: center; }
.ftp-circle i { width: 72%; aspect-ratio: 1; border-radius: 50%; background: var(--second); }
.ftp-tri { grid-column: span 2; background: var(--bg); position: relative; box-shadow: inset 1px 0 0 var(--ink), inset -1px 0 0 var(--ink); }
.ftp-tri i { position: absolute; inset: 14% 10% 0; background: var(--third); clip-path: polygon(50% 0, 100% 100%, 0 100%); }
.ftp-top { grid-column: span 2; background: var(--bg); color: var(--ink); text-decoration: none; display: flex; flex-direction: column; justify-content: space-between; padding: var(--swiss-pad); font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; }
.ftp-arrow { font-family: var(--font-body); font-size: clamp(3rem, 7vw, 7rem); font-weight: 900; line-height: .8; }
.ftp-top:hover { background: var(--block); color: var(--block-fg); }
.ft-grid { display: grid; grid-template-columns: 3fr 3fr 3fr 3fr; }
.ft-col { padding: var(--swiss-pad); padding-left: 0; border-right: 1px solid var(--ink); }
.ft-col + .ft-col { padding-left: var(--swiss-pad); }
.ft-col:last-child { border-right: none; }
.ft-brand b { font-size: clamp(1.3rem, 1rem + 1.2vw, 2rem); font-weight: 900; letter-spacing: -0.04em; display: block; line-height: 1; }
.ft-brand p { margin: .75rem 0 0; max-width: 32ch; color: var(--muted); }
.ft-title { font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; margin: 0 0 .6rem; padding-bottom: .45rem; border-bottom: 1px solid var(--ink); }
.ft-col ul { list-style: none; margin: 0; padding: 0; display: grid; gap: .3rem; }
.ft-col a { color: inherit; text-decoration: none; }
.ft-col a:hover { color: var(--main-text); text-decoration: underline; }
.ft-bottom { display: flex; flex-wrap: wrap; justify-content: space-between; gap: .5rem; border-top: 1px solid var(--ink); padding: .75rem 0 1.5rem; font-family: var(--font-mono); font-size: .7rem; text-transform: uppercase; letter-spacing: .06em; color: var(--muted); }
@media (max-width: 1024px) {
  .ft-grid { grid-template-columns: 1fr 1fr; }
  .ft-poster { height: auto; grid-template-columns: repeat(3, minmax(0, 1fr)); }
  .ftp-word { grid-column: 1 / -1; min-height: 160px; }
  .ftp-circle, .ftp-tri, .ftp-top { grid-column: span 1; aspect-ratio: 1; }
  .ft-col, .ft-col + .ft-col { padding-left: 0; border-right: none; border-bottom: 1px solid var(--ink); }
  .ft-brand { grid-column: 1 / -1; }
}

/* section layouts (homepage only) */
.oss-grid, .hood-grid, .homelab-grid, .homelab-arch-grid { gap: 0 !important; }
#opensource [data-component="heatmap-wrap"] svg { max-width: 100%; height: auto; }
@media (max-width: 1024px) {
  .oss-grid, .hood-grid, .homelab-grid, .homelab-arch-grid, .skills-grid { grid-template-columns: 1fr !important; }
  .hood-grid > * { min-width: 0 !important; }
  .hood-grid > div:first-child { position: static !important; }
  .activity-metrics { grid-template-columns: repeat(2, minmax(0, 1fr)) !important; }
  .activity-metrics > * { grid-column: span 1 !important; }
  .activity-list-grid { grid-template-columns: 1fr !important; }
}
`

// homeCSS styles the homepage sections built on the swiss grid.
const homeCSS = `
.hero .hero-logo-slot { flex: 1; min-height: 240px; }
/* first screen: name spans two rows beside the block and the photo; the
   three facts share one short row below, so the hero ends at the fold */
.hero-head { justify-content: space-between; gap: 1.25rem; min-height: max(440px, calc(100svh - 3.3rem - 11.75rem)); }
.hero-status { color: var(--main-text); }
.hero-head { container-type: inline-size; }
.hero-name { margin: .5rem 0 0; font-size: min(15.5cqi, 9rem); font-weight: 900; letter-spacing: -0.045em; line-height: .9; text-transform: uppercase; overflow-wrap: anywhere; }
.hero-type { margin: 1rem 0 0; font-family: var(--font-mono); font-weight: 700; font-size: clamp(.95rem, .85rem + .5vw, 1.3rem); min-height: 1.5em; }
.hero-lead { margin: 0; font-size: clamp(1.05rem, .95rem + .6vw, 1.5rem); font-weight: 500; letter-spacing: -0.015em; line-height: 1.32; max-width: 40ch; }
.hero-ctas, .cat-links { display: flex; flex-wrap: wrap; gap: .4rem; }
.hero-ctas { margin-top: 1.25rem; }
.hero-side { display: flex; flex-direction: column; gap: 0; }
.hero-block { flex: 1; display: flex; flex-direction: column; gap: .6rem; padding: var(--swiss-pad); background: var(--block); color: var(--block-fg); }
.hero-block [data-component="swiss-kicker"] [data-slot="aside"] { color: inherit; }
.hero-credo { margin: 0; font-weight: 800; font-size: clamp(1.05rem, .95rem + .6vw, 1.5rem); line-height: 1.15; }
.hero-credo-sub { margin: 0; opacity: .85; font-size: .9rem; }
.hero-scroll { border-top: 1px solid currentColor; padding-top: .5rem; margin-top: auto; }
.hero-photo { position: relative; aspect-ratio: 1 / 1; overflow: hidden; border-top: 1px solid var(--ink); }
.hero-photo img { width: 100%; height: 100%; object-fit: cover; object-position: 50% 30%; display: block; position: absolute; inset: 0; }
.hero-photo-cap { position: absolute; left: 0; right: 0; bottom: 0; display: flex; justify-content: space-between; gap: .5rem; padding: .45rem .6rem; background: var(--bg); border-top: 1px solid var(--ink); }
.hero-fact { margin: 0; font-size: clamp(1.1rem, .95rem + .7vw, 1.6rem); font-weight: 800; letter-spacing: -0.02em; line-height: 1.15; }
.hero-fact-sub, .cell-text { margin: 0; color: var(--muted); font-size: .9rem; line-height: 1.45; }
.cell-title { margin: 0; font-size: clamp(1.1rem, 1rem + .5vw, 1.45rem); font-weight: 800; letter-spacing: -0.025em; line-height: 1.1; }
.cell-foot { margin-top: auto; }
.lab-cell { position: relative; }
.lab-cell:hover { background: var(--ink); color: var(--bg); }
.lab-cell:hover .cell-text { color: inherit; opacity: .75; }
.lab-link { color: inherit; text-decoration: none; }
.lab-link::after { content: ""; position: absolute; inset: 0; }
.lab-go { color: var(--main-text); }
.lab-cell:hover .lab-go { color: inherit; }
.paper-head { display: flex; gap: .6rem; align-items: center; }
.paper-logo { width: 36px; height: 36px; object-fit: contain; background: #fff; border: 1px solid var(--line); flex: none; }
.contact-list { margin: 0; padding-left: 1.1rem; display: grid; gap: .35rem; font-size: .92rem; }

.plate { display: grid; grid-template-columns: repeat(9, minmax(0, 1fr)); box-shadow: 0 1px 0 var(--ink); }
.plate-img { grid-column: span 4; box-shadow: 1px 0 0 var(--ink); display: flex; flex-direction: column; min-width: 0; }
.plate-img [data-component="gallery"] { flex: 1; }
.plate-img [data-component="gallery"] [data-slot="stage"] { flex: 1; }
.plate-body { grid-column: span 5; padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .75rem; }
.plate-body p { margin: 0; line-height: 1.5; }
.plate-title { margin: .25rem 0 0; font-size: clamp(1.65rem, 1.35rem + 1.6vw, 3rem); font-weight: 800; letter-spacing: -0.035em; line-height: 1.02; }
.cat-head { padding-block: .75rem; }
.mini { min-height: 100%; }
.mini-body { padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .6rem; flex: 1; }
.mini [data-component="gallery"] { border-bottom: 1px solid var(--ink); }
.mono-field { display: flex; flex-direction: column; justify-content: space-between; aspect-ratio: 16 / 10; padding: var(--swiss-pad); border-bottom: 1px solid var(--ink); }
.plate-img .mono-field { aspect-ratio: auto; height: 100%; border-bottom: none; }
.mono-field > span:first-child { font-size: clamp(3rem, 2rem + 4vw, 6rem); font-weight: 900; letter-spacing: -0.05em; line-height: .85; }
.mono-block { background: var(--block); color: var(--block-fg); }
.mono-second { background: var(--second); color: var(--second-fg); }
.mono-third { background: var(--third); color: var(--third-fg); }
.mono-ink { background: var(--ink); color: var(--bg); }
.poster-field { display: flex; flex-direction: column; justify-content: space-between; gap: 1rem; min-height: clamp(240px, 26vw, 400px); padding: var(--swiss-pad); }
.poster-top, .poster-org { display: flex; justify-content: space-between; align-items: center; gap: .5rem; }
.poster-top { border-bottom: 1px solid currentColor; padding-bottom: .45rem; }
.poster-mark { font-size: clamp(3.5rem, 2rem + 6vw, 8.5rem); font-weight: 900; letter-spacing: -0.06em; line-height: .8; }
.poster-sub { display: block; margin-top: .6rem; font-family: var(--font-mono); font-size: .72rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; }
.paper-stack { display: flex; flex-direction: column; }
.paper + .paper { border-top: 3px solid var(--ink); }
.paper-meta { color: var(--main-text); }
.paper { flex: 1; padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .5rem; background: var(--surface-2); background-image: repeating-linear-gradient(to bottom, transparent 0 23px, var(--rule-subtle) 23px 24px); }
.paper-mark { font-size: 2.6rem; font-weight: 900; letter-spacing: -0.05em; line-height: .85; color: var(--main-text); }
.paper-title { margin: 0; font-size: 1rem; line-height: 1.2; letter-spacing: -0.01em; }
.paper .cell-text { font-size: .82rem; }
.sec-inv > [data-slot="body"] { background: var(--board); color: var(--tile-fg); }
.sec-inv [data-component="swiss-btn"]:not([data-variant="primary"]) { color: inherit; border-color: currentColor; background: transparent; }
.sec-inv [data-component="swiss-btn"]:hover { background: var(--tile-fg); color: var(--board); }
.poster-org img { width: 44px; height: 44px; object-fit: contain; background: #fff; padding: 3px; }
.poster-body { padding: var(--swiss-pad); display: flex; flex-direction: column; gap: .6rem; flex: 1; border-top: 1px solid var(--ink); }

@media (max-width: 1024px) {
  .hero [data-slot="anchor"] { display: grid; grid-template-columns: 1fr auto; column-gap: 1rem; align-items: center; }
  .hero .hero-logo-slot { grid-column: 2; grid-row: 1 / span 2; width: 96px; height: 96px; min-height: 0; flex: none; }
  .hero [data-slot="note"] { grid-column: 1; text-align: left; }
  .hero-head { min-height: 0; }
  .plate-img, .plate-body { grid-column: 1 / -1; box-shadow: none; }
  .plate-body { padding-inline: 0; }
  .plate-img { border-bottom: 1px solid var(--ink); }
}
`
