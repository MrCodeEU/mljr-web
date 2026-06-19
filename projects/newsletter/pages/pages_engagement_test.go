package pages

import (
	"net/http"
	"testing"

	"mljr-web/projects/newsletter/internal/testutil"
)

func TestSuggestionCreateVoteApprove(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "sugowner@example.com", "password123")
	member := testutil.CreateUser(t, app.app, "sugmember@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Suggestion Crew", "suggestion-crew", owner.Id)
	testutil.CreateMembership(t, app.app, group.Id, member.Id, "member")
	slug := group.GetString("slug")

	body, ct := formBody(map[string]string{"prompt": "Best meal this week?", "type": "text"})
	res := app.do(t, http.MethodPost, "/g/"+slug+"/suggestions", body,
		mergeHeaders(cookieHeader(t, member), map[string]string{"content-type": ct}))
	expectStatus(t, res, http.StatusSeeOther)

	suggestions, err := app.app.FindRecordsByFilter("question_suggestions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d (err=%v)", len(suggestions), err)
	}
	suggestion := suggestions[0]
	if suggestion.GetString("status") != "pending" {
		t.Fatalf("expected status=pending, got %q", suggestion.GetString("status"))
	}

	res = app.do(t, http.MethodPost, "/g/"+slug+"/suggestions/"+suggestion.Id+"/vote", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	votes, err := app.app.FindRecordsByFilter("question_suggestion_votes", "suggestion = {:suggestion}", "", 0, 0,
		map[string]any{"suggestion": suggestion.Id})
	if err != nil || len(votes) != 1 {
		t.Fatalf("expected 1 vote, got %d (err=%v)", len(votes), err)
	}

	// Voting again toggles the vote back off.
	res = app.do(t, http.MethodPost, "/g/"+slug+"/suggestions/"+suggestion.Id+"/vote", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)
	votes, err = app.app.FindRecordsByFilter("question_suggestion_votes", "suggestion = {:suggestion}", "", 0, 0,
		map[string]any{"suggestion": suggestion.Id})
	if err != nil || len(votes) != 0 {
		t.Fatalf("expected vote to be removed, got %d (err=%v)", len(votes), err)
	}

	// A non-admin member cannot approve.
	res = app.do(t, http.MethodPost, "/g/"+slug+"/suggestions/"+suggestion.Id+"/approve", nil, cookieHeader(t, member))
	expectStatus(t, res, 403)

	res = app.do(t, http.MethodPost, "/g/"+slug+"/suggestions/"+suggestion.Id+"/approve", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	suggestion, err = app.app.FindRecordById("question_suggestions", suggestion.Id)
	if err != nil {
		t.Fatal(err)
	}
	if suggestion.GetString("status") != "approved" {
		t.Fatalf("expected status=approved, got %q", suggestion.GetString("status"))
	}

	questions, err := app.app.FindRecordsByFilter("question_bank",
		"group = {:group} && scope = \"group\" && prompt = {:prompt}", "", 0, 0,
		map[string]any{"group": group.Id, "prompt": "Best meal this week?"})
	if err != nil || len(questions) != 1 {
		t.Fatalf("expected approved suggestion to create 1 question_bank row, got %d (err=%v)", len(questions), err)
	}
	if !questions[0].GetBool("is_active") {
		t.Error("expected promoted question to be active")
	}
}

func TestReactionToggleAndComments(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "reactowner@example.com", "password123")
	member := testutil.CreateUser(t, app.app, "reactmember@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Reaction Crew", "reaction-crew", owner.Id)
	testutil.CreateMembership(t, app.app, group.Id, member.Id, "member")
	slug := group.GetString("slug")

	res := app.do(t, http.MethodPost, "/g/"+slug+"/editions", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)
	editions, err := app.app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(editions) != 1 {
		t.Fatalf("expected 1 edition, got %d (err=%v)", len(editions), err)
	}
	edition := editions[0]

	eqs, err := app.app.FindRecordsByFilter("edition_questions", "edition = {:edition}", "", 0, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(eqs) == 0 {
		t.Fatalf("expected edition_questions to be populated, got %d (err=%v)", len(eqs), err)
	}
	var textQuestionID string
	for _, eq := range eqs {
		q, err := app.app.FindRecordById("question_bank", eq.GetString("question"))
		if err == nil && q.GetString("type") == "text" {
			textQuestionID = q.Id
			break
		}
	}
	if textQuestionID == "" {
		t.Fatal("expected at least one seeded text question")
	}

	mpBody, mpCT := multipartBody(t, map[string]string{"q_" + textQuestionID: "Tacos"})
	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id, mpBody,
		mergeHeaders(cookieHeader(t, owner), map[string]string{"content-type": mpCT}))
	expectStatus(t, res, http.StatusSeeOther)

	answer, err := app.app.FindFirstRecordByFilter("answers",
		"edition = {:edition} && question = {:question} && user = {:user}",
		map[string]any{"edition": edition.Id, "question": textQuestionID, "user": owner.Id})
	if err != nil {
		t.Fatalf("expected answer to exist: %v", err)
	}

	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id+"/close", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	reactBody, reactCT := formBody(map[string]string{"emoji": "👍"})
	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id+"/answers/"+answer.Id+"/react", reactBody,
		mergeHeaders(cookieHeader(t, member), map[string]string{"content-type": reactCT}))
	expectStatus(t, res, http.StatusSeeOther)

	reactions, err := app.app.FindRecordsByFilter("emoji_reactions", "answer = {:answer}", "", 0, 0,
		map[string]any{"answer": answer.Id})
	if err != nil || len(reactions) != 1 {
		t.Fatalf("expected 1 reaction, got %d (err=%v)", len(reactions), err)
	}

	notifs, err := app.app.FindRecordsByFilter("notifications", "kind = \"emoji_reaction\" && user = {:user}", "", 0, 0,
		map[string]any{"user": owner.Id})
	if err != nil || len(notifs) != 1 {
		t.Fatalf("expected 1 reaction notification for answer owner, got %d (err=%v)", len(notifs), err)
	}

	// Toggle off.
	reactBody2, reactCT2 := formBody(map[string]string{"emoji": "👍"})
	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id+"/answers/"+answer.Id+"/react", reactBody2,
		mergeHeaders(cookieHeader(t, member), map[string]string{"content-type": reactCT2}))
	expectStatus(t, res, http.StatusSeeOther)
	reactions, err = app.app.FindRecordsByFilter("emoji_reactions", "answer = {:answer}", "", 0, 0,
		map[string]any{"answer": answer.Id})
	if err != nil || len(reactions) != 0 {
		t.Fatalf("expected reaction to be removed, got %d (err=%v)", len(reactions), err)
	}

	commentBody, commentCT := formBody(map[string]string{"body": "Tacos are great!"})
	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id+"/answers/"+answer.Id+"/comments", commentBody,
		mergeHeaders(cookieHeader(t, member), map[string]string{"content-type": commentCT}))
	expectStatus(t, res, http.StatusSeeOther)

	comments, err := app.app.FindRecordsByFilter("comments", "answer = {:answer}", "", 0, 0,
		map[string]any{"answer": answer.Id})
	if err != nil || len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d (err=%v)", len(comments), err)
	}
	topComment := comments[0]

	replyBody, replyCT := formBody(map[string]string{"body": "Agreed!", "parent": topComment.Id})
	res = app.do(t, http.MethodPost, "/g/"+slug+"/editions/"+edition.Id+"/answers/"+answer.Id+"/comments", replyBody,
		mergeHeaders(cookieHeader(t, owner), map[string]string{"content-type": replyCT}))
	expectStatus(t, res, http.StatusSeeOther)

	allComments, err := app.app.FindRecordsByFilter("comments", "answer = {:answer}", "", 0, 0,
		map[string]any{"answer": answer.Id})
	if err != nil || len(allComments) != 2 {
		t.Fatalf("expected 2 comments (1 top-level + 1 reply), got %d (err=%v)", len(allComments), err)
	}

	res = app.do(t, http.MethodGet, "/g/"+slug+"/editions/"+edition.Id+"/view", nil, cookieHeader(t, owner))
	expectStatus(t, res, 200)
	expectContains(t, res, "Tacos are great!", "Agreed!")

	res = app.do(t, http.MethodGet, "/g/"+slug+"/recap", nil, cookieHeader(t, owner))
	expectStatus(t, res, 200)
	expectContains(t, res, "Tacos")
}
