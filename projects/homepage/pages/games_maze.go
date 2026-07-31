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

func Maze(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Maze - Michael Reinegger",
			Description: "A maze that carves itself out with a randomized depth-first search, then solves itself with breadth-first search.",
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
				h.H1(g.Text("Maze")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.maze.lede"))),
				h.Div(h.ID("maze-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("control-row"),
					h.Button(h.ID("maze-regenerate"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.maze.regenerate"))),
					h.Label(h.For("maze-speed"), g.Text(i18n.T(lang, "games.maze.label_speed"))),
					h.Input(h.Type("range"), h.ID("maze-speed"), h.Min("1"), h.Max("50"), h.Value("3")),
				),
			),
		),
		mazeDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/maze/sketch.js")),
	)
}

func mazeDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.maze.desc_title",
		[]string{"games.maze.desc_p1", "games.maze.desc_p2", "games.maze.desc_p3"},
		mazeDiagram(), "")
}

func mazeDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 220 180"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		gridLines(220, 160, 20),
		g.El("path", g.Attr("d", "M10 10 L10 90 L50 90 L50 50 L90 50 L90 130 L150 130 L150 70 L190 70"),
			g.Attr("style", "stroke:#78b4ff;stroke-width:3;fill:none;stroke-linecap:round;stroke-linejoin:round")),
		g.El("circle", g.Attr("cx", "10"), g.Attr("cy", "10"), g.Attr("r", "5"), g.Attr("style", "fill:#5ac878")),
		g.El("circle", g.Attr("cx", "190"), g.Attr("cy", "70"), g.Attr("r", "5"), g.Attr("style", "fill:#e6505a")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "175"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("BFS frontier expands ring by ring, guaranteeing the shortest path")),
	)
}
