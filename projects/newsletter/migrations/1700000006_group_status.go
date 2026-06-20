package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		groups, err := app.FindCollectionByNameOrId("groups")
		if err != nil {
			return err
		}
		groups.Fields.Add(
			&core.SelectField{Name: "status", Values: []string{"active", "disabled"}, MaxSelect: 1},
			&core.NumberField{Name: "consecutive_unanswered_editions", OnlyInt: true},
		)
		if err := app.Save(groups); err != nil {
			return err
		}

		records, err := app.FindRecordsByFilter("groups", "status = \"\"", "", 0, 0, nil)
		if err != nil {
			return err
		}
		for _, r := range records {
			r.Set("status", "active")
			if err := app.Save(r); err != nil {
				return err
			}
		}

		notifications, err := app.FindCollectionByNameOrId("notifications")
		if err != nil {
			return err
		}
		if kindField, ok := notifications.Fields.GetByName("kind").(*core.SelectField); ok {
			kindField.Values = append(kindField.Values, "group_disabled")
			if err := app.Save(notifications); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		groups, err := app.FindCollectionByNameOrId("groups")
		if err != nil {
			return err
		}
		groups.Fields.RemoveByName("status")
		groups.Fields.RemoveByName("consecutive_unanswered_editions")
		if err := app.Save(groups); err != nil {
			return err
		}

		notifications, err := app.FindCollectionByNameOrId("notifications")
		if err != nil {
			return err
		}
		if kindField, ok := notifications.Fields.GetByName("kind").(*core.SelectField); ok {
			values := make([]string, 0, len(kindField.Values))
			for _, v := range kindField.Values {
				if v != "group_disabled" {
					values = append(values, v)
				}
			}
			kindField.Values = values
			return app.Save(notifications)
		}
		return nil
	})
}
