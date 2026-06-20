package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// protectedFileTokenDuration is how long a minted file token stays valid.
// Avatars/answer images are linked directly (in-app <img> tags and the
// edition_sent email), so the window needs to outlive a single page view —
// long enough that someone opening the edition_sent email a few weeks late
// still sees the photos, but still bounded instead of forever.
const protectedFileTokenDuration = 30 * 24 * 3600 // 30 days

// Marks answer_images.image and users.avatar as Protected so PocketBase's
// /api/files/... handler actually evaluates each collection's ViewRule
// before serving the file. Without Protected, PocketBase serves any file at
// a known URL unconditionally, regardless of ViewRule — meaning the
// "only this group can see this photo" rule written on answer_images was
// never enforced at the file-serving layer, only at the JSON record API.
func init() {
	m.Register(func(app core.App) error {
		images, err := app.FindCollectionByNameOrId("answer_images")
		if err != nil {
			return err
		}
		if f := images.Fields.GetByName("image"); f != nil {
			f.(*core.FileField).Protected = true
		}
		if err := app.Save(images); err != nil {
			return err
		}

		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		if f := users.Fields.GetByName("avatar"); f != nil {
			f.(*core.FileField).Protected = true
		}
		users.FileToken.Duration = protectedFileTokenDuration
		return app.Save(users)
	}, func(app core.App) error {
		images, err := app.FindCollectionByNameOrId("answer_images")
		if err != nil {
			return err
		}
		if f := images.Fields.GetByName("image"); f != nil {
			f.(*core.FileField).Protected = false
		}
		if err := app.Save(images); err != nil {
			return err
		}

		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		if f := users.Fields.GetByName("avatar"); f != nil {
			f.(*core.FileField).Protected = false
		}
		users.FileToken.Duration = 180
		return app.Save(users)
	})
}
