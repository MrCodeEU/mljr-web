// Package testutil provides shared fixtures for projects/newsletter's test
// suite (scheduler and pages packages), so each new feature's tests can reuse
// the same app-bootstrap/user/group helpers instead of re-deriving them.
package testutil

import (
	"testing"

	_ "mljr-web/projects/newsletter/migrations"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
)

// NewApp boots a fresh PocketBase test app with the newsletter's own
// migrations applied (via the blank import above) and registers automatic
// cleanup of its temp data dir via tb.Cleanup.
func NewApp(tb testing.TB) *tests.TestApp {
	tb.Helper()
	app, err := tests.NewTestApp()
	if err != nil {
		tb.Fatalf("testutil: failed to create test app: %v", err)
	}
	tb.Cleanup(app.Cleanup)
	return app
}

// CreateUser creates a verified users-collection record.
func CreateUser(tb testing.TB, app core.App, email, password string) *core.Record {
	tb.Helper()
	col, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		tb.Fatalf("testutil: find users collection: %v", err)
	}
	user := core.NewRecord(col)
	user.SetEmail(email)
	user.SetPassword(password)
	user.SetVerified(true)
	if err := app.Save(user); err != nil {
		tb.Fatalf("testutil: save user: %v", err)
	}
	return user
}

// AuthCookie mints a session cookie value (matching pages.setSession's
// "nl_session" cookie) for the given user, for use as an ApiScenario
// "Cookie" header so tests can act as an already-logged-in user.
func AuthCookie(tb testing.TB, user *core.Record) string {
	tb.Helper()
	token, err := user.NewAuthToken()
	if err != nil {
		tb.Fatalf("testutil: mint auth token: %v", err)
	}
	return "nl_session=" + token
}

// CreateGroup creates a group with the same schedule defaults
// pages.HandleCreateGroup applies (weekly, Friday 18:00 UTC, 48h reminder
// lead, 24h grace, UTC tz) plus an owner membership for ownerID.
func CreateGroup(tb testing.TB, app core.App, name, slug, ownerID string) *core.Record {
	tb.Helper()
	col, err := app.FindCollectionByNameOrId("groups")
	if err != nil {
		tb.Fatalf("testutil: find groups collection: %v", err)
	}
	group := core.NewRecord(col)
	group.Set("name", name)
	group.Set("slug", slug)
	group.Set("owner", ownerID)
	group.Set("schedule_period", "weekly")
	group.Set("schedule_anchor_weekday", 5)
	group.Set("schedule_send_hour_utc", 18)
	group.Set("reminder_lead_hours", 48)
	group.Set("grace_period_hours", 24)
	group.Set("timezone", "UTC")
	if err := app.Save(group); err != nil {
		tb.Fatalf("testutil: save group: %v", err)
	}
	CreateMembership(tb, app, group.Id, ownerID, "owner")
	return group
}

// CreateMembership creates a group_memberships row directly (bypassing the
// notification hooks that real joins go through, since most tests only need
// the membership relation to exist).
func CreateMembership(tb testing.TB, app core.App, groupID, userID, role string) *core.Record {
	tb.Helper()
	col, err := app.FindCollectionByNameOrId("group_memberships")
	if err != nil {
		tb.Fatalf("testutil: find group_memberships collection: %v", err)
	}
	m := core.NewRecord(col)
	m.Set("group", groupID)
	m.Set("user", userID)
	m.Set("role", role)
	if err := app.Save(m); err != nil {
		tb.Fatalf("testutil: save membership: %v", err)
	}
	return m
}
