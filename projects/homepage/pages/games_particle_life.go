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

func ParticleLife(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Particle Life - Michael Reinegger",
			Description: "Colored particles that attract and repel by simple pairwise rules, organizing into life-like emergent clusters.",
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
				h.H1(g.Text("Particle Life")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.particle_life.lede"))),
				h.Div(h.ID("particle-life-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("control-row"),
					h.Label(h.For("particle-life-count"), g.Text(i18n.T(lang, "games.particle_life.label_count"))),
					h.Input(h.Type("number"), h.ID("particle-life-count"), h.Min("20"), h.Max("500"), h.Value("250")),
					h.Label(h.For("particle-life-force"), g.Text(i18n.T(lang, "games.particle_life.label_force"))),
					h.Input(h.Type("range"), h.ID("particle-life-force"), h.Min("0.1"), h.Max("2"), h.Step("0.05"), h.Value("0.6")),
					h.Label(h.For("particle-life-radius"), g.Text(i18n.T(lang, "games.particle_life.label_radius"))),
					h.Input(h.Type("range"), h.ID("particle-life-radius"), h.Min("30"), h.Max("160"), h.Value("80")),
				),
				h.Div(h.Class("control-row"),
					h.Button(h.ID("particle-life-randomize"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.particle_life.randomize"))),
					h.Button(h.ID("particle-life-reset"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.particle_life.reset_positions"))),
				),
			),
		),
		particleLifeDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/particle-life/sketch.js")),
	)
}

func particleLifeDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.particle_life.desc_title",
		[]string{"games.particle_life.desc_p1", "games.particle_life.desc_p2", "games.particle_life.desc_p3"},
		particleLifeDiagram(), "")
}

func particleLifeDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 200"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.El("circle", g.Attr("cx", "70"), g.Attr("cy", "90"), g.Attr("r", "6"), g.Attr("style", "fill:#e6505a")),
		g.El("circle", g.Attr("cx", "95"), g.Attr("cy", "70"), g.Attr("r", "6"), g.Attr("style", "fill:#5ac878")),
		g.El("circle", g.Attr("cx", "110"), g.Attr("cy", "105"), g.Attr("r", "6"), g.Attr("style", "fill:#5a96e6")),
		g.El("circle", g.Attr("cx", "60"), g.Attr("cy", "120"), g.Attr("r", "6"), g.Attr("style", "fill:#e6c85a")),
		g.El("line", g.Attr("x1", "70"), g.Attr("y1", "90"), g.Attr("x2", "95"), g.Attr("y2", "70"), g.Attr("style", "stroke:#5ac878;stroke-width:1.5;stroke-dasharray:2,2")),
		g.El("line", g.Attr("x1", "70"), g.Attr("y1", "90"), g.Attr("x2", "110"), g.Attr("y2", "105"), g.Attr("style", "stroke:#e6505a;stroke-width:1.5;stroke-dasharray:2,2")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "20"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("one attract/repel value per color pair")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "190"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("...applied to every particle, every frame")),
	)
}
