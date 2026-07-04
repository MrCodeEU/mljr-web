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

func AISnake(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "AI Snake - Michael Reinegger",
			Description: "A self-playing snake that follows a fixed Hamiltonian cycle so it can never trap itself, ported from an old side project.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS + gameCSS + boidsCSS + aiSnakeCSS + gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("AI Snake")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.ai_snake.lede"))),
				h.Div(h.ID("ai-snake-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("control-row"),
					h.Label(h.For("ai-snake-strategy"), g.Text(i18n.T(lang, "games.ai_snake.label_strategy"))),
					h.Select(h.ID("ai-snake-strategy"),
						h.Option(h.Value("cycle_pure"), h.Selected(), g.Text(i18n.T(lang, "games.ai_snake.strategy_cycle_pure"))),
						h.Option(h.Value("greedy_safe"), g.Text(i18n.T(lang, "games.ai_snake.strategy_greedy_safe"))),
						h.Option(h.Value("cycle_shortcuts"), g.Text(i18n.T(lang, "games.ai_snake.strategy_cycle_shortcuts"))),
						h.Option(h.Value("greedy"), g.Text(i18n.T(lang, "games.ai_snake.strategy_greedy"))),
					),
					h.Label(h.For("ai-snake-speed"), g.Text(i18n.T(lang, "games.ai_snake.label_speed"))),
					h.Input(h.Type("range"), h.ID("ai-snake-speed"), h.Min("1"), h.Max("20000"), h.Value("300")),
				),
			),
		),
		aiSnakeDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/ai-snake/sketch.js")),
	)
}

const aiSnakeCSS = `
.control-row select {
  font: inherit;
  box-sizing: border-box;
  padding: var(--sp-1) var(--sp-2);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
}
`

func aiSnakeDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.ai_snake.desc_title",
		[]string{"games.ai_snake.desc_p1", "games.ai_snake.desc_p2", "games.ai_snake.desc_p3", "games.ai_snake.desc_p4"},
		aiSnakeDiagram(), "")
}

func aiSnakeDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		gridLines(260, 220, 24),
		g.El("path", g.Attr("d", "M20 20 L212 20 L212 68 L164 68 L164 44 L116 44 L116 68 L68 68 L68 164 L212 164 L212 116 L20 116 Z"),
			g.Attr("style", "stroke:#78b4ff;stroke-width:2.5;fill:none;stroke-linejoin:round")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "200"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("one fixed loop through every cell, always followed forward")),
	)
}
