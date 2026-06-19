package pages

import (
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// commentThread is a top-level comment plus its (single-level) replies.
type commentThread struct {
	Comment *core.Record
	Replies []*core.Record
}

// commentThreads loads every comment on an answer and groups replies (one
// level of nesting, per the project plan) under their parent.
func commentThreads(re *core.RequestEvent, answerID string) []commentThread {
	all, err := re.App.FindRecordsByFilter(
		"comments", "answer = {:answer}", "created", 0, 0,
		map[string]any{"answer": answerID},
	)
	if err != nil {
		return nil
	}

	var threads []commentThread
	byID := map[string]int{}
	for _, c := range all {
		if c.GetString("parent") != "" {
			continue
		}
		byID[c.Id] = len(threads)
		threads = append(threads, commentThread{Comment: c})
	}
	for _, c := range all {
		parent := c.GetString("parent")
		if parent == "" {
			continue
		}
		if idx, ok := byID[parent]; ok {
			threads[idx].Replies = append(threads[idx].Replies, c)
		}
	}
	return threads
}

// HandleCreateComment adds a comment on an answer, optionally as a reply
// (parent must itself be a top-level comment on the same answer — no
// deeper nesting). Notifies the answer's author and, for replies, the
// parent comment's author, excluding the commenter themselves.
func HandleCreateComment(re *core.RequestEvent) error {
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

	body := strings.TrimSpace(re.Request.FormValue("body"))
	if body == "" {
		return redirect(re, "/g/"+slug+"/editions/"+edition.Id+"/view")
	}

	parentID := re.Request.FormValue("parent")
	var parent *core.Record
	if parentID != "" {
		parent, err = re.App.FindRecordById("comments", parentID)
		if err != nil || parent.GetString("answer") != answer.Id || parent.GetString("parent") != "" {
			return re.BadRequestError("invalid reply target", nil)
		}
	}

	col, err := re.App.FindCollectionByNameOrId("comments")
	if err != nil {
		return err
	}
	comment := core.NewRecord(col)
	comment.Set("answer", answer.Id)
	comment.Set("author", user.Id)
	comment.Set("body", body)
	if parent != nil {
		comment.Set("parent", parent.Id)
	}
	if err := re.App.Save(comment); err != nil {
		return err
	}

	link := "/g/" + slug + "/editions/" + edition.Id + "/view"
	notified := map[string]bool{user.Id: true}
	if answerOwner := answer.GetString("user"); !notified[answerOwner] {
		_ = createNotification(re.App, answerOwner, "comment_reply", group.Id, "", user.Id,
			displayName(user)+" commented on your answer", link)
		notified[answerOwner] = true
	}
	if parent != nil {
		if parentAuthor := parent.GetString("author"); !notified[parentAuthor] {
			_ = createNotification(re.App, parentAuthor, "comment_reply", group.Id, "", user.Id,
				displayName(user)+" replied to your comment", link)
		}
	}

	return redirect(re, link)
}
