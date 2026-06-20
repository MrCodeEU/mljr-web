package pages

import (
	"github.com/pocketbase/pocketbase/core"
)

const defaultReactionEmojis = "👍,❤️,😂,😮,😢,🎉"

func findReaction(re *core.RequestEvent, answerID, userID, emoji string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter(
		"emoji_reactions", "answer = {:answer} && user = {:user} && emoji = {:emoji}",
		map[string]any{"answer": answerID, "user": userID, "emoji": emoji},
	)
}

type reactionCount struct {
	Emoji   string
	Count   int
	Reacted bool
}

// reactionCounts groups an answer's reactions by emoji, in a stable order,
// noting whether the current user is one of the reactors.
func reactionCounts(re *core.RequestEvent, answerID, userID string) []reactionCount {
	reactions, err := re.App.FindRecordsByFilter(
		"emoji_reactions", "answer = {:answer}", "created", 0, 0,
		map[string]any{"answer": answerID},
	)
	if err != nil {
		return nil
	}
	order := []string{}
	counts := map[string]int{}
	reacted := map[string]bool{}
	for _, r := range reactions {
		emoji := r.GetString("emoji")
		if _, ok := counts[emoji]; !ok {
			order = append(order, emoji)
		}
		counts[emoji]++
		if r.GetString("user") == userID {
			reacted[emoji] = true
		}
	}
	out := make([]reactionCount, 0, len(order))
	for _, emoji := range order {
		out = append(out, reactionCount{Emoji: emoji, Count: counts[emoji], Reacted: reacted[emoji]})
	}
	return out
}

// HandleToggleReaction toggles the current user's reaction with a given
// emoji on an answer on/off, and notifies the answer's author (unless
// they're reacting to their own answer).
func HandleToggleReaction(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}
	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := findMembership(re, group.Id, user.Id); err != nil {
		return re.ForbiddenError("not a member of this group", nil)
	}
	edition, err := findEditionInGroup(re, group.Id, re.Request.PathValue("id"))
	if err != nil {
		return re.NotFoundError("edition not found", err)
	}

	answer, err := re.App.FindRecordById("answers", re.Request.PathValue("answerID"))
	if err != nil || answer.GetString("edition") != edition.Id {
		return re.NotFoundError("answer not found", err)
	}

	link := "/g/" + slug + "/editions/" + edition.Id + "/view"

	emoji := re.Request.FormValue("emoji")
	if emoji == "" {
		return redirect(re, link+"?flash=reaction_invalid")
	}

	if existing, err := findReaction(re, answer.Id, user.Id, emoji); err == nil {
		if err := re.App.Delete(existing); err != nil {
			return err
		}
		return redirect(re, link)
	}

	col, err := re.App.FindCollectionByNameOrId("emoji_reactions")
	if err != nil {
		return err
	}
	reaction := core.NewRecord(col)
	reaction.Set("answer", answer.Id)
	reaction.Set("user", user.Id)
	reaction.Set("emoji", emoji)
	if err := re.App.Save(reaction); err != nil {
		return err
	}

	if answerOwner := answer.GetString("user"); answerOwner != user.Id {
		_ = createNotification(re.App, answerOwner, "emoji_reaction", group.Id, "", user.Id,
			displayName(user)+" reacted "+emoji+" to your answer",
			"/g/"+slug+"/editions/"+edition.Id+"/view")
	}

	return redirect(re, link)
}
