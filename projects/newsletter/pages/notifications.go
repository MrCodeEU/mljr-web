package pages

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
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
