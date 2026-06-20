package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Adds bilingual prompt/options storage to question_bank (prompt/options stay
// the canonical English fallback — only global questions get translated) and
// extends the type enum with 5 new answer types, each mapping onto an
// existing ui/form component (switch, slider, numberinput, dateinput, and the
// existing tone swatches reused as buttons).
func init() {
	m.Register(func(app core.App) error {
		questions, err := app.FindCollectionByNameOrId("question_bank")
		if err != nil {
			return err
		}
		questions.Fields.Add(
			&core.JSONField{Name: "prompt_i18n"},  // map[string]string, e.g. {"de": "..."} — en lives in prompt
			&core.JSONField{Name: "options_i18n"}, // map[string][]string, same idea, only for select/emoji types
		)
		if f := questions.Fields.GetByName("type"); f != nil {
			f.(*core.SelectField).Values = []string{
				"text", "single_select", "multi_select", "image", "rating", "emoji_reaction",
				"yes_no", "scale", "number", "date", "color_pick",
			}
		}
		return app.Save(questions)
	}, func(app core.App) error {
		questions, err := app.FindCollectionByNameOrId("question_bank")
		if err != nil {
			return err
		}
		questions.Fields.RemoveByName("prompt_i18n")
		questions.Fields.RemoveByName("options_i18n")
		if f := questions.Fields.GetByName("type"); f != nil {
			f.(*core.SelectField).Values = []string{
				"text", "single_select", "multi_select", "image", "rating", "emoji_reaction",
			}
		}
		return app.Save(questions)
	})
}
