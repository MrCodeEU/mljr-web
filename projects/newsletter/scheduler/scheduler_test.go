package scheduler

import (
	"testing"
	"time"

	"mljr-web/internal/config"
	"mljr-web/projects/newsletter/internal/testutil"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

func testConfig() config.Config {
	return config.Config{
		SMTP:       config.SMTPConfig{From: "Newsletter <noreply@example.com>"},
		Newsletter: config.NewsletterConfig{PublicAppURL: "http://localhost:8096"},
	}
}

func editionStatus(t *testing.T, app core.App, id string) string {
	t.Helper()
	rec, err := app.FindRecordById("newsletter_editions", id)
	if err != nil {
		t.Fatalf("find edition: %v", err)
	}
	return rec.GetString("status")
}

func openedAtForWindow(t *testing.T, group *core.Record, want func(time.Time, time.Time) bool) time.Time {
	t.Helper()
	gracePeriod := time.Duration(group.GetInt("grace_period_hours")) * time.Hour
	for daysAgo := 0; daysAgo < 60; daysAgo++ {
		openedAt := time.Now().UTC().Add(-time.Duration(daysAgo) * 24 * time.Hour)
		closesAt, err := closeTimeForOpenEdition(group, openedAt)
		if err != nil {
			t.Fatal(err)
		}
		if want(closesAt, closesAt.Add(gracePeriod)) {
			return openedAt
		}
	}
	t.Fatal("could not find opened_at matching requested scheduler window")
	return time.Time{}
}

// TestRunScanFullLifecycle drives a single edition through every state
// transition (scheduled -> open -> reminder_sent -> sent) by precomputing its
// timestamps in the past, the same way a real cron tick would once a group's
// window has passed. It mirrors the manual scratch-instance walkthrough used
// to validate Phase 4 before this suite existed.
func TestRunScanFullLifecycle(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "owner@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Test Crew", "test-crew", owner.Id)

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (create): %v", err)
	}

	editions, err := app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(editions) != 1 {
		t.Fatalf("expected exactly 1 edition after first scan, got %d (err=%v)", len(editions), err)
	}
	edition := editions[0]
	if got := edition.GetString("status"); got != "scheduled" {
		t.Fatalf("expected status=scheduled, got %q", got)
	}

	// edition_questions should be populated from the seeded global question bank.
	eqs, err := app.FindRecordsByFilter("edition_questions", "edition = {:edition}", "", 0, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(eqs) == 0 {
		t.Fatalf("expected edition_questions to be populated, got %d (err=%v)", len(eqs), err)
	}
	questionCount := len(eqs)

	now := time.Now().UTC()
	edition.Set("opens_at", now.Add(-time.Hour))
	if err := app.Save(edition); err != nil {
		t.Fatalf("save edition opens_at: %v", err)
	}
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (open): %v", err)
	}
	if got := editionStatus(t, app, edition.Id); got != "open" {
		t.Fatalf("expected status=open, got %q", got)
	}

	edition, err = app.FindRecordById("newsletter_editions", edition.Id)
	if err != nil {
		t.Fatal(err)
	}
	edition.Set("reminder_at", now.Add(-time.Hour))
	if err := app.Save(edition); err != nil {
		t.Fatalf("save edition reminder_at: %v", err)
	}
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (reminder): %v", err)
	}
	if got := editionStatus(t, app, edition.Id); got != "reminder_sent" {
		t.Fatalf("expected status=reminder_sent, got %q", got)
	}
	if mailer.TotalSend() != 1 {
		t.Fatalf("expected 1 reminder email, got %d", mailer.TotalSend())
	}
	if to := mailer.LastMessage().To; len(to) != 1 || to[0].Address != owner.Email() {
		t.Fatalf("expected reminder to %s, got %v", owner.Email(), to)
	}

	logs, err := app.FindRecordsByFilter("email_log", "kind = \"reminder\"", "", 0, 0, nil)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected 1 reminder email_log row, got %d (err=%v)", len(logs), err)
	}
	wantKey := "reminder:" + edition.Id + ":" + owner.Id
	if logs[0].GetString("dedupe_key") != wantKey {
		t.Fatalf("expected dedupe_key=%q, got %q", wantKey, logs[0].GetString("dedupe_key"))
	}

	// Re-running the scan before grace_until passes must not double-send.
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (no-op): %v", err)
	}
	if mailer.TotalSend() != 1 {
		t.Fatalf("expected reminder send count to stay at 1, got %d", mailer.TotalSend())
	}

	edition, err = app.FindRecordById("newsletter_editions", edition.Id)
	if err != nil {
		t.Fatal(err)
	}
	edition.Set("closes_at", now.Add(-2*time.Hour))
	edition.Set("grace_until", now.Add(-time.Hour))
	if err := app.Save(edition); err != nil {
		t.Fatalf("save edition grace_until: %v", err)
	}
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (close): %v", err)
	}
	if got := editionStatus(t, app, edition.Id); got != "sent" {
		t.Fatalf("expected status=sent, got %q", got)
	}
	if mailer.TotalSend() != 2 {
		t.Fatalf("expected 2 total emails (reminder + edition_sent), got %d", mailer.TotalSend())
	}

	skipped, err := app.FindRecordsByFilter("answers", "edition = {:edition} && skipped = true", "", 0, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(skipped) != questionCount {
		t.Fatalf("expected %d skipped answers (one per question, owner answered none), got %d (err=%v)",
			questionCount, len(skipped), err)
	}

	notifs, err := app.FindRecordsByFilter("notifications", "kind = \"edition_sent\"", "", 0, 0, nil)
	if err != nil || len(notifs) != 1 {
		t.Fatalf("expected 1 edition_sent notification, got %d (err=%v)", len(notifs), err)
	}

	// Idempotency: re-running after sent must not re-send or re-mark skipped.
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (post-sent): %v", err)
	}
	if mailer.TotalSend() != 2 {
		t.Fatalf("expected send count to stay at 2 after a post-sent scan, got %d", mailer.TotalSend())
	}
}

// TestRunScanSkipsMembersWhoAnswered ensures the reminder pass only emails
// members who haven't submitted any answer yet for the open edition.
func TestRunScanSkipsMembersWhoAnswered(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "owner2@example.com", "password123")
	answered := testutil.CreateUser(t, app, "answered@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Crew Two", "crew-two", owner.Id)
	testutil.CreateMembership(t, app, group.Id, answered.Id, "member")

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (create): %v", err)
	}
	editions, err := app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(editions) != 1 {
		t.Fatalf("expected 1 edition, got %d (err=%v)", len(editions), err)
	}
	edition := editions[0]

	eqs, err := app.FindRecordsByFilter("edition_questions", "edition = {:edition}", "", 1, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(eqs) == 0 {
		t.Fatalf("expected at least 1 edition_question, got %d (err=%v)", len(eqs), err)
	}
	answersCol, err := app.FindCollectionByNameOrId("answers")
	if err != nil {
		t.Fatal(err)
	}
	ans := core.NewRecord(answersCol)
	ans.Set("edition", edition.Id)
	ans.Set("question", eqs[0].GetString("question"))
	ans.Set("user", answered.Id)
	ans.Set("value", `"yes"`)
	if err := app.Save(ans); err != nil {
		t.Fatalf("save answer: %v", err)
	}

	now := time.Now().UTC()
	edition.Set("opens_at", now.Add(-time.Hour))
	edition.Set("reminder_at", now.Add(-time.Hour))
	if err := app.Save(edition); err != nil {
		t.Fatal(err)
	}
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (reminder): %v", err)
	}

	if mailer.TotalSend() != 1 {
		t.Fatalf("expected exactly 1 reminder (to the member who hasn't answered), got %d", mailer.TotalSend())
	}
	if to := mailer.LastMessage().To; len(to) != 1 || to[0].Address != owner.Email() {
		t.Fatalf("expected reminder to go to owner (no answer yet), got %v", to)
	}
}

// TestRunScanGraceReminderOnlyTargetsIncompleteMembers checks the
// grace-period last-call email: once an edition is past its soft closes_at
// but still inside grace_until, only members who haven't answered every
// question should get a (single, deduped) grace reminder.
func TestRunScanGraceReminderOnlyTargetsIncompleteMembers(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "owner4@example.com", "password123")
	complete := testutil.CreateUser(t, app, "complete@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Crew Four", "crew-four", owner.Id)
	testutil.CreateMembership(t, app, group.Id, complete.Id, "member")

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (create): %v", err)
	}
	editions, err := app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(editions) != 1 {
		t.Fatalf("expected 1 edition, got %d (err=%v)", len(editions), err)
	}
	edition := editions[0]

	eqs, err := app.FindRecordsByFilter("edition_questions", "edition = {:edition}", "", 0, 0,
		map[string]any{"edition": edition.Id})
	if err != nil || len(eqs) == 0 {
		t.Fatalf("expected edition_questions, got %d (err=%v)", len(eqs), err)
	}
	answersCol, err := app.FindCollectionByNameOrId("answers")
	if err != nil {
		t.Fatal(err)
	}
	// `complete` answers every question; `owner` answers none.
	for _, eq := range eqs {
		ans := core.NewRecord(answersCol)
		ans.Set("edition", edition.Id)
		ans.Set("question", eq.GetString("question"))
		ans.Set("user", complete.Id)
		ans.Set("value", `"yes"`)
		ans.Set("skipped", false)
		if err := app.Save(ans); err != nil {
			t.Fatalf("save answer: %v", err)
		}
	}

	now := time.Now().UTC()
	edition.Set("opens_at", now.Add(-3*time.Hour))
	edition.Set("reminder_at", now.Add(-2*time.Hour))
	edition.Set("closes_at", now.Add(-time.Hour))
	edition.Set("grace_until", now.Add(time.Hour))
	if err := app.Save(edition); err != nil {
		t.Fatal(err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (grace): %v", err)
	}
	if got := editionStatus(t, app, edition.Id); got != "grace" {
		t.Fatalf("expected status=grace, got %q", got)
	}

	logs, err := app.FindRecordsByFilter("email_log", "kind = \"grace_reminder\"", "", 0, 0, nil)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected exactly 1 grace_reminder email_log row, got %d (err=%v)", len(logs), err)
	}
	wantKey := "grace_reminder:" + edition.Id + ":" + owner.Id
	if logs[0].GetString("dedupe_key") != wantKey {
		t.Fatalf("expected grace reminder to go to owner (incomplete), got dedupe_key=%q", logs[0].GetString("dedupe_key"))
	}

	sendsAfterFirst := mailer.TotalSend()
	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (grace no-op): %v", err)
	}
	if mailer.TotalSend() != sendsAfterFirst {
		t.Fatalf("expected grace reminder send count to stay at %d, got %d", sendsAfterFirst, mailer.TotalSend())
	}
}

func TestRunScanBackfillsManualOpenEditionAndSendsGraceReminder(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "manual-grace-owner@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Manual Grace", "manual-grace", owner.Id)
	now := time.Now().UTC()
	group.Set("schedule_anchor_weekday", int(now.Add(24*time.Hour).Weekday()))
	group.Set("schedule_send_hour_utc", now.Hour())
	if err := app.Save(group); err != nil {
		t.Fatal(err)
	}
	openedAt := openedAtForWindow(t, group, func(closesAt, graceUntil time.Time) bool {
		return closesAt.Before(now) && graceUntil.After(now)
	})

	editionsCol, err := app.FindCollectionByNameOrId("newsletter_editions")
	if err != nil {
		t.Fatal(err)
	}
	edition := core.NewRecord(editionsCol)
	edition.Set("group", group.Id)
	edition.Set("status", "open")
	edition.Set("opens_at", openedAt)
	if err := app.Save(edition); err != nil {
		t.Fatalf("save manual edition: %v", err)
	}
	if err := populateEditionQuestions(app, edition); err != nil {
		t.Fatalf("populate questions: %v", err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan: %v", err)
	}

	edition, err = app.FindRecordById("newsletter_editions", edition.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got := edition.GetString("status"); got != "grace" {
		t.Fatalf("expected manual edition to move into grace, got %q", got)
	}
	if edition.GetString("closes_at") == "" || edition.GetString("reminder_at") == "" || edition.GetString("grace_until") == "" {
		t.Fatalf("expected scheduler to backfill deadlines, got closes_at=%q reminder_at=%q grace_until=%q",
			edition.GetString("closes_at"), edition.GetString("reminder_at"), edition.GetString("grace_until"))
	}

	logs, err := app.FindRecordsByFilter("email_log", "kind = \"grace_reminder\"", "", 0, 0, nil)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected one grace_reminder email_log row, got %d (err=%v)", len(logs), err)
	}
}

func TestRunScanBackfillsManualOpenEditionAndClosesPastGrace(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "manual-close-owner@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Manual Close", "manual-close", owner.Id)
	now := time.Now().UTC()
	openedAt := openedAtForWindow(t, group, func(closesAt, graceUntil time.Time) bool {
		return closesAt.Before(now) && graceUntil.Before(now)
	})

	editionsCol, err := app.FindCollectionByNameOrId("newsletter_editions")
	if err != nil {
		t.Fatal(err)
	}
	edition := core.NewRecord(editionsCol)
	edition.Set("group", group.Id)
	edition.Set("status", "open")
	edition.Set("opens_at", openedAt)
	if err := app.Save(edition); err != nil {
		t.Fatalf("save manual edition: %v", err)
	}
	if err := populateEditionQuestions(app, edition); err != nil {
		t.Fatalf("populate questions: %v", err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan: %v", err)
	}
	if got := editionStatus(t, app, edition.Id); got != "sent" {
		t.Fatalf("expected past-grace manual edition to be sent, got %q", got)
	}

	logs, err := app.FindRecordsByFilter("email_log", "kind = \"edition_sent\"", "", 0, 0, nil)
	if err != nil || len(logs) != 1 {
		t.Fatalf("expected one edition_sent email_log row, got %d (err=%v)", len(logs), err)
	}
}

// TestRunScanExpiresStaleInvites checks the cron-side invite expiry sweep.
func TestRunScanExpiresStaleInvites(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "owner3@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Crew Three", "crew-three", owner.Id)

	col, err := app.FindCollectionByNameOrId("group_invites")
	if err != nil {
		t.Fatal(err)
	}
	invite := core.NewRecord(col)
	invite.Set("group", group.Id)
	invite.Set("invited_by", owner.Id)
	invite.Set("email", "late@example.com")
	invite.Set("token", "test-token-123")
	invite.Set("role", "member")
	invite.Set("status", "pending")
	invite.Set("expires_at", time.Now().UTC().Add(-time.Hour))
	if err := app.Save(invite); err != nil {
		t.Fatalf("save invite: %v", err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan: %v", err)
	}

	invite, err = app.FindRecordById("group_invites", invite.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got := invite.GetString("status"); got != "expired" {
		t.Fatalf("expected invite status=expired, got %q", got)
	}
}

// TestRunScanDisablesGroupAfterThreeUnansweredEditions checks the
// zero-respondent streak: closing 3 fully-unanswered editions for the same
// group in one tick must flip groups.status to "disabled", and a further
// tick must not create a 4th edition for it.
func TestRunScanDisablesGroupAfterThreeUnansweredEditions(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "ghost-owner@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Ghost Crew", "ghost-crew", owner.Id)

	editionsCol, err := app.FindCollectionByNameOrId("newsletter_editions")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	for i := range disableAfterUnanswered {
		edition := core.NewRecord(editionsCol)
		edition.Set("group", group.Id)
		edition.Set("status", "open")
		edition.Set("opens_at", now.Add(-time.Duration(i+2)*time.Hour))
		edition.Set("closes_at", now.Add(-time.Hour))
		edition.Set("grace_until", now.Add(-time.Minute))
		if err := app.Save(edition); err != nil {
			t.Fatalf("save edition %d: %v", i, err)
		}
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (close all 3): %v", err)
	}

	group, err = app.FindRecordById("groups", group.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got := group.GetInt("consecutive_unanswered_editions"); got != disableAfterUnanswered {
		t.Fatalf("expected streak=%d, got %d", disableAfterUnanswered, got)
	}
	if got := group.GetString("status"); got != "disabled" {
		t.Fatalf("expected group status=disabled, got %q", got)
	}

	sentEditions, err := app.FindRecordsByFilter("newsletter_editions", "group = {:group} && status = \"sent\"", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(sentEditions) != disableAfterUnanswered {
		t.Fatalf("expected %d sent editions, got %d (err=%v)", disableAfterUnanswered, len(sentEditions), err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan (post-disable): %v", err)
	}
	allEditions, err := app.FindRecordsByFilter("newsletter_editions", "group = {:group}", "", 0, 0,
		map[string]any{"group": group.Id})
	if err != nil || len(allEditions) != disableAfterUnanswered {
		t.Fatalf("expected no new edition created for disabled group, still want %d, got %d (err=%v)",
			disableAfterUnanswered, len(allEditions), err)
	}
}

// TestRunScanResetsStreakWhenSomeoneAnswers ensures a single answered
// edition resets the zero-respondent streak instead of letting it accumulate
// across unrelated editions.
func TestRunScanResetsStreakWhenSomeoneAnswers(t *testing.T) {
	app := testutil.NewApp(t)
	cfg := testConfig()
	mailer := &tests.TestMailer{}

	owner := testutil.CreateUser(t, app, "streak-owner@example.com", "password123")
	group := testutil.CreateGroup(t, app, "Streak Crew", "streak-crew", owner.Id)
	group.Set("consecutive_unanswered_editions", 2)
	if err := app.Save(group); err != nil {
		t.Fatal(err)
	}

	editionsCol, err := app.FindCollectionByNameOrId("newsletter_editions")
	if err != nil {
		t.Fatal(err)
	}
	answersCol, err := app.FindCollectionByNameOrId("answers")
	if err != nil {
		t.Fatal(err)
	}
	questionsCol, err := app.FindCollectionByNameOrId("question_bank")
	if err != nil {
		t.Fatal(err)
	}
	question := core.NewRecord(questionsCol)
	question.Set("scope", "global")
	question.Set("type", "text")
	question.Set("prompt", "How was your week?")
	question.Set("is_active", true)
	if err := app.Save(question); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	edition := core.NewRecord(editionsCol)
	edition.Set("group", group.Id)
	edition.Set("status", "open")
	edition.Set("opens_at", now.Add(-2*time.Hour))
	edition.Set("closes_at", now.Add(-time.Hour))
	edition.Set("grace_until", now.Add(-time.Minute))
	if err := app.Save(edition); err != nil {
		t.Fatal(err)
	}
	answer := core.NewRecord(answersCol)
	answer.Set("edition", edition.Id)
	answer.Set("question", question.Id)
	answer.Set("user", owner.Id)
	answer.Set("value", `"good"`)
	answer.Set("skipped", false)
	if err := app.Save(answer); err != nil {
		t.Fatal(err)
	}

	if err := RunScan(app, mailer, cfg); err != nil {
		t.Fatalf("RunScan: %v", err)
	}

	group, err = app.FindRecordById("groups", group.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got := group.GetInt("consecutive_unanswered_editions"); got != 0 {
		t.Fatalf("expected streak reset to 0 after an answered edition, got %d", got)
	}
	if got := group.GetString("status"); got != "active" {
		t.Fatalf("expected group to stay active, got %q", got)
	}
}
