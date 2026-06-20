package pages

import (
	"strings"

	"mljr-web/ui/layout"
	"mljr-web/ui/overlay"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
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

// answerTone returns a user's chosen favorite_color tone, falling back to
// the same deterministic id-derived tone the avatar already uses when no
// color has been picked — keeps answer color-coding consistent with the
// avatar's existing fallback behavior.
func answerTone(user *core.Record) token.Tone {
	if c := user.GetString("favorite_color"); c != "" {
		return token.Tone(c)
	}
	return avatarTone(user.Id)
}

// userAvatarSrc returns the protected file URL for a user's uploaded
// avatar, or "" if none is set (callers fall back to the color+initials
// avatar) or if minting a file token fails. The avatar field is Protected,
// so the URL needs a file token scoped to the viewing user themselves —
// only ever called for a user's own avatar (profile nav, profile page), so
// minting the token off the same record being viewed always satisfies the
// users collection's "id = @request.auth.id" ViewRule.
func userAvatarSrc(user *core.Record) string {
	filename := user.GetString("avatar")
	if filename == "" {
		return ""
	}
	token, err := user.NewFileToken()
	if err != nil {
		return ""
	}
	return "/api/files/" + user.Collection().Id + "/" + user.Id + "/" + filename + "?token=" + token
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

func brand(t func(string, ...any) string) g.Node {
	return h.A(h.Href("/"),
		h.Style("font-family:var(--font-display);font-weight:700;font-size:var(--t-lg);color:var(--ink);text-decoration:none"),
		g.Text(t("newsletter.nav.brand")),
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
func groupsNav(t func(string, ...any) string, groups []groupSummary, activeSlug string) g.Node {
	triggerLabel := t("newsletter.nav.groups_label")
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
			g.If(gr.Slug == activeSlug, h.Span(h.Style("margin-left:auto;font-size:var(--t-xs);color:var(--accent);font-weight:700;white-space:nowrap"), g.Text(t("newsletter.nav.current_badge")))),
		)
		items = append(items, overlay.DropdownItem{Content: row, Href: "/g/" + gr.Slug})
	}
	if len(items) == 0 {
		items = append(items, overlay.DropdownItem{Label: t("newsletter.nav.groups_empty"), Href: "/"})
	}
	items = append(items, overlay.DropdownItem{Divider: true, Label: t("newsletter.nav.groups_back"), Href: "/", Icon: "lucide:layout-grid"})

	return overlay.Dropdown(overlay.DropdownProps{Signal: "navGroups", Align: "right"}, trigger, items...)
}

// profileNav renders the current user's avatar with a menu for the profile
// page and logging out. The logout item submits a hidden plain <form> via
// JS so the existing POST /logout route stays untouched.
func profileNav(t func(string, ...any) string, user *core.Record) g.Node {
	trigger := primitive.Avatar(primitive.AvatarProps{
		Src:      userAvatarSrc(user),
		Initials: initials(displayName(user)),
		Tone:     avatarTone(user.Id),
		Size:     token.SizeSM,
	})

	return h.Div(
		h.Style("display:flex;align-items:center"),
		h.Form(h.ID("nl-logout-form"), h.Method("post"), h.Action("/logout"), h.Style("display:none")),
		overlay.Dropdown(overlay.DropdownProps{Signal: "navUser", Align: "right"}, trigger,
			overlay.DropdownItem{Label: displayName(user), Icon: "lucide:user", Href: "/profile"},
			overlay.DropdownItem{Divider: true, Label: t("newsletter.nav.logout"), Icon: "lucide:log-out", Variant: "danger",
				OnClick: "document.getElementById('nl-logout-form').requestSubmit()"},
		),
	)
}

// appHeader is the shared top nav for every page: brand, dashboard link, a
// groups switcher, and the user's profile menu. activeGroupSlug highlights
// the currently-open group, if any ("" on pages with no single group
// context, like the dashboard, or when unauthenticated).
func appHeader(re *core.RequestEvent, activeGroupSlug string) g.Node {
	t := translator(re)
	langToggle := special.LanguageToggle(special.LanguageToggleProps{
		Languages: []special.Language{
			{Code: "en", Label: "EN", Title: "English", Flag: "circle-flags:gb"},
			{Code: "de", Label: "DE", Title: "Deutsch", Flag: "circle-flags:de"},
		},
		Current:        currentLang(re),
		ReloadOnChange: true,
	})

	user := currentUser(re)
	if user == nil {
		return layout.Navbar(layout.NavbarProps{}, brand(t), g.Group(nil), langToggle)
	}

	groups := userGroups(re, user)
	nav := h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-5)"),
		h.A(h.Href("/"), h.Style("color:inherit;text-decoration:none;font-weight:600"), g.Text(t("newsletter.nav.dashboard"))),
		groupsNav(t, groups, activeGroupSlug),
	)

	right := h.Div(h.Style("display:flex;align-items:center;gap:var(--sp-3)"),
		langToggle,
		notificationsNav(re, user.Id),
		profileNav(t, user),
	)
	return layout.Navbar(layout.NavbarProps{}, brand(t), nav, right)
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

// groupSubNav renders the Editions/Recap/Questions/Suggestions/Invites/
// Settings pill nav shared by every group-scoped page, deriving the active
// tab from the request path so call sites don't each have to say which tab
// is current.
func groupSubNav(re *core.RequestEvent, slug string) g.Node {
	t := translator(re)
	path := re.Request.URL.Path
	items := []struct {
		label, suffix string
	}{
		{t("newsletter.subnav.editions"), "/editions"},
		{t("newsletter.subnav.recap"), "/recap"},
		{t("newsletter.subnav.questions"), "/questions"},
		{t("newsletter.subnav.suggestions"), "/suggestions"},
		{t("newsletter.subnav.invites"), "/invites"},
		{t("newsletter.subnav.settings"), "/settings"},
	}
	base := "/g/" + slug
	pills := make([]layout.NavPillItem, len(items))
	for i, it := range items {
		pills[i] = layout.NavPillItem{
			Label:  it.label,
			Href:   base + it.suffix,
			Active: strings.Contains(path, it.suffix),
		}
	}
	return layout.NavPills(pills)
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
				g.If(activeGroupSlug != "", groupSubNav(re, activeGroupSlug)),
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
