package data

import (
	"fmt"

	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SortableProps struct {
	// ID must be unique per page (default "sortable").
	ID string
	// Handle shows a drag-handle icon instead of the full row being draggable.
	Handle bool
	// OnChange is a JS expression called after a reorder with `(newOrder: string[])`.
	OnChange string
}

// Sortable renders a drag-to-reorder list using the HTML5 drag-and-drop API.
// Each item must be wrapped with SortableItem.
func Sortable(p SortableProps, items ...g.Node) g.Node {
	if p.ID == "" {
		p.ID = "sortable"
	}

	script := fmt.Sprintf(`(function(){
  var list=document.getElementById('%s');
  if(!list) return;
  var dragging=null;
  function getItems(){ return Array.from(list.querySelectorAll('[data-slot="sortable-item"]')); }
  list.addEventListener('dragstart',function(e){
    dragging=e.target.closest('[data-slot="sortable-item"]');
    if(!dragging) return;
    dragging.style.opacity='0.4';
    e.dataTransfer.effectAllowed='move';
  });
  list.addEventListener('dragend',function(){
    if(dragging) dragging.style.opacity='';
    dragging=null;
    getItems().forEach(function(i){ i.classList.remove('drag-over'); });
  });
  list.addEventListener('dragover',function(e){
    e.preventDefault();
    var target=e.target.closest('[data-slot="sortable-item"]');
    if(!target||target===dragging) return;
    var rect=target.getBoundingClientRect();
    var after=(e.clientY-rect.top)>rect.height/2;
    if(after) target.after(dragging); else target.before(dragging);
  });
  list.addEventListener('drop',function(e){
    e.preventDefault();
    %s
  });
})();`,
		p.ID,
		func() string {
			if p.OnChange == "" {
				return ""
			}
			return fmt.Sprintf(`var order=getItems().map(function(i){return i.dataset.value;});(%s)(order);`, p.OnChange)
		}(),
	)

	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "sortable"),
		g.If(p.Handle, g.Attr("data-handle", "true")),
		g.Group(items),
		h.Script(g.Raw(script)),
	)
}

// SortableItem wraps one draggable row. Value is used in the onchange order array.
func SortableItem(value string, handle bool, children ...g.Node) g.Node {
	draggable := "true"
	if handle {
		draggable = "false"
	}

	var handleNode g.Node
	if handle {
		handleNode = h.Span(
			g.Attr("data-slot", "handle"),
			g.Attr("draggable", "true"),
			h.Style("cursor:grab;display:flex;align-items:center;color:var(--muted);padding:0 var(--sp-2)"),
			icon.Icon("lucide:grip-vertical"),
		)
	}

	return h.Div(
		g.Attr("data-slot", "sortable-item"),
		g.Attr("data-value", value),
		g.Attr("draggable", draggable),
		h.Style("display:flex;align-items:center;user-select:none"),
		handleNode,
		h.Div(h.Style("flex:1"), g.Group(children)),
	)
}

// SortableRow is a convenience builder for a simple text row.
func SortableRow(value, label string, handle bool) g.Node {
	return SortableItem(value, handle,
		h.Div(
			h.Style("display:flex;align-items:center;justify-content:space-between;padding:var(--sp-3) var(--sp-4);background:var(--surface);border:var(--bw-1) solid var(--line);border-radius:var(--radius);margin-bottom:var(--sp-1);font-size:var(--t-sm);font-weight:600"),
			g.Text(label),
			icon.Icon("lucide:grip-horizontal", icon.Props{Size: "1rem"}),
		),
	)
}
