package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/internal/i18n"
)

// gameDescriptionSection renders a "how it works" write-up (title + a few
// paragraphs pulled from i18n) alongside an illustrative diagram, reused
// across all four playground pages.
func gameDescriptionSection(lang, titleKey string, paragraphKeys []string, diagram g.Node, diagramCaptionKey string) g.Node {
	paras := make([]g.Node, 0, len(paragraphKeys))
	for _, k := range paragraphKeys {
		paras = append(paras, h.P(g.Text(i18n.T(lang, k))))
	}

	var caption g.Node
	if diagramCaptionKey != "" {
		caption = h.P(h.Class("game-desc-caption"), g.Text(i18n.T(lang, diagramCaptionKey)))
	}

	return subRow("game-desc-section", i18n.T(lang, titleKey), "",
		h.Div(h.Class("game-desc-grid"),
			h.Div(h.Class("game-desc-text"), g.Group(paras)),
			h.Div(h.Class("game-desc-diagram"),
				diagram,
				caption,
			),
		),
	)
}

const gameDescCSS = `
.game-desc-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.3fr) minmax(0, 1fr);
  gap: var(--sp-8);
  align-items: start;
}
.game-desc-text p {
  color: var(--muted);
  line-height: 1.7;
  margin: 0 0 var(--sp-4);
  max-width: 60ch;
}
.game-desc-diagram {
  border: 1px solid var(--ink);
  padding: var(--sp-4);
}
.game-desc-diagram svg {
  width: 100%;
  height: auto;
  display: block;
}
.game-desc-caption {
  margin: var(--sp-3) 0 0;
  font-size: var(--t-sm);
  color: var(--muted);
  text-align: center;
}
@media (max-width: 800px) {
  .game-desc-grid { grid-template-columns: 1fr; }
}
`
