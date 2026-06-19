package pages

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func findEditionInGroup(re *core.RequestEvent, groupID, id string) (*core.Record, error) {
	edition, err := re.App.FindRecordById("newsletter_editions", id)
	if err != nil || edition.GetString("group") != groupID {
		return nil, fmt.Errorf("edition not found")
	}
	return edition, nil
}

type editionQuestion struct {
	EQ       *core.Record
	Question *core.Record
}

func editionQuestions(re *core.RequestEvent, editionID string) ([]editionQuestion, error) {
	eqs, err := re.App.FindRecordsByFilter(
		"edition_questions", "edition = {:edition}", "order", 0, 0,
		map[string]any{"edition": editionID},
	)
	if err != nil {
		return nil, err
	}
	out := make([]editionQuestion, 0, len(eqs))
	for _, eq := range eqs {
		q, err := re.App.FindRecordById("question_bank", eq.GetString("question"))
		if err != nil {
			continue
		}
		out = append(out, editionQuestion{EQ: eq, Question: q})
	}
	return out, nil
}

// ListEditions shows the edition archive plus a "start new edition" action
// (admin only) when no edition is currently open.
func ListEditions(re *core.RequestEvent) error {
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

	editions, err := re.App.FindRecordsByFilter(
		"newsletter_editions", "group = {:group}", "-created", 0, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	hasOpen := false
	var rows []g.Node
	for _, ed := range editions {
		status := ed.GetString("status")
		if status == "open" {
			hasOpen = true
		}
		href := "/g/" + slug + "/editions/" + ed.Id
		if status == "sent" || status == "archived" {
			href += "/view"
		}
		rows = append(rows, h.A(
			h.Href(href),
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line);text-decoration:none;color:var(--ink)"),
			h.Span(g.Text(ed.GetString("created")[:10])),
			h.Span(h.Style("color:var(--muted);font-size:var(--t-sm);text-transform:uppercase;letter-spacing:.04em;white-space:nowrap"), g.Text(status)),
		))
	}

	var newButton g.Node
	if isAdmin && !hasOpen {
		newButton = h.Form(h.Method("post"), h.Action("/g/"+slug+"/editions"), h.Style("margin-bottom:var(--sp-6)"),
			primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit"}, g.Text("Start new edition")),
		)
	}

	return renderPage(re, 200, appPage(re, slug, "Editions — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Editions"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Editions")),
		h.Div(h.Style("margin-top:var(--sp-4)"), newButton),
		g.If(len(rows) == 0, h.P(h.Style("color:var(--muted)"), g.Text("No editions yet."))),
		g.If(len(rows) > 0, primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("padding:var(--sp-2) var(--sp-4)")}}, g.Group(rows))),
	))
}

// HandleCreateEdition starts a new edition with every active question
// (global + this group's custom ones) currently in the bank, admin only.
func HandleCreateEdition(re *core.RequestEvent) error {
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

	existingOpen, err := re.App.FindFirstRecordByFilter(
		"newsletter_editions", "group = {:group} && status = \"open\"",
		map[string]any{"group": group.Id},
	)
	if err == nil {
		return redirect(re, "/g/"+slug+"/editions/"+existingOpen.Id)
	}

	questions, err := re.App.FindRecordsByFilter(
		"question_bank", "is_active = true && (scope = \"global\" || (scope = \"group\" && group = {:group}))", "created", 0, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	editionsCol, err := re.App.FindCollectionByNameOrId("newsletter_editions")
	if err != nil {
		return err
	}
	edition := core.NewRecord(editionsCol)
	edition.Set("group", group.Id)
	edition.Set("status", "open")
	edition.Set("opens_at", types.NowDateTime())
	if err := re.App.Save(edition); err != nil {
		return err
	}

	eqCol, err := re.App.FindCollectionByNameOrId("edition_questions")
	if err != nil {
		return err
	}
	for i, q := range questions {
		eq := core.NewRecord(eqCol)
		eq.Set("edition", edition.Id)
		eq.Set("question", q.Id)
		eq.Set("order", i)
		if err := re.App.Save(eq); err != nil {
			return err
		}
	}

	return redirect(re, "/g/"+slug+"/editions/"+edition.Id)
}

func findAnswer(re *core.RequestEvent, editionID, questionID, userID string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter(
		"answers", "edition = {:edition} && question = {:question} && user = {:user}",
		map[string]any{"edition": editionID, "question": questionID, "user": userID},
	)
}

// EditionAnswer renders the answer form for every question in the edition,
// pre-filled with whatever the current user already submitted.
func EditionAnswer(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	if edition.GetString("status") == "sent" || edition.GetString("status") == "archived" {
		return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"/view")
	}

	eqs, err := editionQuestions(re, edition.Id)
	if err != nil {
		return err
	}

	var fields []g.Node
	for _, eq := range eqs {
		q := eq.Question
		existing, _ := findAnswer(re, edition.Id, q.Id, user.Id)
		fields = append(fields, h.Div(h.Style("margin-bottom:var(--sp-6)"),
			renderAnswerInput(q, existing),
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Answer — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "This edition"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("This edition's questions")),
		h.Form(
			h.Method("post"), h.Action("/g/"+slug+"/editions/"+edition.Id),
			g.Attr("enctype", "multipart/form-data"),
			h.Style("margin-top:var(--sp-4)"),
			g.Group(fields),
			primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit"}, g.Text("Save answers")),
		),
	))
}

// answerValue decodes a JSONField "value" column into a plain Go value
// (string, []any, etc.) — Record.Get returns the raw types.JSONRaw bytes,
// not the decoded value, for JSON-typed fields.
func answerValue(rec *core.Record, key string) any {
	raw, ok := rec.Get(key).(types.JSONRaw)
	if !ok || len(raw) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil
	}
	return v
}

// valueAsString coerces a decoded JSON answer value to a string — ratings
// round-trip as JSON numbers (PocketBase's JSONField normalizes numeric-
// looking strings like "4" into the bare JSON number 4 on write), so a plain
// string type assertion misses them.
func valueAsString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	}
	return ""
}

func answerFieldName(q *core.Record) string { return "q_" + q.Id }

func renderAnswerInput(q *core.Record, existing *core.Record) g.Node {
	name := answerFieldName(q)
	prompt := primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("padding:var(--sp-4)")}},
		h.P(h.Style("font-weight:600;margin-bottom:var(--sp-3)"), g.Text(q.GetString("prompt"))),
		answerControl(q, name, existing),
	)
	return prompt
}

func answerControl(q *core.Record, name string, existing *core.Record) g.Node {
	var existingStr string
	var existingList []string
	if existing != nil {
		switch v := answerValue(existing, "value").(type) {
		case []any:
			for _, item := range v {
				if s, ok := item.(string); ok {
					existingList = append(existingList, s)
				}
			}
		default:
			existingStr = valueAsString(v)
		}
	}
	contains := func(list []string, v string) bool {
		for _, item := range list {
			if item == v {
				return true
			}
		}
		return false
	}

	switch q.GetString("type") {
	case "text":
		return form.Textarea(form.TextareaProps{Name: name, Placeholder: "Your answer…"}, g.Text(existingStr))
	case "single_select":
		opts := make([]form.RadioOption, 0)
		for _, v := range questionOptions(q) {
			opts = append(opts, form.RadioOption{Value: v, Label: v, Checked: v == existingStr})
		}
		return form.RadioGroup(form.RadioGroupProps{Name: name, Options: opts})
	case "multi_select":
		var boxes []g.Node
		for _, v := range questionOptions(q) {
			boxes = append(boxes, form.Checkbox(form.CheckboxProps{
				Label: v, Name: name + "[]", Checked: contains(existingList, v),
				Attrs: []g.Node{h.Value(v)},
			}))
		}
		return h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"), g.Group(boxes))
	case "rating":
		opts := make([]form.RadioOption, 0, 5)
		for i := 1; i <= 5; i++ {
			val := strconv.Itoa(i)
			opts = append(opts, form.RadioOption{Value: val, Label: strings.Repeat("★", i), Checked: val == existingStr})
		}
		return form.RadioGroup(form.RadioGroupProps{Name: name, Options: opts, Attrs: []g.Node{h.Style("flex-direction:row;gap:var(--sp-3)")}})
	case "emoji_reaction":
		opts := make([]form.RadioOption, 0)
		for _, v := range questionOptions(q) {
			opts = append(opts, form.RadioOption{Value: v, Label: v, Checked: v == existingStr})
		}
		return form.RadioGroup(form.RadioGroupProps{Name: name, Options: opts, Attrs: []g.Node{h.Style("flex-direction:row;gap:var(--sp-3);font-size:1.5rem")}})
	case "image":
		return form.FileInput(form.FileInputProps{Name: "f_" + answerFieldNameSuffix(q), Accept: "image/*"})
	default:
		return g.Text("")
	}
}

func answerFieldNameSuffix(q *core.Record) string { return q.Id }

func upsertAnswer(re *core.RequestEvent, editionID, questionID, userID string, value any, skipped bool) (*core.Record, error) {
	answer, err := findAnswer(re, editionID, questionID, userID)
	if err != nil {
		col, cErr := re.App.FindCollectionByNameOrId("answers")
		if cErr != nil {
			return nil, cErr
		}
		answer = core.NewRecord(col)
		answer.Set("edition", editionID)
		answer.Set("question", questionID)
		answer.Set("user", userID)
	}
	if value != nil {
		answer.Set("value", value)
	}
	answer.Set("skipped", skipped)
	if err := re.App.Save(answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// HandleSubmitAnswers upserts one answer row per question from the
// submitted form, including image uploads into answer_images.
func HandleSubmitAnswers(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	if edition.GetString("status") != "open" {
		return re.BadRequestError("this edition is no longer accepting answers", nil)
	}

	if err := re.Request.ParseMultipartForm(20 << 20); err != nil {
		return err
	}

	eqs, err := editionQuestions(re, edition.Id)
	if err != nil {
		return err
	}

	for _, eq := range eqs {
		q := eq.Question
		name := answerFieldName(q)

		switch q.GetString("type") {
		case "text":
			val := strings.TrimSpace(re.Request.FormValue(name))
			if _, err := upsertAnswer(re, edition.Id, q.Id, user.Id, val, val == ""); err != nil {
				return err
			}
		case "single_select", "emoji_reaction", "rating":
			val := re.Request.FormValue(name)
			if _, err := upsertAnswer(re, edition.Id, q.Id, user.Id, val, val == ""); err != nil {
				return err
			}
		case "multi_select":
			vals := re.Request.Form[name+"[]"]
			if _, err := upsertAnswer(re, edition.Id, q.Id, user.Id, vals, len(vals) == 0); err != nil {
				return err
			}
		case "image":
			file, header, ferr := re.Request.FormFile("f_" + answerFieldNameSuffix(q))
			existing, _ := findAnswer(re, edition.Id, q.Id, user.Id)
			hasExistingImage := false
			if existing != nil {
				if imgs, _ := re.App.FindRecordsByFilter("answer_images", "answer = {:answer}", "", 1, 0, map[string]any{"answer": existing.Id}); len(imgs) > 0 {
					hasExistingImage = true
				}
			}
			answer, aerr := upsertAnswer(re, edition.Id, q.Id, user.Id, nil, !hasExistingImage && ferr != nil)
			if aerr != nil {
				return aerr
			}
			if ferr == nil {
				_ = file.Close()
				f, ffErr := filesystem.NewFileFromMultipart(header)
				if ffErr != nil {
					return ffErr
				}
				imagesCol, icErr := re.App.FindCollectionByNameOrId("answer_images")
				if icErr != nil {
					return icErr
				}
				img := core.NewRecord(imagesCol)
				img.Set("answer", answer.Id)
				img.Set("image", f)
				if err := re.App.Save(img); err != nil {
					return err
				}
				answer.Set("skipped", false)
				_ = re.App.Save(answer)
			}
		}
	}

	return redirect(re, "/g/"+slug+"/editions/"+edition.Id)
}

// HandleCloseEdition manually flips an edition to "sent" so members can see
// the compiled answers. Admin only — scheduling automates this in Phase 4.
func HandleCloseEdition(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	edition.Set("status", "sent")
	edition.Set("sent_at", types.NowDateTime())
	if err := re.App.Save(edition); err != nil {
		return err
	}
	return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"/view")
}

// EditionView shows every member's answers, question by question.
func EditionView(re *core.RequestEvent) error {
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
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}
	if edition.GetString("status") != "sent" && edition.GetString("status") != "archived" {
		return redirect(re, "/g/"+slug+"/editions/"+edition.Id)
	}

	eqs, err := editionQuestions(re, edition.Id)
	if err != nil {
		return err
	}
	resolveUser := userDisplayNameCache(re)

	var sections []g.Node
	for _, eq := range eqs {
		q := eq.Question
		answers, _ := re.App.FindRecordsByFilter(
			"answers", "edition = {:edition} && question = {:question} && skipped = false", "created",
			0, 0, map[string]any{"edition": edition.Id, "question": q.Id},
		)
		var rows []g.Node
		for _, a := range answers {
			u, uErr := re.App.FindRecordById("users", a.GetString("user"))
			if uErr != nil {
				continue
			}
			rows = append(rows, h.Div(h.Style("padding:var(--sp-2) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
				h.Span(h.Style("font-weight:600;margin-right:var(--sp-2)"), g.Text(displayName(u)+":")),
				renderAnswerValue(re, q, a),
				renderReactionBar(slug, edition.Id, a, reactionCounts(re, a.Id, user.Id)),
				renderCommentThreads(slug, edition.Id, a, commentThreads(re, a.Id), resolveUser),
			))
		}
		sections = append(sections, primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("padding:var(--sp-4);margin-bottom:var(--sp-4)")}},
			h.P(h.Style("font-weight:600;margin-bottom:var(--sp-3)"), g.Text(q.GetString("prompt"))),
			g.If(len(rows) == 0, h.P(h.Style("color:var(--muted)"), g.Text("No answers."))),
			g.Group(rows),
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Edition — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Edition"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Edition results")),
		h.Div(h.Style("margin-top:var(--sp-4)"), g.Group(sections)),
	))
}

func renderAnswerValue(re *core.RequestEvent, q *core.Record, a *core.Record) g.Node {
	switch q.GetString("type") {
	case "text":
		return g.Text(valueAsString(answerValue(a, "value")))
	case "single_select", "emoji_reaction":
		return g.Text(valueAsString(answerValue(a, "value")))
	case "rating":
		s := valueAsString(answerValue(a, "value"))
		n, _ := strconv.Atoi(s)
		return g.Text(strings.Repeat("★", n) + strings.Repeat("☆", 5-n))
	case "multi_select":
		list, _ := answerValue(a, "value").([]any)
		var parts []string
		for _, v := range list {
			if s, ok := v.(string); ok {
				parts = append(parts, s)
			}
		}
		return g.Text(strings.Join(parts, ", "))
	case "image":
		imgs, _ := re.App.FindRecordsByFilter("answer_images", "answer = {:answer}", "order", 0, 0, map[string]any{"answer": a.Id})
		var thumbs []g.Node
		for _, img := range imgs {
			url := "/api/files/" + img.Collection().Id + "/" + img.Id + "/" + img.GetString("image")
			thumbs = append(thumbs, h.Img(h.Src(url), h.Style("max-width:160px;max-height:160px;border-radius:var(--radius-md);margin-right:var(--sp-2)")))
		}
		return h.Div(g.Group(thumbs))
	default:
		return g.Text("")
	}
}
