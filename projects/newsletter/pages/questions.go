package pages

import (
	"encoding/json"
	"strings"

	"mljr-web/internal/i18n"
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

var questionTypeLabels = map[string]string{
	"text":           "Text",
	"single_select":  "Single choice",
	"multi_select":   "Multiple choice",
	"image":          "Image",
	"rating":         "Rating",
	"emoji_reaction": "Mood emoji",
	"yes_no":         "Yes/No",
	"scale":          "Scale (1-10)",
	"number":         "Number",
	"date":           "Date",
	"color_pick":     "Color",
}

func questionOptions(q *core.Record) []string {
	raw, ok := q.Get("options").(types.JSONRaw)
	if !ok || len(raw) == 0 {
		return nil
	}
	var opts []string
	if err := json.Unmarshal(raw, &opts); err != nil {
		return nil
	}
	return opts
}

// questionPrompt returns a global question's prompt in lang, falling back to
// the canonical (English) prompt if no translation was stored — group/user-
// authored questions never have prompt_i18n set, so they always fall through.
func questionPrompt(q *core.Record, lang string) string {
	if lang != i18n.DefaultLang {
		raw, ok := q.Get("prompt_i18n").(types.JSONRaw)
		if ok && len(raw) > 0 {
			var m map[string]string
			if json.Unmarshal(raw, &m) == nil && m[lang] != "" {
				return m[lang]
			}
		}
	}
	return q.GetString("prompt")
}

// groupQuestions returns the global + this group's custom questions,
// active first.
func groupQuestions(re *core.RequestEvent, groupID string) ([]*core.Record, error) {
	return re.App.FindRecordsByFilter(
		"question_bank", "scope = \"global\" || (scope = \"group\" && group = {:group})", "-is_active,created", 0, 0,
		map[string]any{"group": groupID},
	)
}

// ListQuestions shows the question bank for a group plus a form to add a
// custom group question.
func ListQuestions(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}
	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	membership, err := findMembership(re, group.Id, user.Id)
	if err != nil {
		return re.ForbiddenError("not a member of this group", nil)
	}
	isAdmin := membership.GetString("role") == "owner" || membership.GetString("role") == "admin"
	t := translator(re)

	questions, err := groupQuestions(re, group.Id)
	if err != nil {
		return err
	}

	var rows []g.Node
	for _, q := range questions {
		badge := t("newsletter.questions.global")
		if q.GetString("scope") == "group" {
			badge = t("newsletter.questions.custom")
		}
		var toggle g.Node
		if q.GetString("scope") == "group" && isAdmin {
			label := t("newsletter.questions.deactivate")
			if !q.GetBool("is_active") {
				label = t("newsletter.questions.activate")
			}
			toggle = h.Form(h.Method("post"), h.Action("/g/"+slug+"/questions/"+q.Id+"/toggle"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit"}, g.Text(label)),
			)
		}
		metaPrefix := badge
		if !q.GetBool("is_active") {
			metaPrefix += " · " + t("newsletter.questions.inactive")
		}
		rows = append(rows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Div(h.Style("min-width:0"),
				h.Span(h.Style("min-width:0;overflow-wrap:anywhere"), g.Text(questionPrompt(q, currentLang(re)))),
				h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"), g.Text(metaPrefix)),
				renderQuestionSummary(q),
			),
			toggle,
		))
	}

	return renderPage(re, 200, appPage(re, slug, t("newsletter.subnav.questions")+" — "+group.GetString("name"),
		[]breadcrumbItem{{Label: t("newsletter.nav.dashboard"), Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: t("newsletter.subnav.questions")}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(t("newsletter.questions.bank_heading"))),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				g.Attr("data-signals", `{q_prompt:'',q_type:'text',q_options:''}`),
				h.Method("post"), h.Action("/g/"+slug+"/questions"),
				form.Field(form.FieldProps{Label: t("newsletter.suggestions.prompt_label")},
					form.Input(form.InputProps{Type: "text", Name: "prompt", Required: true, Placeholder: "What's a question you'd like to ask the group?", Signal: "q_prompt"}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.suggestions.type_label"), Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Select(form.SelectProps{
						Name:   "type",
						Signal: "q_type",
						Options: []form.SelectOption{
							{Value: "text", Label: t("newsletter.suggestions.type_text"), Selected: true},
							{Value: "single_select", Label: t("newsletter.suggestions.type_single")},
							{Value: "multi_select", Label: t("newsletter.suggestions.type_multi")},
							{Value: "image", Label: t("newsletter.suggestions.type_image")},
							{Value: "rating", Label: t("newsletter.suggestions.type_rating")},
							{Value: "emoji_reaction", Label: t("newsletter.suggestions.type_emoji")},
							{Value: "yes_no", Label: t("newsletter.suggestions.type_yes_no")},
							{Value: "scale", Label: t("newsletter.suggestions.type_scale")},
							{Value: "number", Label: t("newsletter.suggestions.type_number")},
							{Value: "date", Label: t("newsletter.suggestions.type_date")},
							{Value: "color_pick", Label: t("newsletter.suggestions.type_color")},
						},
					}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.suggestions.options_label"), Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Input(form.InputProps{Type: "text", Name: "options", Placeholder: "Great, Good, Okay, Rough", Signal: "q_options"}),
				),
				renderQuestionPreviewCard(t, "q_prompt", "q_type", "q_options"),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}}, g.Text(t("newsletter.questions.add_button"))),
			),
		),
		primitive.Heading(primitive.HeadingProps{Level: 2, Attrs: []g.Node{h.Style("margin-top:var(--sp-8)")}}, g.Text(t("newsletter.questions.all_heading"))),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-3);padding:var(--sp-2) var(--sp-4)")}}, g.Group(rows)),
	))
}

func HandleCreateQuestion(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}
	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := findMembership(re, group.Id, user.Id); err != nil {
		return re.ForbiddenError("not a member of this group", nil)
	}

	prompt := strings.TrimSpace(re.Request.FormValue("prompt"))
	qtype := re.Request.FormValue("type")
	if prompt == "" || questionTypeLabels[qtype] == "" {
		return redirect(re, "/g/"+slug+"/questions")
	}

	col, err := re.App.FindCollectionByNameOrId("question_bank")
	if err != nil {
		return err
	}
	q := core.NewRecord(col)
	q.Set("scope", "group")
	q.Set("group", group.Id)
	q.Set("author", user.Id)
	q.Set("type", qtype)
	q.Set("prompt", prompt)
	q.Set("is_active", true)

	if qtype == "single_select" || qtype == "multi_select" || qtype == "emoji_reaction" {
		var opts []string
		for _, part := range strings.Split(re.Request.FormValue("options"), ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				opts = append(opts, part)
			}
		}
		q.Set("options", opts)
	}

	if err := re.App.Save(q); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/questions")
}

func HandleToggleQuestion(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}
	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if err := requireAdminMembership(re, group, user); err != nil {
		return err
	}

	q, err := re.App.FindRecordById("question_bank", re.Request.PathValue("id"))
	if err != nil || q.GetString("group") != group.Id {
		return re.NotFoundError("question not found", err)
	}
	q.Set("is_active", !q.GetBool("is_active"))
	if err := re.App.Save(q); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/questions")
}
