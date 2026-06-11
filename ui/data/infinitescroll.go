package data

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type InfiniteScrollProps struct {
	// FetchURL is the Datastar @get URL. Must return HTML to append/replace.
	FetchURL string
	// PageSignal is the Datastar signal name for the current page (default "_isPage").
	// The URL receives it as a query param automatically via Datastar signals.
	PageSignal string
	// ContainerID is the id of the element that receives new items (default "is-container").
	ContainerID string
	// LoadingText shown while fetching (default "Loading…").
	LoadingText string
}

// InfiniteScroll renders a sentinel div that triggers a Datastar @get request
// when it scrolls into the viewport. The server returns HTML to append.
// Uses data-on:intersect__once — fires once per sentinel appearance.
func InfiniteScroll(p InfiniteScrollProps, initialItems ...g.Node) g.Node {
	if p.PageSignal == "" {
		p.PageSignal = "_isPage"
	}
	if p.ContainerID == "" {
		p.ContainerID = "is-container"
	}
	if p.LoadingText == "" {
		p.LoadingText = "Loading…"
	}

	// Expression: increment page, fetch URL
	fetchExpr := "$" + p.PageSignal + "=$" + p.PageSignal + "+1;@get('" + p.FetchURL + "')"
	loadingSignal := p.PageSignal + "Loading"

	return h.Div(
		g.Attr("data-component", "infinite-scroll"),
		g.Attr("data-signals", `{"`+p.PageSignal+`":1,"`+loadingSignal+`":false}`),
		h.Div(h.ID(p.ContainerID), g.Attr("data-slot", "container"),
			g.Group(initialItems),
		),
		// Sentinel — triggers fetch when it enters viewport.
		// __once fires once per rendered sentinel: the server response must
		// include a fresh sentinel (or remove it on the last page).
		h.Div(
			g.Attr("data-slot", "sentinel"),
			g.Attr("data-on:intersect__once", fetchExpr),
			g.Attr("data-indicator:"+loadingSignal),
			h.Style("height:1px;margin-top:var(--sp-4)"),
		),
		// Loading indicator (shown while fetching)
		h.Div(
			g.Attr("data-slot", "loading"),
			g.Attr("data-show", "$"+loadingSignal),
			h.Style("text-align:center;padding:var(--sp-4);color:var(--muted);font-size:var(--t-sm)"),
			g.Text(p.LoadingText),
		),
	)
}
