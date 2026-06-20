package pages

import (
	"mljr-web/ui/form"
	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// Dashboard lists the groups the current user belongs to, plus a
// create-group form. Auth is required; unauthenticated users are sent to
// /login.
func Dashboard(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	t := translator(re)
	memberships, err := re.App.FindRecordsByFilter(
		"group_memberships", "user = {:user}", "-created", 0, 0,
		map[string]any{"user": user.Id},
	)
	if err != nil {
		return err
	}

	var groupRows []g.Node
	for _, m := range memberships {
		group, err := re.App.FindRecordById("groups", m.GetString("group"))
		if err != nil {
			continue
		}
		groupRows = append(groupRows, primitive.Card(primitive.CardProps{
			Interactive: true,
			Attrs:       []g.Node{h.Style("padding:var(--sp-4)")},
		},
			h.A(
				h.Href("/g/"+group.GetString("slug")),
				h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;text-decoration:none;color:var(--ink)"),
				h.Span(h.Style("font-weight:700;min-width:0;overflow-wrap:anywhere"), g.Text(group.GetString("name"))),
				h.Span(h.Style("color:var(--muted);font-size:var(--t-sm);text-transform:uppercase;letter-spacing:.04em;white-space:nowrap"), g.Text(m.GetString("role"))),
			),
		))
	}

	return renderPage(re, 200, appPage(re, "", t("newsletter.nav.dashboard"), nil,
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(t("newsletter.dashboard.heading"))),
		flashAlert(re),
		g.If(len(groupRows) == 0,
			h.P(h.Style("color:var(--muted);margin:var(--sp-3) 0 var(--sp-6)"), g.Text(t("newsletter.dashboard.empty"))),
		),
		g.If(len(groupRows) > 0,
			h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-3);margin:var(--sp-4) 0 var(--sp-8)"), g.Group(groupRows)),
		),

		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(t("newsletter.dashboard.create_heading"))),
		flashAlert(re),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-3)")}},
			h.Form(
				h.Method("post"), h.Action("/groups"),
				form.Field(form.FieldProps{Label: t("newsletter.dashboard.name_label")},
					form.Input(form.InputProps{Type: "text", Name: "name", Required: true, Placeholder: "e.g. The Weekly Crew"}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-3)")},
				}, icon.Icon("lucide:plus"), g.Text(t("newsletter.dashboard.create_button"))),
			),
		),
	))
}
