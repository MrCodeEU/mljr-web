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
