package pages

import (
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// HandleCreateGroup creates a group owned by the current user and adds them
// as an "owner" member. Schedule fields default to a sane preset (weekly,
// Friday, 18:00 UTC) — editable later from group settings.
func HandleCreateGroup(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	name := re.Request.FormValue("name")
	if name == "" {
		return redirect(re, "/")
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

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := findMembership(re, group.Id, user.Id); err != nil {
		return re.ForbiddenError("not a member of this group", nil)
	}

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
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name")}},
		h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:baseline;margin-bottom:var(--sp-6)"),
			primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(group.GetString("name"))),
			h.Div(h.Style("display:flex;gap:var(--sp-3);white-space:nowrap"),
				h.A(h.Href("/g/"+slug+"/editions"), g.Text("Editions")),
				h.A(h.Href("/g/"+slug+"/questions"), g.Text("Questions")),
				h.A(h.Href("/g/"+slug+"/invites"), g.Text("Invites")),
				h.A(h.Href("/g/"+slug+"/settings"), g.Text("Settings")),
			),
		),
		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Members")),
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

	return renderPage(re, 200, appPage(re, slug, "Settings — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Settings"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(group.GetString("name")+" settings")),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/g/"+slug+"/settings"),
				form.Field(form.FieldProps{Label: "Name"},
					form.Input(form.InputProps{Type: "text", Name: "name", Required: true, Attrs: []g.Node{h.Value(group.GetString("name"))}}),
				),
				form.Field(form.FieldProps{Label: "Schedule period", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Select(form.SelectProps{
						Name: "schedule_period",
						Options: []form.SelectOption{
							periodOption("weekly", "Weekly"),
							periodOption("biweekly", "Biweekly"),
							periodOption("monthly", "Monthly"),
							periodOption("quarterly", "Quarterly"),
						},
					}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text("Save")),
			),
		),
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

	group.Set("name", re.Request.FormValue("name"))
	group.Set("schedule_period", re.Request.FormValue("schedule_period"))
	if err := re.App.Save(group); err != nil {
		return err
	}

	return redirect(re, "/g/"+slug+"/settings")
}
