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
		answers, err := app.FindCollectionByNameOrId("answers")
		if err != nil {
			return err
		}

		suggestions, err := createQuestionSuggestionsCollection(app, groups, users)
		if err != nil {
			return err
		}
		if err := createQuestionSuggestionVotesCollection(app, suggestions, users); err != nil {
			return err
		}
		if err := createEmojiReactionsCollection(app, answers, users); err != nil {
			return err
		}
		return createCommentsCollection(app, answers, users)
	}, func(app core.App) error {
		for _, name := range []string{"comments", "emoji_reactions", "question_suggestion_votes", "question_suggestions"} {
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

func createQuestionSuggestionsCollection(app core.App, groups, users *core.Collection) (*core.Collection, error) {
	suggestions := core.NewBaseCollection("question_suggestions")
	suggestions.Fields.Add(
		&core.RelationField{Name: "group", Required: true, CollectionId: groups.Id, MaxSelect: 1},
		&core.RelationField{Name: "suggested_by", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.SelectField{Name: "type", Required: true, Values: []string{
			"text", "single_select", "multi_select", "image", "rating", "emoji_reaction",
		}, MaxSelect: 1},
		&core.TextField{Name: "prompt", Required: true, Max: 300},
		&core.JSONField{Name: "options"}, // []string for select/emoji types
		&core.SelectField{Name: "status", Required: true, Values: []string{"pending", "approved", "rejected"}, MaxSelect: 1},
		&core.AutodateField{Name: "created", OnCreate: true},
	)

	groupMember := "group.group_memberships_via_group.user ?= @request.auth.id"
	suggestions.ListRule = types.Pointer(groupMember)
	suggestions.ViewRule = types.Pointer(groupMember)

	if err := app.Save(suggestions); err != nil {
		return nil, err
	}
	return suggestions, nil
}

func createQuestionSuggestionVotesCollection(app core.App, suggestions, users *core.Collection) error {
	votes := core.NewBaseCollection("question_suggestion_votes")
	votes.Fields.Add(
		&core.RelationField{Name: "suggestion", Required: true, CollectionId: suggestions.Id, MaxSelect: 1},
		&core.RelationField{Name: "user", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.AutodateField{Name: "created", OnCreate: true},
	)
	votes.AddIndex("idx_question_suggestion_votes_unique", true, "suggestion, user", "")

	groupMember := "suggestion.group.group_memberships_via_group.user ?= @request.auth.id"
	votes.ListRule = types.Pointer(groupMember)
	votes.ViewRule = types.Pointer(groupMember)

	return app.Save(votes)
}

func createEmojiReactionsCollection(app core.App, answers, users *core.Collection) error {
	reactions := core.NewBaseCollection("emoji_reactions")
	reactions.Fields.Add(
		&core.RelationField{Name: "answer", Required: true, CollectionId: answers.Id, MaxSelect: 1},
		&core.RelationField{Name: "user", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.TextField{Name: "emoji", Required: true, Max: 16},
		&core.AutodateField{Name: "created", OnCreate: true},
	)
	reactions.AddIndex("idx_emoji_reactions_unique", true, "answer, user, emoji", "")

	groupMember := "answer.edition.group.group_memberships_via_group.user ?= @request.auth.id"
	reactions.ListRule = types.Pointer(groupMember)
	reactions.ViewRule = types.Pointer(groupMember)

	return app.Save(reactions)
}

func createCommentsCollection(app core.App, answers, users *core.Collection) error {
	comments := core.NewBaseCollection("comments")
	comments.Fields.Add(
		&core.RelationField{Name: "answer", Required: true, CollectionId: answers.Id, MaxSelect: 1},
		&core.RelationField{Name: "author", Required: true, CollectionId: users.Id, MaxSelect: 1},
		&core.TextField{Name: "body", Required: true, Max: 1000},
		&core.AutodateField{Name: "created", OnCreate: true},
	)
	if err := app.Save(comments); err != nil {
		return err
	}

	// Self-relation for single-level reply; added after creation since it
	// needs the collection's own (now-assigned) Id.
	comments.Fields.Add(&core.RelationField{Name: "parent", CollectionId: comments.Id, MaxSelect: 1})

	groupMember := "answer.edition.group.group_memberships_via_group.user ?= @request.auth.id"
	comments.ListRule = types.Pointer(groupMember)
	comments.ViewRule = types.Pointer(groupMember)

	return app.Save(comments)
}
