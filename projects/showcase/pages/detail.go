//go:build showcase

package pages

import (
	"fmt"
	"net/url"
	"strings"

	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/special"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func ComponentDetail(c *registry.Component) g.Node {
	return layout.PageShell(
		layout.PageProps{
			Title: c.Name + " — mljr-ui showcase",
			Theme: token.ThemeSwissBrut,
			Mode:  token.ModeLight,
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),
		h.Div(g.Attr("data-signals", buildDetailSignals(c))),

		layout.Navbar(layout.NavbarProps{},
			h.A(h.Href("/"), g.Text("mljr-ui · showcase")),
			g.Group{h.A(h.Href("/"), g.Text("← Catalogue"))},
			g.Group{special.ThemeToggle(), special.ModeToggle()},
		),

		h.Div(
			g.Attr("class", "detail-outer-grid"),
			g.Attr("style", "display:grid;grid-template-columns:220px 1fr;min-height:calc(100vh - 56px)"),
			componentSidebar(c.Slug),
			h.Main(
				g.Attr("style", "min-width:0;padding:var(--sp-6) clamp(var(--sp-3),4vw,var(--sp-6)) var(--sp-9)"),
				detailHeader(c),
				detailBody(c),
				codeSection(c),
				combinationsSection(c),
				propsTable(c),
			),
		),

		layout.Footer(layout.FooterProps{},
			h.Div(g.Text("mljr-ui showcase · "+c.Category+" / "+c.Slug)),
		),

		// Keyboard navigation: ← prev component, → next component
		h.Script(g.Raw(`(function(){
  var links=Array.from(document.querySelectorAll('nav[aria-label="Component navigation"] a'));
  var cur=location.pathname;
  var idx=links.findIndex(function(l){ return l.pathname===cur; });
  document.addEventListener('keydown',function(e){
    if(e.target.tagName==='INPUT'||e.target.tagName==='TEXTAREA'||e.target.tagName==='SELECT') return;
    if(e.key==='ArrowLeft'&&idx>0){ location.href=links[idx-1].href; }
    if(e.key==='ArrowRight'&&idx<links.length-1){ location.href=links[idx+1].href; }
  });
})()`)),
	)
}

func componentSidebar(activeSlug string) g.Node {
	all := registry.All()
	cats := registry.Categories()

	catSections := make([]g.Node, 0, len(cats))
	for _, cat := range cats {
		items := make([]g.Node, 0)
		for _, comp := range all {
			if comp.Category != cat {
				continue
			}
			isActive := comp.Slug == activeSlug
			itemStyle := "display:block;padding:var(--sp-1) var(--sp-3);font-size:var(--t-sm);border-radius:var(--radius);text-decoration:none;color:var(--muted);transition:background var(--fade),color var(--fade)"
			if isActive {
				itemStyle += ";background:var(--accent);color:var(--accent-ink);font-weight:700"
			}
			items = append(items, h.A(
				h.Href("/components/"+comp.Slug),
				g.Attr("style", itemStyle),
				g.Text(comp.Name),
			))
		}
		catSections = append(catSections,
			h.Div(
				g.Attr("style", "margin-bottom:var(--sp-4)"),
				h.Div(
					g.Attr("style", "font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:var(--tracking-wide);opacity:.45;padding:var(--sp-2) var(--sp-3);margin-bottom:var(--sp-1)"),
					g.Text(strings.ToUpper(cat[:1])+cat[1:]),
				),
				g.Group(items),
			),
		)
	}

	return h.Nav(
		g.Attr("aria-label", "Component navigation"),
		g.Attr("style", "border-right:var(--bw-1) solid var(--line);padding:var(--sp-5) var(--sp-2);overflow-y:auto;position:sticky;top:56px;height:calc(100vh - 56px);background:var(--bg)"),
		g.Group(catSections),
	)
}

func detailHeader(c *registry.Component) g.Node {
	return h.Div(
		g.Attr("style", "margin-bottom:var(--sp-6)"),
		h.Nav(
			g.Attr("style", "font-size:var(--t-sm);opacity:.6;margin-bottom:var(--sp-3)"),
			h.A(h.Href("/"), g.Text("Catalogue")),
			g.Text(" › "+c.Category+" › "),
			h.Strong(g.Text(c.Name)),
		),
		primitive.Heading(primitive.HeadingProps{Level: 1}, g.Text(c.Name)),
		h.P(
			g.Attr("style", "color:var(--muted);margin-top:var(--sp-2)"),
			g.Text(c.Summary),
		),
		h.Div(
			g.Attr("style", "margin-top:var(--sp-3)"),
			primitive.Tag(primitive.TagProps{}, g.Text(c.Category)),
			g.Text(" "),
			primitive.Tag(primitive.TagProps{}, g.Text(c.Slug)),
		),
	)
}

func detailBody(c *registry.Component) g.Node {
	srcExpr := buildSrcExpr(c)
	defaultSrc := buildDefaultSrcURL(c)
	hasControls := len(c.Controls) > 0
	iframeH := c.PreviewHeight
	if iframeH == "" {
		iframeH = "480px"
	}

	previewPanel := h.Div(
		sizeTabs(),
		// Reactive iframe: data-effect updates src when any signal changes.
		// data-attr handles responsive width changes.
		h.Div(
			g.Attr("data-effect", fmt.Sprintf("(function(){var f=document.getElementById('preview-frame');if(f)f.src=%s;})()", srcExpr)),
		),
		h.Div(
			g.Attr("style", "background:var(--surface-2);border:var(--bw-1) var(--border-style) var(--line);border-radius:var(--radius);overflow:auto;display:flex;justify-content:center;align-items:flex-start;min-height:496px"),
			g.El("iframe",
				h.ID("preview-frame"),
				h.Src(defaultSrc),
				g.Attr("scrolling", "auto"),
				g.Attr("data-attr", fmt.Sprintf(`{"style":'border:none;height:%s;display:block;background:var(--bg);width:'+$previewWidth+';overflow:auto'}`, iframeH)),
				g.Attr("style", fmt.Sprintf("border:none;height:%s;display:block;width:100%%;overflow:auto", iframeH)),
			),
		),
		g.If(hasControls, codeSnippet(c)),
	)

	if !hasControls {
		return h.Div(
			g.Attr("style", "margin-bottom:var(--sp-7)"),
			previewPanel,
		)
	}

	return h.Div(
		g.Attr("class", "detail-controls-grid"),
		g.Attr("style", "display:grid;grid-template-columns:260px 1fr;gap:var(--sp-6);margin-bottom:var(--sp-7);align-items:start"),
		controlsPanel(c),
		previewPanel,
	)
}

func sizeTabs() g.Node {
	return h.Div(
		g.Attr("style", "display:flex;gap:var(--sp-2);margin-bottom:var(--sp-3);align-items:center;flex-wrap:wrap;overflow-x:auto"),
		sizeTabBtn("Desktop", "100%"),
		sizeTabBtn("Tablet", "768px"),
		sizeTabBtn("Mobile", "375px"),
		h.Span(
			g.Attr("style", "margin-left:auto;font-family:var(--font-mono);font-size:var(--t-xs);opacity:.45;padding:var(--sp-1) var(--sp-2);border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
			g.Attr("data-text", `$previewWidth==="100%"?"Full width":$previewWidth`),
		),
	)
}

func sizeTabBtn(label, width string) g.Node {
	return h.Button(
		g.Attr("data-component", "button"),
		g.Attr("data-size", "sm"),
		g.Attr("data-variant", "outline"),
		g.Attr("data-attr", fmt.Sprintf(`{"data-variant":$previewWidth==='%s'?'primary':'outline'}`, width)),
		g.Attr("data-on:click", fmt.Sprintf("$previewWidth='%s'", width)),
		g.Text(label),
	)
}

func controlsPanel(c *registry.Component) g.Node {
	rows := make([]g.Node, len(c.Controls))
	for i, ctl := range c.Controls {
		rows[i] = controlRow(ctl)
	}
	return h.Div(
		g.Attr("style", "display:flex;flex-direction:column;gap:var(--sp-4);padding:var(--sp-4);border:var(--bw-1) var(--border-style) var(--line);border-radius:var(--radius);background:var(--surface);position:sticky;top:var(--sp-4)"),
		h.Div(
			g.Attr("style", "font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:var(--tracking-wide);opacity:.5"),
			g.Text("Controls"),
		),
		g.Group(rows),
	)
}

func controlRow(ctl registry.Control) g.Node {
	return h.Div(
		g.Attr("style", "display:flex;flex-direction:column;gap:var(--sp-1)"),
		h.Span(
			g.Attr("style", "font-size:var(--t-sm);font-weight:600"),
			g.Text(ctl.Name),
		),
		controlInput(ctl),
	)
}

func controlInput(ctl registry.Control) g.Node {
	switch ctl.Type {
	case registry.ControlEnum:
		opts := make([]g.Node, len(ctl.Options))
		for i, o := range ctl.Options {
			label := o
			if label == "" {
				label = "(none)"
			}
			opts[i] = h.Option(h.Value(o), g.Text(label))
		}
		return h.Select(
			g.Attr("data-component", "select"),
			g.Attr("data-bind:"+ctl.Name),
			g.Group(opts),
		)
	case registry.ControlBool:
		return h.Label(
			g.Attr("data-component", "checkbox"),
			h.Input(
				h.Type("checkbox"),
				g.Attr("data-bind:"+ctl.Name),
			),
			h.Span(g.Attr("data-slot", "box")),
			h.Span(g.Attr("data-slot", "label"), g.Text(ctl.Name)),
		)
	default:
		return h.Input(
			g.Attr("data-component", "input"),
			h.Type("text"),
			g.Attr("data-bind:"+ctl.Name),
		)
	}
}

func codeSnippet(c *registry.Component) g.Node {
	rows := make([]g.Node, len(c.Controls))
	for i, ctl := range c.Controls {
		var valNode g.Node
		if ctl.Type == registry.ControlBool {
			valNode = h.Span(g.Attr("data-text", "$"+ctl.Name))
		} else {
			valNode = g.Group{g.Text(`"`), h.Span(g.Attr("data-text", "$"+ctl.Name)), g.Text(`"`)}
		}
		rows[i] = h.Div(
			h.Span(g.Attr("style", "opacity:.45"), g.Text(ctl.Name+": ")),
			valNode,
		)
	}
	return g.El("details",
		g.Attr("style", "margin-top:var(--sp-4)"),
		g.El("summary",
			g.Attr("style", "cursor:pointer;font-size:var(--t-sm);font-weight:600;padding:var(--sp-2) 0;user-select:none"),
			g.Text("Current props"),
		),
		h.Div(
			g.Attr("style", "position:relative;margin-top:var(--sp-2)"),
			h.Pre(
				h.ID("code-snippet"),
				g.Attr("style", "font-family:var(--font-mono);font-size:var(--t-sm);padding:var(--sp-4);background:var(--surface);border:var(--bw-1) solid var(--line);border-radius:var(--radius);margin:0;overflow-x:auto"),
				g.Group(rows),
			),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-size", "sm"),
				g.Attr("data-variant", "outline"),
				g.Attr("style", "position:absolute;top:var(--sp-2);right:var(--sp-2)"),
				g.Attr("data-on:click", "navigator.clipboard.writeText(document.getElementById('code-snippet').innerText)"),
				g.Text("Copy"),
			),
		),
	)
}

func combinationsSection(c *registry.Component) g.Node {
	combos := registry.Combinations(c)
	if len(combos) <= 1 {
		return nil
	}
	cards := make([]g.Node, len(combos))
	for i, combo := range combos {
		label := registry.ComboLabel(c, combo)
		searchable := strings.ToLower(label)
		cards[i] = h.Div(
			g.Attr("data-show", fmt.Sprintf("$cf===''||'%s'.includes($cf.toLowerCase())", searchable)),
			g.Attr("style", "display:flex;flex-direction:column;gap:var(--sp-2)"),
			h.Div(
				g.Attr("style", "font-size:var(--t-xs);opacity:.5;font-family:var(--font-mono);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;padding:0 var(--sp-1)"),
				g.Text(label),
			),
			h.Div(
				g.Attr("style", "padding:var(--sp-3);border:var(--bw-1) dashed var(--line);border-radius:var(--radius);display:flex;flex-wrap:wrap;gap:var(--sp-2);align-items:flex-start;justify-content:center;min-height:80px;background:var(--surface);overflow:hidden;max-height:200px"),
				c.Render(combo),
			),
		)
	}
	return h.Section(
		g.Attr("data-signals", `{"cf":""}`),
		g.Attr("style", "margin-top:var(--sp-7)"),
		h.Div(
			g.Attr("style", "display:flex;align-items:baseline;gap:var(--sp-4);flex-wrap:wrap;margin-bottom:var(--sp-5)"),
			primitive.Heading(primitive.HeadingProps{Level: 2},
				g.Text("All combinations "),
				h.Span(
					g.Attr("style", "font-size:var(--t-base);opacity:.5"),
					g.Text(fmt.Sprintf("(%d)", len(combos))),
				),
			),
			h.Input(
				g.Attr("data-component", "input"),
				h.Type("search"),
				h.Placeholder("Filter…"),
				g.Attr("data-bind:cf"),
				g.Attr("style", "width:160px;flex-shrink:0"),
			),
		),
		h.P(
			g.Attr("style", "color:var(--muted);font-size:var(--t-sm);margin-bottom:var(--sp-5)"),
			g.Text("Enum controls only · bool and text use defaults"),
		),
		h.Div(
			g.Attr("style", "display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:var(--sp-4)"),
			g.Group(cards),
		),
	)
}

func codeSection(c *registry.Component) g.Node {
	if c.Code == "" {
		return nil
	}
	return h.Section(
		g.Attr("style", "margin-top:var(--sp-7)"),
		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Usage")),
		h.Div(
			g.Attr("style", "position:relative;margin-top:var(--sp-4)"),
			h.Pre(
				h.ID("usage-snippet"),
				g.Attr("style", "font-family:var(--font-mono);font-size:var(--t-sm);padding:var(--sp-5);background:var(--surface);border:var(--bw-1) solid var(--line);border-radius:var(--radius);margin:0;overflow-x:auto;line-height:1.6"),
				h.Code(g.Text(c.Code)),
			),
			h.Button(
				g.Attr("data-component", "button"),
				g.Attr("data-size", "sm"),
				g.Attr("data-variant", "outline"),
				g.Attr("style", "position:absolute;top:var(--sp-2);right:var(--sp-2)"),
				g.Attr("data-on:click", "navigator.clipboard.writeText(document.getElementById('usage-snippet').innerText)"),
				g.Text("Copy"),
			),
		),
	)
}

func propsTable(c *registry.Component) g.Node {
	if len(c.Controls) == 0 {
		return nil
	}
	rows := make([]g.Node, len(c.Controls))
	for i, ctl := range c.Controls {
		opts := strings.Join(ctl.Options, ", ")
		if opts == "" {
			opts = "—"
		}
		rows[i] = h.Tr(
			h.Td(g.Attr("style", "padding:var(--sp-2) var(--sp-3)"), h.Code(g.Text(ctl.Name))),
			h.Td(g.Attr("style", "padding:var(--sp-2) var(--sp-3)"), g.Text(string(ctl.Type))),
			h.Td(g.Attr("style", "padding:var(--sp-2) var(--sp-3);max-width:280px"), h.Code(g.Attr("style", "font-size:var(--t-xs);word-break:break-word"), g.Text(opts))),
			h.Td(g.Attr("style", "padding:var(--sp-2) var(--sp-3)"), h.Code(g.Text(ctl.Default))),
		)
	}
	thStyle := "text-align:left;padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-2) solid var(--line)"
	return h.Section(
		g.Attr("style", "margin-top:var(--sp-7);padding-bottom:var(--sp-7)"),
		primitive.Heading(primitive.HeadingProps{Level: 2}, g.Text("Props")),
		h.Div(
			g.Attr("style", "overflow-x:auto;margin-top:var(--sp-4)"),
			h.Table(
				g.Attr("style", "width:100%;border-collapse:collapse;font-size:var(--t-sm)"),
				h.THead(h.Tr(
					h.Th(g.Attr("style", thStyle), g.Text("Name")),
					h.Th(g.Attr("style", thStyle), g.Text("Type")),
					h.Th(g.Attr("style", thStyle), g.Text("Options")),
					h.Th(g.Attr("style", thStyle), g.Text("Default")),
				)),
				h.TBody(g.Group(rows)),
			),
		),
	)
}

func buildDetailSignals(c *registry.Component) string {
	defaults := registry.DefaultProps(c)
	parts := make([]string, 0, len(c.Controls)+1)
	for _, ctl := range c.Controls {
		v := defaults[ctl.Name]
		if ctl.Type == registry.ControlBool {
			parts = append(parts, fmt.Sprintf(`%q:%s`, ctl.Name, v))
		} else {
			parts = append(parts, fmt.Sprintf(`%q:%q`, ctl.Name, v))
		}
	}
	parts = append(parts, `"previewWidth":"100%"`)
	return "{" + strings.Join(parts, ",") + "}"
}

func buildSrcExpr(c *registry.Component) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("'/components/%s/preview", c.Slug))
	first := true
	for _, ctl := range c.Controls {
		sep := "&"
		if first {
			sep = "?"
			first = false
		}
		if ctl.Type == registry.ControlText {
			sb.WriteString(fmt.Sprintf("%s%s='+encodeURIComponent($%s)+'", sep, ctl.Name, ctl.Name))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s='+$%s+'", sep, ctl.Name, ctl.Name))
		}
	}
	if first {
		sb.WriteString("?theme='+$theme+'&mode='+$mode")
	} else {
		sb.WriteString("&theme='+$theme+'&mode='+$mode")
	}
	return sb.String()
}

// buildDefaultSrcURL computes the initial iframe src from default props (server-side).
// This is used as the static src attribute so the iframe loads immediately on page render,
// before Datastar runs the data-effect.
func buildDefaultSrcURL(c *registry.Component) string {
	defaults := registry.DefaultProps(c)
	var parts []string
	for _, ctl := range c.Controls {
		parts = append(parts, ctl.Name+"="+url.QueryEscape(defaults[ctl.Name]))
	}
	parts = append(parts, "theme=swissbrut", "mode=light")
	base := "/components/" + c.Slug + "/preview"
	if len(parts) == 0 {
		return base
	}
	return base + "?" + strings.Join(parts, "&")
}
