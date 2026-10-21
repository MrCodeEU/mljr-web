package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

// AISnakeStreamPath is the SSE endpoint feeding the shared AI snake game.
const AISnakeStreamPath = "/api/snake-ai/stream"

// AISnakeLive is the shared, server-side AI snake: one game, every visitor watches the same state.
func AISnakeLive(lang string, a AnalyticsConfig) g.Node {
	t := func(key string) string { return i18n.T(lang, key) }
	stat := func(labelKey, expr string) g.Node {
		return h.Div(h.Class("snake-live-stat"),
			h.Span(h.Class("snake-live-label"), g.Text(t(labelKey))),
			h.Strong(ui.Text(expr)),
		)
	}
	status := func(value, key string) g.Node {
		return h.Span(h.Class("snake-live-status"), g.Attr("data-status", value),
			ui.Show("$snake.status === '"+value+"'"), g.Text(t(key)))
	}
	return layout.PageShell(
		layout.PageProps{
			Title:       "AI Snake Live - Michael Reinegger",
			Description: "One snake game running on the server, played by a tiny neural network, watched live by everyone on the page.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS+gameCSS+aiSnakeLiveCSS+gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(t("legal.back"))),
				h.H1(g.Text("AI Snake Live")),
				h.P(h.Class("game-lede"), g.Text(t("games.ai_snake_live.lede"))),
				h.Div(h.Class("snake-live"),
					ui.Signals(`{snake:{}}`),
					g.Attr("data-init", "@get('"+AISnakeStreamPath+"')"),
					g.Attr("data-effect", "window.snakeLiveRender && window.snakeLiveRender($snake)"),
					h.Div(h.Class("snake-live-stats"),
						stat("games.ai_snake_live.score", "$snake.score"),
						stat("games.ai_snake_live.best", "$snake.best"),
						stat("games.ai_snake_live.games", "$snake.games"),
						stat("games.ai_snake_live.viewers", "$snake.viewers"),
						status("playing", "games.ai_snake_live.status_playing"),
						status("dead", "games.ai_snake_live.status_dead"),
						status("complete", "games.ai_snake_live.status_complete"),
					),
					h.Canvas(h.ID("ai-snake-live-canvas")),
				),
				h.P(h.Class("snake-live-note"), g.Text(t("games.ai_snake_live.note"))),
				h.P(h.A(h.Href("/games/ai-snake"), g.Text(t("games.ai_snake_live.classic_link")))),
			),
		),
		aiSnakeLiveDescriptionSection(lang),
		siteFooter(lang),
		// Loaded before the stream's first frame; defined as a global the Datastar effect calls.
		h.Script(h.Src("/static/games/ai-snake-live/renderer.js")),
	)
}

const aiSnakeLiveCSS = `
.snake-live {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-4);
  margin-bottom: var(--sp-4);
}
.snake-live-stats {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  justify-content: center;
  gap: var(--sp-2) var(--sp-6);
  width: min(100%, 800px);
}
.snake-live-label {
  color: var(--muted);
  font-size: .8rem;
  text-transform: uppercase;
  letter-spacing: .08em;
  margin-right: var(--sp-2);
}
.snake-live-status {
  margin-left: auto;
  font-weight: 800;
}
.snake-live-status[data-status="playing"] { color: var(--accent); }
.snake-live canvas {
  width: min(100%, 800px);
  aspect-ratio: 1;
  display: block;
  background: var(--surface);
  border: var(--bw-2) solid var(--line);
  image-rendering: pixelated;
}
.snake-live-note { color: var(--muted); }
`

func aiSnakeLiveDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.ai_snake_live.desc_title",
		[]string{"games.ai_snake_live.desc_p1", "games.ai_snake_live.desc_p2", "games.ai_snake_live.desc_p3"},
		aiSnakeDiagram(), "")
}
