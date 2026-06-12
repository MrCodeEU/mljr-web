package overlay

import (
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ContextMenuItem struct {
	Label     string
	Icon      string // lucide icon name
	OnClick   string // JS expression
	Href      string
	Separator bool // renders a divider instead of an item
	Disabled  bool
}

type ContextMenuProps struct {
	ID    string // unique id for this menu (default "ctx-menu")
	Items []ContextMenuItem
}

// ContextMenu renders a right-click context menu.
// Call ContextMenuTrigger(menuID, element) to attach a trigger.
// Or use: data-on:contextmenu__prevent="window._ctxShow('menuID',evt)"
func ContextMenu(p ContextMenuProps) g.Node {
	if p.ID == "" {
		p.ID = "ctx-menu"
	}

	itemNodes := make([]g.Node, 0, len(p.Items))
	for _, item := range p.Items {
		if item.Separator {
			itemNodes = append(itemNodes, h.Div(g.Attr("data-slot", "separator")))
			continue
		}
		inner := []g.Node{}
		if item.Icon != "" {
			inner = append(inner, icon.Icon(item.Icon, icon.Props{Size: "1rem"}))
		}
		inner = append(inner, h.Span(g.Text(item.Label)))

		var el g.Node
		if item.Href != "" {
			el = h.A(
				g.Attr("data-slot", "item"),
				h.Href(item.Href),
				g.Group(inner),
			)
		} else {
			clickExpr := item.OnClick
			if clickExpr == "" {
				clickExpr = ""
			}
			el = h.Button(
				g.Attr("data-slot", "item"),
				h.Type("button"),
				g.If(item.Disabled, g.Attr("disabled", "")),
				g.If(clickExpr != "", g.Attr("onclick", "window._ctxClose(this);"+clickExpr)),
				g.Group(inner),
			)
		}
		itemNodes = append(itemNodes, h.Li(
			g.Attr("role", "menuitem"),
			g.If(item.Disabled, g.Attr("aria-disabled", "true")),
			el,
		))
	}

	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "context-menu"),
		h.Role("menu"),
		h.Style("display:none"),
		h.Ul(g.Attr("data-slot", "list"), g.Group(itemNodes)),
		h.Script(g.Raw(`(function(){
  var menu=document.getElementById('`+p.ID+`');
  if(!menu) return;
  function hide(){ menu.style.display='none'; }
  window._ctxClose=window._ctxClose||function(){ hide(); };
  // Attach to any [data-ctx="`+p.ID+`"] elements
  document.querySelectorAll('[data-ctx="`+p.ID+`"]').forEach(function(el){
    el.addEventListener('contextmenu',function(e){
      e.preventDefault();
      menu.style.display='block';
      var x=Math.min(e.clientX,window.innerWidth-menu.offsetWidth-8);
      var y=Math.min(e.clientY,window.innerHeight-menu.offsetHeight-8);
      menu.style.left=x+'px'; menu.style.top=y+'px';
    });
  });
  document.addEventListener('click',function(e){ if(!menu.contains(e.target)) hide(); });
  document.addEventListener('keydown',function(e){ if(e.key==='Escape') hide(); });
  document.addEventListener('contextmenu',function(e){ if(!e.target.closest('[data-ctx="`+p.ID+`"]')) hide(); });
})()`)),
	)
}
