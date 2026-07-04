package pages

import (
	"encoding/json"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

// fourierI18NJSON hands the panel-label strings (drawn onto the canvas by
// the JS sketch, not rendered as HTML) to the client so the canvas labels
// match the page's language too.
func fourierI18NJSON(lang string) string {
	labels := map[string]string{
		"drawHint":       i18n.T(lang, "games.fourier.draw_hint"),
		"yourDrawing":    i18n.T(lang, "games.fourier.your_drawing"),
		"reconstruction": i18n.T(lang, "games.fourier.reconstruction"),
		"waveLabel":      i18n.T(lang, "games.fourier.wave_label"),
	}
	b, _ := json.Marshal(labels)
	return string(b)
}

func Fourier(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Fourier Drawing - Michael Reinegger",
			Description: "Draw a shape and watch a chain of rotating circles reconstruct it from sine waves via the Discrete Fourier Transform.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS + gameCSS + fourierCSS + gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell game-shell-wide"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("Fourier Drawing")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.fourier.lede"))),
				h.Div(h.ID("fourier-game"), h.Class("game-canvas-wrap")),
				h.Div(h.Class("game-controls fourier-controls"),
					h.Button(h.ID("fourier-play"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.fourier.play"))),
					h.Button(h.ID("fourier-clear"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.fourier.clear"))),
					h.Label(h.For("fourier-waves"), g.Text(i18n.T(lang, "games.fourier.label_waves"))),
					h.Input(h.Type("range"), h.ID("fourier-waves"), h.Min("1"), h.Max("150"), h.Value("30")),
					h.Span(h.ID("fourier-waves-value"), g.Text("30")),
				),
			),
		),
		fourierDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		g.El("script", g.Raw(
			"window.FOURIER_I18N="+fourierI18NJSON(lang)+";",
		)),
		h.Script(h.Src("/static/games/fourier/sketch.js")),
	)
}

const fourierCSS = `
.fourier-controls {
  gap: var(--sp-3);
  flex-wrap: wrap;
}
.fourier-controls input[type="range"] {
  flex: 1 1 12rem;
  min-width: 8rem;
}
#fourier-waves-value {
  font-weight: 900;
  min-width: 2.5ch;
  text-align: right;
}
`

func fourierDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.fourier.desc_title",
		[]string{"games.fourier.desc_p1", "games.fourier.desc_p2", "games.fourier.desc_p3", "games.fourier.desc_p4"},
		fourierDiagram(), "")
}

func fourierDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 260 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.El("circle", g.Attr("cx", "80"), g.Attr("cy", "110"), g.Attr("r", "55"), g.Attr("style", "stroke:var(--muted);stroke-width:1;fill:none")),
		g.El("circle", g.Attr("cx", "135"), g.Attr("cy", "80"), g.Attr("r", "28"), g.Attr("style", "stroke:var(--muted);stroke-width:1;fill:none")),
		g.El("circle", g.Attr("cx", "163"), g.Attr("cy", "95"), g.Attr("r", "12"), g.Attr("style", "stroke:var(--muted);stroke-width:1;fill:none")),
		g.El("line", g.Attr("x1", "80"), g.Attr("y1", "110"), g.Attr("x2", "135"), g.Attr("y2", "80"), g.Attr("style", "stroke:var(--ink);stroke-width:2")),
		g.El("line", g.Attr("x1", "135"), g.Attr("y1", "80"), g.Attr("x2", "163"), g.Attr("y2", "95"), g.Attr("style", "stroke:var(--ink);stroke-width:2")),
		g.El("line", g.Attr("x1", "163"), g.Attr("y1", "95"), g.Attr("x2", "175"), g.Attr("y2", "72"), g.Attr("style", "stroke:var(--ink);stroke-width:2")),
		g.El("circle", g.Attr("cx", "175"), g.Attr("cy", "72"), g.Attr("r", "3"), g.Attr("style", "fill:#e56b6f")),
		g.El("line", g.Attr("x1", "175"), g.Attr("y1", "72"), g.Attr("x2", "230"), g.Attr("y2", "72"), g.Attr("style", "stroke:var(--accent, #78b4ff);stroke-width:1;stroke-dasharray:3,3")),
		g.El("path", g.Attr("d", "M230 40 L235 100 L240 55 L245 90 L250 65"), g.Attr("style", "stroke:#78b4ff;stroke-width:2;fill:none")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "210"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("circles chained tip-to-tail; the final tip traces the shape and feeds the wave graph")),
	)
}
