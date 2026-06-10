package primitive

import (
	"fmt"

	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ToggleOption struct {
	Value string
	Label string
	Icon  string // optional icon name
}

type ToggleGroupProps struct {
	Signal  string
	Default string
	Attrs   []g.Node
}

// ToggleGroup renders a segmented button control bound to a Datastar signal.
// Only one option can be active at a time.
func ToggleGroup(p ToggleGroupProps, options ...ToggleOption) g.Node {
	sig := p.Signal
	initialSig := fmt.Sprintf(`{%s:%q}`, sig, p.Default)

	buttons := make([]g.Node, len(options))
	for i, opt := range options {
		val := opt.Value
		var iconNode g.Node
		if opt.Icon != "" {
			iconNode = icon.Icon(opt.Icon)
		}
		buttons[i] = h.Button(
			g.Attr("data-component", "toggle-btn"),
			h.Type("button"),
			g.Attr("data-on:click", fmt.Sprintf("$%s=%q", sig, val)),
			g.Attr("data-attr", fmt.Sprintf(`{"data-state":$%s===%q?"active":""}`, sig, val)),
			iconNode,
			g.If(opt.Label != "", g.Text(opt.Label)),
		)
	}

	return h.Div(
		g.Attr("data-component", "toggle-group"),
		g.Attr("data-signals", initialSig),
		g.Group(p.Attrs),
		g.Group(buttons),
	)
}
