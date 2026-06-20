package pages

import (
	"sort"
	"strconv"

	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// EditionQuestions lets a group admin curate which active questions go into
// an edition and in what order, but only while it's still "scheduled" — once
// openScheduledEditions flips it to "open" members may already be answering,
// so editing the set out from under them would be confusing.
func EditionQuestions(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	if edition.GetString("status") != "scheduled" {
		return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"?flash=edition_questions_locked")
	}

	t := translator(re)
	candidates, err := groupQuestions(re, group.Id)
	if err != nil {
		return err
	}
	eqs, err := editionQuestions(re, edition.Id)
	if err != nil {
		return err
	}
	order := map[string]int{}
	for _, eq := range eqs {
		order[eq.Question.Id] = eq.EQ.GetInt("order")
	}

	var rows []g.Node
	for _, q := range candidates {
		if !q.GetBool("is_active") {
			continue
		}
		_, checked := order[q.Id]
		orderValue := order[q.Id]
		rows = append(rows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-3);align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Input(h.Type("checkbox"), h.Name("q_"+q.Id), h.Value("1"), g.If(checked, h.Checked())),
			h.Span(h.Style("flex:1;min-width:0;overflow-wrap:anywhere"), g.Text(questionPrompt(q, currentLang(re)))),
			h.Input(
				h.Type("number"), h.Name("order_"+q.Id), h.Value(strconv.Itoa(orderValue)),
				h.Min("0"), h.Max("999"), h.Step("1"), h.Style("width:5em"),
			),
		))
	}

	return renderPage(re, 200, appPage(re, slug, t("newsletter.edition_questions.title")+" — "+group.GetString("name"),
		[]breadcrumbItem{
			{Label: t("newsletter.nav.dashboard"), Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug},
			{Label: t("newsletter.subnav.editions"), Href: "/g/" + slug + "/editions"}, {Label: t("newsletter.edition_questions.title")},
		},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(t("newsletter.edition_questions.heading"))),
		flashAlert(re),
		h.P(h.Style("color:var(--muted);margin-top:var(--sp-2)"),
			g.Text(t("newsletter.edition_questions.note"))),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4);padding:var(--sp-2) var(--sp-4)")}},
			h.Form(h.Method("post"), h.Action("/g/"+slug+"/editions/"+edition.Id+"/questions"),
				g.Group(rows),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}}, g.Text(t("newsletter.groups.save_button"))),
			),
		),
	))
}

// HandleEditionQuestions replaces an edition's edition_questions rows with
// the curated checked set + order values, the same shape
// populateEditionQuestions builds when it randomly assigns them.
func HandleEditionQuestions(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	if edition.GetString("status") != "scheduled" {
		return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"?flash=edition_questions_locked")
	}

	candidates, err := groupQuestions(re, group.Id)
	if err != nil {
		return err
	}

	type selected struct {
		questionID string
		order      int
	}
	var picks []selected
	for _, q := range candidates {
		if re.Request.FormValue("q_"+q.Id) == "" {
			continue
		}
		order, _ := strconv.Atoi(re.Request.FormValue("order_" + q.Id))
		picks = append(picks, selected{questionID: q.Id, order: order})
	}
	sort.SliceStable(picks, func(i, j int) bool { return picks[i].order < picks[j].order })

	existing, err := re.App.FindRecordsByFilter(
		"edition_questions", "edition = {:edition}", "", 0, 0, map[string]any{"edition": edition.Id},
	)
	if err != nil {
		return err
	}
	for _, eq := range existing {
		if err := re.App.Delete(eq); err != nil {
			return err
		}
	}

	eqCol, err := re.App.FindCollectionByNameOrId("edition_questions")
	if err != nil {
		return err
	}
	for i, p := range picks {
		eq := core.NewRecord(eqCol)
		eq.Set("edition", edition.Id)
		eq.Set("question", p.questionID)
		eq.Set("order", i)
		if err := re.App.Save(eq); err != nil {
			return err
		}
	}

	return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"/questions?flash=edition_questions_saved")
}
