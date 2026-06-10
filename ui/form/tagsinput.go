package form

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TagsInputProps struct {
	Name        string
	Placeholder string
	Default     []string // pre-filled tags
	Max         int      // max tags; 0 = unlimited
}

// TagsInput renders a multi-value tag input. Tags are added on Enter or comma,
// removed with ×. The hidden input named p.Name holds comma-separated values.
func TagsInput(p TagsInputProps) g.Node {
	if p.Placeholder == "" {
		p.Placeholder = "Add tag…"
	}
	defaultVal := strings.Join(p.Default, ",")
	maxAttr := ""
	if p.Max > 0 {
		maxAttr = fmt.Sprint(p.Max)
	}

	prebuiltTags := make([]g.Node, len(p.Default))
	for i, t := range p.Default {
		prebuiltTags[i] = tagChip(t)
	}

	return h.Div(
		g.Attr("data-component", "tags-input"),
		g.If(maxAttr != "", g.Attr("data-max", maxAttr)),
		h.Div(
			g.Attr("data-slot", "tags"),
			g.Group(prebuiltTags),
			h.Input(
				g.Attr("data-slot", "text-input"),
				h.Type("text"),
				h.Placeholder(p.Placeholder),
				g.Attr("autocomplete", "off"),
			),
		),
		h.Input(
			g.Attr("data-slot", "value"),
			h.Type("hidden"),
			h.Name(p.Name),
			h.Value(defaultVal),
		),
		h.Script(g.Raw(tagsInputScript)),
	)
}

func tagChip(label string) g.Node {
	return h.Span(
		g.Attr("data-slot", "tag"),
		g.Attr("data-tag", label),
		g.Text(label),
		h.Button(
			h.Type("button"),
			g.Attr("data-slot", "remove"),
			g.Attr("aria-label", "Remove "+label),
			g.Raw("×"),
		),
	)
}

const tagsInputScript = `(function(){
  document.querySelectorAll('[data-component="tags-input"]').forEach(function(root){
    var tagsSlot=root.querySelector('[data-slot="tags"]');
    var textInput=root.querySelector('[data-slot="text-input"]');
    var hidden=root.querySelector('[data-slot="value"]');
    var max=parseInt(root.dataset.max)||0;

    function tags(){ return Array.from(root.querySelectorAll('[data-slot="tag"]')).map(function(t){return t.dataset.tag;}); }
    function sync(){ hidden.value=tags().join(','); }

    function addTag(val){
      val=val.trim();
      if(!val) return;
      if(tags().includes(val)) return;
      if(max>0 && tags().length>=max) return;
      var chip=document.createElement('span');
      chip.dataset.slot='tag';
      chip.dataset.tag=val;
      chip.textContent=val+' ';
      var btn=document.createElement('button');
      btn.type='button'; btn.dataset.slot='remove'; btn.setAttribute('aria-label','Remove '+val);
      btn.textContent='×';
      btn.addEventListener('click',function(){ chip.remove(); sync(); });
      chip.appendChild(btn);
      tagsSlot.insertBefore(chip, textInput);
      sync();
    }

    // Existing remove buttons
    root.querySelectorAll('[data-slot="remove"]').forEach(function(btn){
      btn.addEventListener('click',function(){ btn.closest('[data-slot="tag"]').remove(); sync(); });
    });

    textInput.addEventListener('keydown',function(e){
      if(e.key==='Enter'||e.key===','){
        e.preventDefault();
        addTag(textInput.value.replace(',',''));
        textInput.value='';
      } else if(e.key==='Backspace'&&textInput.value===''){
        var last=root.querySelectorAll('[data-slot="tag"]');
        if(last.length) { last[last.length-1].remove(); sync(); }
      }
    });
    textInput.addEventListener('blur',function(){
      if(textInput.value.trim()){ addTag(textInput.value); textInput.value=''; }
    });
  });
})();`
