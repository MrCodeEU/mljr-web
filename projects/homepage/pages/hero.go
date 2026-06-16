package pages

import (
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func cvPath(lang string) string {
	if lang == "de" {
		return "/static/cv/lebenslauf.pdf"
	}
	return "/static/cv/lebenslauf_eng.pdf"
}

func heroSection(d hpdata.SiteData, lang string, projectCount int) g.Node {
	li := d.LinkedIn
	hero := d.ContentFor(lang).Hero
	return h.Section(
		h.ID("hero"),
		h.Style("min-height:90vh;display:flex;align-items:center;padding:var(--sp-12) 0 var(--sp-8)"),
		layout.Container(layout.ContainerProps{},
			h.Div(
				h.Class("hero-grid"),
				h.Style("display:grid;grid-template-columns:minmax(0,1fr) minmax(0,1fr);gap:var(--sp-10);align-items:stretch"),
				// Left column
				h.Div(
					h.ID("hero-content"),
					h.Style("display:flex;flex-direction:column;justify-content:center;gap:var(--sp-5);position:relative;z-index:2"),
					// Availability tag
					h.Div(
						h.Style("display:inline-flex;align-items:center;gap:var(--sp-2);padding:var(--sp-1) var(--sp-3);border:var(--bw-2) solid var(--ink);border-radius:calc(var(--radius)*2);font-size:var(--t-xs);font-weight:700;width:fit-content;background:var(--surface)"),
						h.Span(h.Style("width:8px;height:8px;border-radius:50%;background:#22c55e;flex-shrink:0;animation:pulse-dot 2s ease infinite")),
						g.Text(hero.StatusTag),
					),
					// Main headline
					h.H1(
						h.Style("font-size:clamp(2.5rem,6vw,4rem);font-weight:900;line-height:1.05;margin:0"),
						g.Text(i18n.T(lang, "hero.greeting") + " "),
						primitive.GradientText(primitive.GradientTextProps{
							From:  "var(--accent)",
							To:    "var(--ink)",
							Angle: "135deg",
						}, g.Text(i18n.T(lang, "hero.name"))),
					),
					// Typewriter tagline
					h.P(
						h.Style("font-size:clamp(var(--t-lg),3vw,var(--t-xl));font-weight:700;margin:0;line-height:1.3"),
						g.Text("I build "),
						primitive.Typewriter(primitive.TypewriterProps{
							Lines:       hero.TaglineLines,
							Speed:       55,
							DeleteSpeed: 28,
							Pause:       2200,
							ID:          "hero-tw",
						}),
					),
					// Description
					h.P(
						h.Style("font-size:var(--t-base);color:var(--muted);max-width:46ch;margin:0;line-height:1.6"),
						g.Text(hero.Description),
					),
					// CTA buttons
					h.Div(
						h.Class("hero-ctas"),
						h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
						h.A(h.Href("#projects"),
							primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeLG},
								g.Text(i18n.T(lang, "hero.cta_projects")),
								icon.Icon("lucide:arrow-right", icon.Props{Size: "1.1rem"}),
							),
						),
						h.A(h.Href("#contact"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeLG},
								icon.Icon("lucide:mail", icon.Props{Size: "1.1rem"}),
								g.Text(i18n.T(lang, "hero.cta_contact")),
							),
						),
						h.A(
							h.Href("https://github.com/MrCodeEU"),
							g.Attr("target", "_blank"),
							g.Attr("rel", "noopener"),
							primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeLG},
								icon.Icon("simple-icons:github", icon.Props{Size: "1.1rem"}),
								g.Text(i18n.T(lang, "hero.cta_github")),
							),
						),
						h.A(
							h.Href(cvPath(lang)),
							g.Attr("download", ""),
							g.Attr("target", "_blank"),
							g.Attr("rel", "noopener"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeLG},
								icon.Icon("lucide:download", icon.Props{Size: "1.1rem"}),
								g.Text(i18n.T(lang, "hero.cta_download_cv")),
							),
						),
					),
				),
				// Right column: Bento Grid
				h.Div(
					h.Class("hero-bento"),
					h.Style("position:relative;z-index:2"),
					heroBento(li, hero, lang, projectCount),
				),
			),
		),
	)
}

func heroBento(li hpdata.LinkedInData, hero hpdata.HeroContent, lang string, projectCount int) g.Node {
	bento := hero.Bento
	// Neo-brutalist bento: context cards first, simple counters second.
	card := func(tone token.Tone, children ...g.Node) g.Node {
		return primitive.Card(primitive.CardProps{Tone: tone, Attrs: []g.Node{h.Style("height:100%")}}, children...)
	}

	kickerStyle := "font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.12em;opacity:.72"
	numStyle := "font-size:clamp(2.2rem,4vw,3.6rem);font-weight:950;line-height:1;font-variant-numeric:tabular-nums;margin-bottom:var(--sp-2)"
	labelStyle := "font-size:var(--t-xs);font-weight:800;text-transform:uppercase;letter-spacing:.12em;opacity:.7"

	photoCell := h.Div(
		g.Attr("data-component", "bento-item"),
		h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		h.Div(
			h.Class("bento-photo"),
			h.Style("height:100%;min-height:160px;aspect-ratio:4/3;width:100%;overflow:hidden;border-radius:var(--radius);border:var(--bw-2) solid var(--ink);position:relative"),
			h.Img(
				h.Src(li.Profile.PhotoURL),
				h.Alt(i18n.T(lang, "hero.photo_name")),
				g.Attr("loading", "eager"),
				h.Style("width:100%;height:100%;object-fit:cover;object-position:center 28%"),
			),
			h.Div(
				h.Style("position:absolute;bottom:0;left:0;right:0;padding:var(--sp-3) var(--sp-4);background:linear-gradient(transparent,rgba(0,0,0,.75));color:white"),
				h.Div(h.Style("font-weight:900;font-size:var(--t-md)"), g.Text(i18n.T(lang, "hero.photo_name"))),
				h.Div(h.Style("font-size:var(--t-xs);opacity:.85"), g.Text(i18n.T(lang, "hero.photo_location"))),
			),
		),
	)

	focusCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		card(token.ToneLime,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);height:100%"),
				icon.Icon("lucide:building-2", icon.Props{Size: "1.8rem"}),
				h.Div(h.Style(kickerStyle+";margin-top:var(--sp-1)"), g.Text(bento.Focus.Label)),
				h.Div(h.Style("font-size:clamp(1.2rem,2vw,1.8rem);font-weight:950;line-height:1.1"), g.Text(bento.Focus.Value)),
				h.Div(h.Style(labelStyle+";margin-top:auto"), g.Text(bento.Focus.Sub)),
			),
		),
	)
	statusCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		card(token.ToneMint,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);height:100%"),
				icon.Icon("lucide:map-pin", icon.Props{Size: "1.6rem"}),
				h.Div(h.Style(kickerStyle+";margin-top:var(--sp-1)"), g.Text(bento.Status.Label)),
				h.Div(h.Style("font-size:var(--t-lg);font-weight:950;line-height:1.15"), g.Text(bento.Status.Value)),
				h.Div(h.Style(labelStyle+";margin-top:auto"), g.Text(bento.Status.Sub)),
			),
		),
	)
	mscCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		card(token.ToneCyan,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);height:100%"),
				icon.Icon("lucide:graduation-cap", icon.Props{Size: "1.6rem"}),
				h.Div(h.Style(kickerStyle+";margin-top:var(--sp-1)"), g.Text(bento.Education.Label)),
				h.Div(h.Style("font-size:var(--t-base);font-weight:900;line-height:1.25"), g.Text(bento.Education.Value)),
				h.Div(h.Style(labelStyle+";margin-top:auto"), g.Text(bento.Education.Sub)),
			),
		),
	)
	infraCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 2;grid-row:span 1;min-width:0"),
		card(token.ToneSky,
			h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-4)"),
				icon.Icon("lucide:server", icon.Props{Size: "2rem"}),
				h.Div(
					h.Div(h.Style(kickerStyle), g.Text(bento.Homelab.Label)),
					h.Div(h.Style("font-size:var(--t-base);font-weight:900;line-height:1.25;margin-top:var(--sp-1)"), g.Text(bento.Homelab.Value)),
					h.Div(h.Style(labelStyle+";margin-top:var(--sp-2)"), g.Text(bento.Homelab.Sub)),
				),
			),
		),
	)
	projectsCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		card(token.ToneViolet,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);height:100%"),
				icon.Icon("lucide:folder-git-2", icon.Props{Size: "1.4rem"}),
				h.Div(h.Style(numStyle+";margin-top:auto"),
					primitive.NumberTicker(primitive.NumberTickerProps{Value: float64(projectCount), From: float64(projectCount), TriggerOnView: true, ID: "nt-proj", Duration: 3200}),
				),
				h.Div(h.Style(labelStyle), g.Text(i18n.T(lang, "hero.projects_label"))),
			),
		),
	)
	yearsCell := h.Div(g.Attr("data-component", "bento-item"), h.Style("grid-column:span 1;grid-row:span 1;min-width:0"),
		card(token.ToneYellow,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2);height:100%"),
				icon.Icon("lucide:code-2", icon.Props{Size: "1.4rem"}),
				h.Div(h.Style(numStyle+";margin-top:auto"),
					primitive.NumberTicker(primitive.NumberTickerProps{Value: float64(time.Now().Year() - 2015), From: float64(time.Now().Year() - 2015), Suffix: "+", TriggerOnView: true, ID: "nt-yrs", Duration: 2800}),
				),
				h.Div(h.Style(labelStyle), g.Text(i18n.T(lang, "hero.years_label"))),
			),
		),
	)

	return h.Div(
		g.Attr("data-component", "bento-grid"),
		h.Style("display:grid;grid-template-columns:repeat(2,minmax(0,1fr));grid-auto-rows:minmax(150px,auto);gap:var(--sp-3);height:100%"),
		photoCell, focusCell, statusCell, mscCell, infraCell, projectsCell, yearsCell,
	)
}
