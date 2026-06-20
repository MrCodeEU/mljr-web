package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		groups, err := app.FindCollectionByNameOrId("groups")
		if err != nil {
			return err
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		invites, err := createGroupInvitesCollection(app, groups, users)
		if err != nil {
			return err
		}
		return createNotificationsCollection(app, users, groups, invites)
	}, func(app core.App) error {
		for _, name := range []string{"notifications", "group_invites"} {
			c, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			if err := app.Delete(c); err != nil {
				return err
			}
		}
		return nil
	})
}

func createGroupInvitesCollection(app core.App, groups, users *core.Collection) (*core.Collection, error) {
	invites := core.NewBaseCollection("group_invites")
	invites.Fields.Add(
		&core.RelationField{Name: "group", Required: true, CollectionId: groups.Id, MaxSelect: 1},
		&core.RelationField{Name: "invited_by", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.EmailField{Name: "email", Required: true},
		&core.RelationField{Name: "invited_user", CollectionId: users.Id, MaxSelect: 1},
		&core.TextField{Name: "token", Required: true, Max: 64},
		&core.SelectField{
			Name:      "role",
			Required:  true,
			Values:    []string{"admin", "member"},
			MaxSelect: 1,
		},
		&core.SelectField{
			Name:      "status",
			Required:  true,
			Values:    []string{"pending", "accepted", "expired", "revoked"},
			MaxSelect: 1,
		},
		&core.DateField{Name: "expires_at", Required: true},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)
	invites.AddIndex("idx_group_invites_token", true, "token", "")

	ownerOnly := "group.owner = @request.auth.id"
	invites.ListRule = types.Pointer(ownerOnly)
	invites.ViewRule = types.Pointer(ownerOnly + " || invited_user = @request.auth.id")
	invites.CreateRule = types.Pointer(ownerOnly)
	invites.UpdateRule = types.Pointer(ownerOnly)
	invites.DeleteRule = types.Pointer(ownerOnly)

	if err := app.Save(invites); err != nil {
		return nil, err
	}
	return invites, nil
}

func createNotificationsCollection(app core.App, users, groups, invites *core.Collection) error {
	notifications := core.NewBaseCollection("notifications")
	notifications.Fields.Add(
		&core.RelationField{Name: "user", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.SelectField{
			Name:     "kind",
			Required: true,
			Values: []string{
				"invite", "group_joined", "group_left",
				"edition_open", "edition_reminder", "edition_sent",
				"comment_reply", "emoji_reaction",
			},
			MaxSelect: 1,
		},
		&core.RelationField{Name: "group", CollectionId: groups.Id, MaxSelect: 1},
		&core.RelationField{Name: "invite", CollectionId: invites.Id, MaxSelect: 1},
		&core.RelationField{Name: "actor", CollectionId: users.Id, MaxSelect: 1},
		&core.TextField{Name: "body", Required: true, Max: 500},
		&core.TextField{Name: "link", Max: 200},
		&core.DateField{Name: "read_at"},
		&core.AutodateField{Name: "created", OnCreate: true},
	)
	notifications.AddIndex("idx_notifications_user", false, "user, created", "")

	selfOnly := "user = @request.auth.id"
	notifications.ListRule = types.Pointer(selfOnly)
	notifications.ViewRule = types.Pointer(selfOnly)
	notifications.UpdateRule = types.Pointer(selfOnly)
	notifications.DeleteRule = types.Pointer(selfOnly)

	return app.Save(notifications)
}
