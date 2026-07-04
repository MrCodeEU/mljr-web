package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

const pendulumMaxSegments = 5

func DoublePendulum(lang string, a AnalyticsConfig) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title:       "Double Pendulum - Michael Reinegger",
			Description: "A chaotic n-segment pendulum physics simulation, ported and rebuilt from an old side project.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			Lang:        lang,
			HeadExtra:   append([]g.Node{g.El("style", g.Raw(homepageCSS + gameCSS + pendulumCSS + gameDescCSS))}, AnalyticsHead(a)...),
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		siteNavbar(lang),
		h.Main(h.Class("game-page"),
			h.Div(h.Class("game-shell game-shell-wide"),
				h.A(h.Href("/"), h.Class("legal-back"), icon.Icon("lucide:arrow-left"), g.Text(i18n.T(lang, "legal.back"))),
				h.H1(g.Text("Double Pendulum")),
				h.P(h.Class("game-lede"), g.Text(i18n.T(lang, "games.pendulum.lede"))),
				h.Div(h.ID("double-pendulum-game"), h.Class("game-canvas-wrap")),
				pendulumControls(lang),
			),
		),
		pendulumDescriptionSection(lang),
		siteFooter(lang),
		h.Script(h.Src("/static/games/vendor/p5.min.js")),
		h.Script(h.Src("/static/games/double-pendulum/sketch.js")),
	)
}

func pendulumControls(lang string) g.Node {
	rows := make([]g.Node, 0, pendulumMaxSegments)
	for i := 0; i < pendulumMaxSegments; i++ {
		display := ""
		if i >= 3 {
			display = "display:none"
		}
		rows = append(rows, h.Div(
			h.ID(fmt.Sprintf("pendulum-row-%d", i)),
			h.Class("pendulum-row"),
			h.Style(display),
			h.Span(h.Class("pendulum-row-label"), g.Text(i18n.T(lang, "games.pendulum.label_segment", i+1))),
			h.Label(h.For(fmt.Sprintf("pendulum-len-%d", i)), g.Text(i18n.T(lang, "games.pendulum.label_length"))),
			h.Input(h.Type("number"), h.ID(fmt.Sprintf("pendulum-len-%d", i)), h.Min("20"), h.Max("200"), h.Value("100")),
			h.Label(h.For(fmt.Sprintf("pendulum-mass-%d", i)), g.Text(i18n.T(lang, "games.pendulum.label_mass"))),
			h.Input(h.Type("number"), h.ID(fmt.Sprintf("pendulum-mass-%d", i)), h.Min("1"), h.Max("40"), h.Value("15")),
		))
	}

	return h.Div(h.Class("pendulum-controls"),
		h.Div(h.Class("pendulum-row"),
			h.Label(h.For("pendulum-segments"), g.Text(i18n.T(lang, "games.pendulum.label_segments"))),
			h.Select(h.ID("pendulum-segments"),
				h.Option(h.Value("2"), g.Text("2")),
				h.Option(h.Value("3"), h.Selected(), g.Text("3")),
				h.Option(h.Value("4"), g.Text("4")),
				h.Option(h.Value("5"), g.Text("5")),
			),
			h.Label(h.For("pendulum-gravity"), g.Text(i18n.T(lang, "games.pendulum.label_gravity"))),
			h.Input(h.Type("number"), h.ID("pendulum-gravity"), h.Min("0.1"), h.Max("5"), h.Step("0.1"), h.Value("1")),
			h.Label(h.For("pendulum-speed"), g.Text(i18n.T(lang, "games.pendulum.label_speed"))),
			h.Input(h.Type("range"), h.ID("pendulum-speed"), h.Min("1"), h.Max("100"), h.Value("25")),
		),
		g.Group(rows),
		h.Button(h.ID("pendulum-reset"), h.Class("game-btn"), g.Text(i18n.T(lang, "games.pendulum.apply_reset"))),
	)
}

const pendulumCSS = `
.game-shell-wide {
  max-width: 1080px;
}
.pendulum-controls {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
  margin-bottom: var(--sp-4);
}
.pendulum-row {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  flex-wrap: wrap;
  font-size: var(--t-sm);
}
.pendulum-row-label {
  font-weight: 900;
  min-width: 6rem;
}
.pendulum-row input[type="number"] {
  font: inherit;
  box-sizing: border-box;
  width: 4.5rem;
  text-align: center;
  padding: var(--sp-1) var(--sp-2);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
}
.pendulum-row input[type="range"] {
  flex: 1 1 12rem;
  min-width: 8rem;
}
.pendulum-row select {
  font: inherit;
  box-sizing: border-box;
  padding: var(--sp-1) var(--sp-2);
  border: var(--bw-2) solid var(--line);
  background: var(--surface);
  color: var(--ink);
}
`

func pendulumDescriptionSection(lang string) g.Node {
	return gameDescriptionSection(lang, "games.pendulum.desc_title",
		[]string{"games.pendulum.desc_p1", "games.pendulum.desc_p2", "games.pendulum.desc_p3"},
		pendulumDiagram(), "")
}

func pendulumDiagram() g.Node {
	return g.El("svg",
		g.Attr("viewBox", "0 0 280 220"),
		g.Attr("xmlns", "http://www.w3.org/2000/svg"),
		g.El("circle", g.Attr("cx", "140"), g.Attr("cy", "20"), g.Attr("r", "4"), g.Attr("style", "fill:var(--ink)")),
		g.El("path", g.Attr("d", "M140 20 L140 60"), g.Attr("style", "stroke:var(--muted);stroke-width:1;stroke-dasharray:3,3;fill:none")),
		g.El("path", g.Attr("d", "M140 20 A40 40 0 0 1 165 55"), g.Attr("style", "stroke:var(--accent, #e56b6f);stroke-width:1.5;fill:none")),
		g.El("text", g.Attr("x", "150"), g.Attr("y", "45"), g.Attr("style", "font-size:11px;fill:var(--accent, #e56b6f)"), g.Text("θ₁")),
		g.El("line", g.Attr("x1", "140"), g.Attr("y1", "20"), g.Attr("x2", "185"), g.Attr("y2", "95"),
			g.Attr("style", "stroke:var(--ink);stroke-width:2")),
		g.El("circle", g.Attr("cx", "185"), g.Attr("cy", "95"), g.Attr("r", "9"), g.Attr("style", "fill:#78b4ff")),
		g.El("text", g.Attr("x", "195"), g.Attr("y", "90"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("m₁, l₁")),

		g.El("path", g.Attr("d", "M185 95 L185 135"), g.Attr("style", "stroke:var(--muted);stroke-width:1;stroke-dasharray:3,3;fill:none")),
		g.El("path", g.Attr("d", "M185 95 A40 40 0 0 1 210 130"), g.Attr("style", "stroke:var(--accent, #e56b6f);stroke-width:1.5;fill:none")),
		g.El("text", g.Attr("x", "195"), g.Attr("y", "120"), g.Attr("style", "font-size:11px;fill:var(--accent, #e56b6f)"), g.Text("θ₂")),
		g.El("line", g.Attr("x1", "185"), g.Attr("y1", "95"), g.Attr("x2", "150"), g.Attr("y2", "170"),
			g.Attr("style", "stroke:var(--ink);stroke-width:2")),
		g.El("circle", g.Attr("cx", "150"), g.Attr("cy", "170"), g.Attr("r", "9"), g.Attr("style", "fill:#ff78a0")),
		g.El("text", g.Attr("x", "160"), g.Attr("y", "185"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("m₂, l₂")),
		g.El("text", g.Attr("x", "10"), g.Attr("y", "210"), g.Attr("style", "font-size:11px;fill:var(--muted)"), g.Text("...continues for each added segment")),
	)
}
