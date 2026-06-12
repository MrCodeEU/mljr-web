package overlay

import (
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AlertDialogProps struct {
	SignalName  string // default "_alertOpen"
	Title       string
	Description string
	ConfirmText string        // default "Confirm"
	CancelText  string        // default "Cancel"
	Variant     token.Variant // button variant for confirm (default Danger)
	OnConfirm   string        // Datastar expression to run on confirm
}

// AlertDialog renders a confirmation dialog driven by a Datastar signal.
// Open it by setting the signal to true: data-on:click="$_alertOpen=true"
func AlertDialog(p AlertDialogProps) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_alertOpen"
	}
	if p.Title == "" {
		p.Title = "Are you sure?"
	}
	if p.ConfirmText == "" {
		p.ConfirmText = "Confirm"
	}
	if p.CancelText == "" {
		p.CancelText = "Cancel"
	}
	if p.Variant == "" {
		p.Variant = token.Danger
	}

	sig := p.SignalName
	closeExpr := "$" + sig + "=false"
	confirmExpr := closeExpr
	if p.OnConfirm != "" {
		confirmExpr = p.OnConfirm + ";" + closeExpr
	}

	return h.Div(
		g.Attr("data-component", "alert-dialog"),
		g.Attr("data-signals", `{"`+sig+`":false}`),
		g.Attr("data-show", "$"+sig),
		h.Style("display:none"),

		h.Div(
			g.Attr("data-slot", "backdrop"),
			g.Attr("data-on:click", closeExpr),
		),

		h.Div(
			g.Attr("data-slot", "panel"),
			g.Attr("role", "alertdialog"),
			g.Attr("aria-modal", "true"),
			g.Attr("aria-labelledby", "alert-title"),

			h.Div(
				g.Attr("data-slot", "icon"),
				icon.Icon("lucide:alert-triangle"),
			),

			h.H2(
				h.ID("alert-title"),
				g.Attr("data-slot", "title"),
				g.Text(p.Title),
			),
			g.If(p.Description != "", h.P(
				g.Attr("data-slot", "description"),
				g.Text(p.Description),
			)),

			h.Div(
				g.Attr("data-slot", "actions"),
				primitive.Button(primitive.ButtonProps{Variant: token.Outline},
					g.Attr("data-on:click", closeExpr),
					g.Text(p.CancelText),
				),
				primitive.Button(primitive.ButtonProps{Variant: p.Variant},
					g.Attr("data-on:click", confirmExpr),
					g.Text(p.ConfirmText),
				),
			),
		),
	)
}
