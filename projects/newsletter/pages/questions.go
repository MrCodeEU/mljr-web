package pages

import (
	"encoding/json"
	"strings"

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

	questions, err := groupQuestions(re, group.Id)
	if err != nil {
		return err
	}

	var rows []g.Node
	for _, q := range questions {
		badge := "Global"
		if q.GetString("scope") == "group" {
			badge = "Custom"
		}
		var toggle g.Node
		if q.GetString("scope") == "group" && isAdmin {
			label := "Deactivate"
			if !q.GetBool("is_active") {
				label = "Activate"
			}
			toggle = h.Form(h.Method("post"), h.Action("/g/"+slug+"/questions/"+q.Id+"/toggle"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit"}, g.Text(label)),
			)
		}
		meta := badge + " · " + questionTypeLabels[q.GetString("type")]
		if !q.GetBool("is_active") {
			meta += " · inactive"
		}
		rows = append(rows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Div(h.Style("min-width:0"),
				h.Span(h.Style("min-width:0;overflow-wrap:anywhere"), g.Text(q.GetString("prompt"))),
				h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"), g.Text(meta)),
			),
			toggle,
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Questions — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Questions"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Question bank")),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/g/"+slug+"/questions"),
				form.Field(form.FieldProps{Label: "Prompt"},
					form.Input(form.InputProps{Type: "text", Name: "prompt", Required: true, Placeholder: "What's a question you'd like to ask the group?"}),
				),
				form.Field(form.FieldProps{Label: "Answer type", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Select(form.SelectProps{
						Name: "type",
						Options: []form.SelectOption{
							{Value: "text", Label: "Text", Selected: true},
							{Value: "single_select", Label: "Single choice"},
							{Value: "multi_select", Label: "Multiple choice"},
							{Value: "image", Label: "Image"},
							{Value: "rating", Label: "Rating"},
							{Value: "emoji_reaction", Label: "Mood emoji"},
						},
					}),
				),
				form.Field(form.FieldProps{Label: "Options (comma-separated, only for choice/mood types)", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Input(form.InputProps{Type: "text", Name: "options", Placeholder: "Great, Good, Okay, Rough"}),
				),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}}, g.Text("Add question")),
			),
		),
		primitive.Heading(primitive.HeadingProps{Level: 2, Attrs: []g.Node{h.Style("margin-top:var(--sp-8)")}}, g.Text("All questions")),
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
	if _, err := requireAdminMembership(re, group, user); err != nil {
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
