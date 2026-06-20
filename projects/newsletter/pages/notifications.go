package pages

import (
	"mljr-web/ui/feedback"
	"mljr-web/ui/overlay"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// createNotification inserts an in-app notification row. link, group, invite
// and actor are all optional — pass "" to omit.
func createNotification(app core.App, userID, kind, group, invite, actor, body, link string) error {
	col, err := app.FindCollectionByNameOrId("notifications")
	if err != nil {
		return err
	}
	n := core.NewRecord(col)
	n.Set("user", userID)
	n.Set("kind", kind)
	if group != "" {
		n.Set("group", group)
	}
	if invite != "" {
		n.Set("invite", invite)
	}
	if actor != "" {
		n.Set("actor", actor)
	}
	n.Set("body", body)
	n.Set("link", link)
	return app.Save(n)
}

func recentNotifications(re *core.RequestEvent, userID string) []*core.Record {
	list, err := re.App.FindRecordsByFilter(
		"notifications", "user = {:user}", "-created", 8, 0,
		map[string]any{"user": userID},
	)
	if err != nil {
		return nil
	}
	return list
}

func countUnreadNotifications(re *core.RequestEvent, userID string) int {
	list, err := re.App.FindRecordsByFilter(
		"notifications", "user = {:user} && read_at = \"\"", "-created", 0, 0,
		map[string]any{"user": userID},
	)
	if err != nil {
		return 0
	}
	return len(list)
}

// NotifyGroupAdmins notifies every owner/admin of a group except excludeUserID
// (typically the user who triggered the event). Exported for use from
// record hooks, which only have a core.App, not a *core.RequestEvent.
func NotifyGroupAdmins(app core.App, groupID, excludeUserID, kind, actorID, body, link string) {
	members, err := app.FindRecordsByFilter(
		"group_memberships", "group = {:group} && (role = \"owner\" || role = \"admin\")", "", 0, 0,
		map[string]any{"group": groupID},
	)
	if err != nil {
		return
	}
	for _, m := range members {
		if m.GetString("user") == excludeUserID {
			continue
		}
		_ = createNotification(app, m.GetString("user"), kind, groupID, "", actorID, body, link)
	}
}

// notificationsNav renders the notification bell: an unread-count badge plus
// a dropdown of the most recent notifications and a "mark all read" action.
// Previously this data existed (recentNotifications/countUnreadNotifications)
// but nothing in the header rendered it, so invites/comments/reactions/
// edition events were silently invisible to the recipient.
func notificationsNav(re *core.RequestEvent, userID string) g.Node {
	t := translator(re)
	unread := countUnreadNotifications(re, userID)
	list := recentNotifications(re, userID)

	trigger := feedback.NotificationBadge(feedback.NotificationBadgeProps{Count: unread},
		h.Button(h.Type("button"), h.Style("background:none;border:none;cursor:pointer;font-size:1.2em;line-height:1;padding:var(--sp-1)"),
			g.Text("🔔")),
	)

	var items []overlay.DropdownItem
	if len(list) == 0 {
		items = append(items, overlay.DropdownItem{Label: t("newsletter.notifications.empty"), Href: "#"})
	}
	for _, n := range list {
		weight := "400"
		if n.GetString("read_at") == "" {
			weight = "700"
		}
		row := h.Div(h.Style("min-width:0;max-width:280px"),
			h.Span(h.Style("display:block;font-weight:"+weight+";overflow-wrap:anywhere"), g.Text(n.GetString("body"))),
			h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);margin-top:2px"), g.Text(n.GetString("created")[:10])),
		)
		href := n.GetString("link")
		if href == "" {
			href = "#"
		}
		items = append(items, overlay.DropdownItem{Content: row, Href: href})
	}
	items = append(items, overlay.DropdownItem{
		Divider: true, Label: t("newsletter.notifications.mark_all_read"), Icon: "lucide:check-check",
		OnClick: "document.getElementById('nl-mark-read-form').requestSubmit()",
	})

	return h.Div(
		h.Style("display:flex;align-items:center"),
		h.Form(h.ID("nl-mark-read-form"), h.Method("post"), h.Action("/notifications/read-all"), h.Style("display:none")),
		overlay.Dropdown(overlay.DropdownProps{Signal: "navNotifs", Align: "right"}, trigger, items...),
	)
}

// HandleMarkAllNotificationsRead marks every unread notification for the
// current user as read, then bounces back to wherever the request came from.
func HandleMarkAllNotificationsRead(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	unread, err := re.App.FindRecordsByFilter(
		"notifications", "user = {:user} && read_at = \"\"", "", 0, 0,
		map[string]any{"user": user.Id},
	)
	if err == nil {
		now := types.NowDateTime()
		for _, n := range unread {
			n.Set("read_at", now)
			_ = re.App.Save(n)
		}
	}

	back := re.Request.Referer()
	if back == "" {
		back = "/"
	}
	return redirect(re, back)
}
