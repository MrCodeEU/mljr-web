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
			Description: "A self-playing snake that pathfinds to food while avoiding trapping itself, ported from an old side project.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS + gameCSS + gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("AI Snake")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.ai_snake.lede"))),
				h.Div(h.ID("ai-snake-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("game-controls"),
					h.Label(h.For("ai-snake-speed"), g.Text(i18n.T(lang, "games.ai_snake.label_speed"))),
					h.Input(h.Type("range"), h.ID("ai-snake-speed"), h.Min("1"), h.Max("300"), h.Value("60")),
				),
			),
		),
		aiSnakeDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/ai-snake/sketch.js")),
	)
}

func aiSnakeDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.ai_snake.desc_title",
		[]string{"games.ai_snake.desc_p1", "games.ai_snake.desc_p2", "games.ai_snake.desc_p3"},
		aiSnakeDiagram(), "games.ai_snake.desc_diagram_caption")
}

func aiSnakeDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		gridLines(260, 220, 24),
		// head
		g.El("rect", g.Attr("x", "20"), g.Attr("y", "92"), g.Attr("width", "22"), g.Attr("height", "22"), g.Attr("style", "fill:#78c878")),
		// tail
		g.El("rect", g.Attr("x", "20"), g.Attr("y", "20"), g.Attr("width", "22"), g.Attr("height", "22"), g.Attr("style", "fill:var(--muted)")),
		// food
		g.El("rect", g.Attr("x", "212"), g.Attr("y", "92"), g.Attr("width", "22"), g.Attr("height", "22"), g.Attr("style", "fill:#e56b6f")),
		// path to food (green dashed)
		g.El("path", g.Attr("d", "M31 103 L223 103"), g.Attr("style", "stroke:#78c878;stroke-width:2;stroke-dasharray:5,4;fill:none")),
		// path to tail (blue dashed)
		g.El("path", g.Attr("d", "M31 103 L31 31"), g.Attr("style", "stroke:#78b4ff;stroke-width:2;stroke-dasharray:5,4;fill:none")),
		g.El("text", g.Attr("x", "60"), g.Attr("y", "95"), g.Attr("style", "font-size:11px;fill:#78c878"), g.Text("path to food")),
		g.El("text", g.Attr("x", "36"), g.Attr("y", "65"), g.Attr("style", "font-size:11px;fill:#78b4ff"), g.Text("path to tail")),
	)
}
