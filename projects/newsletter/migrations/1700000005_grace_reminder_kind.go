package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		emailLog, err := app.FindCollectionByNameOrId("email_log")
		if err != nil {
			return err
		}
		field, ok := emailLog.Fields.GetByName("kind").(*core.SelectField)
		if !ok {
			return nil
		}
		field.Values = append(field.Values, "grace_reminder")
		return app.Save(emailLog)
	}, func(app core.App) error {
		emailLog, err := app.FindCollectionByNameOrId("email_log")
		if err != nil {
			return err
		}
		field, ok := emailLog.Fields.GetByName("kind").(*core.SelectField)
		if !ok {
			return nil
		}
		values := make([]string, 0, len(field.Values))
		for _, v := range field.Values {
			if v != "grace_reminder" {
				values = append(values, v)
			}
		}
		field.Values = values
		return app.Save(emailLog)
	})
}
