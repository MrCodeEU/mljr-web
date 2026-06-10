package form

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SearchInputProps struct {
	Name        string
	Placeholder string
	// Target is the URL for @get (Datastar SSE fetch). Empty = no auto-fetch.
	Target string
	// Debounce is the debounce modifier value (default "300ms").
	Debounce string
	// IndicatorSignal is the name of the loading indicator signal (default "_searching").
	IndicatorSignal string
	// Clearable adds an × button to reset the input.
	Clearable bool
}

// SearchInput renders a search input with icon prefix, optional debounced @get,
// and optional clear button.
func SearchInput(p SearchInputProps) g.Node {
	if p.Placeholder == "" {
		p.Placeholder = "Search…"
	}
	if p.Debounce == "" {
		p.Debounce = "300ms"
	}
	if p.IndicatorSignal == "" {
		p.IndicatorSignal = "_searching"
	}
	if p.Name == "" {
		p.Name = "q"
	}

	inputAttrs := []g.Node{
		g.Attr("data-component", "input"),
		h.Type("search"),
		h.Name(p.Name),
		h.Placeholder(p.Placeholder),
		g.Attr("data-bind:"+p.Name),
		g.Attr("autocomplete", "off"),
	}
	if p.Target != "" {
		inputAttrs = append(inputAttrs,
			g.Attr("data-on:input__debounce."+p.Debounce, "@get('"+p.Target+"')"),
			g.Attr("data-indicator:"+p.IndicatorSignal),
		)
	}

	children := []g.Node{
		h.Div(g.Attr("data-slot", "icon"), icon.Icon("lucide:search")),
		h.Input(inputAttrs...),
	}

	if p.Clearable {
		children = append(children, h.Button(
			g.Attr("data-slot", "clear"),
			h.Type("button"),
			g.Attr("data-show", "$"+p.Name+"!==''"),
			h.Style("display:none"),
			g.Attr("data-on:click", "$"+p.Name+"=''"),
			g.Attr("aria-label", "Clear search"),
			icon.Icon("lucide:x", icon.Props{Size: "1rem"}),
		))
	}

	if p.Target != "" {
		children = append(children, h.Div(
			g.Attr("data-slot", "spinner"),
			g.Attr("data-show", "$"+p.IndicatorSignal),
			h.Style("display:none"),
		))
	}

	return h.Div(
		g.Attr("data-component", "search-input"),
		g.Attr("data-signals", `{"`+p.Name+`":""}`),
		g.Group(children),
	)
}
