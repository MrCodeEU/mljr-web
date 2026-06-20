package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		if err := addUserTimezone(app); err != nil {
			return err
		}
		groups, err := createGroupsCollection(app)
		if err != nil {
			return err
		}
		if err := createGroupMembershipsCollection(app, groups); err != nil {
			return err
		}
		return nil
	}, func(app core.App) error {
		for _, name := range []string{"group_memberships", "groups"} {
			c, err := app.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			if err := app.Delete(c); err != nil {
				return err
			}
		}
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Fields.RemoveByName("timezone")
		return app.Save(users)
	})
}

func addUserTimezone(app core.App) error {
	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}
	users.Fields.Add(&core.TextField{
		Name: "timezone",
		Max:  64,
	})
	return app.Save(users)
}

func createGroupsCollection(app core.App) (*core.Collection, error) {
	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return nil, err
	}

	groups := core.NewBaseCollection("groups")
	groups.Fields.Add(
		&core.TextField{Name: "name", Required: true, Max: 120},
		&core.TextField{Name: "slug", Required: true, Max: 120, Pattern: "^[a-z0-9-]+$"},
		&core.RelationField{Name: "owner", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.SelectField{
			Name:      "schedule_period",
			Required:  true,
			Values:    []string{"weekly", "biweekly", "monthly", "quarterly"},
			MaxSelect: 1,
		},
		&core.NumberField{Name: "schedule_anchor_weekday", OnlyInt: true}, // 0=Sunday..6=Saturday
		&core.NumberField{Name: "schedule_anchor_day_of_month", OnlyInt: true},
		&core.DateField{Name: "schedule_epoch_date"}, // biweekly parity anchor
		&core.NumberField{Name: "schedule_send_hour_utc", OnlyInt: true},
		&core.NumberField{Name: "reminder_lead_hours", OnlyInt: true},
		&core.NumberField{Name: "grace_period_hours", OnlyInt: true},
		&core.TextField{Name: "timezone", Required: true, Max: 64},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)
	groups.AddIndex("idx_groups_slug", true, "slug", "")

	authRule := "owner = @request.auth.id || @request.auth.group_memberships_via_user.group ?= id"
	groups.ListRule = types.Pointer(authRule)
	groups.ViewRule = types.Pointer(authRule)
	groups.CreateRule = types.Pointer("@request.auth.id != \"\"")
	groups.UpdateRule = types.Pointer("owner = @request.auth.id")
	groups.DeleteRule = types.Pointer("owner = @request.auth.id")

	if err := app.Save(groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func createGroupMembershipsCollection(app core.App, groups *core.Collection) error {
	users, err := app.FindCollectionByNameOrId("users")
	if err != nil {
		return err
	}

	memberships := core.NewBaseCollection("group_memberships")
	memberships.Fields.Add(
		&core.RelationField{Name: "group", Required: true, CollectionId: groups.Id, MaxSelect: 1},
		&core.RelationField{Name: "user", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.SelectField{
			Name:      "role",
			Required:  true,
			Values:    []string{"owner", "admin", "member"},
			MaxSelect: 1,
		},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)
	memberships.AddIndex("idx_group_memberships_unique", true, "group, user", "")

	selfOrGroupmate := "user = @request.auth.id || group.group_memberships_via_group.user ?= @request.auth.id"
	ownerOnly := "group.owner = @request.auth.id"
	memberships.ListRule = types.Pointer(selfOrGroupmate)
	memberships.ViewRule = types.Pointer(selfOrGroupmate)
	memberships.CreateRule = types.Pointer(ownerOnly)
	memberships.UpdateRule = types.Pointer(ownerOnly)
	memberships.DeleteRule = types.Pointer(ownerOnly)

	return app.Save(memberships)
}
