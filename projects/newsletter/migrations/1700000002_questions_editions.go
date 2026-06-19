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

		questions, err := createQuestionBankCollection(app, groups, users)
		if err != nil {
			return err
		}
		editions, err := createEditionsCollection(app, groups)
		if err != nil {
			return err
		}
		if err := createEditionQuestionsCollection(app, editions, questions); err != nil {
			return err
		}
		answers, err := createAnswersCollection(app, editions, questions, users)
		if err != nil {
			return err
		}
		if err := createAnswerImagesCollection(app, answers); err != nil {
			return err
		}
		return seedGlobalQuestions(app, questions)
	}, func(app core.App) error {
		for _, name := range []string{"answer_images", "answers", "edition_questions", "newsletter_editions", "question_bank"} {
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

const questionTypes = "text,single_select,multi_select,image,rating,emoji_reaction"

func createQuestionBankCollection(app core.App, groups, users *core.Collection) (*core.Collection, error) {
	questions := core.NewBaseCollection("question_bank")
	questions.Fields.Add(
		&core.SelectField{Name: "scope", Required: true, Values: []string{"global", "group", "user"}, MaxSelect: 1},
		&core.RelationField{Name: "group", CollectionId: groups.Id, MaxSelect: 1},
		&core.RelationField{Name: "author", CollectionId: users.Id, MaxSelect: 1},
		&core.SelectField{Name: "type", Required: true, Values: []string{
			"text", "single_select", "multi_select", "image", "rating", "emoji_reaction",
		}, MaxSelect: 1},
		&core.TextField{Name: "prompt", Required: true, Max: 300},
		&core.JSONField{Name: "options"}, // []string for select/emoji types; unused for text/image/rating
		&core.BoolField{Name: "is_active"},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)

	groupMember := "group.group_memberships_via_group.user ?= @request.auth.id"
	visible := "scope = \"global\" || (scope = \"group\" && " + groupMember + ") || author = @request.auth.id"
	manageable := "author = @request.auth.id || (scope = \"group\" && " + groupMember + ")"
	questions.ListRule = types.Pointer(visible)
	questions.ViewRule = types.Pointer(visible)
	questions.CreateRule = types.Pointer("@request.auth.id != \"\"")
	questions.UpdateRule = types.Pointer(manageable)
	questions.DeleteRule = types.Pointer(manageable)

	if err := app.Save(questions); err != nil {
		return nil, err
	}
	return questions, nil
}

func createEditionsCollection(app core.App, groups *core.Collection) (*core.Collection, error) {
	editions := core.NewBaseCollection("newsletter_editions")
	editions.Fields.Add(
		&core.RelationField{Name: "group", Required: true, CollectionId: groups.Id, MaxSelect: 1},
		&core.DateField{Name: "opens_at"},
		&core.DateField{Name: "closes_at"},
		&core.DateField{Name: "reminder_at"},
		&core.DateField{Name: "grace_until"},
		&core.SelectField{Name: "status", Required: true, Values: []string{
			"scheduled", "open", "reminder_sent", "grace", "sent", "archived",
		}, MaxSelect: 1},
		&core.DateField{Name: "sent_at"},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)

	groupMember := "group.group_memberships_via_group.user ?= @request.auth.id"
	editions.ListRule = types.Pointer(groupMember)
	editions.ViewRule = types.Pointer(groupMember)

	if err := app.Save(editions); err != nil {
		return nil, err
	}
	return editions, nil
}

func createEditionQuestionsCollection(app core.App, editions, questions *core.Collection) error {
	eq := core.NewBaseCollection("edition_questions")
	eq.Fields.Add(
		&core.RelationField{Name: "edition", Required: true, CollectionId: editions.Id, MaxSelect: 1},
		&core.RelationField{Name: "question", Required: true, CollectionId: questions.Id, MaxSelect: 1},
		&core.NumberField{Name: "order", OnlyInt: true},
		&core.NumberField{Name: "vote_count", OnlyInt: true},
	)
	eq.AddIndex("idx_edition_questions_unique", true, "edition, question", "")

	groupMember := "edition.group.group_memberships_via_group.user ?= @request.auth.id"
	eq.ListRule = types.Pointer(groupMember)
	eq.ViewRule = types.Pointer(groupMember)

	return app.Save(eq)
}

func createAnswersCollection(app core.App, editions, questions, users *core.Collection) (*core.Collection, error) {
	answers := core.NewBaseCollection("answers")
	answers.Fields.Add(
		&core.RelationField{Name: "edition", Required: true, CollectionId: editions.Id, MaxSelect: 1},
		&core.RelationField{Name: "question", Required: true, CollectionId: questions.Id, MaxSelect: 1},
		&core.RelationField{Name: "user", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.JSONField{Name: "value"}, // shape depends on question.type
		&core.BoolField{Name: "skipped"},
		&core.AutodateField{Name: "created", OnCreate: true},
		&core.AutodateField{Name: "updated", OnCreate: true, OnUpdate: true},
	)
	answers.AddIndex("idx_answers_unique", true, "edition, question, user", "")

	groupMember := "edition.group.group_memberships_via_group.user ?= @request.auth.id"
	answers.ListRule = types.Pointer(groupMember)
	answers.ViewRule = types.Pointer(groupMember)

	if err := app.Save(answers); err != nil {
		return nil, err
	}
	return answers, nil
}

func createAnswerImagesCollection(app core.App, answers *core.Collection) error {
	images := core.NewBaseCollection("answer_images")
	images.Fields.Add(
		&core.RelationField{Name: "answer", Required: true, CollectionId: answers.Id, MaxSelect: 1},
		&core.FileField{Name: "image", Required: true, MaxSelect: 1, MaxSize: 5 << 20, MimeTypes: []string{
			"image/jpeg", "image/png", "image/webp", "image/gif",
		}},
		&core.NumberField{Name: "order", OnlyInt: true},
	)

	groupMember := "answer.edition.group.group_memberships_via_group.user ?= @request.auth.id"
	images.ListRule = types.Pointer(groupMember)
	images.ViewRule = types.Pointer(groupMember)

	return app.Save(images)
}

// seedGlobalQuestions adds a starter set of built-in questions covering all
// 6 answer types, so a fresh install has content out of the box.
func seedGlobalQuestions(app core.App, questions *core.Collection) error {
	type seed struct {
		qtype   string
		prompt  string
		options []string
	}
	seeds := []seed{
		{"text", "What's the best thing that happened to you this week?", nil},
		{"text", "What are you currently obsessed with (show, game, hobby, etc.)?", nil},
		{"text", "Any plans coming up you're excited about?", nil},
		{"single_select", "How's your week been overall?", []string{"Great", "Good", "Okay", "Rough", "Terrible"}},
		{"single_select", "Favorite season right now?", []string{"Spring", "Summer", "Autumn", "Winter"}},
		{"multi_select", "What kept you busy this week?", []string{"Work", "Family", "Friends", "Travel", "Gaming", "Fitness", "Reading", "Cooking"}},
		{"multi_select", "What are you grateful for lately?", []string{"Health", "Family", "Friends", "Job", "A win", "Free time"}},
		{"image", "Share a photo from your week.", nil},
		{"image", "Show us something you made or found.", nil},
		{"rating", "Rate your week out of 5 stars.", nil},
		{"rating", "How relaxed do you feel right now?", nil},
		{"emoji_reaction", "Pick the emoji that matches your mood.", []string{"😄", "🙂", "😐", "😕", "😢", "😡", "🥳", "😴"}},
	}

	for _, s := range seeds {
		q := core.NewRecord(questions)
		q.Set("scope", "global")
		q.Set("type", s.qtype)
		q.Set("prompt", s.prompt)
		q.Set("is_active", true)
		if s.options != nil {
			q.Set("options", s.options)
		}
		if err := app.Save(q); err != nil {
			return err
		}
	}
	return nil
}
