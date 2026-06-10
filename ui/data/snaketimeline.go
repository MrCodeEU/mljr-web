package data

import (
	"fmt"

	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// SnakeTimelineItem is one event in the snake timeline.
type SnakeTimelineItem struct {
	Period  string
	Title   string
	Org     string
	OrgLogo string
	Desc    string
	Tags    []string
	Tone    token.Tone
}

type SnakeTimelineProps struct {
	// Cols is the number of items per row before snaking (default 2, max 4).
	Cols int
}

// SnakeTimeline renders items in a serpentine pattern.
// Rows alternate LTR/RTL. The connector rail is row-level so the line,
// row turn, dots, and cards all share the same geometry.
func SnakeTimeline(p SnakeTimelineProps, items ...SnakeTimelineItem) g.Node {
	if p.Cols < 1 {
		p.Cols = 2
	}
	if p.Cols > 4 {
		p.Cols = 4
	}

	var rows [][]SnakeTimelineItem
	for i := 0; i < len(items); i += p.Cols {
		end := i + p.Cols
		if end > len(items) {
			end = len(items)
		}
		rows = append(rows, items[i:end])
	}

	var segments []g.Node
	globalIdx := 0

	for ri, row := range rows {
		isRTL := ri%2 == 1
		dir := "ltr"
		if isRTL {
			dir = "rtl"
		}

		// Dot rail and cards use the same column count so numbered points
		// stay aligned with their cards.
		rowCols := len(row)
		rowStyle := fmt.Sprintf("--snake-cols:%d", rowCols)
		if ri == len(rows)-1 {
			rowStyle += ";padding-bottom:0"
		}

		dotCols := make([]g.Node, len(row))
		for ci := range row {
			num := globalIdx + ci + 1
			dotCols[ci] = snakeDot(num)
		}

		cards := make([]g.Node, len(row))
		for ci, item := range row {
			cards[ci] = snakeCard(item)
		}

		var turn g.Node
		if ri < len(rows)-1 {
			side := "right"
			if isRTL {
				side = "left"
			}
			turn = h.Div(g.Attr("data-slot", "turn"), g.Attr("data-side", side))
		}

		segments = append(segments, h.Div(
			g.Attr("data-slot", "row"),
			g.Attr("data-dir", dir),
			h.Style(rowStyle),
			h.Div(
				g.Attr("data-slot", "path"),
				h.Div(g.Attr("data-slot", "rail")),
				h.Div(g.Attr("data-slot", "dots"), g.Group(dotCols)),
			),
			h.Div(g.Attr("data-slot", "items"), g.Group(cards)),
			turn,
		))

		globalIdx += len(row)
	}

	return h.Div(
		g.Attr("data-component", "snake-timeline"),
		g.Attr("data-cols", fmt.Sprintf("%d", p.Cols)),
		g.Group(segments),
	)
}

func snakeDot(num int) g.Node {
	return h.Div(
		g.Attr("data-slot", "dot-cell"),
		h.Div(g.Attr("data-slot", "dot"), g.Text(fmt.Sprintf("%d", num))),
	)
}

func snakeCard(item SnakeTimelineItem) g.Node {
	var logoNode g.Node
	if item.OrgLogo != "" {
		logoNode = h.Img(
			h.Src(item.OrgLogo),
			h.Alt(item.Org),
			h.Style("width:52px;height:52px;object-fit:contain;border-radius:var(--radius);border:var(--bw-2) solid var(--line);flex-shrink:0"),
		)
	}

	tagNodes := make([]g.Node, 0, len(item.Tags))
	for _, t := range item.Tags {
		if t != "" {
			tagNodes = append(tagNodes, primitive.Tag(primitive.TagProps{}, g.Text(t)))
		}
	}

	return h.Div(
		g.Attr("data-slot", "item"),
		h.Style("direction:ltr"), // ensure LTR even in RTL rows
		primitive.Card(primitive.CardProps{Tone: item.Tone},
			h.Div(
				h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
				h.Div(
					h.Div(h.Style("font-size:var(--t-xs);font-weight:800;opacity:.6;text-transform:uppercase;letter-spacing:.06em"), g.Text(item.Period)),
					g.If(item.Org != "", h.Div(h.Style("font-size:var(--t-xs);opacity:.65;margin-top:2px"), g.Text(item.Org))),
				),
				logoNode,
			),
			h.H4(h.Style("font-weight:900;font-size:var(--t-base);margin:0 0 var(--sp-2);line-height:1.3"), g.Text(item.Title)),
			g.If(item.Desc != "", h.P(h.Style("font-size:var(--t-sm);opacity:.8;margin:0 0 var(--sp-2);line-height:1.5"), g.Text(item.Desc))),
			g.If(len(tagNodes) > 0, h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-1)"), g.Group(tagNodes))),
		),
	)
}
