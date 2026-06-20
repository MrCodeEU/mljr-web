package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"mljr-web/internal/config"
	"mljr-web/internal/i18n"
	"mljr-web/projects/newsletter/internal/calendar"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// questionPrompt mirrors pages.questionPrompt (duplicated rather than
// imported — pages imports scheduler indirectly via routes wiring, so
// importing pages from here would risk a cycle) — returns a global
// question's prompt in lang, falling back to the canonical English prompt.
func questionPrompt(q *core.Record, lang string) string {
	if lang != i18n.DefaultLang {
		raw, ok := q.Get("prompt_i18n").(types.JSONRaw)
		if ok && len(raw) > 0 {
			var m map[string]string
			if json.Unmarshal(raw, &m) == nil && m[lang] != "" {
				return m[lang]
			}
		}
	}
	return q.GetString("prompt")
}

// toneHex maps token.Tone values to email-safe hex colors, mirroring the
// app's CSS tone palette (ui/css/_base.css) since email clients can't load
// the site stylesheet and need literal colors.
var toneHex = map[token.Tone]string{
	token.ToneYellow: "#ffd23f",
	token.ToneLime:   "#a3e635",
	token.ToneCyan:   "#15c8d4",
	token.ToneViolet: "#8b5cf6",
	token.TonePink:   "#ff5da2",
	token.ToneSky:    "#dbeafe",
	token.ToneMint:   "#d8f5e3",
	token.ToneBlush:  "#ffe4ec",
}

// editionAvatarTones mirrors pages' avatarTone fallback (same fixed
// 8-tone palette, same deterministic id-sum derivation) so a member who
// hasn't picked a favorite_color still gets a stable, consistent color
// between the in-app view and the email.
var editionAvatarTones = []token.Tone{
	token.ToneYellow, token.ToneLime, token.ToneCyan, token.ToneViolet,
	token.TonePink, token.ToneSky, token.ToneMint, token.ToneBlush,
}

func editionAvatarTone(seed string) token.Tone {
	var sum int
	for _, b := range []byte(seed) {
		sum += int(b)
	}
	return editionAvatarTones[sum%len(editionAvatarTones)]
}

func userTone(user *core.Record) token.Tone {
	if c := user.GetString("favorite_color"); c != "" {
		return token.Tone(c)
	}
	return editionAvatarTone(user.Id)
}

type editionAnswerRow struct {
	UserName  string
	Tone      token.Tone
	Text      string
	IsImage   bool
	ImageURLs []string
}

type editionQuestionSection struct {
	Prompt string
	Rows   []editionAnswerRow
}

type memberBirthday struct {
	Name  string
	Month int
	Day   int
}

type editionComposeData struct {
	GroupName         string
	OpensAtLabel      string
	ViewURL           string
	Sections          []editionQuestionSection
	Horoscopes        map[string]string // zodiac sign -> blurb, only signs present among answering members
	UpcomingBirthdays []memberBirthday
}

// gatherEditionComposeData fetches everything the edition_sent email needs:
// the same per-question/per-answer data compileEditionText used to fetch,
// plus each answering member's color/starsign, distinct-sign horoscope
// blurbs, and upcoming birthdays before the group's next edition. Both
// renderEditionText and renderEditionHTML consume this one struct so there
// is exactly one query/parsing path for the two output formats.
//
// forUser is whoever will actually receive the resulting email — answer
// images are a Protected file field, so each image URL needs a file token
// minted for that specific recipient, which means this (and the email
// built from it) can no longer be shared across multiple recipients.
func gatherEditionComposeData(app core.App, cfg config.Config, group, edition, forUser *core.Record, horoscopes HoroscopeProvider) (editionComposeData, error) {
	data := editionComposeData{
		GroupName: group.GetString("name"),
		ViewURL:   cfg.Newsletter.PublicAppURL + "/g/" + group.GetString("slug") + "/editions/" + edition.Id + "/view",
	}
	if opensAt := edition.GetString("opens_at"); len(opensAt) >= 10 {
		data.OpensAtLabel = opensAt[:10]
	}

	eqs, err := app.FindRecordsByFilter(
		"edition_questions", "edition = {:edition}", "order", 0, 0, map[string]any{"edition": edition.Id},
	)
	if err != nil {
		return data, err
	}

	signsPresent := map[string]bool{}

	for _, eq := range eqs {
		question, err := app.FindRecordById("question_bank", eq.GetString("question"))
		if err != nil {
			continue
		}
		section := editionQuestionSection{Prompt: questionPrompt(question, forUser.GetString("language"))}

		answers, err := app.FindRecordsByFilter(
			"answers", "edition = {:edition} && question = {:question} && skipped = false", "", 0, 0,
			map[string]any{"edition": edition.Id, "question": question.Id},
		)
		if err != nil {
			continue
		}
		for _, answer := range answers {
			user, err := app.FindRecordById("users", answer.GetString("user"))
			if err != nil {
				continue
			}
			text := valueAsText(decodeJSONField(answer, "value"))
			isImage := question.GetString("type") == "image"
			var imageURLs []string
			if isImage {
				text = "[shared a photo]"
				imageURLs = answerImageURLs(app, cfg, answer.Id, forUser)
			}
			if text == "" {
				continue
			}
			section.Rows = append(section.Rows, editionAnswerRow{
				UserName:  displayNameOrEmail(user),
				Tone:      userTone(user),
				Text:      text,
				IsImage:   isImage,
				ImageURLs: imageURLs,
			})

			if bday := user.GetDateTime("birthday").Time(); !bday.IsZero() {
				if sign := calendar.SignForDate(int(bday.Month()), bday.Day()); sign != "" {
					signsPresent[sign] = true
				}
			}
		}
		data.Sections = append(data.Sections, section)
	}

	if len(signsPresent) > 0 && horoscopes != nil {
		data.Horoscopes = map[string]string{}
		ctx := context.Background()
		now := time.Now().UTC()
		for sign := range signsPresent {
			blurb, err := horoscopes.Daily(ctx, sign, now)
			if err != nil || blurb == "" {
				continue // graceful degradation: omit this sign's blurb, never fail the send
			}
			data.Horoscopes[sign] = blurb
		}
	}

	upcoming, err := upcomingBirthdays(app, group, edition)
	if err == nil {
		data.UpcomingBirthdays = upcoming
	}

	return data, nil
}

// answerImageURLs resolves an image answer's uploaded files to absolute
// URLs (PublicAppURL + the PocketBase file-serving path) so they render in
// an email client, which can't resolve relative paths the way a browser
// viewing the app would. answer_images.image is Protected, so each URL
// also carries a file token minted for forUser (the email's one recipient)
// — without it the image would 404 for everyone, since PocketBase only
// serves protected files to a request whose token resolves to a record the
// collection's ViewRule actually grants access to.
func answerImageURLs(app core.App, cfg config.Config, answerID string, forUser *core.Record) []string {
	imgs, err := app.FindRecordsByFilter(
		"answer_images", "answer = {:answer}", "order", 0, 0, map[string]any{"answer": answerID},
	)
	if err != nil {
		return nil
	}
	token, err := forUser.NewFileToken()
	if err != nil {
		return nil
	}
	urls := make([]string, 0, len(imgs))
	for _, img := range imgs {
		urls = append(urls, cfg.Newsletter.PublicAppURL+"/api/files/"+img.Collection().Id+"/"+img.Id+"/"+img.GetString("image")+"?token="+token)
	}
	return urls
}

// upcomingBirthdays returns every group member whose birthday falls
// between now and the group's next expected edition opens_at — a
// cadence-aware window (rather than a fixed N days) so weekly groups don't
// repeat the same birthday across several editions and quarterly groups
// don't miss one entirely.
func upcomingBirthdays(app core.App, group, edition *core.Record) ([]memberBirthday, error) {
	from := time.Now().UTC()
	to, err := nextWindowForGroup(group, edition.GetDateTime("opens_at").Time())
	if err != nil || !to.After(from) {
		return nil, nil
	}

	members, err := app.FindRecordsByFilter(
		"group_memberships", "group = {:group}", "", 0, 0, map[string]any{"group": group.Id},
	)
	if err != nil {
		return nil, err
	}

	var out []memberBirthday
	for _, membership := range members {
		user, err := app.FindRecordById("users", membership.GetString("user"))
		if err != nil {
			continue
		}
		bday := user.GetDateTime("birthday").Time()
		if bday.IsZero() {
			continue
		}
		month, day := int(bday.Month()), bday.Day()
		if calendar.InRange(month, day, from, to) {
			out = append(out, memberBirthday{Name: displayNameOrEmail(user), Month: month, Day: day})
		}
	}
	return out, nil
}

// Inline-style design tokens hand-ported from ui/css/_base.css — email
// clients can't load the site stylesheet, and scheduler/ deliberately
// doesn't import ui/ (those components emit data-* attrs that depend on
// it), so the actual values are duplicated here as literals instead.
const (
	emailFontBody    = `-apple-system,BlinkMacSystemFont,"Inter",system-ui,sans-serif`
	emailFontHeading = `"Archivo",-apple-system,BlinkMacSystemFont,system-ui,sans-serif`
	emailInk         = "#141414"
	emailMuted       = "#5a5a5a"
	emailSurface2    = "#fff7e6"
	emailAccent      = "#8b5cf6"
	emailLine        = "#e5e1d8"
	emailRadius      = "4px"
)

// transactionalEmailHTML builds the inline-styled HTML body shared by the
// simple one-paragraph-and-a-link emails (reminder, grace reminder) — same
// email-safe inline-style approach as renderEditionHTML, just smaller.
func transactionalEmailHTML(greetingName, lead, linkURL, linkLabel string) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<div style="font-family:%s;max-width:640px;margin:0 auto;color:%s">`, emailFontBody, emailInk)
	fmt.Fprintf(&b, `<p style="font-size:14px">Hi %s,</p>`, html.EscapeString(greetingName))
	fmt.Fprintf(&b, `<p style="font-size:14px;line-height:1.5">%s</p>`, html.EscapeString(lead))
	if linkURL != "" {
		fmt.Fprintf(&b,
			`<p style="margin-top:20px"><a href="%s" style="display:inline-block;background:%s;color:#fff;padding:10px 18px;border-radius:%s;text-decoration:none;font-size:14px;font-family:%s">%s</a></p>`,
			html.EscapeString(linkURL), emailAccent, emailRadius, emailFontBody, html.EscapeString(linkLabel),
		)
	}
	b.WriteString(`</div>`)
	return b.String()
}

// renderEditionText replaces compileEditionText's body: the same plain-text
// shape as before, plus a horoscopes section and an upcoming-birthdays
// section when there's anything to show.
func renderEditionText(lang string, data editionComposeData) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s newsletter\n\n", data.GroupName)
	if data.ViewURL != "" {
		fmt.Fprintf(&b, "%s: %s\n\n", i18n.T(lang, "newsletter.email.view_online"), data.ViewURL)
	}

	for _, section := range data.Sections {
		fmt.Fprintf(&b, "%s\n", section.Prompt)
		for _, row := range section.Rows {
			fmt.Fprintf(&b, "  - %s: %s\n", row.UserName, row.Text)
		}
		b.WriteString("\n")
	}

	if len(data.UpcomingBirthdays) > 0 {
		fmt.Fprintf(&b, "%s:\n", i18n.T(lang, "newsletter.email.upcoming_birthdays"))
		for _, bd := range data.UpcomingBirthdays {
			fmt.Fprintf(&b, "  - %s (%02d/%02d)\n", bd.Name, bd.Month, bd.Day)
		}
		b.WriteString("\n")
	}

	if len(data.Horoscopes) > 0 {
		fmt.Fprintf(&b, "%s:\n", i18n.T(lang, "newsletter.email.horoscopes"))
		for sign, blurb := range data.Horoscopes {
			fmt.Fprintf(&b, "  - %s: %s\n", sign, blurb)
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderEditionHTML builds an inline-styled, table-based HTML email body —
// standard email-safe practice. Deliberately not using the ui/ gomponents
// components: those emit data-* attributes that depend on the site
// stylesheet, which email clients never load, and scheduler/ doesn't
// otherwise depend on ui/ — keeping that separation intact.
func renderEditionHTML(lang string, data editionComposeData) string {
	var b strings.Builder

	fmt.Fprintf(&b, `<div style="font-family:%s;max-width:640px;margin:0 auto;color:%s">`, emailFontBody, emailInk)
	fmt.Fprintf(&b, `<h1 style="font-family:%s;font-size:20px;margin-bottom:4px">%s newsletter</h1>`, emailFontHeading, html.EscapeString(data.GroupName))
	if data.OpensAtLabel != "" {
		fmt.Fprintf(&b, `<p style="color:%s;font-size:13px;margin-top:0">%s</p>`, emailMuted, html.EscapeString(data.OpensAtLabel))
	}
	if data.ViewURL != "" {
		fmt.Fprintf(&b, `<p style="margin-top:0"><a href="%s" style="color:%s;font-size:13px;text-decoration:underline">%s</a></p>`,
			html.EscapeString(data.ViewURL), emailAccent, html.EscapeString(i18n.T(lang, "newsletter.email.view_online")))
	}

	if len(data.UpcomingBirthdays) > 0 {
		fmt.Fprintf(&b, `<div style="background:%s;border:1px solid %s;border-radius:%s;padding:12px 16px;margin:16px 0"><strong>🎂 %s</strong><ul style="margin:8px 0 0;padding-left:20px">`,
			emailSurface2, emailLine, emailRadius, html.EscapeString(i18n.T(lang, "newsletter.email.upcoming_birthdays")))
		for _, bd := range data.UpcomingBirthdays {
			fmt.Fprintf(&b, `<li>%s — %02d/%02d</li>`, html.EscapeString(bd.Name), bd.Month, bd.Day)
		}
		b.WriteString(`</ul></div>`)
	}

	for _, section := range data.Sections {
		fmt.Fprintf(&b, `<h2 style="font-family:%s;font-size:16px;margin-top:24px;margin-bottom:8px">%s</h2>`, emailFontHeading, html.EscapeString(section.Prompt))
		if len(section.Rows) == 0 {
			fmt.Fprintf(&b, `<p style="color:%s">%s</p>`, emailMuted, html.EscapeString(i18n.T(lang, "newsletter.email.no_answers")))
			continue
		}
		for _, row := range section.Rows {
			hex := toneHex[row.Tone]
			if hex == "" {
				hex = "#cccccc"
			}
			fmt.Fprintf(&b, `<div style="border-left:4px solid %s;padding:6px 0 6px 10px;margin-bottom:6px">`, hex)
			if len(row.ImageURLs) > 0 {
				fmt.Fprintf(&b, `<strong>%s:</strong><div style="margin-top:6px">`, html.EscapeString(row.UserName))
				for _, url := range row.ImageURLs {
					fmt.Fprintf(&b, `<img src="%s" alt="photo from %s" style="max-width:280px;max-height:280px;border-radius:%s;margin:0 8px 8px 0">`,
						html.EscapeString(url), html.EscapeString(row.UserName), emailRadius)
				}
				b.WriteString(`</div>`)
			} else {
				fmt.Fprintf(&b, `<strong>%s:</strong> %s`, html.EscapeString(row.UserName), html.EscapeString(row.Text))
			}
			b.WriteString(`</div>`)
		}
	}

	if len(data.Horoscopes) > 0 {
		fmt.Fprintf(&b, `<div style="background:%s;border:1px solid %s;border-radius:%s;padding:12px 16px;margin:24px 0"><strong>✨ Horoscopes</strong><ul style="margin:8px 0 0;padding-left:20px">`,
			emailSurface2, emailLine, emailRadius)
		for sign, blurb := range data.Horoscopes {
			fmt.Fprintf(&b, `<li><strong>%s:</strong> %s</li>`, html.EscapeString(sign), html.EscapeString(blurb))
		}
		b.WriteString(`</ul></div>`)
	}

	b.WriteString(`</div>`)
	return b.String()
}
