package pages

import (
	"fmt"
	"strconv"
	"strings"

	"mljr-web/ui"
	"mljr-web/ui/feedback"
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// currentEdition returns the group's open edition, or else its most
// recently created one (sent/archived/scheduled), or nil if none exists yet.
func currentEdition(re *core.RequestEvent, groupID string) *core.Record {
	if open, err := re.App.FindFirstRecordByFilter(
		"newsletter_editions", "group = {:group} && status = \"open\"",
		map[string]any{"group": groupID},
	); err == nil {
		return open
	}
	editions, err := re.App.FindRecordsByFilter(
		"newsletter_editions", "group = {:group}", "-created", 1, 0,
		map[string]any{"group": groupID},
	)
	if err != nil || len(editions) == 0 {
		return nil
	}
	return editions[0]
}

// editionHeroCard renders the group's current edition as a large status
// card with a countdown to its hard deadline (grace_until, falling back to
// closes_at) when one is set — manually-created editions never get either
// field populated, so it degrades to an "admin closes manually" message.
// editionRespondentCount counts the distinct members who've answered at
// least one (non-skipped) question on an edition.
func editionRespondentCount(re *core.RequestEvent, editionID string) int {
	answers, err := re.App.FindRecordsByFilter(
		"answers", "edition = {:edition} && skipped = false", "", 0, 0,
		map[string]any{"edition": editionID},
	)
	if err != nil {
		return 0
	}
	users := map[string]bool{}
	for _, a := range answers {
		users[a.GetString("user")] = true
	}
	return len(users)
}

func editionHeroCard(re *core.RequestEvent, slug string, isAdmin bool, ed *core.Record, totalMembers int) g.Node {
	t := translator(re)
	if ed == nil {
		var cta g.Node
		if isAdmin {
			cta = h.Form(h.Method("post"), h.Action("/g/"+slug+"/editions"), h.Style("margin-top:var(--sp-4)"),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit"}, g.Text(t("newsletter.groups.start_edition"))),
			)
		}
		return primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-bottom:var(--sp-6)")}},
			h.P(h.Style("color:var(--muted)"), g.Text(t("newsletter.groups.no_edition_yet"))),
			cta,
		)
	}

	status := ed.GetString("status")
	eqs, _ := editionQuestions(re, ed.Id)

	respondentCount := 0
	var closeDeadline, graceDeadline g.Node
	if status == "open" {
		respondentCount = editionRespondentCount(re, ed.Id)
		allAnswered := totalMembers > 0 && respondentCount >= totalMembers

		closesAt := ed.GetString("closes_at")
		graceUntil := ed.GetString("grace_until")

		if closesAt != "" {
			closeDeadline = h.Div(
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text(t("newsletter.groups.answers_close_in"))),
				primitive.Countdown(primitive.CountdownProps{Target: closesAt, ID: "edition-hero-cd-close"}),
			)
		}
		// The grace period is the actual hard cutoff the scheduler enforces
		// (closeEditions watches grace_until, not closes_at) — but once
		// every member has answered there's nothing left to wait for, so
		// showing it would just be a confusing second timer.
		if graceUntil != "" && graceUntil != closesAt && !allAnswered {
			graceDeadline = h.Div(h.Style("margin-top:var(--sp-3)"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm)"), g.Text(t("newsletter.groups.final_cutoff_in"))),
				primitive.Countdown(primitive.CountdownProps{Target: graceUntil, ID: "edition-hero-cd-grace"}),
			)
		}
		if closesAt == "" && graceUntil == "" {
			closeDeadline = h.P(h.Style("color:var(--muted)"), g.Text(t("newsletter.groups.open_no_deadline")))
		}
	}

	href := "/g/" + slug + "/editions/" + ed.Id
	ctaLabel := t("newsletter.groups.answer_now")
	if status == "sent" || status == "archived" {
		href += "/view"
		ctaLabel = t("newsletter.groups.view_results")
	}

	return primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-bottom:var(--sp-6);padding:var(--sp-6)")}},
		h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:baseline"),
			primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(t("newsletter.groups.current_edition"))),
			h.Span(h.Style("color:var(--muted);font-size:var(--t-sm);text-transform:uppercase;letter-spacing:.04em"), g.Text(status)),
		),
		h.P(h.Style("color:var(--muted);margin-top:var(--sp-2)"), g.Text(t("newsletter.groups.question_count", len(eqs)))),
		g.If(status == "open", h.P(h.Style("color:var(--muted);margin-top:var(--sp-1)"),
			g.Text(t("newsletter.groups.answered_count", respondentCount, totalMembers)),
		)),
		g.If(closeDeadline != nil, h.Div(h.Style("margin-top:var(--sp-4)"), closeDeadline)),
		graceDeadline,
		h.Div(h.Style("margin-top:var(--sp-5)"),
			h.A(
				g.Attr("data-component", "button"),
				g.Attr("data-variant", string(token.Primary)),
				g.Attr("data-size", string(token.SizeMD)),
				h.Href(href), h.Style("text-decoration:none;display:inline-block"),
				g.Text(ctaLabel),
			),
		),
	)
}

// HandleCreateGroup creates a group owned by the current user and adds them
// as an "owner" member. Schedule fields default to a sane preset (weekly,
// Friday, 18:00 UTC) — editable later from group settings.
func HandleCreateGroup(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	name := strings.TrimSpace(re.Request.FormValue("name"))
	if name == "" {
		return redirect(re, "/?flash=group_name_required")
	}

	groupsCol, err := re.App.FindCollectionByNameOrId("groups")
	if err != nil {
		return err
	}
	membershipsCol, err := re.App.FindCollectionByNameOrId("group_memberships")
	if err != nil {
		return err
	}

	group := core.NewRecord(groupsCol)
	group.Set("name", name)
	group.Set("slug", slugify(name)+"-"+randomSuffix())
	group.Set("owner", user.Id)
	group.Set("schedule_period", "weekly")
	group.Set("schedule_anchor_weekday", 5) // Friday
	group.Set("schedule_send_hour_utc", 18)
	group.Set("reminder_lead_hours", 48)
	group.Set("grace_period_hours", 24)
	group.Set("timezone", "UTC")
	group.Set("status", "active")
	if err := re.App.Save(group); err != nil {
		return err
	}

	membership := core.NewRecord(membershipsCol)
	membership.Set("group", group.Id)
	membership.Set("user", user.Id)
	membership.Set("role", "owner")
	if err := re.App.Save(membership); err != nil {
		return err
	}

	return redirect(re, "/g/"+group.GetString("slug"))
}

func findGroupBySlug(re *core.RequestEvent, slug string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter("groups", "slug = {:slug}", map[string]any{"slug": slug})
}

func findMembership(re *core.RequestEvent, groupID, userID string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter(
		"group_memberships", "group = {:group} && user = {:user}",
		map[string]any{"group": groupID, "user": userID},
	)
}

// GroupHome shows group info and the member list.
func GroupHome(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	t := translator(re)
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

	members, err := re.App.FindRecordsByFilter(
		"group_memberships", "group = {:group}", "created", 0, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	var memberRows []g.Node
	for _, m := range members {
		u, err := re.App.FindRecordById("users", m.GetString("user"))
		if err != nil {
			continue
		}
		memberRows = append(memberRows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Span(h.Style("min-width:0;overflow-wrap:anywhere"), g.Text(displayName(u))),
			h.Span(h.Style("color:var(--muted);font-size:var(--t-sm);text-transform:uppercase;letter-spacing:.04em;white-space:nowrap"), g.Text(m.GetString("role"))),
		))
	}

	return renderPage(re, 200, appPage(re, slug, group.GetString("name"),
		[]breadcrumbItem{{Label: t("newsletter.nav.dashboard"), Href: "/"}, {Label: group.GetString("name")}},
		primitive.Heading(primitive.HeadingProps{Level: 1, Attrs: []g.Node{h.Style("margin-bottom:var(--sp-6)")}}, g.Text(group.GetString("name"))),
		editionHeroCard(re, slug, isAdmin, currentEdition(re, group.Id), len(members)),
		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(t("newsletter.groups.members_heading"))),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-3);padding:var(--sp-2) var(--sp-4)")}},
			g.Group(memberRows),
		),
	))
}

// GroupSettings shows the schedule + member-list editor (owner/admin only).
func GroupSettings(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	t := translator(re)
	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	membership, err := findMembership(re, group.Id, user.Id)
	if err != nil || (membership.GetString("role") != "owner" && membership.GetString("role") != "admin") {
		return re.ForbiddenError("only group owners/admins can edit settings", nil)
	}

	currentPeriod := group.GetString("schedule_period")
	periodOption := func(value, label string) form.SelectOption {
		return form.SelectOption{Value: value, Label: label, Selected: value == currentPeriod}
	}

	isOwner := group.GetString("owner") == user.Id

	var disabledBanner g.Node
	if group.GetString("status") == "disabled" {
		var reactivateBtn g.Node
		if isOwner {
			reactivateBtn = h.Form(h.Method("post"), h.Action("/g/"+slug+"/reactivate"), h.Style("margin-top:var(--sp-3)"),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit"}, g.Text(t("newsletter.groups.reactivate_button"))),
			)
		}
		disabledBanner = feedback.Alert(
			feedback.AlertProps{Variant: feedback.AlertWarning, Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.P(g.Text(t("newsletter.groups.auto_disabled_notice", group.GetInt("consecutive_unanswered_editions")))),
			reactivateBtn,
		)
	}

	var leaveSection g.Node
	if !isOwner {
		leaveSection = primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-6);border-color:var(--danger)")}},
			primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(t("newsletter.groups.leave_heading"))),
			h.P(h.Style("color:var(--muted);margin:var(--sp-2) 0 var(--sp-4)"), g.Text(t("newsletter.groups.leave_body"))),
			h.Form(h.Method("post"), h.Action("/g/"+slug+"/leave"),
				primitive.Button(primitive.ButtonProps{Variant: token.Danger, Type: "submit"}, g.Text(t("newsletter.groups.leave_button"))),
			),
		)
	}

	return renderPage(re, 200, appPage(re, slug, t("newsletter.subnav.settings")+" — "+group.GetString("name"),
		[]breadcrumbItem{{Label: t("newsletter.nav.dashboard"), Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: t("newsletter.subnav.settings")}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(t("newsletter.groups.settings_heading", group.GetString("name")))),
		flashAlert(re),
		disabledBanner,
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/g/"+slug+"/settings"),
				ui.Signals(fmt.Sprintf(`{reminderLead:%d,gracePeriod:%d}`,
					group.GetInt("reminder_lead_hours"), group.GetInt("grace_period_hours"))),
				form.Field(form.FieldProps{Label: t("newsletter.groups.name_label")},
					form.Input(form.InputProps{Type: "text", Name: "name", Required: true, Attrs: []g.Node{h.Value(group.GetString("name"))}}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.groups.schedule_period_label"), Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Select(form.SelectProps{
						Name: "schedule_period",
						Options: []form.SelectOption{
							periodOption("weekly", t("newsletter.groups.period_weekly")),
							periodOption("biweekly", t("newsletter.groups.period_biweekly")),
							periodOption("monthly", t("newsletter.groups.period_monthly")),
							periodOption("quarterly", t("newsletter.groups.period_quarterly")),
						},
					}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.groups.reminder_lead_label"), Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.NumberInput(form.NumberInputProps{
						Signal: "reminderLead", Name: "reminder_lead_hours",
						Min: 0, Max: 336, Step: 1, Value: group.GetInt("reminder_lead_hours"),
					}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.groups.grace_period_label"), Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.NumberInput(form.NumberInputProps{
						Signal: "gracePeriod", Name: "grace_period_hours",
						Min: 0, Max: 336, Step: 1, Value: group.GetInt("grace_period_hours"),
					}),
				),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin-top:var(--sp-2)"),
					g.Text(t("newsletter.groups.settings_note"))),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text(t("newsletter.groups.save_button"))),
			),
		),
		leaveSection,
	))
}

func HandleGroupSettings(re *core.RequestEvent) error {
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
	if err != nil || (membership.GetString("role") != "owner" && membership.GetString("role") != "admin") {
		return re.ForbiddenError("only group owners/admins can edit settings", nil)
	}

	name := strings.TrimSpace(re.Request.FormValue("name"))
	if name == "" {
		return redirect(re, "/g/"+slug+"/settings?flash=group_name_required")
	}
	group.Set("name", name)
	group.Set("schedule_period", re.Request.FormValue("schedule_period"))
	if v, err := strconv.Atoi(re.Request.FormValue("reminder_lead_hours")); err == nil && v >= 0 {
		group.Set("reminder_lead_hours", v)
	}
	if v, err := strconv.Atoi(re.Request.FormValue("grace_period_hours")); err == nil && v >= 0 {
		group.Set("grace_period_hours", v)
	}
	if err := re.App.Save(group); err != nil {
		return err
	}

	return redirect(re, "/g/"+slug+"/settings?flash=group_settings_saved")
}

// HandleLeaveGroup removes the current user's membership. The group's
// recorded owner can't leave this way — there's no ownership-transfer flow
// yet, so letting them leave would orphan the group with no owner.
func HandleLeaveGroup(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if group.GetString("owner") == user.Id {
		return re.BadRequestError("the group owner can't leave — transfer ownership first", nil)
	}
	membership, err := findMembership(re, group.Id, user.Id)
	if err != nil {
		return re.ForbiddenError("not a member of this group", nil)
	}
	if err := re.App.Delete(membership); err != nil {
		return err
	}

	return redirect(re, "/?flash=left_group")
}

// HandleReactivateGroup resets an auto-disabled group's status and
// zero-response streak so the scheduler resumes creating editions for it.
// Owner-only — admins can manage settings but reactivation is a bigger call.
func HandleReactivateGroup(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if group.GetString("owner") != user.Id {
		return re.ForbiddenError("only the group owner can reactivate it", nil)
	}

	group.Set("status", "active")
	group.Set("consecutive_unanswered_editions", 0)
	if err := re.App.Save(group); err != nil {
		return err
	}

	return redirect(re, "/g/"+slug+"/settings?flash=group_settings_saved")
}
