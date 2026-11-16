package pages

import (
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

func Minenfeld(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Minenfeld - Michael Reinegger",
			Description: "A minefield-crossing game ported from an old side project.",
			Theme:       token.ThemeSwiss,
			AutoMode:    true,
			LockTheme:   true,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS+gameCSS+gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwiss, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.ID("main"), h.Class("game-page sub-page"),
			subHead("/#playground", tr(lang, "Labor", "Lab"), 8,
				h.H1(g.Text("Minenfeld")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.minenfeld.lede"))),
			),
			h.Div(h.Class("game-shell sub-body"),
				h.Div(h.ID("minenfeld-game"),
					h.Canvas(h.ID("minenfeld-canvas"), h.Width("640"), h.Height("640")),
					h.Div(h.ID("minenfeld-input-wrap"), h.Class("game-input-wrap"),
						h.P(h.ID("minenfeld-input-label"), g.Text("Wo wollen Sie starten? [0/19]")),
						h.Input(h.Type("text"), h.ID("minenfeld-input-box"), g.Attr("aria-labelledby", "minenfeld-input-label")),
						h.Button(h.ID("minenfeld-input-btn"), h.Class("game-btn"), g.Text("Go!")),
					),
					h.P(h.ID("minenfeld-error"), h.Class("game-error"), h.Style("display:none")),
					h.P(h.ID("minenfeld-message"), h.Class("game-message"), h.Style("display:none")),
					h.Button(h.ID("minenfeld-restart"), h.Class("game-restart game-btn"), h.Style("display:none"), g.Text(i18n.T(lang, "games.minenfeld.play_again"))),
				),
			),
		),
		minenfeldDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/minenfeld/game.js"), h.Defer()),
	)
}

func minenfeldDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.minenfeld.desc_title",
		[]string{"games.minenfeld.desc_p1", "games.minenfeld.desc_p2", "games.minenfeld.desc_p3"},
		minenfeldDiagram(), "")
}

func minenfeldDiagram() g.Node {
	minePositions := []int{3, 1, 4, 2, 0}
	pathCols := []int{2, 2, 3, 2, 1}
	var mines []g.Node
	for row, col := range minePositions {
		mines = append(mines, g.El("rect",
			g.Attr("x", strconv.Itoa(20+col*40)), g.Attr("y", strconv.Itoa(20+row*36)),
			g.Attr("width", "36"), g.Attr("height", "32"),
			g.Attr("style", "fill:#e56b6f;opacity:.5")))
	}
	var pathPoints []string
	for row, col := range pathCols {
		pathPoints = append(pathPoints, strconv.Itoa(20+col*40+18)+","+strconv.Itoa(20+row*36+16))
	}
	polyline := strings.Join(pathPoints, " ")
	return g.El("svg",
		g.Attr("viewBox", "0 0 220 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		gridLines(220, 200, 40),
		g.Group(mines),
		g.El("polyline", g.Attr("points", polyline), g.Attr("style", "stroke:#78b4ff;stroke-width:2.5;fill:none")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "215"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("one mine per row, chosen path in blue")),
	)
}

const gameCSS = `
/* embedded in the homepage lab: game only, no page chrome */
html.embedded body > :not(main):not(script):not(style) { display: none !important; }
html.embedded .game-page { min-height: 0; }
.game-page { min-height: 70vh; }
#minenfeld-canvas {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 0 auto var(--sp-4);
  border: var(--bw-2) solid var(--line);
}
.game-canvas-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: var(--sp-4);
}
.game-canvas-wrap canvas {
  max-width: 100%;
  height: auto;
  border: var(--bw-2) solid var(--line);
}
.game-controls {
  display: flex;
  justify-content: center;
}
.game-input-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--sp-3);
  flex-wrap: wrap;
  margin-bottom: var(--sp-3);
}
.game-error {
  color: var(--danger, #e5484d);
  font-weight: 700;
  text-align: center;
}
.game-message {
  font-weight: 700;
  text-align: center;
}
.game-btn {
  font: inherit;
  font-weight: 900;
  padding: var(--sp-2) var(--sp-4);
  background: var(--ink);
  color: var(--surface);
  border: var(--bw-2) solid var(--ink);
  cursor: pointer;
}
.game-btn:hover {
  opacity: .85;
}
.game-input-wrap input[type="text"] {
  font: inherit;
  padding: var(--sp-2) var(--sp-3);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
  width: 5ch;
}
.game-restart {
  display: block;
  margin: var(--sp-3) auto 0;
}
`
