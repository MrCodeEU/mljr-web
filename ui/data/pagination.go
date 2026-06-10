package data

import (
	"fmt"
	"math"

	"mljr-web/ui"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PaginationProps struct {
	ID      string
	Total   int
	PerPage int
	Attrs   []g.Node
}

func PaginationSignals(id string, perPage int) g.Node {
	if id == "" {
		id = "pg"
	}
	sig := id + "Page"
	return h.Div(ui.Signals(fmt.Sprintf(`{%s:0}`, sig)))
}

func Pagination(p PaginationProps) g.Node {
	id := p.ID
	if id == "" {
		id = "pg"
	}
	perPage := p.PerPage
	if perPage <= 0 {
		perPage = 6
	}
	sig := id + "Page"

	pages := int(math.Ceil(float64(p.Total) / float64(perPage)))
	if pages < 1 {
		pages = 1
	}
	maxPage := pages - 1

	prevExpr := fmt.Sprintf("$%s = Math.max(0,$%s-1)", sig, sig)
	nextExpr := fmt.Sprintf("$%s = Math.min(%d,$%s+1)", sig, maxPage, sig)

	// prev: disabled when on first page
	prevDisabled := fmt.Sprintf(`{"data-state": $%s === 0 ? "disabled" : ""}`, sig)
	// next: disabled when on last page
	nextDisabled := fmt.Sprintf(`{"data-state": $%s === %d ? "disabled" : ""}`, sig, maxPage)

	btns := make([]g.Node, 0, pages+2)
	btns = append(btns,
		h.Button(
			g.Attr("data-slot", "prev"),
			g.Attr("data-attr", prevDisabled),
			ui.On("click", fmt.Sprintf("if($%s>0){%s}", sig, prevExpr)),
			g.Text("←"),
		),
	)
	for i := range pages {
		idx := i
		// active state driven by Datastar — data-attr sets data-state="active" when current page
		activeAttr := fmt.Sprintf(`{"data-state": $%s === %d ? "active" : ""}`, sig, idx)
		btns = append(btns,
			h.Button(
				g.Attr("data-slot", "btn"),
				g.Attr("data-attr", activeAttr),
				ui.On("click", fmt.Sprintf("$%s = %d", sig, idx)),
				g.Text(fmt.Sprintf("%d", idx+1)),
			),
		)
	}
	btns = append(btns,
		h.Button(
			g.Attr("data-slot", "next"),
			g.Attr("data-attr", nextDisabled),
			ui.On("click", fmt.Sprintf("if($%s<%d){%s}", sig, maxPage, nextExpr)),
			g.Text("→"),
		),
	)

	return h.Div(
		g.Attr("data-component", "pagination"),
		g.Group(p.Attrs),
		g.Group(btns),
	)
}
