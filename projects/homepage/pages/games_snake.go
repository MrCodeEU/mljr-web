package pages

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

func Snake(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Snake - Michael Reinegger",
			Description: "A classic keyboard-controlled snake game ported from an old side project.",
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
				h.H1(g.Text("Snake")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.snake.lede"))),
				h.Div(h.ID("snake-game"), h.Class("game-canvas-wrap")),
			),
		),
		snakeDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/snake/sketch.js"), h.Defer()),
	)
}

func snakeDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.snake.desc_title",
		[]string{"games.snake.desc_p1", "games.snake.desc_p2", "games.snake.desc_p3"},
		snakeDiagram(), "")
}

func snakeDiagram() g.Node {
	cells := []struct{ x, y int }{{2, 3}, {3, 3}, {4, 3}, {5, 3}}
	var body []g.Node
	for i, c := range cells {
		fillColor := "var(--muted)"
		if i == len(cells)-1 {
			fillColor = "#78c878"
		}
		body = append(body, g.El("rect",
			g.Attr("x", strconv.Itoa(20+c.x*24)), g.Attr("y", strconv.Itoa(20+c.y*24)),
			g.Attr("width", "22"), g.Attr("height", "22"),
			g.Attr("style", "fill:"+fillColor)))
	}
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		gridLines(260, 220, 24),
		g.Group(body),
		g.El("rect", g.Attr("x", "212"), g.Attr("y", "116"), g.Attr("width", "22"), g.Attr("height", "22"), g.Attr("style", "fill:#e56b6f")),
		g.El("text", g.Attr("x", "20"), g.Attr("y", "200"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("head highlighted, apple in red")),
	)
}

func gridLines(w, h, step int) g.Node {
	var lines []g.Node
	for x := 0; x <= w; x += step {
		lines = append(lines, g.El("line",
			g.Attr("x1", strconv.Itoa(x)), g.Attr("y1", "0"), g.Attr("x2", strconv.Itoa(x)), g.Attr("y2", strconv.Itoa(h)),
			g.Attr("style", "stroke:var(--line);stroke-width:0.5;opacity:.4")))
	}
	for y := 0; y <= h; y += step {
		lines = append(lines, g.El("line",
			g.Attr("x1", "0"), g.Attr("y1", strconv.Itoa(y)), g.Attr("x2", strconv.Itoa(w)), g.Attr("y2", strconv.Itoa(y)),
			g.Attr("style", "stroke:var(--line);stroke-width:0.5;opacity:.4")))
	}
	return g.Group(lines)
}
