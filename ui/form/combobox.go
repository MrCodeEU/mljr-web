package form

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ComboboxOption struct {
	Value string
	Label string // display label; defaults to Value
}

type ComboboxProps struct {
	Name        string
	Placeholder string
	Options     []ComboboxOption
	Default     string // pre-selected value
}

// Combobox renders a filterable select input. Options are filtered as the user
// types; arrow keys and Enter navigate/select; Escape closes.
// The hidden input named p.Name carries the selected value for form submission.
func Combobox(p ComboboxProps) g.Node {
	if p.Placeholder == "" {
		p.Placeholder = "Select…"
	}
	defaultLabel := ""
	for _, o := range p.Options {
		lbl := o.Label
		if lbl == "" {
			lbl = o.Value
		}
		if o.Value == p.Default {
			defaultLabel = lbl
		}
	}

	optNodes := make([]g.Node, len(p.Options))
	for i, o := range p.Options {
		lbl := o.Label
		if lbl == "" {
			lbl = o.Value
		}
		optNodes[i] = h.Li(
			g.Attr("role", "option"),
			g.Attr("data-value", o.Value),
			g.Attr("data-label", lbl),
			g.Text(lbl),
		)
	}

	return h.Div(
		g.Attr("data-component", "combobox"),
		g.Attr("role", "combobox"),
		g.Attr("aria-haspopup", "listbox"),
		g.Attr("aria-expanded", "false"),
		h.Input(
			g.Attr("data-slot", "input"),
			g.Attr("data-component", "input"),
			h.Type("text"),
			h.Placeholder(p.Placeholder),
			h.Value(defaultLabel),
			g.Attr("autocomplete", "off"),
			g.Attr("aria-autocomplete", "list"),
		),
		h.Input(g.Attr("data-slot", "value"), h.Type("hidden"), h.Name(p.Name), h.Value(p.Default)),
		h.Ul(
			g.Attr("data-slot", "options"),
			h.Role("listbox"),
			g.Attr("hidden", ""),
			g.Group(optNodes),
		),
		h.Script(g.Raw(comboboxScript)),
	)
}

const comboboxScript = `(function(){
  document.querySelectorAll('[data-component="combobox"]').forEach(function(root){
    var input=root.querySelector('[data-slot="input"]');
    var hidden=root.querySelector('[data-slot="value"]');
    var list=root.querySelector('[data-slot="options"]');
    var items=Array.from(list.querySelectorAll('[data-value]'));
    var cursor=-1;

    function open(){
      list.hidden=false;
      root.setAttribute('aria-expanded','true');
    }
    function close(){
      list.hidden=true;
      root.setAttribute('aria-expanded','false');
      cursor=-1;
    }
    function select(li){
      hidden.value=li.dataset.value;
      input.value=li.dataset.label;
      close();
    }
    function filter(q){
      q=q.toLowerCase();
      var vis=[];
      items.forEach(function(li){
        var match=!q||li.dataset.label.toLowerCase().includes(q)||li.dataset.value.toLowerCase().includes(q);
        li.hidden=!match;
        if(match) vis.push(li);
      });
      if(vis.length) open(); else close();
      cursor=-1;
    }
    function highlight(i){
      items.forEach(function(li){ li.removeAttribute('data-active'); });
      cursor=i;
      if(items[cursor]){ items[cursor].setAttribute('data-active',''); items[cursor].scrollIntoView({block:'nearest'}); }
    }
    function visibleItems(){ return items.filter(function(li){ return !li.hidden; }); }

    input.addEventListener('input',function(){ filter(input.value); });
    input.addEventListener('focus',function(){ filter(input.value); });
    input.addEventListener('keydown',function(e){
      var vis=visibleItems();
      if(e.key==='ArrowDown'){ e.preventDefault(); highlight(Math.min(cursor+1,vis.length-1)); }
      else if(e.key==='ArrowUp'){ e.preventDefault(); highlight(Math.max(cursor-1,0)); }
      else if(e.key==='Enter'){ e.preventDefault(); if(cursor>=0&&vis[cursor]) select(vis[cursor]); }
      else if(e.key==='Escape'){ close(); input.blur(); }
    });
    list.addEventListener('mousedown',function(e){
      var li=e.target.closest('[data-value]');
      if(li){ e.preventDefault(); select(li); }
    });
    document.addEventListener('click',function(e){
      if(!root.contains(e.target)) close();
    });
  });
})();`
