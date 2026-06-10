package pages

import (
	"fmt"
	"strings"

	"mljr-web/ui/layout"
	"mljr-web/ui/overlay"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func Catalogue() g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title: "mljr-ui — showcase",
			Theme: token.ThemeSwissBrut,
			Mode:  token.ModeLight,
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),

		layout.Navbar(layout.NavbarProps{},
			g.Text("mljr-ui · showcase"),
			g.Group{
				h.A(h.Href("/"), g.Text("Catalogue")),
				h.A(h.Href("/patterns"), g.Text("Patterns")),
			},
			g.Group{
				special.ThemeToggle(),
				special.ModeToggle(),
			},
		),

		// Scroll restoration: save position on navigate, restore on back
		h.Script(g.Raw(`(function(){
  var k='mljr-cat-scroll';
  var saved=sessionStorage.getItem(k);
  if(saved){window.scrollTo(0,parseInt(saved));sessionStorage.removeItem(k);}
  document.querySelectorAll('a[href^="/components/"]').forEach(function(a){
    a.addEventListener('click',function(){sessionStorage.setItem(k,window.scrollY);});
  });
})();`)),

		h.Main(
			layout.Container(layout.ContainerProps{},
				primitive.Display(primitive.DisplayProps{},
					g.Text("Component "),
					h.Em(g.Text("catalogue")),
				),
				h.Div(
					g.Attr("style", "display:flex;align-items:center;gap:var(--sp-4);flex-wrap:wrap;margin-bottom:var(--sp-6)"),
					h.P(
						g.Attr("style", "margin:0;flex:1"),
						g.Text("Every registered component, rendered with default props. Cycle theme and mode to verify all four skins."),
					),
					h.Input(
						g.Attr("data-component", "input"),
						h.Type("search"),
						h.Placeholder("Search components…"),
						g.Attr("data-bind:q"),
						g.Attr("style", "width:220px"),
					),
				),
				h.Div(g.Attr("data-signals", `{"q":""}`)),
				stackIntro(),
				categories(),
			),
		),

		layout.Footer(layout.FooterProps{},
			h.Div(g.Text("mljr-ui showcase · build tag: showcase")),
		),

		overlay.Toaster(overlay.ToasterProps{}),
		overlay.Portal("portal"),
	)
}

func stackIntro() g.Node {
	type pill struct {
		name    string
		version string
		role    string
		detail  string
		color   string
	}
	pills := []pill{
		{"Go + gomponents", "1.23", "Server rendering", "Type-safe HTML components. No templates, no JSX. Server owns the HTML.", "#00ADD8"},
		{"Datastar", "1.0.2", "Hypermedia reactivity", "Signal-based interactivity over SSE. No virtual DOM, no hydration step.", "var(--accent)"},
		{"Tailwind", "v4", "Design system", "4 themes × 2 modes. CSS-only, zero runtime. Rebuilt per component category.", "var(--success)"},
		{"Motion", "v10", "Animations", "24 KB WAAPI-based library. Spring physics, stagger, inView, timeline. No React.", "var(--warning)"},
	}

	pillNodes := make([]g.Node, len(pills))
	for i, p := range pills {
		pillNodes[i] = h.Div(
			h.Style("flex:1;min-width:180px;padding:var(--sp-4);border:var(--bw-2) solid var(--line);border-radius:var(--radius);border-top:4px solid "+p.color+";background:var(--surface-2);display:flex;flex-direction:column;gap:var(--sp-2)"),
			h.Div(
				h.Style("display:flex;align-items:baseline;gap:var(--sp-2)"),
				h.Span(h.Style("font-weight:900;font-family:var(--font-display)"), g.Text(p.name)),
				h.Code(h.Style("font-size:var(--t-xs);opacity:.5"), g.Text(p.version)),
			),
			h.Span(h.Style("font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:.06em;color:"+p.color), g.Text(p.role)),
			h.P(h.Style("font-size:var(--t-sm);color:var(--muted);margin:0"), g.Text(p.detail)),
		)
	}

	return h.Div(
		h.Style("margin:var(--sp-8) 0;border:var(--bw-2) solid var(--line);border-radius:var(--radius);overflow:hidden"),

		// Header strip
		h.Div(
			h.Style("background:var(--line);padding:var(--sp-3) var(--sp-5);display:flex;align-items:center;gap:var(--sp-3)"),
			h.Span(h.Style("font-size:var(--t-xs);font-weight:900;text-transform:uppercase;letter-spacing:.1em"), g.Text("Stack")),
			h.Div(h.Style("width:1px;height:1em;background:currentColor;opacity:.3")),
			h.Span(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text("What this is, what it does, and why")),
		),

		h.Div(
			h.Style("padding:var(--sp-5) var(--sp-5) var(--sp-4)"),

			// Philosophy
			h.Div(
				h.Style("margin-bottom:var(--sp-5);padding-bottom:var(--sp-5);border-bottom:var(--bw-1) dashed var(--line)"),
				h.P(
					h.Style("font-size:var(--t-lg);font-weight:700;margin:0 0 var(--sp-3);line-height:1.3"),
					g.Text("Server renders HTML. Datastar adds signals. Motion handles animations. "),
					h.Span(h.Style("opacity:.5"), g.Text("No JS framework. No hydration. No CDN.")),
				),
				h.P(
					h.Style("font-size:var(--t-sm);color:var(--muted);margin:0;max-width:64ch;line-height:1.6"),
					g.Text("mljr-ui is a production Go component library built for speed and minimal client-side footprint. "+
						"Each component is a pure Go function that renders typed HTML. "+
						"Reactivity comes from Datastar's 14 KB SSE-driven signal layer — not a SPA framework. "+
						"The total client JS budget is under 40 KB (Datastar 14 KB + Motion 24 KB), and every byte earns its place."),
				),
			),

			// Stack cards
			h.Div(
				h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-3);margin-bottom:var(--sp-5)"),
				g.Group(pillNodes),
			),

			// Flow
			h.Div(
				h.Style("padding:var(--sp-4);background:var(--surface-2);border-radius:var(--radius);border:var(--bw-1) dashed var(--line)"),
				h.P(h.Style("font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"), g.Text("Request flow")),
				h.Div(
					h.Style("display:flex;align-items:center;gap:var(--sp-2);flex-wrap:wrap;font-size:var(--t-sm);font-weight:700"),
					flowStep("Browser", "var(--line)"),
					flowArrow(),
					flowStep("Go handler", "var(--primary)"),
					flowArrow(),
					flowStep("gomponents HTML", "var(--primary)"),
					flowArrow(),
					flowStep("Datastar signals", "var(--accent)"),
					flowArrow(),
					flowStep("SSE patch", "var(--accent)"),
					flowArrow(),
					flowStep("Motion animate", "var(--warning)"),
				),
			),
		),
	)
}

func flowStep(label, color string) g.Node {
	return h.Span(
		h.Style("padding:var(--sp-1) var(--sp-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius);font-size:var(--t-xs);white-space:nowrap;border-left:3px solid "+color),
		g.Text(label),
	)
}

func flowArrow() g.Node {
	return h.Span(h.Style("opacity:.4;font-size:var(--t-xs)"), g.Text("→"))
}

func categories() g.Node {
	cats := registry.Categories()
	sections := make([]g.Node, 0, len(cats))
	for _, cat := range cats {
		sections = append(sections, categorySection(cat))
	}
	return g.Group(sections)
}

func categorySection(cat string) g.Node {
	all := registry.All()
	cards := []g.Node{}
	for _, c := range all {
		if c.Category != cat {
			continue
		}
		cards = append(cards, componentCard(c))
	}
	title := strings.ToUpper(cat[:1]) + cat[1:]
	return g.El("details",
		g.Attr("open", ""),
		g.Attr("style", "margin-top:var(--sp-6)"),
		g.El("summary",
			g.Attr("style", "cursor:pointer;list-style:none;display:flex;align-items:center;gap:var(--sp-3);margin-bottom:var(--sp-4)"),
			primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text(title)),
			h.Span(
				g.Attr("style", "font-size:var(--t-sm);opacity:.45"),
				g.Text(fmt.Sprintf("%d components", len(cards))),
			),
		),
		layout.Grid(layout.GridProps{}, cards...),
	)
}

func componentCard(c *registry.Component) g.Node {
	searchable := strings.ToLower(c.Name + " " + c.Category + " " + c.Slug)
	return layout.Col(layout.ColProps{Span: 6},
		g.Attr("data-show", fmt.Sprintf("$q===''||'%s'.includes($q.toLowerCase())", searchable)),
		primitive.Card(primitive.CardProps{},
			primitive.Heading(primitive.HeadingProps{Level: 3}, g.Text(c.Name)),
			h.P(g.Attr("style", "color:var(--muted)"), g.Text(c.Summary)),
			// Render each component in its own iframe so Motion, Datastar, and CSS
			// are fully isolated — avoids missing Motion, style bleed, and scroll issues.
			// data-attr reactively rebuilds src when $theme/$mode signals change.
			g.El("iframe",
				g.Attr("data-attr", fmt.Sprintf(`{"src":"/components/%s/preview?theme="+$theme+"&mode="+$mode}`, c.Slug)),
				g.Attr("style", "width:100%;height:220px;border:none;border-radius:var(--radius);display:block;background:var(--bg)"),
				g.Attr("loading", "lazy"),
				g.Attr("scrolling", "no"),
				g.Attr("tabindex", "-1"),
			),
			h.Div(
				g.Attr("style", "display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:var(--sp-2)"),
				primitive.Tag(primitive.TagProps{}, g.Text(c.Category+" · "+c.Slug)),
				h.A(
					h.Href("/components/"+c.Slug),
					g.Attr("data-component", "button"),
					g.Attr("data-variant", "outline"),
					g.Attr("data-size", "sm"),
					g.Text("View →"),
				),
			),
		),
	)
}
