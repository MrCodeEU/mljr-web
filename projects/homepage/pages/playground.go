package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

type playgroundItem struct {
	Label   string
	DescKey string
	Href    string
	Icon    string
}

var playgroundItems = []playgroundItem{
	{Label: "AI Snake", DescKey: "games.ai_snake.card_desc", Href: "/games/ai-snake", Icon: "lucide:bot"},
	{Label: "Snake", DescKey: "games.snake.card_desc", Href: "/games/snake", Icon: "lucide:gamepad-2"},
	{Label: "Minenfeld", DescKey: "games.minenfeld.card_desc", Href: "/games/minenfeld", Icon: "lucide:bomb"},
	{Label: "Double Pendulum", DescKey: "games.pendulum.card_desc", Href: "/games/double-pendulum", Icon: "lucide:orbit"},
	{Label: "Fourier Drawing", DescKey: "games.fourier.card_desc", Href: "/games/fourier", Icon: "lucide:audio-waveform"},
	{Label: "Boids", DescKey: "games.boids.card_desc", Href: "/games/boids", Icon: "lucide:bird"},
	{Label: "Particle Life", DescKey: "games.particle_life.card_desc", Href: "/games/particle-life", Icon: "lucide:atom"},
	{Label: "Maze", DescKey: "games.maze.card_desc", Href: "/games/maze", Icon: "lucide:route"},
	{Label: "Sorting Visualizer", DescKey: "games.sorting.card_desc", Href: "/games/sorting", Icon: "lucide:bar-chart-3"},
}

// PlaygroundURLPaths returns the route paths of all playground items, for sitemap generation.
func PlaygroundURLPaths() []string {
	paths := make([]string, len(playgroundItems))
	for i, it := range playgroundItems {
		paths[i] = it.Href
	}
	return paths
}

func playgroundSection(num, lang string) g.Node {
	tones := []token.Tone{token.ToneViolet, token.ToneCyan, token.ToneLime, token.TonePink}
	cards := make([]g.Node, 0, len(playgroundItems))
	for i, it := range playgroundItems {
		cards = append(cards, playgroundCard(it, lang, tones[i%len(tones)]))
	}

	return h.Section(
		h.ID("playground"),
		h.Style("padding:var(--sp-12) 0;border-top:var(--bw-2) solid var(--ink)"),
		layout.Container(layout.ContainerProps{},
			sectionHeader(num, i18n.T(lang, "sections.playground.title"), i18n.T(lang, "sections.playground.sub"), token.ToneAccent),
			h.Div(h.Class("tools-grid"), g.Group(cards)),
		),
	)
}

func playgroundCard(it playgroundItem, lang string, tone token.Tone) g.Node {
	return primitive.Card(primitive.CardProps{
		Tone:  tone,
		Attrs: []g.Node{h.Style("display:flex;flex-direction:column;gap:var(--sp-4);height:100%")},
	},
		h.Div(
			h.Style("display:flex;align-items:flex-start;gap:var(--sp-3)"),
			icon.Icon(it.Icon, icon.Props{Size: "1.4rem"}),
			h.H3(h.Style("font-size:var(--t-xl);font-weight:900;line-height:1.1;margin:0;flex:1"), g.Text(it.Label)),
		),
		h.P(h.Style("margin:0;font-size:var(--t-sm);line-height:1.6;opacity:.85;flex:1"), g.Text(i18n.T(lang, it.DescKey))),
		h.A(
			h.Href(it.Href),
			h.Style("margin-top:auto"),
			primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeMD},
				g.Text(i18n.T(lang, "sections.playground.open")),
			),
		),
	)
}
