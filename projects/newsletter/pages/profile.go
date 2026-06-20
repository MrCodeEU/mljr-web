package pages

import (
	"net/http"
	"slices"
	"strings"

	"mljr-web/internal/i18n"
	"mljr-web/projects/newsletter/internal/calendar"
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// favoriteColorOptions mirrors the migration's favoriteColorValues — kept
// local since pages doesn't import the migrations package.
var favoriteColorOptions = []string{
	"yellow", "lime", "cyan", "violet", "pink", "sky", "mint", "blush",
}

// Profile shows the current user's avatar, display name, email, and the
// rest of their profile attributes (birthday, favorite animal/food/color),
// with forms to update each. Avatar falls back to a deterministic
// color+initials placeholder until a real image is uploaded.
func Profile(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	t := translator(re)
	var starsignNode g.Node
	if bday := user.GetDateTime("birthday").Time(); !bday.IsZero() {
		if sign := calendar.SignForDate(int(bday.Month()), bday.Day()); sign != "" {
			starsignNode = h.P(h.Style("color:var(--muted)"), g.Text(t("newsletter.profile.starsign", sign)))
		}
	}

	colorOpts := make([]form.SelectOption, 0, len(favoriteColorOptions)+1)
	colorOpts = append(colorOpts, form.SelectOption{Value: "", Label: t("newsletter.profile.not_set"), Selected: user.GetString("favorite_color") == ""})
	for _, c := range favoriteColorOptions {
		colorOpts = append(colorOpts, form.SelectOption{Value: c, Label: c, Selected: user.GetString("favorite_color") == c})
	}

	birthdayVal := ""
	if bday := user.GetDateTime("birthday").Time(); !bday.IsZero() {
		birthdayVal = bday.Format("2006-01-02")
	}

	currentLanguage := user.GetString("language")
	languageOpts := []form.SelectOption{
		{Value: "en", Label: "English", Selected: currentLanguage == "" || currentLanguage == "en"},
		{Value: "de", Label: "Deutsch", Selected: currentLanguage == "de"},
	}

	return renderPage(re, 200, appPage(re, "", t("newsletter.profile.title"),
		[]breadcrumbItem{{Label: t("newsletter.nav.dashboard"), Href: "/"}, {Label: t("newsletter.profile.title")}},
		h.Div(
			h.Style("display:flex;flex-wrap:wrap;align-items:center;gap:var(--sp-4);margin-bottom:var(--sp-6)"),
			primitive.Avatar(primitive.AvatarProps{
				Src:      userAvatarSrc(user),
				Initials: initials(displayName(user)),
				Tone:     avatarTone(user.Id),
				Size:     token.SizeLG,
			}),
			h.Div(
				h.Style("min-width:0;overflow-wrap:anywhere"),
				primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(displayName(user))),
				h.P(h.Style("color:var(--muted)"), g.Text(user.Email())),
				starsignNode,
			),
		),
		flashAlert(re),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/profile/avatar"), h.EncType("multipart/form-data"),
				form.Field(form.FieldProps{Label: t("newsletter.profile.avatar_label")},
					form.FileInput(form.FileInputProps{Name: "avatar", Accept: "image/*", Signal: "profileAvatarFilename"}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Secondary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text(t("newsletter.profile.upload_avatar"))),
			),
		),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/profile"),
				form.Field(form.FieldProps{Label: t("newsletter.profile.name_label")},
					form.Input(form.InputProps{Type: "text", Name: "name", Required: true, Attrs: []g.Node{h.Value(user.GetString("name"))}}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.profile.birthday_label")},
					form.DateInput(form.DateInputProps{Name: "birthday", Value: birthdayVal}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.profile.animal_label")},
					form.Input(form.InputProps{Type: "text", Name: "favorite_animal", Attrs: []g.Node{h.Value(user.GetString("favorite_animal"))}}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.profile.food_label")},
					form.Input(form.InputProps{Type: "text", Name: "favorite_food", Attrs: []g.Node{h.Value(user.GetString("favorite_food"))}}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.profile.color_label")},
					form.Select(form.SelectProps{Name: "favorite_color", Options: colorOpts}),
				),
				form.Field(form.FieldProps{Label: t("newsletter.profile.language_label")},
					form.Select(form.SelectProps{Name: "language", Options: languageOpts}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text(t("newsletter.profile.save_button"))),
			),
		),
	))
}

func HandleProfile(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	name := strings.TrimSpace(re.Request.FormValue("name"))
	if name == "" {
		return redirect(re, "/profile?flash=name_required")
	}
	user.Set("name", name)

	if birthday := strings.TrimSpace(re.Request.FormValue("birthday")); birthday != "" {
		user.Set("birthday", birthday)
	} else {
		user.Set("birthday", "")
	}
	user.Set("favorite_animal", strings.TrimSpace(re.Request.FormValue("favorite_animal")))
	user.Set("favorite_food", strings.TrimSpace(re.Request.FormValue("favorite_food")))

	favoriteColor := re.Request.FormValue("favorite_color")
	if favoriteColor != "" && !slices.Contains(favoriteColorOptions, favoriteColor) {
		return redirect(re, "/profile?flash=invalid_color")
	}
	user.Set("favorite_color", favoriteColor)

	if lang := re.Request.FormValue("language"); i18n.IsSupported(lang) {
		user.Set("language", lang)
	}

	if err := re.App.Save(user); err != nil {
		return err
	}
	return redirect(re, "/profile?flash=profile_saved")
}

// HandleProfileAvatar uploads/replaces the current user's avatar image,
// reusing the same downscale-and-re-encode pipeline already used for answer
// images so stored avatars stay small regardless of the original upload.
func HandleProfileAvatar(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	re.Request.Body = http.MaxBytesReader(re.Response, re.Request.Body, maxUploadImageSize)
	if err := re.Request.ParseMultipartForm(maxUploadImageSize); err != nil {
		return redirect(re, "/profile?flash=image_too_large")
	}
	file, _, err := re.Request.FormFile("avatar")
	if err != nil {
		return redirect(re, "/profile?flash=image_invalid")
	}
	data, _, procErr := processUploadedImage(file)
	_ = file.Close()
	if procErr != nil {
		return redirect(re, "/profile?flash=image_invalid")
	}
	f, ffErr := filesystem.NewFileFromBytes(data, "avatar.jpg")
	if ffErr != nil {
		return ffErr
	}
	user.Set("avatar", f)
	if err := re.App.Save(user); err != nil {
		return err
	}
	return redirect(re, "/profile?flash=avatar_saved")
}
