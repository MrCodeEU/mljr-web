package pages

import (
	"mljr-web/ui/layout"
	"mljr-web/ui/overlay"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	"github.com/pocketbase/pocketbase/core"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func displayName(u *core.Record) string {
	if name := u.GetString("name"); name != "" {
		return name
	}
	return u.Email()
}

// avatarTones is the fixed palette random avatars are drawn from until real
// uploads exist — picked deterministically per id/name so the same
// user/group always gets the same color.
var avatarTones = []token.Tone{
	token.ToneYellow, token.ToneLime, token.ToneCyan, token.ToneViolet,
	token.TonePink, token.ToneSky, token.ToneMint, token.ToneBlush,
}

func avatarTone(seed string) token.Tone {
	var sum int
	for _, b := range []byte(seed) {
		sum += int(b)
	}
	return avatarTones[sum%len(avatarTones)]
}

// initials computes a 1-2 letter avatar fallback from a display name.
func initials(name string) string {
	var letters []rune
	word := true
	for _, r := range name {
		if r == ' ' {
			word = true
			continue
		}
		if word {
			letters = append(letters, r)
			word = false
			if len(letters) == 2 {
				break
			}
		}
	}
	if len(letters) == 0 {
		return "?"
	}
	return string(letters)
}

func brand() g.Node {
	return h.A(h.Href("/"),
		h.Style("font-family:var(--font-display);font-weight:700;font-size:var(--t-lg);color:var(--ink);text-decoration:none"),
		g.Text("Newsletter"),
	)
}

// groupSummary is the minimal info needed to render a group in the nav.
type groupSummary struct {
	Name string
	Slug string
	Role string
}

func userGroups(re *core.RequestEvent, user *core.Record) []groupSummary {
	memberships, err := re.App.FindRecordsByFilter(
		"group_memberships", "user = {:user}", "-created", 0, 0,
		map[string]any{"user": user.Id},
	)
	if err != nil {
		return nil
	}
	var out []groupSummary
	for _, m := range memberships {
		group, err := re.App.FindRecordById("groups", m.GetString("group"))
		if err != nil {
			continue
		}
		out = append(out, groupSummary{
			Name: group.GetString("name"),
			Slug: group.GetString("slug"),
			Role: m.GetString("role"),
		})
	}
	return out
}

func groupAvatar(seed, name string) g.Node {
	return primitive.Avatar(primitive.AvatarProps{
		Initials: initials(name),
		Tone:     avatarTone(seed),
		Size:     token.SizeSM,
	})
}

// groupsNav renders the "switch group" dropdown: every group the user is in,
// each with its avatar, name and role, with the currently-open one marked.
func groupsNav(groups []groupSummary, activeSlug string) g.Node {
	triggerLabel := "Groups"
	for _, gr := range groups {
		if gr.Slug == activeSlug {
			triggerLabel = gr.Name
		}
	}

	trigger := h.Button(
		h.Type("button"),
		h.Style("display:flex;align-items:center;gap:var(--sp-2);background:none;border:none;cursor:pointer;font:inherit;color:inherit;padding:0"),
		g.Text(triggerLabel),
		h.Span(h.Style("font-size:.7em"), g.Text("▾")),
	)

	var items []overlay.DropdownItem
	for _, gr := range groups {
		row := h.Div(
			h.Style("display:flex;align-items:center;gap:var(--sp-3);min-width:0;max-width:240px"),
			groupAvatar(gr.Slug, gr.Name),
			h.Div(
				h.Style("min-width:0;overflow:hidden"),
				h.Span(h.Style("display:block;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis"), g.Text(gr.Name)),
				h.Span(h.Style("display:block;font-size:var(--t-xs);color:var(--muted);text-transform:uppercase;letter-spacing:.04em"), g.Text(gr.Role)),
			),
			g.If(gr.Slug == activeSlug, h.Span(h.Style("margin-left:auto;font-size:var(--t-xs);color:var(--accent);font-weight:700;white-space:nowrap"), g.Text("CURRENT"))),
		)
		items = append(items, overlay.DropdownItem{Content: row, Href: "/g/" + gr.Slug})
	}
	if len(items) == 0 {
		items = append(items, overlay.DropdownItem{Label: "You're not in any groups yet", Href: "/"})
	}
	items = append(items, overlay.DropdownItem{Divider: true, Label: "Back to dashboard", Href: "/", Icon: "lucide:layout-grid"})

	return overlay.Dropdown(overlay.DropdownProps{Signal: "navGroups", Align: "right"}, trigger, items...)
}

// profileNav renders the current user's avatar with a menu for the profile
// page and logging out. The logout item submits a hidden plain <form> via
// JS so the existing POST /logout route stays untouched.
func profileNav(user *core.Record) g.Node {
	trigger := primitive.Avatar(primitive.AvatarProps{
		Initials: initials(displayName(user)),
		Tone:     avatarTone(user.Id),
		Size:     token.SizeSM,
	})

	return h.Div(
		h.Style("display:flex;align-items:center"),
		h.Form(h.ID("nl-logout-form"), h.Method("post"), h.Action("/logout"), h.Style("display:none")),
		overlay.Dropdown(overlay.DropdownProps{Signal: "navUser", Align: "right"}, trigger,
			overlay.DropdownItem{Label: displayName(user), Icon: "lucide:user", Href: "/profile"},
			overlay.DropdownItem{Divider: true, Label: "Log out", Icon: "lucide:log-out", Variant: "danger",
				OnClick: "document.getElementById('nl-logout-form').requestSubmit()"},
		),
	)
}

// appHeader is the shared top nav for every page: brand, dashboard link, a
// groups switcher, and the user's profile menu. activeGroupSlug highlights
// the currently-open group, if any ("" on pages with no single group
// context, like the dashboard, or when unauthenticated).
func appHeader(re *core.RequestEvent, activeGroupSlug string) g.Node {
	user := currentUser(re)
	if user == nil {
		return layout.Navbar(layout.NavbarProps{}, brand(), g.Group(nil), g.Group(nil))
	}

	groups := userGroups(re, user)
	nav := h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-5)"),
		h.A(h.Href("/"), h.Style("color:inherit;text-decoration:none;font-weight:600"), g.Text("Dashboard")),
		groupsNav(groups, activeGroupSlug),
	)

	return layout.Navbar(layout.NavbarProps{}, brand(), nav, profileNav(user))
}

// breadcrumbItem is one link in a breadcrumb trail; the last item rendered
// is plain text (the current page), the rest are links.
type breadcrumbItem struct {
	Label string
	Href  string
}

// breadcrumbs renders a secondary, always-visible path back up the page
// hierarchy (e.g. "Dashboard / Family Crew / Settings") — a second way to
// navigate alongside the nav links, since header chrome alone isn't always
// noticed.
func breadcrumbs(items []breadcrumbItem) g.Node {
	if len(items) == 0 {
		return g.Group(nil)
	}
	var nodes []g.Node
	for i, item := range items {
		if i > 0 {
			nodes = append(nodes, h.Span(h.Style("color:var(--muted)"), g.Text(" / ")))
		}
		if i == len(items)-1 {
			nodes = append(nodes, h.Span(h.Style("color:var(--muted)"), g.Text(item.Label)))
		} else {
			nodes = append(nodes, h.A(h.Href(item.Href), h.Style("color:var(--ink);text-decoration:none;font-weight:600"), g.Text(item.Label)))
		}
	}
	return h.Nav(
		g.Attr("aria-label", "Breadcrumb"),
		h.Style("max-width:640px;margin:0 auto;padding:var(--sp-4) var(--sp-4) 0;font-size:var(--t-sm)"),
		g.Group(nodes),
	)
}

// appPage wraps page content with the shared header and an optional
// breadcrumb trail inside the common PageShell, so every authenticated
// route gets consistent chrome. activeGroupSlug highlights the open group
// in the nav; pass "" for pages with no single group context.
func appPage(re *core.RequestEvent, activeGroupSlug, title string, crumbs []breadcrumbItem, content ...g.Node) g.Node {
	return layout.PageShell(
		layout.PageProps{Title: title + " — Newsletter", Theme: token.ThemeSwissBrut, Mode: token.ModeLight},
		h.Div(
			appHeader(re, activeGroupSlug),
			breadcrumbs(crumbs),
			h.Main(
				h.Style("max-width:640px;margin:0 auto;padding:var(--sp-8) var(--sp-4)"),
				g.Group(content),
			),
		),
	)
}

// authPage wraps unauthenticated pages (login/signup) with the same brand
// header, minus any nav that requires a logged-in user.
func authPage(re *core.RequestEvent, title string, content ...g.Node) g.Node {
	return layout.PageShell(
		layout.PageProps{Title: title + " — Newsletter", Theme: token.ThemeSwissBrut, Mode: token.ModeLight},
		h.Div(
			appHeader(re, ""),
			h.Main(
				h.Style("max-width:640px;margin:0 auto;padding:var(--sp-8) var(--sp-4)"),
				g.Group(content),
			),
		),
	)
}
