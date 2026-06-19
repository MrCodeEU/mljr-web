# agents.md — `projects/newsletter`

Friend-group recurring newsletter. Members join groups, each group runs a
recurring "edition" (weekly/biweekly/monthly/quarterly) of questions that
members answer (text/single-select/multi-select/image/rating/emoji), with
reminders, a grace period, and a compiled send. See the root [`agents.md`](../../agents.md)
for repo-wide conventions (gomponents, Datastar, CSS, CI). **This project does
not follow those Echo-based patterns** — read this file before touching
anything here.

## Why this project looks different from the rest of the repo

Every other `projects/*` is a small Echo handler set rendering gomponents.
This one embeds **PocketBase** (`github.com/pocketbase/pocketbase`) as the
actual HTTP server/process, because it's the first project in the repo that
needs real persistent storage, user accounts, file uploads, and a scheduled
job — and PocketBase gives all four for free (SQLite, auth, `FileField`,
`app.Cron()`) instead of hand-rolling them on Echo.

Consequences:
- No Echo anywhere in this directory. Routes are PocketBase
  `core.RequestEvent` handlers, bound via `e.Router.GET/POST(...)` inside an
  `app.OnServe().BindFunc(...)` hook.
- Rendering uses `internal/web.RenderPB(e, status, node)` (in
  `internal/web/pocketbase.go`), the PocketBase-flavored sibling of the
  Echo-targeting `Render` used elsewhere. `pages/render.go`'s `renderPage`
  wraps it.
- Auth, password hashing, and sessions are PocketBase's built-in `users`
  collection — no custom crypto in this codebase.
- Still reuses the shared `ui/` component library (`ui/form`, `ui/primitive`,
  `ui/layout`, `ui/overlay`, `ui/special`, `ui/feedback`, `ui/icon`,
  `ui/token`) exactly like every other project — the `class=`-ban and other
  root `agents.md` guard rules still apply to any `ui/**.go` files touched
  from here.

## Running it

```
make dev-newsletter        # port 8096 (config.Load() defaults port "8090"->8096
                            # for this binary in main.go when PORT is unset)
```

Config knobs (`internal/config/config.go`, `NewsletterConfig`):
- `NEWSLETTER_DATA_DIR` (default `pb_data`) — PocketBase's SQLite + file
  storage directory. **This is the first project in the repo with a real
  persistent-data requirement** — a prod deploy needs a mounted volume here
  or every redeploy loses all groups/editions/answers/images. Not yet
  coordinated with homelab-automation (see Known gaps).
- `NEWSLETTER_PUBLIC_URL` (default `http://localhost:8095`) — base URL used
  to build links inside emails (invite links, edition links). Must match
  whatever the deployed origin actually is, or emailed links will be wrong.
- SMTP comes from the repo-wide `SMTP_*` env vars (shared with the homepage
  contact form), bootstrapped into PocketBase's own `app.Settings().SMTP` on
  boot by `scheduler.BootstrapMailer` — see Mail below.

PocketBase's own admin UI is available at `/_/` on the same port (create the
first superuser via `go run . superuser create <email> <password>` from this
directory, or through the CLI prompt on first run) — useful for inspecting
collections/records directly and for manually editing a group's schedule
fields to force the scanner to act sooner during manual testing.

## Data model (`migrations/`, one file per schema slice)

Defined as Go migrations registered via `m.Register(up, down)` — there is no
hand-written SQL schema file; the migration source *is* the schema, and the
SQLite file in `pb_data/` is derived/runtime state (gitignored, not source of
truth).

| File | Collections added |
|---|---|
| `1700000000_init.go` | `users` (+`timezone` field on the built-in collection), `groups`, `group_memberships` |
| `1700000001_invites_notifications.go` | `group_invites`, `notifications` |
| `1700000002_questions_editions.go` | `question_bank`, `newsletter_editions`, `edition_questions`, `answers`, `answer_images` (+ seeds ~12 global questions, one/two per type) |
| `1700000003_email_log.go` | `email_log` |

Key fields per collection (see the migration files for the authoritative
field list/types/API rules — this is a summary, not a substitute):

- **`groups`**: `name`, `slug` (unique, `^[a-z0-9-]+$`), `owner` (relation),
  `schedule_period` (weekly/biweekly/monthly/quarterly),
  `schedule_anchor_weekday` (0=Sun..6=Sat), `schedule_anchor_day_of_month`,
  `schedule_epoch_date` (biweekly parity anchor), `schedule_send_hour_utc`,
  `reminder_lead_hours`, `grace_period_hours`, `timezone`. All schedule math
  is precomputed once per edition into stored timestamps — the scanner only
  ever compares `now` against `opens_at`/`reminder_at`/`grace_until`, never
  recomputes the group's recurrence rule on every tick.
- **`group_memberships`**: `group`, `user`, `role` (owner/admin/member),
  unique on (`group`,`user`).
- **`group_invites`**: `group`, `invited_by`, `email`, `invited_user`
  (nullable — only set if the email matched an existing account at creation
  time), `token` (unique), `role` (admin/member), `status`
  (pending/accepted/expired/revoked), `expires_at`.
- **`notifications`**: `user`, `kind` (invite/group_joined/group_left/
  edition_open/edition_reminder/edition_sent/comment_reply/emoji_reaction),
  `group`, `invite`, `actor`, `body`, `link`, `read_at`.
- **`question_bank`**: `scope` (global/group/user), `group`, `author`,
  `type` (text/single_select/multi_select/image/rating/emoji_reaction),
  `prompt`, `options` (JSONField, `[]string`, used by select/emoji types
  only), `is_active`.
- **`newsletter_editions`**: `group`, `opens_at`, `closes_at`, `reminder_at`,
  `grace_until`, `status` (scheduled→open→reminder_sent→grace→sent→archived),
  `sent_at`.
- **`edition_questions`**: `edition`, `question`, `order`, `vote_count`
  (unused until Phase 5 voting lands), unique on (`edition`,`question`).
- **`answers`**: `edition`, `question`, `user`, `value` (JSONField, shape
  depends on `question.type`), `skipped`, unique on
  (`edition`,`question`,`user`) — one row per user per question per edition,
  upserted by `editions.go`'s `upsertAnswer`.
- **`answer_images`**: `answer`, `image` (FileField, 5MB max,
  jpeg/png/webp/gif), `order` — image answers are N rows here, not a
  `MaxSelect` tuning knob on a single field, so multi-image answers are just
  multiple rows.
- **`email_log`**: `kind` (invite/reminder/edition_sent), `dedupe_key`
  (unique — `"reminder:{editionID}:{userID}"` / `"send:{editionID}"`),
  `recipient_email`, `status` (sent/failed), `error`. No API rules — only
  ever touched from backend code (the cron scanner), never exposed over
  PocketBase's REST API.

API access rules on group-scoped collections gate via membership filters
like `group.group_memberships_via_group.user ?= @request.auth.id`, with
admin-only mutation rules layered on top where relevant (see each
migration's `*Rule` assignments). These REST rules exist for completeness /
defense-in-depth, but in practice almost everything goes through this
project's own handlers in `pages/`, which re-check membership/role
server-side rather than relying solely on the collection rules — see
Authorization patterns below.

## JSONField decode gotcha (read this before touching `answers`/`question_bank`/`options`)

`core.Record.Get(key)` on a `JSONField` column returns the **raw**
`types.JSONRaw` bytes, not a decoded Go value. Every read site must
`json.Unmarshal` manually:

```go
// pages/questions.go
func questionOptions(q *core.Record) []string {
    raw, ok := q.Get("options").(types.JSONRaw)
    if !ok || len(raw) == 0 {
        return nil
    }
    var opts []string
    json.Unmarshal(raw, &opts)
    return opts
}
```

`pages/editions.go`'s `answerValue`/`valueAsString` do the same for
`answers.value`, with an extra wrinkle: ratings round-trip through JSON as a
bare number, so `valueAsString` must handle both `string` and `float64` —
a plain `.(string)` assertion silently drops every rating answer. The
scheduler has its own copy of this pattern (`decodeJSONField` in
`scheduler/scheduler.go`) since it can't import `pages` (would create an
import cycle); don't try to deduplicate the two without checking that first.

## Routes (`pages/routes.go`, `RegisterRoutes(e *core.ServeEvent)`)

| Method/Path | Handler | Notes |
|---|---|---|
| `GET/POST /login` | `pages.Login` / `HandleLogin` | sets `nl_session` cookie |
| `GET/POST /signup` | `pages.Signup` / `HandleSignup` | supports `?invite=` token |
| `POST /logout` | `pages.HandleLogout` | |
| `GET/POST /profile` | `pages.Profile` / `HandleProfile` | |
| `GET /` | `pages.Dashboard` | lists user's groups + create-group form |
| `GET /g/{slug}` | `pages.GroupHome` | member list, 403 if not a member |
| `GET/POST /g/{slug}/settings` | `pages.GroupSettings` / `HandleGroupSettings` | owner/admin only |
| `POST /groups` | `pages.HandleCreateGroup` | creates group + owner membership |
| `GET/POST /g/{slug}/invites` | `pages.ListInvites` / `HandleCreateInvite` | owner/admin only |
| `POST /g/{slug}/invites/{id}/revoke` | `pages.HandleRevokeInvite` | owner/admin only |
| `GET /invites/{token}` | `pages.InviteAccept` | works logged-out (shows signup/login CTAs) |
| `POST /invites/{token}/accept` | `pages.HandleAcceptInvite` | requires login; validates invite targets the logged-in user |
| `POST /notifications/read-all` | `pages.HandleMarkAllNotificationsRead` | redirects back to `Referer` |
| `GET/POST /g/{slug}/questions` | `pages.ListQuestions` / `HandleCreateQuestion` | any member can list/add; only admins can toggle |
| `POST /g/{slug}/questions/{id}/toggle` | `pages.HandleToggleQuestion` | owner/admin only |
| `GET/POST /g/{slug}/editions` | `pages.ListEditions` / `HandleCreateEdition` | create is owner/admin only, no-ops if an edition is already open |
| `GET/POST /g/{slug}/editions/{id}` | `pages.EditionAnswer` / `HandleSubmitAnswers` | multipart form (image uploads); 400 if edition isn't `open` |
| `POST /g/{slug}/editions/{id}/close` | `pages.HandleCloseEdition` | manual close path, owner/admin only — the real path is the scheduler |
| `GET /g/{slug}/editions/{id}/view` | `pages.EditionView` | only for `sent`/`archived` editions |
| `GET /healthz` | inline 200 | liveness probe |

`pages.RegisterHooks(app)` (in `pages/hooks.go`) is bound separately from
routes (both called from `main.go`, both also called from the test harness —
see Testing below) and fires `group_joined`/`group_left` notifications to
group admins from `OnRecordAfterCreateSuccess`/`AfterDeleteSuccess` on
`group_memberships`.

## Authorization patterns

There's no middleware layer — every handler that needs auth starts with:

```go
user := currentUser(re)
if user == nil {
    return redirect(re, "/login")
}
```

Group-scoped handlers then look up the group by slug and check membership:

```go
group, err := findGroupBySlug(re, slug)
if err != nil { return re.NotFoundError("group not found", err) }
if _, err := findMembership(re, group.Id, user.Id); err != nil {
    return re.ForbiddenError("not a member of this group", nil)
}
```

Admin-gated actions use the shared `requireAdminMembership(re, group, user)`
(`pages/invites.go`) instead of repeating the owner/admin role check inline.

**Two historical IDOR fixes worth knowing about** (already patched, but the
pattern matters for any new invite/signup-adjacent code): `inviteTargetsUser`
(`pages/invites.go`) must be checked before honoring an accept — both
`HandleAcceptInvite` and the signup-via-invite path in `HandleSignup` call it
— otherwise any logged-in user could accept an invite token meant for
someone else's email, or a new signup could self-assign membership via a
guessed/leaked token regardless of which email it was actually sent to.
`findEditionInGroup` similarly re-validates `edition.group == group.Id`
rather than trusting the path param alone, since `newsletter_editions` IDs
are global, not scoped per group in the URL.

## Mail (`scheduler/mailer.go`)

`Mailer` is a small interface (`Send(*pbmailer.Message) error`).
`BootstrapMailer(app, cfg.SMTP)` writes the repo's `.env`-sourced
`SMTPConfig` into PocketBase's own `app.Settings().SMTP` on every boot (so
PocketBase's native `app.NewMailClient()` picks it up), and returns:
- `pbMailer{app}` wrapping `app.NewMailClient()` if `cfg.SMTP.Host != ""`.
- `logMailer{}` (just logs the message, dev fallback) if SMTP host is unset
  — this is why local dev works without real credentials, and why `make
  dev-newsletter` is silent about email unless you set `SMTP_*` in `.env`.

`pages.SetMailer(mailerClient, cfg)` (called once from `main.go`'s
`OnServe`) stashes the mailer + config as package-level state so handlers
like `sendInviteEmail` (`pages/invites.go`) can send without needing it
threaded through every call. Tests don't go through `BootstrapMailer` at all
— they inject `&tests.TestMailer{}` directly into `scheduler.RunScan`.

## Scheduling (`scheduler/`)

`app.Cron().MustAdd("newsletter_scan", "*/5 * * * *", func() {
scheduler.RunScan(app, mailerClient, cfg) })` — registered once in
`main.go`'s `OnServe`. `RunScan` runs five sub-steps every tick, **each one
logs and continues past its own error rather than aborting the rest** (a
problem with one group/edition must not block the whole scan):

1. `createDueEditions` — for every group with no currently-open edition,
   compute the next window via `anchor.NextWindow` and create one if it's
   due; populates `edition_questions` from the active global+group question
   bank (`populateEditionQuestions`).
2. `openScheduledEditions` — flips `scheduled→open` once `opens_at` passes.
3. `sendReminders` — flips `open→reminder_sent` once `reminder_at` passes;
   emails members who haven't submitted *any* answer yet for that edition
   (one reminder per member, deduped via `email_log.dedupe_key =
   "reminder:{editionID}:{userID}"`).
4. `closeEditions` — once `grace_until` passes: marks every still-missing
   answer `skipped=true` (`markMissingAnswersSkipped`), compiles and sends
   the combined edition email (`dedupe_key = "send:{editionID}"`),
   stamps `sent_at`, flips status to `sent`.
5. `expireInvites` — flips stale `pending` `group_invites` past
   `expires_at` to `expired`.

**Idempotency**: `sendOnce` (in `scheduler.go`) checks `email_log` for an
existing row with the target `dedupe_key` *before* building/sending, then
always writes a result row afterward (`status`: sent/failed + `error`).
This is an application-level dedupe check against `email_log`, not a
database-level unique-index rejection race — acceptable at this scan
cadence (5 min, single process) but would need locking if the scanner ever
ran concurrently or more frequently.

`scheduler/anchor.go` holds the pure date math —
`NextWindow(period, anchorWeekday, anchorDayOfMonth, epochDate, sendHourUTC,
tz, after)` — including day-of-month clamping (day 31 in Feb → last valid
day) and biweekly parity via `epochDate`. It has no DB dependency, which is
why it has its own fast unit test file separate from the DB-backed scheduler
tests.

**Known gap**: `compileEditionText` renders image-type answers as a literal
`"[shared a photo]"` text placeholder in the digest email — there's no
actual image attachment/embedding in the sent email yet. Relevant if Phase 5
recap features are expected to show images in an email context, not just on
the web view.

## Slugs & tokens (`pages/slug.go`)

- `slugify(name)` + `randomSuffix()` (4 random bytes hex) → group slugs, so
  two groups named "Family" don't collide.
- `randomToken()` (24 random bytes hex) → invite tokens
  (`group_invites.token`). Long enough that brute-forcing one isn't
  practical; treated as a bearer credential (whoever has the link can act as
  the invited account up to the IDOR checks described above).

## Session (`pages/session.go`)

Cookie name `nl_session`, set via `user.NewAuthToken()`
(PocketBase-native auth token, not custom JWT/session logic),
`HttpOnly`, `SameSite=Lax`, 30-day expiry. `currentUser(re)` resolves it via
`re.App.FindAuthRecordByToken(cookie.Value, core.TokenTypeAuth)`. There is no
session table — validity is whatever `FindAuthRecordByToken` accepts
(PocketBase handles token expiry/revocation internally).

## Testing — the one non-obvious lesson worth preserving

`tests.ApiScenario` (PocketBase's own test helper) **does not work for
multi-request flows against one shared `*tests.TestApp`**. `apis.NewRouter(app)`
internally calls `app.OnServe().BindFunc(...)` to wire its own extensions
(e.g. `/_/extensions.js`); calling it more than once against the same app —
which `ApiScenario.Test()` does every time it's invoked — accumulates
duplicate `OnServe`-bound route registrations and panics on the second
request with something like `pattern GET /_/extensions.js conflicts with
pattern GET /_/extensions.js` (this is a known upstream issue, see
PocketBase's own `discussions/7267`). Any test that needs >1 HTTP call
against the same app/db state (e.g. "create group, then invite, then accept,
then assert membership") hits this immediately.

**Fix used here, in `pages/pages_test.go`**: build the router once
(`apis.NewRouter`), fire `OnServe()` exactly once via a manually-constructed
`core.ServeEvent{App, Router}`, build the `http.Handler` mux exactly once via
`baseRouter.BuildMux()`, then reuse that single mux across as many
`httptest.NewRequest`/`mux.ServeHTTP` calls as the test needs:

```go
type httpApp struct {
    app *tests.TestApp
    mux http.Handler
}

func newTestApp(t *testing.T) *httpApp {
    app := testutil.NewApp(t)
    baseRouter, _ := apis.NewRouter(app)
    app.OnServe().Trigger(&core.ServeEvent{App: app, Router: baseRouter},
        func(e *core.ServeEvent) error { return RegisterRoutes(e) })
    RegisterHooks(app)
    mux, _ := baseRouter.BuildMux()
    return &httpApp{app: app, mux: mux}
}
```

Any new HTTP test file in `pages/` should reuse `newTestApp`/`(*httpApp).do`
from `pages_test.go` rather than reaching for `tests.ApiScenario` directly.
Single-request smoke tests *could* use `ApiScenario` safely, but there's no
reason to mix two patterns in one package — stick with the `httpApp` harness
for consistency.

Fixtures live in `internal/testutil/testutil.go`: `NewApp` (boots a fresh
app with this project's migrations applied via blank import), `CreateUser`,
`AuthCookie` (mints an `nl_session=...` cookie value), `CreateGroup` (applies
the same weekly/Friday/18:00/48h/24h defaults as
`HandleCreateGroup`), `CreateMembership`.

Current coverage: `pages/pages_http_test.go` (10 tests, every `pages/*.go`
handler that's reachable gets at least one case — login, dashboard auth
gate, group create, group-home 403-for-non-member, invite create/permission,
invite accept (existing user), question create+toggle, edition
create→answer→close→view); `scheduler/scheduler_test.go` (full lifecycle
scheduled→open→reminder_sent→sent, reminder-skips-members-who-answered,
invite expiry) + `scheduler/anchor_test.go` (pure date-math cases). All
green as of Phase 4.5; run with `go test ./projects/newsletter/...`.

## Known gaps / explicitly deferred

- **Phase 5 (not started)**: `question_suggestions` + 
  `question_suggestion_votes` (voting on next-edition questions — note
  `edition_questions.vote_count` already exists in the schema but is unused
  until this lands), past-answer recap, `emoji_reactions` (toggle
  reactions on an answer), `comments` (flat + single-level reply via a
  self-relation `parent` field), and an edition archive page beyond the
  current bare `ListEditions` list.
- **Phase 6 (not started)**: this project is **not yet in
  `.github/workflows/docker.yml`'s `ALL` build matrix** — CI doesn't build
  or publish a newsletter image yet. Also not yet coordinated with
  homelab-automation for a persistent volume on `NEWSLETTER_DATA_DIR` (see
  Running it above) — deploying without that volume loses all data on every
  redeploy.
- Image answers have no email-digest representation (placeholder text only,
  see Scheduling above).
- No avatar upload yet — `pages/shell.go`'s `avatarTone`/`initials` generate
  a deterministic color+initials avatar from user id/name instead.
