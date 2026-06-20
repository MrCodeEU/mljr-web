package pages

import (
	"mljr-web/ui/feedback"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type flashMsg struct {
	variant feedback.AlertVariant
	text    string
}

// flashCatalog maps every ?flash= code used by a redirect somewhere in this
// project to the banner shown on the page it lands on. This app uses plain
// multipart POST/redirect (no SSE), so a query param is the simplest way to
// carry "it worked"/"it didn't" feedback across the redirect.
var flashCatalog = map[string]flashMsg{
	"saved":   {feedback.AlertSuccess, "Saved — all questions answered."},
	"partial": {feedback.AlertInfo, "Saved — you can keep editing the rest until this edition closes."},

	"invite_sent":          {feedback.AlertSuccess, "Invite sent."},
	"invite_email_failed":  {feedback.AlertWarning, "Invite created, but the email failed to send — share the link from the list below instead."},
	"invite_invalid_email": {feedback.AlertDanger, "Enter a valid email address to send an invite."},
	"invite_revoked":       {feedback.AlertSuccess, "Invite revoked."},

	"profile_saved": {feedback.AlertSuccess, "Profile updated."},
	"name_required": {feedback.AlertDanger, "Name can't be empty."},

	"group_settings_saved": {feedback.AlertSuccess, "Group settings saved."},
	"group_name_required":  {feedback.AlertDanger, "Group name can't be empty."},
	"left_group":           {feedback.AlertSuccess, "You left the group."},

	"comment_empty":    {feedback.AlertDanger, "Comment can't be empty."},
	"reaction_invalid": {feedback.AlertDanger, "Pick a reaction to add."},

	"image_invalid":   {feedback.AlertDanger, "That file isn't a readable image — try a JPEG, PNG, GIF, or WebP."},
	"image_too_large": {feedback.AlertDanger, "That image is too large (max 25MB) — try a smaller file."},

	"edition_questions_saved":  {feedback.AlertSuccess, "Question set saved."},
	"edition_questions_locked": {feedback.AlertWarning, "This edition has already opened — its questions are locked in."},
}

// flashAlert renders the banner for the current ?flash= code, if any.
func flashAlert(re *core.RequestEvent) g.Node {
	msg, ok := flashCatalog[re.Request.URL.Query().Get("flash")]
	if !ok {
		return nil
	}
	return feedback.Alert(
		feedback.AlertProps{Variant: msg.variant, Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
		g.Text(msg.text),
	)
}
