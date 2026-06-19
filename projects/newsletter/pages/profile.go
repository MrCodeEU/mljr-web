package pages

import (
	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Profile shows the current user's avatar, display name, and email, with a
// form to update their name. Avatar upload is deferred — for now everyone
// gets a deterministic color+initials avatar derived from their name/id.
func Profile(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	return renderPage(re, 200, appPage(re, "", "Profile",
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: "Profile"}},
		h.Div(
			h.Style("display:flex;flex-wrap:wrap;align-items:center;gap:var(--sp-4);margin-bottom:var(--sp-6)"),
			primitive.Avatar(primitive.AvatarProps{
				Initials: initials(displayName(user)),
				Tone:     avatarTone(user.Id),
				Size:     token.SizeLG,
			}),
			h.Div(
				h.Style("min-width:0;overflow-wrap:anywhere"),
				primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(displayName(user))),
				h.P(h.Style("color:var(--muted)"), g.Text(user.Email())),
			),
		),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/profile"),
				form.Field(form.FieldProps{Label: "Name"},
					form.Input(form.InputProps{Type: "text", Name: "name", Required: true, Attrs: []g.Node{h.Value(user.GetString("name"))}}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text("Save")),
			),
		),
	))
}

func HandleProfile(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	user.Set("name", re.Request.FormValue("name"))
	if err := re.App.Save(user); err != nil {
		return err
	}
	return redirect(re, "/profile")
}
