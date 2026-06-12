package primitive

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CopyButtonProps struct {
	Text    string        // text to copy to clipboard
	Variant token.Variant // default Outline
	Size    token.Size    // default SizeIcon
	Label   string        // optional visible label after icon
}

// CopyButton renders a clipboard copy button with Datastar-powered copied feedback.
// Wraps in a data-signals scope so each instance has independent state.
func CopyButton(p CopyButtonProps) g.Node {
	if p.Variant == "" {
		p.Variant = token.Outline
	}
	if p.Size == "" {
		p.Size = token.SizeIcon
	}
	copyExpr := `navigator.clipboard.writeText(` + "`" + p.Text + "`" + `).catch(()=>{});$_cpd=true;setTimeout(()=>$_cpd=false,2000)`

	inner := []g.Node{
		h.Span(g.Attr("data-show", "!$_cpd"), icon.Icon("lucide:copy")),
		h.Span(g.Attr("data-show", "$_cpd"), h.Style("display:none"), icon.Icon("lucide:check")),
	}
	if p.Label != "" {
		inner = append(inner,
			h.Span(g.Attr("data-show", "!$_cpd"), g.Text(p.Label)),
			h.Span(g.Attr("data-show", "$_cpd"), h.Style("display:none"), g.Text("Copied!")),
		)
	}

	return h.Span(
		g.Attr("data-component", "copy-button"),
		g.Attr("data-signals", `{"_cpd":false}`),
		h.Button(
			g.Attr("data-component", "button"),
			g.Attr("data-variant", string(p.Variant)),
			g.Attr("data-size", string(p.Size)),
			h.Type("button"),
			h.Title("Copy to clipboard"),
			g.Attr("data-on:click", copyExpr),
			g.Group(inner),
		),
	)
}
