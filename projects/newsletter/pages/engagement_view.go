package pages

import (
	"strings"

	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// userDisplayNameCache returns a memoized userID->displayName resolver, so
// rendering several comment threads on the same page doesn't re-fetch the
// same user record once per comment.
func userDisplayNameCache(re *core.RequestEvent) func(userID string) string {
	cache := map[string]string{}
	return func(userID string) string {
		if name, ok := cache[userID]; ok {
			return name
		}
		name := userID
		if u, err := re.App.FindRecordById("users", userID); err == nil {
			name = displayName(u)
		}
		cache[userID] = name
		return name
	}
}

// renderReactionBar shows existing emoji reaction counts (highlighted if the
// current user reacted) plus a row of one-click buttons for the emojis that
// haven't been used yet, each posting to the toggle-reaction endpoint.
func renderReactionBar(slug, editionID string, answer *core.Record, counts []reactionCount) g.Node {
	used := map[string]bool{}
	var badges []g.Node
	for _, c := range counts {
		used[c.Emoji] = true
		variant := token.Ghost
		if c.Reacted {
			variant = token.Primary
		}
		badges = append(badges, h.Form(
			h.Method("post"), h.Action("/g/"+slug+"/editions/"+editionID+"/answers/"+answer.Id+"/react"),
			h.Style("display:inline-block"),
			h.Input(h.Type("hidden"), h.Name("emoji"), h.Value(c.Emoji)),
			primitive.Button(primitive.ButtonProps{Variant: variant, Tone: token.ToneNone, Type: "submit",
				Attrs: []g.Node{h.Style("padding:var(--sp-1) var(--sp-2);font-size:var(--t-sm)")}},
				g.Textf("%s %d", c.Emoji, c.Count),
			),
		))
	}

	var picker []g.Node
	for _, emoji := range strings.Split(defaultReactionEmojis, ",") {
		if used[emoji] {
			continue
		}
		picker = append(picker, h.Form(
			h.Method("post"), h.Action("/g/"+slug+"/editions/"+editionID+"/answers/"+answer.Id+"/react"),
			h.Style("display:inline-block"),
			h.Input(h.Type("hidden"), h.Name("emoji"), h.Value(emoji)),
			primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit",
				Attrs: []g.Node{h.Style("padding:var(--sp-1) var(--sp-2);font-size:var(--t-sm);opacity:.6")}},
				g.Text(emoji),
			),
		))
	}

	return h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-1);margin-top:var(--sp-2)"),
		g.Group(badges), g.Group(picker),
	)
}

// renderCommentThreads shows every comment (with one level of replies) on an
// answer, plus a reply form per top-level comment and a new-top-level-comment
// form at the bottom.
func renderCommentThreads(re *core.RequestEvent, slug, editionID string, answer *core.Record, threads []commentThread, byUser func(userID string) string) g.Node {
	t := translator(re)
	actionURL := "/g/" + slug + "/editions/" + editionID + "/answers/" + answer.Id + "/comments"

	var nodes []g.Node
	for _, thread := range threads {
		nodes = append(nodes, h.Div(h.Style("margin-top:var(--sp-2);padding-left:var(--sp-3);border-left:var(--border-w) var(--border-style) var(--line)"),
			h.P(h.Style("font-size:var(--t-sm)"),
				h.Span(h.Style("font-weight:600"), g.Text(byUser(thread.Comment.GetString("author"))+": ")),
				g.Text(thread.Comment.GetString("body")),
			),
			g.Group(renderReplies(thread.Replies, byUser)),
			h.Form(h.Method("post"), h.Action(actionURL), h.Style("display:flex;gap:var(--sp-2);margin-top:var(--sp-1)"),
				h.Input(h.Type("hidden"), h.Name("parent"), h.Value(thread.Comment.Id)),
				form.Input(form.InputProps{Type: "text", Name: "body", Placeholder: t("newsletter.comments.reply_placeholder"), Required: true,
					Attrs: []g.Node{h.Style("flex:1;font-size:var(--t-sm)")}}),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit",
					Attrs: []g.Node{h.Style("font-size:var(--t-sm)")}}, g.Text(t("newsletter.comments.reply_button"))),
			),
		))
	}

	nodes = append(nodes, h.Form(h.Method("post"), h.Action(actionURL), h.Style("display:flex;gap:var(--sp-2);margin-top:var(--sp-2)"),
		form.Input(form.InputProps{Type: "text", Name: "body", Placeholder: t("newsletter.comments.add_placeholder"), Required: true,
			Attrs: []g.Node{h.Style("flex:1;font-size:var(--t-sm)")}}),
		primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit",
			Attrs: []g.Node{h.Style("font-size:var(--t-sm)")}}, g.Text(t("newsletter.comments.comment_button"))),
	))

	return h.Div(g.Group(nodes))
}

func renderReplies(replies []*core.Record, byUser func(userID string) string) []g.Node {
	var nodes []g.Node
	for _, r := range replies {
		nodes = append(nodes, h.P(h.Style("font-size:var(--t-sm);margin-top:var(--sp-1)"),
			h.Span(h.Style("font-weight:600"), g.Text(byUser(r.GetString("author"))+": ")),
			g.Text(r.GetString("body")),
		))
	}
	return nodes
}
