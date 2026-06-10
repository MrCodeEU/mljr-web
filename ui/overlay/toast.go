package overlay

import (
	"fmt"

	"mljr-web/ui/icon"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ToasterProps struct {
	ID    string // default "toaster"
	Attrs []g.Node
}

// Toaster mounts the fixed toast region. Server-side SSE handlers patch new
// Toast nodes into this region via sse.PatchElements.
func Toaster(p ToasterProps, children ...g.Node) g.Node {
	if p.ID == "" {
		p.ID = "toaster"
	}
	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "toaster"),
		g.Attr("aria-live", "polite"),
		g.Attr("aria-atomic", "true"),
		g.Group(p.Attrs),
		g.Group(children),
	)
}

// ToastIconName returns the canonical icon for a toast variant.
func ToastIconName(v token.ToastVariant) string {
	switch v {
	case token.ToastSuccess:
		return "lucide:circle-check"
	case token.ToastWarning:
		return "lucide:alert-triangle"
	case token.ToastDanger:
		return "lucide:circle-x"
	default:
		return "lucide:info"
	}
}

type ToastProps struct {
	ID          string             // unique per toast (drives morph/remove)
	Variant     token.ToastVariant // info | success | warning | danger
	Title       string
	Icon        string // icon name; defaults to variant icon if empty
	AutoDismiss int    // seconds before auto-close (0 = manual only)
	Attrs       []g.Node
}

// Toast renders one toast notification. Pair with Toaster.
// For auto-dismiss, set AutoDismiss > 0; a CSS animation handles the timer and
// animationend removes the element.
func Toast(p ToastProps, children ...g.Node) g.Node {
	iconName := p.Icon
	if iconName == "" && p.Variant != "" {
		iconName = ToastIconName(p.Variant)
	}

	attrs := []g.Node{
		g.Attr("data-component", "toast"),
		g.If(p.Variant != "", g.Attr("data-variant", string(p.Variant))),
		g.Attr("role", "status"),
	}
	if p.ID != "" {
		attrs = append(attrs, h.ID(p.ID))
	}
	if p.AutoDismiss > 0 {
		dur := fmt.Sprintf("%ds", p.AutoDismiss)
		attrs = append(attrs,
			g.Attr("data-autodismiss", ""),
			h.Style(fmt.Sprintf("--toast-duration:%s", dur)),
			// remove element when exit animation ends
			g.Attr("data-on:animationend",
				"if(evt.animationName==='mljr-toast-out'){evt.target.remove()}"),
		)
	}
	attrs = append(attrs, p.Attrs...)

	var iconNode g.Node
	if iconName != "" {
		iconNode = h.Div(g.Attr("data-slot", "icon"), icon.Icon(iconName))
	}

	body := h.Div(
		g.Attr("data-slot", "content"),
		g.If(p.Title != "", h.Div(g.Attr("data-slot", "title"), g.Text(p.Title))),
		g.If(len(children) > 0, h.Div(g.Group(children))),
	)

	closeBtn := h.Button(
		g.Attr("data-slot", "close"),
		g.Attr("aria-label", "Dismiss"),
		g.Attr("data-on:click", "evt.target.closest('[data-component=toast]').remove()"),
		g.Text("×"),
	)

	var progressBar g.Node
	if p.AutoDismiss > 0 {
		progressBar = h.Div(g.Attr("data-slot", "progress"))
	}

	return h.Div(append(attrs, iconNode, body, progressBar, closeBtn)...)
}
