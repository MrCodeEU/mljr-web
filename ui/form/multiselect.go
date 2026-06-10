package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MultiSelectOption struct {
	Value string
	Label string
}

type MultiSelectProps struct {
	Name     string
	Options  []MultiSelectOption
	Default  []string // pre-selected values
	Max      int      // max selections (0 = unlimited)
	Placeholder string
}

// MultiSelect renders a chip-based multi-value select.
// Click options to toggle selection; selected values shown as chips.
// Hidden inputs named p.Name[] carry all selected values for form submission.
func MultiSelect(p MultiSelectProps) g.Node {
	if p.Placeholder == "" {
		p.Placeholder = "Select options…"
	}

	defaultSet := make(map[string]bool, len(p.Default))
	for _, v := range p.Default {
		defaultSet[v] = true
	}

	optNodes := make([]g.Node, len(p.Options))
	for i, opt := range p.Options {
		selected := defaultSet[opt.Value]
		optNodes[i] = h.Div(
			h.Class("ms-option"),
			g.Attr("data-value", opt.Value),
			g.Attr("data-selected", func() string {
				if selected {
					return "true"
				}
				return "false"
			}()),
			g.Text(opt.Label),
		)
	}

	hiddenInputs := make([]g.Node, 0, len(p.Default))
	for _, v := range p.Default {
		hiddenInputs = append(hiddenInputs, h.Input(
			h.Type("hidden"),
			h.Class("ms-hidden"),
			h.Name(p.Name+"[]"),
			h.Value(v),
		))
	}

	maxAttr := ""
	if p.Max > 0 {
		maxAttr = fmt.Sprintf("%d", p.Max)
	}

	return h.Div(
		g.Attr("data-component", "multi-select"),
		g.If(maxAttr != "", g.Attr("data-max", maxAttr)),

		h.Div(
			g.Attr("data-slot", "trigger"),
			g.Attr("tabindex", "0"),
			h.Div(g.Attr("data-slot", "chips")),
			h.Span(g.Attr("data-slot", "placeholder"), h.Style("color:var(--muted)"), g.Text(p.Placeholder)),
		),

		h.Div(
			g.Attr("data-slot", "dropdown"),
			h.Style("display:none"),
			g.Group(optNodes),
		),

		g.Group(hiddenInputs),
		h.Script(g.Raw(multiSelectScript)),
	)
}

const multiSelectScript = `(function(){
  document.querySelectorAll('[data-component="multi-select"]').forEach(function(root){
    var trigger=root.querySelector('[data-slot="trigger"]');
    var chipsSlot=root.querySelector('[data-slot="chips"]');
    var placeholder=root.querySelector('[data-slot="placeholder"]');
    var dropdown=root.querySelector('[data-slot="dropdown"]');
    var options=Array.from(root.querySelectorAll('.ms-option'));
    var max=parseInt(root.dataset.max)||0;
    var name=root.querySelector('.ms-hidden');
    var baseName=name?name.name.replace('[]',''):'';

    function selected(){ return options.filter(function(o){ return o.dataset.selected==='true'; }); }
    function sync(){
      // Update chips
      chipsSlot.innerHTML='';
      var sel=selected();
      sel.forEach(function(o){
        var chip=document.createElement('span');
        chip.style.cssText='display:inline-flex;align-items:center;gap:4px;padding:2px 8px;background:var(--accent);color:var(--accent-ink);border-radius:calc(var(--radius)*.6);font-size:var(--t-xs);font-weight:700';
        chip.textContent=o.textContent+' ';
        var x=document.createElement('button');
        x.type='button';x.textContent='×';x.style.cssText='background:none;border:none;cursor:pointer;color:inherit;padding:0;font-size:1em;line-height:1';
        x.onclick=function(e){ e.stopPropagation(); o.dataset.selected='false'; sync(); };
        chip.appendChild(x);
        chipsSlot.appendChild(chip);
      });
      placeholder.style.display=sel.length?'none':'';
      // Sync hidden inputs
      root.querySelectorAll('.ms-hidden').forEach(function(i){ i.remove(); });
      sel.forEach(function(o){
        var inp=document.createElement('input');
        inp.type='hidden'; inp.className='ms-hidden';
        inp.name=baseName+'[]'; inp.value=o.dataset.value;
        root.appendChild(inp);
      });
    }

    // Option toggle
    options.forEach(function(opt){
      opt.addEventListener('click',function(){
        var isSel=opt.dataset.selected==='true';
        if(!isSel&&max>0&&selected().length>=max) return;
        opt.dataset.selected=isSel?'false':'true';
        sync();
      });
    });

    // Open/close dropdown
    trigger.addEventListener('click',function(e){
      e.stopPropagation();
      dropdown.style.display=dropdown.style.display==='none'?'flex':'none';
    });
    document.addEventListener('click',function(){ dropdown.style.display='none'; });

    sync();
  });
})();`
