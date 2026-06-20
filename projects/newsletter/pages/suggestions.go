package pages

import (
	"strconv"
	"strings"

	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func countSuggestionVotes(re *core.RequestEvent, suggestionID string) int {
	votes, err := re.App.FindRecordsByFilter(
		"question_suggestion_votes", "suggestion = {:suggestion}", "", 0, 0,
		map[string]any{"suggestion": suggestionID},
	)
	if err != nil {
		return 0
	}
	return len(votes)
}

func findVote(re *core.RequestEvent, suggestionID, userID string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter(
		"question_suggestion_votes", "suggestion = {:suggestion} && user = {:user}",
		map[string]any{"suggestion": suggestionID, "user": userID},
	)
}

// ListSuggestions shows pending question suggestions for a group, ranked by
// vote count, with a form for any member to propose a new one and a vote
// toggle. Owners/admins additionally get approve/reject actions.
func ListSuggestions(re *core.RequestEvent) error {
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

	pending, err := re.App.FindRecordsByFilter(
		"question_suggestions", "group = {:group} && status = \"pending\"", "-created", 0, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	type ranked struct {
		rec   *core.Record
		votes int
	}
	var items []ranked
	for _, s := range pending {
		items = append(items, ranked{rec: s, votes: countSuggestionVotes(re, s.Id)})
	}
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].votes > items[i].votes {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	var rows []g.Node
	for _, it := range items {
		s := it.rec
		hasVoted := false
		if _, err := findVote(re, s.Id, user.Id); err == nil {
			hasVoted = true
		}
		voteLabel := "Vote"
		if hasVoted {
			voteLabel = "Voted ✓"
		}

		actions := []g.Node{
			h.Form(h.Method("post"), h.Action("/g/"+slug+"/suggestions/"+s.Id+"/vote"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit"}, g.Text(voteLabel)),
			),
		}
		if isAdmin {
			actions = append(actions,
				h.Form(h.Method("post"), h.Action("/g/"+slug+"/suggestions/"+s.Id+"/approve"),
					primitive.Button(primitive.ButtonProps{Variant: token.Primary, Tone: token.ToneNone, Type: "submit"}, g.Text("Approve")),
				),
				h.Form(h.Method("post"), h.Action("/g/"+slug+"/suggestions/"+s.Id+"/reject"),
					primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit"}, g.Text("Reject")),
				),
			)
		}

		rows = append(rows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Div(h.Style("min-width:0"),
				h.Span(h.Style("min-width:0;overflow-wrap:anywhere"), g.Text(s.GetString("prompt"))),
				h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"),
					g.Text(questionTypeLabels[s.GetString("type")]+" · "+strconv.Itoa(it.votes)+" vote(s)")),
			),
			h.Div(h.Style("display:flex;gap:var(--sp-2);flex-wrap:wrap"), g.Group(actions)),
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Suggestions — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Suggestions"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Suggest a question")),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/g/"+slug+"/suggestions"),
				form.Field(form.FieldProps{Label: "Prompt"},
					form.Input(form.InputProps{Type: "text", Name: "prompt", Required: true, Placeholder: "What should we ask next?"}),
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
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}}, g.Text("Suggest")),
			),
		),
		primitive.Heading(primitive.HeadingProps{Level: 2, Attrs: []g.Node{h.Style("margin-top:var(--sp-8)")}}, g.Text("Pending suggestions")),
		g.If(len(rows) == 0, h.P(h.Style("color:var(--muted);margin-top:var(--sp-3)"), g.Text("No pending suggestions."))),
		g.If(len(rows) > 0, primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-3);padding:var(--sp-2) var(--sp-4)")}}, g.Group(rows))),
	))
}

func HandleCreateSuggestion(re *core.RequestEvent) error {
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
		return redirect(re, "/g/"+slug+"/suggestions")
	}

	col, err := re.App.FindCollectionByNameOrId("question_suggestions")
	if err != nil {
		return err
	}
	s := core.NewRecord(col)
	s.Set("group", group.Id)
	s.Set("suggested_by", user.Id)
	s.Set("type", qtype)
	s.Set("prompt", prompt)
	s.Set("status", "pending")

	if qtype == "single_select" || qtype == "multi_select" || qtype == "emoji_reaction" {
		var opts []string
		for _, part := range strings.Split(re.Request.FormValue("options"), ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				opts = append(opts, part)
			}
		}
		s.Set("options", opts)
	}

	if err := re.App.Save(s); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/suggestions")
}

// HandleToggleVote toggles the current user's vote on a suggestion on/off.
func HandleToggleVote(re *core.RequestEvent) error {
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

	suggestion, err := re.App.FindRecordById("question_suggestions", re.Request.PathValue("id"))
	if err != nil || suggestion.GetString("group") != group.Id {
		return re.NotFoundError("suggestion not found", err)
	}

	if existing, err := findVote(re, suggestion.Id, user.Id); err == nil {
		if err := re.App.Delete(existing); err != nil {
			return err
		}
		return redirect(re, "/g/"+slug+"/suggestions")
	}

	votesCol, err := re.App.FindCollectionByNameOrId("question_suggestion_votes")
	if err != nil {
		return err
	}
	vote := core.NewRecord(votesCol)
	vote.Set("suggestion", suggestion.Id)
	vote.Set("user", user.Id)
	if err := re.App.Save(vote); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/suggestions")
}

// HandleApproveSuggestion promotes a suggestion into the group's active
// question bank and marks it approved (owner/admin only).
func HandleApproveSuggestion(re *core.RequestEvent) error {
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

	suggestion, err := re.App.FindRecordById("question_suggestions", re.Request.PathValue("id"))
	if err != nil || suggestion.GetString("group") != group.Id {
		return re.NotFoundError("suggestion not found", err)
	}

	questionsCol, err := re.App.FindCollectionByNameOrId("question_bank")
	if err != nil {
		return err
	}
	q := core.NewRecord(questionsCol)
	q.Set("scope", "group")
	q.Set("group", group.Id)
	q.Set("author", suggestion.GetString("suggested_by"))
	q.Set("type", suggestion.GetString("type"))
	q.Set("prompt", suggestion.GetString("prompt"))
	q.Set("is_active", true)
	if raw := answerValue(suggestion, "options"); raw != nil {
		q.Set("options", raw)
	}
	if err := re.App.Save(q); err != nil {
		return err
	}

	suggestion.Set("status", "approved")
	if err := re.App.Save(suggestion); err != nil {
		return err
	}

	if suggestedBy := suggestion.GetString("suggested_by"); suggestedBy != user.Id {
		_ = createNotification(re.App, suggestedBy, "comment_reply", group.Id, "", user.Id,
			"Your suggested question \""+suggestion.GetString("prompt")+"\" was approved",
			"/g/"+slug+"/questions")
	}

	return redirect(re, "/g/"+slug+"/suggestions")
}

func HandleRejectSuggestion(re *core.RequestEvent) error {
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

	suggestion, err := re.App.FindRecordById("question_suggestions", re.Request.PathValue("id"))
	if err != nil || suggestion.GetString("group") != group.Id {
		return re.NotFoundError("suggestion not found", err)
	}
	suggestion.Set("status", "rejected")
	if err := re.App.Save(suggestion); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/suggestions")
}
