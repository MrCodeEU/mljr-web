package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// favoriteColorValues mirrors ui/token's Tone enum so favorite_color can be
// cast directly to token.Tone with no translation step.
var favoriteColorValues = []string{
	"yellow", "lime", "cyan", "violet", "pink", "sky", "mint", "blush",
}

func init() {
	m.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Fields.Add(
			&core.FileField{
				Name:      "avatar",
				MaxSelect: 1,
				MaxSize:   5 << 20,
				MimeTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
			},
			&core.DateField{Name: "birthday"},
			&core.TextField{Name: "favorite_animal", Max: 80},
			&core.TextField{Name: "favorite_food", Max: 80},
			&core.SelectField{Name: "favorite_color", Values: favoriteColorValues, MaxSelect: 1},
		)
		return app.Save(users)
	}, func(app core.App) error {
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		for _, name := range []string{"avatar", "birthday", "favorite_animal", "favorite_food", "favorite_color"} {
			users.Fields.RemoveByName(name)
		}
		return app.Save(users)
	})
}
