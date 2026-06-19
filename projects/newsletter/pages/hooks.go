package pages

import "github.com/pocketbase/pocketbase/core"

// RegisterHooks binds record-change hooks that fire in-app notifications.
// Split out from main.go for the same reason as RegisterRoutes: so tests can
// bind it onto a tests.TestApp without booting the real server.
func RegisterHooks(app core.App) {
	app.OnRecordAfterCreateSuccess("group_memberships").BindFunc(func(e *core.RecordEvent) error {
		groupID := e.Record.GetString("group")
		userID := e.Record.GetString("user")
		if member, err := e.App.FindRecordById("users", userID); err == nil {
			if group, err := e.App.FindRecordById("groups", groupID); err == nil {
				NotifyGroupAdmins(e.App, groupID, userID, "group_joined", userID,
					displayName(member)+" joined "+group.GetString("name"),
					"/g/"+group.GetString("slug"))
			}
		}
		return e.Next()
	})

	app.OnRecordAfterDeleteSuccess("group_memberships").BindFunc(func(e *core.RecordEvent) error {
		groupID := e.Record.GetString("group")
		userID := e.Record.GetString("user")
		if member, err := e.App.FindRecordById("users", userID); err == nil {
			if group, err := e.App.FindRecordById("groups", groupID); err == nil {
				NotifyGroupAdmins(e.App, groupID, userID, "group_left", userID,
					displayName(member)+" left "+group.GetString("name"),
					"/g/"+group.GetString("slug"))
			}
		}
		return e.Next()
	})
}
