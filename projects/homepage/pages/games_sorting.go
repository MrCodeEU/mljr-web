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

func Sorting(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Sorting Visualizer - Michael Reinegger",
			Description: "Race classic sorting algorithms as animated bars: bubble, selection, insertion, merge, and quick sort.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS + gameCSS + boidsCSS + sortingCSS + gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("Sorting Visualizer")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.sorting.lede"))),
				h.Div(h.ID("sorting-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("control-row"),
					h.Label(h.For("sorting-algo"), g.Text(i18n.T(lang, "games.sorting.label_algorithm"))),
					h.Select(h.ID("sorting-algo"),
						h.Option(h.Value("bubble"), g.Text(i18n.T(lang, "games.sorting.algo_bubble"))),
						h.Option(h.Value("selection"), g.Text(i18n.T(lang, "games.sorting.algo_selection"))),
						h.Option(h.Value("insertion"), g.Text(i18n.T(lang, "games.sorting.algo_insertion"))),
						h.Option(h.Value("merge"), h.Selected(), g.Text(i18n.T(lang, "games.sorting.algo_merge"))),
						h.Option(h.Value("quick"), g.Text(i18n.T(lang, "games.sorting.algo_quick"))),
					),
					h.Label(h.For("sorting-size"), g.Text(i18n.T(lang, "games.sorting.label_size"))),
					h.Input(h.Type("range"), h.ID("sorting-size"), h.Min("10"), h.Max("200"), h.Value("80")),
					h.Label(h.For("sorting-speed"), g.Text(i18n.T(lang, "games.sorting.label_speed"))),
					h.Input(h.Type("range"), h.ID("sorting-speed"), h.Min("1"), h.Max("50"), h.Value("3")),
				),
				h.Div(h.Class("control-row"),
					h.Button(h.ID("sorting-shuffle"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.sorting.shuffle"))),
					h.Button(h.ID("sorting-start"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.sorting.start"))),
					h.Span(h.ID("sorting-status"), h.Class("sorting-status")),
				),
			),
		),
		sortingDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/sorting/sketch.js")),
	)
}

const sortingCSS = `
.sorting-status {
  font-size: var(--t-sm);
  color: var(--muted);
  font-weight: 700;
}
.control-row select {
  font: inherit;
  box-sizing: border-box;
  padding: var(--sp-1) var(--sp-2);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
}
`

func sortingDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.sorting.desc_title",
		[]string{"games.sorting.desc_p1", "games.sorting.desc_p2", "games.sorting.desc_p3"},
		sortingDiagram(), "")
}

func sortingDiagram() g.Node {
	heights := []int{40, 90, 60, 130, 30, 100, 70, 50, 110, 80}
	var bars []g.Node
	for i, h2 := range heights {
		color := "var(--muted)"
		if i == 3 || i == 4 {
			color = "#e6505a"
		}
		bars = append(bars, g.El("rect",
			g.Attr("x", strconv.Itoa(10+i*22)), g.Attr("y", strconv.Itoa(150-h2)),
			g.Attr("width", "18"), g.Attr("height", strconv.Itoa(h2)),
			g.Attr("style", "fill:"+color)))
	}
	return g.El("svg",
		g.Attr("viewBox", "0 0 240 170"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.Group(bars),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "165"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("highlighted pair: currently being compared or swapped")),
	)
}
