package scheduler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"mljr-web/internal/config"

	"github.com/pocketbase/pocketbase/core"
	pbmailer "github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/types"
)

// RunScan drives every group's edition state machine one tick. It is safe to
// call repeatedly (e.g. every 5 minutes via app.Cron()) — every email send is
// guarded by a unique dedupe_key in email_log, and every status transition
// re-queries by status so a half-finished tick just gets picked up again.
func RunScan(app core.App, m Mailer, cfg config.Config) error {
	now := time.Now().UTC()

	if err := createDueEditions(app); err != nil {
		log.Printf("newsletter scan: create editions: %v", err)
	}
	if err := openScheduledEditions(app, now); err != nil {
		log.Printf("newsletter scan: open editions: %v", err)
	}
	if err := sendReminders(app, m, cfg, now); err != nil {
		log.Printf("newsletter scan: reminders: %v", err)
	}
	if err := closeEditions(app, m, cfg, now); err != nil {
		log.Printf("newsletter scan: close editions: %v", err)
	}
	if err := expireInvites(app, now); err != nil {
		log.Printf("newsletter scan: expire invites: %v", err)
	}
	return nil
}

func createDueEditions(app core.App) error {
	groups, err := app.FindRecordsByFilter("groups", "", "", 0, 0, nil)
	if err != nil {
		return err
	}

	for _, group := range groups {
		filter := "group = {:group} && (" +
			"status = \"scheduled\" || status = \"open\" || status = \"reminder_sent\" || status = \"grace\")"
		if _, err := app.FindFirstRecordByFilter("newsletter_editions", filter, map[string]any{"group": group.Id}); err == nil {
			continue // already has an edition in flight
		}

		after := group.GetDateTime("created").Time()
		if last, lerr := app.FindRecordsByFilter(
			"newsletter_editions", "group = {:group}", "-opens_at", 1, 0, map[string]any{"group": group.Id},
		); lerr == nil && len(last) > 0 {
			after = last[0].GetDateTime("opens_at").Time()
		}

		opensAt, err := nextWindowForGroup(group, after)
		if err != nil {
			log.Printf("newsletter scan: group %s: %v", group.Id, err)
			continue
		}
		closesAt, err := nextWindowForGroup(group, opensAt)
		if err != nil {
			log.Printf("newsletter scan: group %s: %v", group.Id, err)
			continue
		}
		closesAt = closesAt.Add(-24 * time.Hour)
		reminderLead := time.Duration(group.GetInt("reminder_lead_hours")) * time.Hour
		gracePeriod := time.Duration(group.GetInt("grace_period_hours")) * time.Hour

		col, err := app.FindCollectionByNameOrId("newsletter_editions")
		if err != nil {
			return err
		}
		edition := core.NewRecord(col)
		edition.Set("group", group.Id)
		edition.Set("status", "scheduled")
		edition.Set("opens_at", opensAt)
		edition.Set("closes_at", closesAt)
		edition.Set("reminder_at", closesAt.Add(-reminderLead))
		edition.Set("grace_until", closesAt.Add(gracePeriod))
		if err := app.Save(edition); err != nil {
			return err
		}

		if err := populateEditionQuestions(app, edition); err != nil {
			return err
		}
	}
	return nil
}

func nextWindowForGroup(group *core.Record, after time.Time) (time.Time, error) {
	epoch := group.GetDateTime("schedule_epoch_date").Time()
	return NextWindow(
		group.GetString("schedule_period"),
		group.GetInt("schedule_anchor_weekday"),
		group.GetInt("schedule_anchor_day_of_month"),
		epoch,
		group.GetInt("schedule_send_hour_utc"),
		group.GetString("timezone"),
		after,
	)
}

func populateEditionQuestions(app core.App, edition *core.Record) error {
	questions, err := app.FindRecordsByFilter(
		"question_bank", "is_active = true && (scope = \"global\" || (scope = \"group\" && group = {:group}))", "created", 0, 0,
		map[string]any{"group": edition.GetString("group")},
	)
	if err != nil {
		return err
	}
	eqCol, err := app.FindCollectionByNameOrId("edition_questions")
	if err != nil {
		return err
	}
	for i, q := range questions {
		eq := core.NewRecord(eqCol)
		eq.Set("edition", edition.Id)
		eq.Set("question", q.Id)
		eq.Set("order", i)
		if err := app.Save(eq); err != nil {
			return err
		}
	}
	return nil
}

func openScheduledEditions(app core.App, now time.Time) error {
	editions, err := app.FindRecordsByFilter(
		"newsletter_editions", "status = \"scheduled\" && opens_at <= {:now}", "", 0, 0,
		map[string]any{"now": now},
	)
	if err != nil {
		return err
	}
	for _, edition := range editions {
		edition.Set("status", "open")
		if err := app.Save(edition); err != nil {
			return err
		}
		group, err := app.FindRecordById("groups", edition.GetString("group"))
		if err != nil {
			continue
		}
		notifyMembers(app, group.Id, "", "edition_open", "",
			"A new edition is open for "+group.GetString("name"),
			"/g/"+group.GetString("slug")+"/editions/"+edition.Id)
	}
	return nil
}

func sendReminders(app core.App, m Mailer, cfg config.Config, now time.Time) error {
	editions, err := app.FindRecordsByFilter(
		"newsletter_editions", "status = \"open\" && reminder_at != \"\" && reminder_at <= {:now}", "", 0, 0,
		map[string]any{"now": now},
	)
	if err != nil {
		return err
	}
	for _, edition := range editions {
		group, err := app.FindRecordById("groups", edition.GetString("group"))
		if err != nil {
			continue
		}
		members, err := app.FindRecordsByFilter(
			"group_memberships", "group = {:group}", "", 0, 0, map[string]any{"group": group.Id},
		)
		if err != nil {
			continue
		}
		for _, membership := range members {
			userID := membership.GetString("user")
			answered, err := app.FindFirstRecordByFilter(
				"answers", "edition = {:edition} && user = {:user}",
				map[string]any{"edition": edition.Id, "user": userID},
			)
			if err == nil && answered != nil {
				continue // already started/submitted
			}
			user, err := app.FindRecordById("users", userID)
			if err != nil || user.Email() == "" {
				continue
			}
			dedupeKey := "reminder:" + edition.Id + ":" + userID
			sendOnce(app, m, dedupeKey, "reminder", user.Email(), func() (*pbmailer.Message, error) {
				return reminderMessage(cfg, group, edition, user), nil
			})
		}
		edition.Set("status", "reminder_sent")
		if err := app.Save(edition); err != nil {
			return err
		}
	}
	return nil
}

func closeEditions(app core.App, m Mailer, cfg config.Config, now time.Time) error {
	editions, err := app.FindRecordsByFilter(
		"newsletter_editions",
		"(status = \"open\" || status = \"reminder_sent\") && grace_until != \"\" && grace_until <= {:now}",
		"", 0, 0, map[string]any{"now": now},
	)
	if err != nil {
		return err
	}
	for _, edition := range editions {
		group, err := app.FindRecordById("groups", edition.GetString("group"))
		if err != nil {
			continue
		}
		if err := markMissingAnswersSkipped(app, edition); err != nil {
			log.Printf("newsletter scan: mark skipped for edition %s: %v", edition.Id, err)
		}

		dedupeKey := "send:" + edition.Id
		recipients, err := memberEmails(app, group.Id)
		if err == nil && len(recipients) > 0 {
			sendOnce(app, m, dedupeKey, "edition_sent", strings.Join(recipients, ","), func() (*pbmailer.Message, error) {
				return editionSentMessage(app, cfg, group, edition, recipients)
			})
		}

		edition.Set("status", "sent")
		edition.Set("sent_at", now)
		if err := app.Save(edition); err != nil {
			return err
		}
		notifyMembers(app, group.Id, "", "edition_sent", "",
			group.GetString("name")+"'s edition is ready to read",
			"/g/"+group.GetString("slug")+"/editions/"+edition.Id+"/view")
	}
	return nil
}

func markMissingAnswersSkipped(app core.App, edition *core.Record) error {
	eqs, err := app.FindRecordsByFilter(
		"edition_questions", "edition = {:edition}", "", 0, 0, map[string]any{"edition": edition.Id},
	)
	if err != nil {
		return err
	}
	members, err := app.FindRecordsByFilter(
		"group_memberships", "group = {:group}", "", 0, 0, map[string]any{"group": edition.GetString("group")},
	)
	if err != nil {
		return err
	}
	answersCol, err := app.FindCollectionByNameOrId("answers")
	if err != nil {
		return err
	}
	for _, eq := range eqs {
		questionID := eq.GetString("question")
		for _, membership := range members {
			userID := membership.GetString("user")
			if _, err := app.FindFirstRecordByFilter(
				"answers", "edition = {:edition} && question = {:question} && user = {:user}",
				map[string]any{"edition": edition.Id, "question": questionID, "user": userID},
			); err == nil {
				continue
			}
			answer := core.NewRecord(answersCol)
			answer.Set("edition", edition.Id)
			answer.Set("question", questionID)
			answer.Set("user", userID)
			answer.Set("skipped", true)
			if err := app.Save(answer); err != nil {
				return err
			}
		}
	}
	return nil
}

func expireInvites(app core.App, now time.Time) error {
	invites, err := app.FindRecordsByFilter(
		"group_invites", "status = \"pending\" && expires_at <= {:now}", "", 0, 0,
		map[string]any{"now": now},
	)
	if err != nil {
		return err
	}
	for _, invite := range invites {
		invite.Set("status", "expired")
		if err := app.Save(invite); err != nil {
			return err
		}
	}
	return nil
}

// notifyMembers creates an in-app notification for every member of a group
// except excludeUserID (pass "" to notify everyone).
func notifyMembers(app core.App, groupID, excludeUserID, kind, actorID, body, link string) {
	col, err := app.FindCollectionByNameOrId("notifications")
	if err != nil {
		return
	}
	members, err := app.FindRecordsByFilter(
		"group_memberships", "group = {:group}", "", 0, 0, map[string]any{"group": groupID},
	)
	if err != nil {
		return
	}
	for _, membership := range members {
		userID := membership.GetString("user")
		if userID == excludeUserID {
			continue
		}
		n := core.NewRecord(col)
		n.Set("user", userID)
		n.Set("kind", kind)
		n.Set("group", groupID)
		if actorID != "" {
			n.Set("actor", actorID)
		}
		n.Set("body", body)
		n.Set("link", link)
		_ = app.Save(n)
	}
}

func memberEmails(app core.App, groupID string) ([]string, error) {
	members, err := app.FindRecordsByFilter(
		"group_memberships", "group = {:group}", "", 0, 0, map[string]any{"group": groupID},
	)
	if err != nil {
		return nil, err
	}
	var emails []string
	for _, membership := range members {
		user, err := app.FindRecordById("users", membership.GetString("user"))
		if err != nil || user.Email() == "" {
			continue
		}
		emails = append(emails, user.Email())
	}
	return emails, nil
}

// sendOnce checks email_log for dedupeKey before building/sending, and logs
// the outcome (success or failure) so a retried cron tick never double-sends.
func sendOnce(app core.App, m Mailer, dedupeKey, kind, recipient string, build func() (*pbmailer.Message, error)) {
	if _, err := app.FindFirstRecordByFilter(
		"email_log", "dedupe_key = {:key}", map[string]any{"key": dedupeKey},
	); err == nil {
		return // already sent
	}

	msg, err := build()
	status, errText := "sent", ""
	if err != nil {
		status, errText = "failed", err.Error()
	} else if sendErr := m.Send(msg); sendErr != nil {
		status, errText = "failed", sendErr.Error()
	}

	col, cErr := app.FindCollectionByNameOrId("email_log")
	if cErr != nil {
		log.Printf("newsletter mail: missing email_log collection: %v", cErr)
		return
	}
	logRec := core.NewRecord(col)
	logRec.Set("kind", kind)
	logRec.Set("dedupe_key", dedupeKey)
	logRec.Set("recipient_email", recipient)
	logRec.Set("status", status)
	logRec.Set("error", errText)
	if err := app.Save(logRec); err != nil {
		log.Printf("newsletter mail: failed to write email_log: %v", err)
	}
	if status == "failed" {
		log.Printf("newsletter mail: send failed kind=%s dedupe=%s: %s", kind, dedupeKey, errText)
	}
}

func reminderMessage(cfg config.Config, group, edition, user *core.Record) *pbmailer.Message {
	link := cfg.Newsletter.PublicAppURL + "/g/" + group.GetString("slug") + "/editions/" + edition.Id
	body := fmt.Sprintf(
		"Hi %s,\n\n%s's newsletter edition is still open and you haven't answered yet.\n\n%s\n",
		displayNameOrEmail(user), group.GetString("name"), link,
	)
	return &pbmailer.Message{
		From:    fromAddress(cfg.SMTP.From),
		To:      []mail.Address{{Address: user.Email()}},
		Subject: group.GetString("name") + ": don't forget to answer this week's questions",
		Text:    body,
	}
}

func editionSentMessage(app core.App, cfg config.Config, group, edition *core.Record, recipients []string) (*pbmailer.Message, error) {
	body, err := compileEditionText(app, group, edition)
	if err != nil {
		return nil, err
	}
	to := make([]mail.Address, len(recipients))
	for i, addr := range recipients {
		to[i] = mail.Address{Address: addr}
	}
	return &pbmailer.Message{
		From:    fromAddress(cfg.SMTP.From),
		Bcc:     to,
		Subject: group.GetString("name") + " newsletter — " + edition.GetString("opens_at")[:10],
		Text:    body,
	}, nil
}

func compileEditionText(app core.App, group, edition *core.Record) (string, error) {
	eqs, err := app.FindRecordsByFilter(
		"edition_questions", "edition = {:edition}", "order", 0, 0, map[string]any{"edition": edition.Id},
	)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s newsletter\n\n", group.GetString("name"))

	for _, eq := range eqs {
		question, err := app.FindRecordById("question_bank", eq.GetString("question"))
		if err != nil {
			continue
		}
		fmt.Fprintf(&b, "%s\n", question.GetString("prompt"))

		answers, err := app.FindRecordsByFilter(
			"answers", "edition = {:edition} && question = {:question} && skipped = false", "", 0, 0,
			map[string]any{"edition": edition.Id, "question": question.Id},
		)
		if err != nil {
			continue
		}
		for _, answer := range answers {
			user, err := app.FindRecordById("users", answer.GetString("user"))
			if err != nil {
				continue
			}
			text := valueAsText(decodeJSONField(answer, "value"))
			if question.GetString("type") == "image" {
				text = "[shared a photo]"
			}
			if text == "" {
				continue
			}
			fmt.Fprintf(&b, "  - %s: %s\n", displayNameOrEmail(user), text)
		}
		b.WriteString("\n")
	}

	return b.String(), nil
}

func displayNameOrEmail(user *core.Record) string {
	if name := user.GetString("name"); name != "" {
		return name
	}
	return user.Email()
}

// decodeJSONField mirrors pages.answerValue: Record.Get on a JSONField
// returns the raw types.JSONRaw bytes, not the decoded value.
func decodeJSONField(rec *core.Record, key string) any {
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

func valueAsText(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []any:
		parts := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, ", ")
	}
	return ""
}
