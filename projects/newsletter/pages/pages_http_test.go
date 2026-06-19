package pages

import (
	"net/http"
	"testing"

	"mljr-web/projects/newsletter/internal/testutil"
)

func TestLoginPageRenders(t *testing.T) {
	app := newTestApp(t)
	res := app.do(t, http.MethodGet, "/login", nil, nil)
	expectStatus(t, res, 200)
	expectContains(t, res, `name="email"`, `name="password"`)
}

func TestHandleLoginSuccess(t *testing.T) {
	app := newTestApp(t)
	testutil.CreateUser(t, app.app, "loginok@example.com", "password123")

	body, ct := formBody(map[string]string{"email": "loginok@example.com", "password": "password123"})
	res := app.do(t, http.MethodPost, "/login", body, map[string]string{"content-type": ct})
	expectStatus(t, res, http.StatusSeeOther)

	if loc := res.Header.Get("Location"); loc != "/" {
		t.Errorf("expected redirect to /, got %q", loc)
	}
	found := false
	for _, c := range res.Cookies() {
		if c.Name == "nl_session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected nl_session cookie to be set")
	}
}

func TestHandleLoginWrongPassword(t *testing.T) {
	app := newTestApp(t)
	testutil.CreateUser(t, app.app, "loginbad@example.com", "password123")

	body, ct := formBody(map[string]string{"email": "loginbad@example.com", "password": "wrongpass"})
	res := app.do(t, http.MethodPost, "/login", body, map[string]string{"content-type": ct})
	expectStatus(t, res, 401)
	expectContains(t, res, "Invalid email or password")
}

func TestDashboardRequiresAuth(t *testing.T) {
	app := newTestApp(t)
	res := app.do(t, http.MethodGet, "/", nil, nil)
	expectStatus(t, res, http.StatusSeeOther)
	if loc := res.Header.Get("Location"); loc != "/login" {
		t.Errorf("expected redirect to /login, got %q", loc)
	}
}

func TestDashboardAuthenticated(t *testing.T) {
	app := newTestApp(t)
	user := testutil.CreateUser(t, app.app, "dash@example.com", "password123")

	res := app.do(t, http.MethodGet, "/", nil, cookieHeader(t, user))
	expectStatus(t, res, 200)
	expectContains(t, res, "Your groups")
}

func TestHandleCreateGroup(t *testing.T) {
	app := newTestApp(t)
	user := testutil.CreateUser(t, app.app, "groupowner@example.com", "password123")

	body, ct := formBody(map[string]string{"name": "My Crew"})
	res := app.do(t, http.MethodPost, "/groups", body,
		mergeHeaders(cookieHeader(t, user), map[string]string{"content-type": ct}))
	expectStatus(t, res, http.StatusSeeOther)

	groups, err := app.app.FindRecordsByFilter("groups", "owner = {:owner}", "", 0, 0, map[string]any{"owner": user.Id})
	if err != nil || len(groups) != 1 {
		t.Fatalf("expected 1 group owned by user, got %d (err=%v)", len(groups), err)
	}
	if groups[0].GetString("name") != "My Crew" {
		t.Errorf("expected name %q, got %q", "My Crew", groups[0].GetString("name"))
	}

	memberships, err := app.app.FindRecordsByFilter("group_memberships",
		"group = {:group} && user = {:user} && role = \"owner\"", "", 0, 0,
		map[string]any{"group": groups[0].Id, "user": user.Id})
	if err != nil || len(memberships) != 1 {
		t.Fatalf("expected owner membership to be created, got %d (err=%v)", len(memberships), err)
	}
}

func TestGroupHomeForbiddenForNonMember(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "ghowner@example.com", "password123")
	outsider := testutil.CreateUser(t, app.app, "outsider@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Closed Crew", "closed-crew", owner.Id)

	res := app.do(t, http.MethodGet, "/g/"+group.GetString("slug"), nil, cookieHeader(t, outsider))
	expectStatus(t, res, 403)
	expectContains(t, res, "Not a member of this group.")

	res = app.do(t, http.MethodGet, "/g/"+group.GetString("slug"), nil, cookieHeader(t, owner))
	expectStatus(t, res, 200)
	expectContains(t, res, "Closed Crew")
}

func TestHandleCreateInvitePermissions(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "invowner@example.com", "password123")
	member := testutil.CreateUser(t, app.app, "invmember@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Invite Crew", "invite-crew", owner.Id)
	testutil.CreateMembership(t, app.app, group.Id, member.Id, "member")

	body, ct := formBody(map[string]string{"email": "newperson@example.com", "role": "member"})
	res := app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/invites", body,
		mergeHeaders(cookieHeader(t, member), map[string]string{"content-type": ct}))
	expectStatus(t, res, 403)

	body2, ct2 := formBody(map[string]string{"email": "newperson@example.com", "role": "member"})
	res = app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/invites", body2,
		mergeHeaders(cookieHeader(t, owner), map[string]string{"content-type": ct2}))
	expectStatus(t, res, http.StatusSeeOther)

	invites, err := app.app.FindRecordsByFilter("group_invites", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(invites) != 1 {
		t.Fatalf("expected 1 invite to be created, got %d (err=%v)", len(invites), err)
	}
	if invites[0].GetString("email") != "newperson@example.com" {
		t.Errorf("unexpected invite email %q", invites[0].GetString("email"))
	}
}

func TestInviteAcceptFlowExistingUser(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "ia-owner@example.com", "password123")
	invitee := testutil.CreateUser(t, app.app, "ia-invitee@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Accept Crew", "accept-crew", owner.Id)

	invitesCol, err := app.app.FindCollectionByNameOrId("group_invites")
	if err != nil {
		t.Fatal(err)
	}
	invite := newRecordHelper(t, app.app, invitesCol, map[string]any{
		"group":        group.Id,
		"invited_by":   owner.Id,
		"email":        invitee.Email(),
		"invited_user": invitee.Id,
		"token":        "accept-test-token",
		"role":         "member",
		"status":       "pending",
		"expires_at":   "2999-01-01 00:00:00.000Z",
	})

	res := app.do(t, http.MethodGet, "/invites/"+invite.GetString("token"), nil, cookieHeader(t, invitee))
	expectStatus(t, res, 200)
	expectContains(t, res, "Accept Crew")

	res = app.do(t, http.MethodPost, "/invites/"+invite.GetString("token")+"/accept", nil, cookieHeader(t, invitee))
	expectStatus(t, res, http.StatusSeeOther)

	memberships, err := app.app.FindRecordsByFilter("group_memberships",
		"group = {:group} && user = {:user}", "", 0, 0,
		map[string]any{"group": group.Id, "user": invitee.Id})
	if err != nil || len(memberships) != 1 {
		t.Fatalf("expected membership to be created on accept, got %d (err=%v)", len(memberships), err)
	}

	invite, err = app.app.FindRecordById("group_invites", invite.Id)
	if err != nil {
		t.Fatal(err)
	}
	if invite.GetString("status") != "accepted" {
		t.Errorf("expected invite status=accepted, got %q", invite.GetString("status"))
	}
}

func TestHandleCreateAndToggleQuestion(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "qowner@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Question Crew", "question-crew", owner.Id)

	body, ct := formBody(map[string]string{"prompt": "Favorite snack?", "type": "text"})
	res := app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/questions", body,
		mergeHeaders(cookieHeader(t, owner), map[string]string{"content-type": ct}))
	expectStatus(t, res, http.StatusSeeOther)

	qs, err := app.app.FindRecordsByFilter("question_bank", "group = {:group} && scope = \"group\"", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(qs) != 1 {
		t.Fatalf("expected 1 custom question, got %d (err=%v)", len(qs), err)
	}
	if !qs[0].GetBool("is_active") {
		t.Fatal("expected new question to be active")
	}

	res = app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/questions/"+qs[0].Id+"/toggle", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	toggled, err := app.app.FindRecordById("question_bank", qs[0].Id)
	if err != nil {
		t.Fatal(err)
	}
	if toggled.GetBool("is_active") {
		t.Error("expected question to be deactivated after toggle")
	}
}

func TestEditionCreateAnswerCloseView(t *testing.T) {
	app := newTestApp(t)
	owner := testutil.CreateUser(t, app.app, "edowner@example.com", "password123")
	group := testutil.CreateGroup(t, app.app, "Edition Crew", "edition-crew", owner.Id)

	res := app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/editions", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	editions, err := app.app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(editions) != 1 {
		t.Fatalf("expected 1 edition, got %d (err=%v)", len(editions), err)
	}
	edition := editions[0]
	if edition.GetString("status") != "open" {
		t.Fatalf("expected status=open, got %q", edition.GetString("status"))
	}

	allEqs, err := app.app.FindRecordsByFilter("edition_questions", "edition = {:edition}", "", 0, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(allEqs) == 0 {
		t.Fatalf("expected edition_questions to be populated, got %d (err=%v)", len(allEqs), err)
	}

	var textQuestionID string
	for _, eq := range allEqs {
		q, err := app.app.FindRecordById("question_bank", eq.GetString("question"))
		if err == nil && q.GetString("type") == "text" {
			textQuestionID = q.Id
			break
		}
	}
	if textQuestionID == "" {
		t.Fatal("expected at least one seeded text question")
	}

	mpBody, mpCT := multipartBody(t, map[string]string{"q_" + textQuestionID: "Pretzels"})
	res = app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/editions/"+edition.Id, mpBody,
		mergeHeaders(cookieHeader(t, owner), map[string]string{"content-type": mpCT}))
	expectStatus(t, res, http.StatusSeeOther)

	ans, err := app.app.FindFirstRecordByFilter("answers",
		"edition = {:edition} && question = {:question} && user = {:user}",
		map[string]any{"edition": edition.Id, "question": textQuestionID, "user": owner.Id})
	if err != nil {
		t.Fatalf("expected answer to be saved: %v", err)
	}
	if ans.GetBool("skipped") {
		t.Error("expected answered text question to not be skipped")
	}

	res = app.do(t, http.MethodPost, "/g/"+group.GetString("slug")+"/editions/"+edition.Id+"/close", nil, cookieHeader(t, owner))
	expectStatus(t, res, http.StatusSeeOther)

	closed, err := app.app.FindRecordById("newsletter_editions", edition.Id)
	if err != nil {
		t.Fatal(err)
	}
	if closed.GetString("status") != "sent" {
		t.Fatalf("expected status=sent after close, got %q", closed.GetString("status"))
	}

	res = app.do(t, http.MethodGet, "/g/"+group.GetString("slug")+"/editions/"+edition.Id+"/view", nil, cookieHeader(t, owner))
	expectStatus(t, res, 200)
	expectContains(t, res, "Pretzels")
}
