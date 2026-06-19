package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		return createEmailLogCollection(app)
	}, func(app core.App) error {
		c, err := app.FindCollectionByNameOrId("email_log")
		if err != nil {
			return nil
		}
		return app.Delete(c)
	})
}

// createEmailLogCollection tracks every send attempt for idempotency: a
// unique dedupe_key (e.g. "reminder:{editionID}:{userID}" or
// "send:{editionID}") guarantees the cron scanner can retry a tick without
// double-sending.
func createEmailLogCollection(app core.App) error {
	log := core.NewBaseCollection("email_log")
	log.Fields.Add(
		&core.SelectField{Name: "kind", Required: true, Values: []string{
			"invite", "reminder", "edition_sent",
		}, MaxSelect: 1},
		&core.TextField{Name: "dedupe_key", Required: true, Max: 200},
		&core.TextField{Name: "recipient_email", Max: 200},
		&core.SelectField{Name: "status", Required: true, Values: []string{"sent", "failed"}, MaxSelect: 1},
		&core.TextField{Name: "error", Max: 2000},
		&core.AutodateField{Name: "created", OnCreate: true},
	)
	log.AddIndex("idx_email_log_dedupe", true, "dedupe_key", "")

	// No API rules — this collection is only ever touched from backend
	// code (cron scanner), never exposed to the REST API.
	return app.Save(log)
}
