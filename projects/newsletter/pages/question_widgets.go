package pages

import (
	"strings"

	"mljr-web/ui/primitive"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// renderQuestionPreviewCard renders a live, client-side-only preview of a
// question being created, driven entirely by the form's own prompt/type/
// options signals — no server round-trip.
func renderQuestionPreviewCard(t func(string, ...any) string, promptSig, typeSig, optionsSig string) g.Node {
	optsExpr := "$" + optionsSig + ".split(',').map(s=>s.trim()).filter(s=>s)"
	return primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4);background:var(--surface-2)")}},
		h.P(h.Style("font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em;margin-bottom:var(--sp-2)"),
			g.Text(t("newsletter.question_widgets.live_preview"))),
		h.P(h.Style("font-weight:600;margin-bottom:var(--sp-3)"),
			g.Attr("data-text", "$"+promptSig+" || '"+t("newsletter.question_widgets.placeholder_prompt")+"'"),
		),
		h.Div(
			g.Attr("data-show", "$"+typeSig+"==='single_select'||$"+typeSig+"==='multi_select'||$"+typeSig+"==='emoji_reaction'"),
			h.Style("padding:var(--sp-2) var(--sp-3);border-radius:var(--radius-full);background:var(--surface-3);font-size:var(--t-sm);display:inline-block"),
			g.Attr("data-text", optsExpr+".join(' · ') || '"+t("newsletter.question_widgets.no_options")+"'"),
		),
		h.Div(
			g.Attr("data-show", "$"+typeSig+"==='rating'"),
			primitive.Rating(primitive.RatingProps{Signal: "__preview_rating", ReadOnly: true}),
		),
		h.Div(
			g.Attr("data-show", "$"+typeSig+"==='image'"),
			h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text(t("newsletter.question_widgets.image_note"))),
		),
		h.Div(
			g.Attr("data-show", "$"+typeSig+"==='text'"),
			h.Div(h.Style("border:var(--border-w) dashed var(--line);border-radius:var(--radius-md);padding:var(--sp-3);color:var(--muted);font-size:var(--t-sm)"),
				g.Text(t("newsletter.question_widgets.free_text"))),
		),
	)
}

// renderQuestionSummary renders a static, full-detail summary of an existing
// question/suggestion record — type plus its full option list / scale /
// upload note, so reviewers see exactly what they're voting on.
func renderQuestionSummary(q *core.Record) g.Node {
	label := questionTypeLabels[q.GetString("type")]
	switch q.GetString("type") {
	case "single_select", "multi_select", "emoji_reaction":
		opts := questionOptions(q)
		if len(opts) > 0 {
			return h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted)"),
				g.Text(label+": "+strings.Join(opts, " · ")))
		}
	case "rating":
		return h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted)"), g.Text(label+": 1-5 stars"))
	case "image":
		return h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted)"), g.Text(label+": photo upload"))
	case "scale":
		return h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted)"), g.Text(label+": 1-10"))
	}
	return h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"), g.Text(label))
}
