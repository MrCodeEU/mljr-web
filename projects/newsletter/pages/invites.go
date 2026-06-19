package pages

import (
	"log"
	"net/mail"
	"strings"

	"mljr-web/ui/form"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	pbmailer "github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/types"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

const inviteValidDays = 7

func requireAdminMembership(re *core.RequestEvent, group *core.Record, user *core.Record) (*core.Record, error) {
	membership, err := findMembership(re, group.Id, user.Id)
	if err != nil || (membership.GetString("role") != "owner" && membership.GetString("role") != "admin") {
		return nil, re.ForbiddenError("only group owners/admins can manage invites", nil)
	}
	return membership, nil
}

// ListInvites shows pending invites + a form to create new ones (owner/admin only).
func ListInvites(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := requireAdminMembership(re, group, user); err != nil {
		return err
	}

	invites, err := re.App.FindRecordsByFilter(
		"group_invites", "group = {:group} && status = \"pending\"", "-created", 0, 0,
		map[string]any{"group": group.Id},
	)
	if err != nil {
		return err
	}

	var rows []g.Node
	for _, inv := range invites {
		rows = append(rows, h.Div(
			h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);justify-content:space-between;align-items:center;padding:var(--sp-3) 0;border-bottom:var(--border-w) var(--border-style) var(--line)"),
			h.Div(
				h.Span(h.Style("font-weight:600;min-width:0;overflow-wrap:anywhere"), g.Text(inv.GetString("email"))),
				h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"), g.Text(inv.GetString("role"))),
			),
			h.Form(h.Method("post"), h.Action("/g/"+slug+"/invites/"+inv.Id+"/revoke"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Tone: token.ToneNone, Type: "submit"}, g.Text("Revoke")),
			),
		))
	}

	return renderPage(re, 200, appPage(re, slug, "Invites — "+group.GetString("name"),
		[]breadcrumbItem{{Label: "Dashboard", Href: "/"}, {Label: group.GetString("name"), Href: "/g/" + slug}, {Label: "Invites"}},
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Invite people")),
		primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
			h.Form(
				h.Method("post"), h.Action("/g/"+slug+"/invites"),
				form.Field(form.FieldProps{Label: "Email"},
					form.Input(form.InputProps{Type: "email", Name: "email", Required: true, Placeholder: "friend@example.com"}),
				),
				form.Field(form.FieldProps{Label: "Role", Attrs: []g.Node{h.Style("margin-top:var(--sp-4)")}},
					form.Select(form.SelectProps{
						Name: "role",
						Options: []form.SelectOption{
							{Value: "member", Label: "Member", Selected: true},
							{Value: "admin", Label: "Admin"},
						},
					}),
				),
				primitive.Button(primitive.ButtonProps{
					Variant: token.Primary,
					Type:    "submit",
					Attrs:   []g.Node{h.Style("margin-top:var(--sp-4)")},
				}, g.Text("Send invite")),
			),
		),
		primitive.Heading(primitive.HeadingProps{Level: 2, Attrs: []g.Node{h.Style("margin-top:var(--sp-8)")}}, g.Text("Pending invites")),
		g.If(len(rows) == 0, h.P(h.Style("color:var(--muted);margin-top:var(--sp-3)"), g.Text("No pending invites."))),
		g.If(len(rows) > 0, primitive.Card(primitive.CardProps{Attrs: []g.Node{h.Style("margin-top:var(--sp-3);padding:var(--sp-2) var(--sp-4)")}}, g.Group(rows))),
	))
}

func HandleCreateInvite(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := requireAdminMembership(re, group, user); err != nil {
		return err
	}

	email := re.Request.FormValue("email")
	role := re.Request.FormValue("role")
	if role != "admin" {
		role = "member"
	}
	if email == "" {
		return redirect(re, "/g/"+slug+"/invites")
	}

	col, err := re.App.FindCollectionByNameOrId("group_invites")
	if err != nil {
		return err
	}

	invite := core.NewRecord(col)
	invite.Set("group", group.Id)
	invite.Set("invited_by", user.Id)
	invite.Set("email", email)
	invite.Set("token", randomToken())
	invite.Set("role", role)
	invite.Set("status", "pending")
	invite.Set("expires_at", types.NowDateTime().Add(inviteValidDays*24*60*60*1e9))
	if invited, err := re.App.FindAuthRecordByEmail("users", email); err == nil {
		invite.Set("invited_user", invited.Id)
	}
	if err := re.App.Save(invite); err != nil {
		return err
	}

	if invitedID := invite.GetString("invited_user"); invitedID != "" {
		_ = createNotification(re.App, invitedID, "invite", group.Id, invite.Id, user.Id,
			displayName(user)+" invited you to join "+group.GetString("name"),
			"/invites/"+invite.GetString("token"))
	}

	sendInviteEmail(invite, group, user)

	return redirect(re, "/g/"+slug+"/invites")
}

func sendInviteEmail(invite, group, invitedBy *core.Record) {
	link := publicAppURL + "/invites/" + invite.GetString("token")
	msg := &pbmailer.Message{
		From:    mailFrom,
		To:      []mail.Address{{Address: invite.GetString("email")}},
		Subject: displayName(invitedBy) + " invited you to join " + group.GetString("name"),
		Text: "Hi,\n\n" + displayName(invitedBy) + " invited you to join \"" + group.GetString("name") +
			"\" on the newsletter.\n\n" + link + "\n",
	}
	if err := sendMail(msg); err != nil {
		log.Printf("newsletter invite mail: failed to send to %s: %v", invite.GetString("email"), err)
	}
}

func HandleRevokeInvite(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login")
	}

	slug := re.Request.PathValue("slug")
	group, err := findGroupBySlug(re, slug)
	if err != nil {
		return re.NotFoundError("group not found", err)
	}
	if _, err := requireAdminMembership(re, group, user); err != nil {
		return err
	}

	id := re.Request.PathValue("id")
	invite, err := re.App.FindRecordById("group_invites", id)
	if err != nil || invite.GetString("group") != group.Id {
		return re.NotFoundError("invite not found", err)
	}
	invite.Set("status", "revoked")
	if err := re.App.Save(invite); err != nil {
		return err
	}

	return redirect(re, "/g/"+slug+"/invites")
}

func findInviteByToken(re *core.RequestEvent, token string) (*core.Record, error) {
	return re.App.FindFirstRecordByFilter("group_invites", "token = {:token}", map[string]any{"token": token})
}

// inviteTargetsUser reports whether invite was sent to user, by account
// (if it was linked at creation time) or by email otherwise.
func inviteTargetsUser(invite, user *core.Record) bool {
	if invitedUser := invite.GetString("invited_user"); invitedUser != "" {
		return invitedUser == user.Id
	}
	return strings.EqualFold(invite.GetString("email"), user.Email())
}

// InviteAccept shows the invite ("you've been invited to join X") and, if
// logged in, an accept button; if not, links to signup/login with the
// invite token carried through.
func InviteAccept(re *core.RequestEvent) error {
	inviteToken := re.Request.PathValue("token")
	invite, err := findInviteByToken(re, inviteToken)
	if err != nil {
		return re.NotFoundError("invite not found", err)
	}
	if invite.GetString("status") != "pending" {
		return renderPage(re, 410, authPage(re, "Invite no longer valid",
			primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("This invite is no longer valid")),
			h.P(h.Style("color:var(--muted);margin-top:var(--sp-2)"), g.Text("It may have already been used or revoked.")),
		))
	}
	if invite.GetDateTime("expires_at").Time().Before(types.NowDateTime().Time()) {
		invite.Set("status", "expired")
		_ = re.App.Save(invite)
		return renderPage(re, 410, authPage(re, "Invite expired",
			primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("This invite has expired")),
			h.P(h.Style("color:var(--muted);margin-top:var(--sp-2)"), g.Text("Ask whoever invited you to send a new one.")),
		))
	}

	group, err := re.App.FindRecordById("groups", invite.GetString("group"))
	if err != nil {
		return err
	}

	user := currentUser(re)
	if user == nil {
		return renderPage(re, 200, authPage(re, "Join "+group.GetString("name"),
			primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Join "+group.GetString("name"))),
			h.P(h.Style("color:var(--muted);margin:var(--sp-2) 0 var(--sp-6)"), g.Text("You've been invited as "+invite.GetString("role")+". Sign in or create an account to accept.")),
			h.Div(h.Style("display:flex;gap:var(--sp-3)"),
				primitive.Button(primitive.ButtonProps{Variant: token.Primary, Attrs: []g.Node{h.Type("button")}},
					h.A(h.Href("/signup?invite="+inviteToken), h.Style("color:inherit;text-decoration:none"), g.Text("Create account")),
				),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Attrs: []g.Node{h.Type("button")}},
					h.A(h.Href("/login?invite="+inviteToken), h.Style("color:inherit;text-decoration:none"), g.Text("Sign in instead")),
				),
			),
		))
	}

	if _, err := findMembership(re, group.Id, user.Id); err == nil {
		return redirect(re, "/g/"+group.GetString("slug"))
	}
	if !inviteTargetsUser(invite, user) {
		return renderPage(re, 403, authPage(re, "Invite not for this account",
			primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("This invite isn't for your account")),
			h.P(h.Style("color:var(--muted);margin-top:var(--sp-2)"), g.Text("It was sent to "+invite.GetString("email")+". Sign in with that account to accept it.")),
		))
	}

	return renderPage(re, 200, appPage(re, "", "Join "+group.GetString("name"), nil,
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text("Join "+group.GetString("name"))),
		h.P(h.Style("color:var(--muted);margin:var(--sp-2) 0 var(--sp-6)"), g.Text("You've been invited as "+invite.GetString("role")+".")),
		h.Form(h.Method("post"), h.Action("/invites/"+inviteToken+"/accept"),
			primitive.Button(primitive.ButtonProps{Variant: token.Primary, Type: "submit"}, g.Text("Accept & join")),
		),
	))
}

// acceptInvite creates the membership, marks the invite accepted, and
// notifies the group's admins. Shared by HandleAcceptInvite and the
// signup-via-invite path.
func acceptInvite(re *core.RequestEvent, invite, user *core.Record) (*core.Record, error) {
	group, err := re.App.FindRecordById("groups", invite.GetString("group"))
	if err != nil {
		return nil, err
	}

	if _, err := findMembership(re, group.Id, user.Id); err == nil {
		return group, nil
	}

	membershipsCol, err := re.App.FindCollectionByNameOrId("group_memberships")
	if err != nil {
		return nil, err
	}
	membership := core.NewRecord(membershipsCol)
	membership.Set("group", group.Id)
	membership.Set("user", user.Id)
	role := invite.GetString("role")
	if role == "" {
		role = "member"
	}
	membership.Set("role", role)
	if err := re.App.Save(membership); err != nil {
		return nil, err
	}

	invite.Set("status", "accepted")
	if err := re.App.Save(invite); err != nil {
		return nil, err
	}

	return group, nil
}

func HandleAcceptInvite(re *core.RequestEvent) error {
	user := currentUser(re)
	if user == nil {
		return redirect(re, "/login?invite="+re.Request.PathValue("token"))
	}

	invite, err := findInviteByToken(re, re.Request.PathValue("token"))
	if err != nil || invite.GetString("status") != "pending" {
		return re.NotFoundError("invite not found or already used", err)
	}
	if !inviteTargetsUser(invite, user) {
		return re.ForbiddenError("this invite was not sent to your account", nil)
	}

	group, err := acceptInvite(re, invite, user)
	if err != nil {
		return err
	}

	return redirect(re, "/g/"+group.GetString("slug"))
}
