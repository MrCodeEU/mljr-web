package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// bilingualSeed describes one global question's bilingual prompt/options.
// Used both to backfill German onto the 12 questions seeded by
// 1700000002_questions_editions.go (matched by their existing English
// prompt, since that migration already ran on existing installs and can't
// be edited in place) and to insert new bilingual global questions.
type bilingualSeed struct {
	qtype     string
	promptEN  string
	promptDE  string
	optionsEN []string
	optionsDE []string
}

// existingSeedTranslations backfills prompt_i18n/options_i18n (German) onto
// the 12 questions 1700000002_questions_editions.go already inserted.
var existingSeedTranslations = []bilingualSeed{
	{qtype: "text", promptEN: "What's the best thing that happened to you this week?", promptDE: "Was war das Beste, das dir diese Woche passiert ist?"},
	{qtype: "text", promptEN: "What are you currently obsessed with (show, game, hobby, etc.)?", promptDE: "Wovon bist du gerade total begeistert (Serie, Spiel, Hobby, etc.)?"},
	{qtype: "text", promptEN: "Any plans coming up you're excited about?", promptDE: "Gibt es Pläne, auf die du dich freust?"},
	{qtype: "single_select", promptEN: "How's your week been overall?", promptDE: "Wie war deine Woche insgesamt?",
		optionsEN: []string{"Great", "Good", "Okay", "Rough", "Terrible"}, optionsDE: []string{"Großartig", "Gut", "Okay", "Anstrengend", "Furchtbar"}},
	{qtype: "single_select", promptEN: "Favorite season right now?", promptDE: "Aktuelle Lieblingsjahreszeit?",
		optionsEN: []string{"Spring", "Summer", "Autumn", "Winter"}, optionsDE: []string{"Frühling", "Sommer", "Herbst", "Winter"}},
	{qtype: "multi_select", promptEN: "What kept you busy this week?", promptDE: "Was hat dich diese Woche beschäftigt?",
		optionsEN: []string{"Work", "Family", "Friends", "Travel", "Gaming", "Fitness", "Reading", "Cooking"},
		optionsDE: []string{"Arbeit", "Familie", "Freunde", "Reisen", "Gaming", "Fitness", "Lesen", "Kochen"}},
	{qtype: "multi_select", promptEN: "What are you grateful for lately?", promptDE: "Wofür bist du in letzter Zeit dankbar?",
		optionsEN: []string{"Health", "Family", "Friends", "Job", "A win", "Free time"},
		optionsDE: []string{"Gesundheit", "Familie", "Freunde", "Job", "Ein Erfolg", "Freizeit"}},
	{qtype: "image", promptEN: "Share a photo from your week.", promptDE: "Teile ein Foto aus deiner Woche."},
	{qtype: "image", promptEN: "Show us something you made or found.", promptDE: "Zeig uns etwas, das du gemacht oder gefunden hast."},
	{qtype: "rating", promptEN: "Rate your week out of 5 stars.", promptDE: "Bewerte deine Woche mit bis zu 5 Sternen."},
	{qtype: "rating", promptEN: "How relaxed do you feel right now?", promptDE: "Wie entspannt fühlst du dich gerade?"},
	{qtype: "emoji_reaction", promptEN: "Pick the emoji that matches your mood.", promptDE: "Wähle das Emoji, das zu deiner Stimmung passt.",
		optionsEN: []string{"😄", "🙂", "😐", "😕", "😢", "😡", "🥳", "😴"}},
}

// newGlobalQuestions are brand-new global questions, inserted bilingual from
// the start, spanning the original 6 types plus the 5 new ones added by
// 1700000010_question_i18n_and_types.go.
var newGlobalQuestions = []bilingualSeed{
	{qtype: "text", promptEN: "What's a small thing that made you smile today?", promptDE: "Was hat dich heute zum Lächeln gebracht?"},
	{qtype: "text", promptEN: "What's something you're looking forward to?", promptDE: "Auf was freust du dich gerade?"},
	{qtype: "text", promptEN: "Describe your week in one sentence.", promptDE: "Beschreibe deine Woche in einem Satz."},
	{qtype: "single_select", promptEN: "How's your energy been this week?", promptDE: "Wie war dein Energielevel diese Woche?",
		optionsEN: []string{"High", "Medium", "Low", "Drained"}, optionsDE: []string{"Hoch", "Mittel", "Niedrig", "Erschöpft"}},
	{qtype: "single_select", promptEN: "What's your current mood?", promptDE: "Wie ist deine aktuelle Stimmung?",
		optionsEN: []string{"Happy", "Calm", "Stressed", "Excited", "Tired"}, optionsDE: []string{"Glücklich", "Ruhig", "Gestresst", "Aufgeregt", "Müde"}},
	{qtype: "multi_select", promptEN: "What are you looking forward to next week?", promptDE: "Auf was freust du dich nächste Woche?",
		optionsEN: []string{"Work", "Travel", "Family time", "A hobby", "Rest", "A celebration"},
		optionsDE: []string{"Arbeit", "Reisen", "Familienzeit", "Ein Hobby", "Erholung", "Eine Feier"}},
	{qtype: "multi_select", promptEN: "Which of these did you do this week?", promptDE: "Was davon hast du diese Woche gemacht?",
		optionsEN: []string{"Cooked something new", "Exercised", "Read a book", "Watched a movie", "Called a friend", "Tried something new"},
		optionsDE: []string{"Etwas Neues gekocht", "Sport gemacht", "Ein Buch gelesen", "Einen Film gesehen", "Einen Freund angerufen", "Etwas Neues ausprobiert"}},
	{qtype: "image", promptEN: "Share a photo of your favorite meal this week.", promptDE: "Teile ein Foto deiner Lieblingsmahlzeit dieser Woche."},
	{qtype: "rating", promptEN: "Rate how productive you felt this week.", promptDE: "Bewerte, wie produktiv du dich diese Woche gefühlt hast."},
	{qtype: "emoji_reaction", promptEN: "Which emoji best captures your week?", promptDE: "Welches Emoji beschreibt deine Woche am besten?",
		optionsEN: []string{"😄", "🙂", "😐", "😕", "😢", "😡", "🥳", "😴"}},
	{qtype: "yes_no", promptEN: "Did anything stress you out this week?", promptDE: "Hat dich diese Woche etwas gestresst?"},
	{qtype: "yes_no", promptEN: "Did you try something new this week?", promptDE: "Hast du diese Woche etwas Neues ausprobiert?"},
	{qtype: "scale", promptEN: "Rate your energy level this week, 1-10.", promptDE: "Bewerte dein Energielevel diese Woche, 1-10."},
	{qtype: "scale", promptEN: "How connected did you feel to friends/family this week, 1-10?", promptDE: "Wie verbunden hast du dich diese Woche mit Freunden/Familie gefühlt, 1-10?"},
	{qtype: "number", promptEN: "How many books did you read this month?", promptDE: "Wie viele Bücher hast du diesen Monat gelesen?"},
	{qtype: "number", promptEN: "How many hours did you sleep on average this week?", promptDE: "Wie viele Stunden hast du diese Woche durchschnittlich geschlafen?"},
	{qtype: "date", promptEN: "When's your next big plan?", promptDE: "Wann ist dein nächstes großes Vorhaben?"},
	{qtype: "color_pick", promptEN: "Pick the color that matches your mood.", promptDE: "Wähle die Farbe, die zu deiner Stimmung passt."},
}

func init() {
	m.Register(func(app core.App) error {
		questions, err := app.FindCollectionByNameOrId("question_bank")
		if err != nil {
			return err
		}

		for _, s := range existingSeedTranslations {
			existing, err := app.FindFirstRecordByFilter(
				"question_bank", "scope = \"global\" && prompt = {:prompt}",
				map[string]any{"prompt": s.promptEN},
			)
			if err != nil {
				continue // not found on this install (e.g. seed migration never ran) — skip, nothing to backfill
			}
			existing.Set("prompt_i18n", map[string]string{"de": s.promptDE})
			if s.optionsDE != nil {
				existing.Set("options_i18n", map[string][]string{"de": s.optionsDE})
			}
			if err := app.Save(existing); err != nil {
				return err
			}
		}

		for _, s := range newGlobalQuestions {
			q := core.NewRecord(questions)
			q.Set("scope", "global")
			q.Set("type", s.qtype)
			q.Set("prompt", s.promptEN)
			q.Set("prompt_i18n", map[string]string{"de": s.promptDE})
			q.Set("is_active", true)
			if s.optionsEN != nil {
				q.Set("options", s.optionsEN)
			}
			if s.optionsDE != nil {
				q.Set("options_i18n", map[string][]string{"de": s.optionsDE})
			}
			if err := app.Save(q); err != nil {
				return err
			}
		}
		return nil
	}, func(app core.App) error {
		for _, s := range newGlobalQuestions {
			existing, err := app.FindFirstRecordByFilter(
				"question_bank", "scope = \"global\" && prompt = {:prompt}",
				map[string]any{"prompt": s.promptEN},
			)
			if err != nil {
				continue
			}
			if err := app.Delete(existing); err != nil {
				return err
			}
		}
		for _, s := range existingSeedTranslations {
			existing, err := app.FindFirstRecordByFilter(
				"question_bank", "scope = \"global\" && prompt = {:prompt}",
				map[string]any{"prompt": s.promptEN},
			)
			if err != nil {
				continue
			}
			existing.Set("prompt_i18n", nil)
			existing.Set("options_i18n", nil)
			if err := app.Save(existing); err != nil {
				return err
			}
		}
		return nil
	})
}
