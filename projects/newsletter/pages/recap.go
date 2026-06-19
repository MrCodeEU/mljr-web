package pages

import (
	"mljr-web/ui/primitive"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

const recapEditionLimit = 5

// Recap shows the current user's own answers across their last few sent
// editions in this group, so they can look back without digging through
// the edition archive one at a time.
func Recap(re *core.RequestEvent) error {
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

	editions, err := re.App.FindRecordsByFilter(
		"newsletter_editions", "group = {:group} && (status = \"sent\" || status = \"archived\")",
		"-sent_at", recapEditionLimit, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	var sections []g.Node
	for _, edition := range editions {
		eqs, err := editionQuestions(re, edition.Id)
		if err != nil {
			continue
		}
		var rows []g.Node
		for _, eq := range eqs {
			answer, err := findAnswer(re, edition.Id, eq.Question.Id, user.Id)
			if err != nil || answer.GetBool("skipped") {
				continue
			}
			rows = append(rows, h.Div(h.Style("padding:var(--sp-2) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
				h.P(h.Style("font-weight:600;font-size:var(--t-sm);color:var(--muted)"), g.Text(eq.Question.GetString("prompt"))),
				renderAnswerValue(re, eq.Question, answer),
			))
		}
		if len(rows) == 0 {
			continue
		}
		created := edition.GetString("created")
		dateLabel := created
		if len(created) >= 10 {
			dateLabel = created[:10]
		}
		sections = append(sections, primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("padding:var(--sp-4);margin-bottom:var(--sp-4)")}},
			h.Div(h.Style("display:flex;justify-content:space-between;align-items:baseline;margin-bottom:var(--sp-3)"),
				h.P(h.Style("font-weight:700"), g.Text(dateLabel)),
				h.A(h.Href("/g/"+slug+"/editions/"+edition.Id+"/view"), h.Style("font-size:var(--t-sm)"), g.Text("View full edition")),
			),
			g.Group(rows),
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Recap — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Recap"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Your recap")),
		h.P(h.Style("color:var(--muted);margin:var(--sp-2) 0 var(--sp-6)"), g.Text("Your own answers from recent editions.")),
		g.If(len(sections) == 0, h.P(h.Style("color:var(--muted)"), g.Text("No past answers yet."))),
		g.Group(sections),
	))
}
