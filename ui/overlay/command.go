package overlay

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CommandItem struct {
	Label    string
	Icon     string // icon name e.g. "lucide:search"
	Shortcut string // display shortcut e.g. "⌘K"
	Href     string // if set, item is a link
	OnClick  string // JS expression on click
	Group    string // group header label
}

type CommandProps struct {
	SignalName  string // default "_cmdOpen"
	Placeholder string
	Items       []CommandItem
}

// Command renders a ⌘K-style command palette overlay. Open it by setting
// the Datastar signal to true, or via the keyboard shortcut (Ctrl/⌘+K).
// Add data-on:click="$_cmdOpen=true" to any trigger button.
func Command(p CommandProps) g.Node {
	if p.SignalName == "" {
		p.SignalName = "_cmdOpen"
	}
	if p.Placeholder == "" {
		p.Placeholder = "Type a command or search…"
	}
	sig := p.SignalName
	closeExpr := "$" + sig + "=false"

	// Group items
	type group struct {
		label string
		items []CommandItem
	}
	groupMap := map[string]*group{}
	groupOrder := []string{}
	for _, item := range p.Items {
		g := item.Group
		if _, ok := groupMap[g]; !ok {
			groupMap[g] = &group{label: g}
			groupOrder = append(groupOrder, g)
		}
		groupMap[g].items = append(groupMap[g].items, item)
	}

	var itemNodes []g.Node
	for _, gk := range groupOrder {
		gr := groupMap[gk]
		if gr.label != "" {
			itemNodes = append(itemNodes, h.Li(
				g.Attr("data-slot", "group-label"),
				g.Attr("role", "presentation"),
				g.Text(gr.label),
			))
		}
		for _, item := range gr.items {
			var inner []g.Node
			if item.Icon != "" {
				inner = append(inner, icon.Icon(item.Icon))
			}
			inner = append(inner, h.Span(g.Attr("data-slot", "cmd-label"), g.Text(item.Label)))
			if item.Shortcut != "" {
				inner = append(inner, h.Kbd(g.Attr("data-component", "kbd"), g.Text(item.Shortcut)))
			}
			var el g.Node
			if item.Href != "" {
				el = h.A(
					g.Attr("data-slot", "item"),
					h.Href(item.Href),
					g.Attr("data-on:click", closeExpr),
					g.Group(inner),
				)
			} else {
				clickExpr := closeExpr
				if item.OnClick != "" {
					clickExpr = item.OnClick + ";" + closeExpr
				}
				el = h.Button(
					g.Attr("data-slot", "item"),
					h.Type("button"),
					g.Attr("data-on:click", clickExpr),
					g.Group(inner),
				)
			}
			itemNodes = append(itemNodes, h.Li(
				g.Attr("role", "option"),
				g.Attr("data-component", "cmd-option"),
				el,
			))
		}
	}

	return h.Div(
		g.Attr("data-component", "command"),
		g.Attr("data-signals", `{"`+sig+`":false,"_cmdQ":""}`),
		g.Attr("data-show", "$"+sig),
		h.Style("display:none"),

		// Backdrop
		h.Div(
			g.Attr("data-slot", "backdrop"),
			g.Attr("data-on:click", closeExpr),
		),

		// Panel
		h.Div(
			g.Attr("data-slot", "panel"),
			h.Div(
				g.Attr("data-slot", "search"),
				icon.Icon("lucide:search"),
				h.Input(
					h.ID("cmd-input"),
					h.Type("search"),
					h.Placeholder(p.Placeholder),
					g.Attr("data-bind:_cmdQ"),
					g.Attr("autocomplete", "off"),
					g.Attr("spellcheck", "false"),
				),
				h.Kbd(g.Attr("data-component", "kbd"), g.Attr("data-on:click", closeExpr), g.Text("Esc")),
			),
			h.Ul(
				g.Attr("data-slot", "list"),
				h.Role("listbox"),
				g.Group(itemNodes),
			),
		),

		// Keyboard shortcut: ⌘K / Ctrl+K
		h.Script(g.Raw(`(function(){
  var sig='`+sig+`';
  document.addEventListener('keydown',function(e){
    if((e.metaKey||e.ctrlKey)&&e.key==='k'){
      e.preventDefault();
      // Toggle via Datastar store if available
      if(window.__ds){ var s={}; s[sig]=!window.__ds.store[sig]; window.__ds.store=Object.assign({},window.__ds.store,s); }
    }
    if(e.key==='Escape'){
      if(window.__ds){ var s={}; s[sig]=false; window.__ds.store=Object.assign({},window.__ds.store,s); }
    }
  });
  // Focus input when opened
  var observer=new MutationObserver(function(){
    var panel=document.querySelector('[data-component="command"] [data-slot="panel"]');
    var input=document.getElementById('cmd-input');
    if(panel&&input&&panel.offsetParent!==null){ input.focus(); }
  });
  observer.observe(document.body,{childList:true,subtree:true,attributes:true,attributeFilter:['style']});
  // Filter list items
  document.addEventListener('input',function(e){
    if(e.target.id!=='cmd-input') return;
    var q=e.target.value.toLowerCase();
    document.querySelectorAll('[data-component="cmd-option"]').forEach(function(li){
      var label=li.querySelector('[data-slot="cmd-label"]');
      if(!label) return;
      li.hidden=q&&!label.textContent.toLowerCase().includes(q);
    });
    // Hide group labels with no visible items
    document.querySelectorAll('[data-slot="group-label"]').forEach(function(gl){
      var next=gl.nextElementSibling;
      var hasVisible=false;
      while(next&&next.getAttribute('data-slot')!=='group-label'){
        if(!next.hidden) hasVisible=true;
        next=next.nextElementSibling;
      }
      gl.hidden=!hasVisible;
    });
  });
})()`)),
	)
}
