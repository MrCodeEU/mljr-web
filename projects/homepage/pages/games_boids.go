package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

func Boids(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Boids - Michael Reinegger",
			Description: "A flocking simulation built from Craig Reynolds' three classic boid rules: separation, alignment, and cohesion.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS+gameCSS+boidsCSS+gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("Boids")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.boids.lede"))),
				h.Div(h.ID("boids-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("control-row"),
					h.Label(h.For("boids-separation"), g.Text(i18n.T(lang, "games.boids.label_separation"))),
					h.Input(h.Type("range"), h.ID("boids-separation"), h.Min("0"), h.Max("3"), h.Step("0.1"), h.Value("1.5")),
					h.Label(h.For("boids-alignment"), g.Text(i18n.T(lang, "games.boids.label_alignment"))),
					h.Input(h.Type("range"), h.ID("boids-alignment"), h.Min("0"), h.Max("3"), h.Step("0.1"), h.Value("1.0")),
					h.Label(h.For("boids-cohesion"), g.Text(i18n.T(lang, "games.boids.label_cohesion"))),
					h.Input(h.Type("range"), h.ID("boids-cohesion"), h.Min("0"), h.Max("3"), h.Step("0.1"), h.Value("1.0")),
					h.Label(h.For("boids-count"), g.Text(i18n.T(lang, "games.boids.label_count"))),
					h.Input(h.Type("number"), h.ID("boids-count"), h.Min("5"), h.Max("300"), h.Value("80")),
				),
			),
		),
		boidsDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/boids/sketch.js")),
	)
}

const boidsCSS = `
.control-row {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  flex-wrap: wrap;
  font-size: var(--t-sm);
  margin-top: var(--sp-4);
}
.control-row input[type="range"] {
  flex: 1 1 8rem;
  min-width: 6rem;
}
.control-row input[type="number"] {
  font: inherit;
  box-sizing: border-box;
  width: 4.5rem;
  text-align: center;
  padding: var(--sp-1) var(--sp-2);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
}
`

func boidsDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.boids.desc_title",
		[]string{"games.boids.desc_p1", "games.boids.desc_p2", "games.boids.desc_p3"},
		boidsDiagram(), "")
}

func boidsDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 200"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.El("circle", g.Attr("cx", "130"), g.Attr("cy", "100"), g.Attr("r", "50"), g.Attr("style", "stroke:var(--muted);stroke-width:1;stroke-dasharray:3,3;fill:none")),
		g.El("polygon", g.Attr("points", "130,90 124,102 136,102"), g.Attr("style", "fill:#e56b6f")),
		g.El("polygon", g.Attr("points", "105,70 100,80 110,80"), g.Attr("style", "fill:var(--muted)")),
		g.El("polygon", g.Attr("points", "160,75 155,85 165,85"), g.Attr("style", "fill:var(--muted)")),
		g.El("polygon", g.Attr("points", "100,120 95,130 105,130"), g.Attr("style", "fill:var(--muted)")),
		g.El("polygon", g.Attr("points", "155,130 150,140 160,140"), g.Attr("style", "fill:var(--muted)")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "20"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("each boid only reacts to neighbours inside its own radius")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "190"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("separation + alignment + cohesion → flock")),
	)
}
