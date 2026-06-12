package form

import (
	"fmt"
	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type FileDropZoneProps struct {
	Name      string
	Accept    string // e.g. "image/*,.pdf"
	Multiple  bool
	MaxSizeMB int // display hint, not enforced client-side
	Label     string
}

// FileDropZone renders a drag-and-drop file upload area backed by a hidden file input.
func FileDropZone(p FileDropZoneProps, children ...g.Node) g.Node {
	if p.Label == "" {
		p.Label = "Drop files here or click to browse"
	}
	hint := ""
	if p.Accept != "" {
		hint += p.Accept
	}
	if p.MaxSizeMB > 0 {
		if hint != "" {
			hint += " · "
		}
		hint += fmt.Sprintf("max %d MB", p.MaxSizeMB)
	}

	return h.Div(
		g.Attr("data-component", "file-drop-zone"),
		h.Label(
			g.Attr("data-slot", "area"),
			h.Input(
				h.Type("file"),
				h.Name(p.Name),
				g.If(p.Accept != "", g.Attr("accept", p.Accept)),
				g.If(p.Multiple, g.Attr("multiple", "")),
				h.Style("position:absolute;inset:0;opacity:0;cursor:pointer;width:100%;height:100%"),
			),
			icon.Icon("lucide:upload", icon.Props{Size: "2rem"}),
			h.Span(g.Attr("data-slot", "label"), g.Text(p.Label)),
			g.If(hint != "", h.Span(g.Attr("data-slot", "hint"), g.Text(hint))),
			h.Div(
				g.Attr("data-slot", "file-list"),
				g.Attr("style", "display:none"),
			),
			g.Group(children),
		),
		h.Script(g.Raw(dropZoneScript)),
	)
}

const dropZoneScript = `(function(){
  document.querySelectorAll('[data-component="file-drop-zone"]').forEach(function(root){
    var area=root.querySelector('[data-slot="area"]');
    var input=area.querySelector('input[type="file"]');
    var list=root.querySelector('[data-slot="file-list"]');
    var label=root.querySelector('[data-slot="label"]');

    function showFiles(files){
      if(!files||!files.length) return;
      list.style.display='flex';
      list.innerHTML='';
      Array.from(files).forEach(function(f){
        var item=document.createElement('span');
        item.style.cssText='font-size:var(--t-xs);padding:2px 8px;background:var(--surface-2);border-radius:var(--radius);border:1px solid var(--line)';
        item.textContent=f.name;
        list.appendChild(item);
      });
      label.textContent=files.length===1?files[0].name:files.length+' files selected';
    }

    input.addEventListener('change',function(){ showFiles(input.files); });

    area.addEventListener('dragover',function(e){
      e.preventDefault(); root.setAttribute('data-drag','');
    });
    area.addEventListener('dragleave',function(){ root.removeAttribute('data-drag'); });
    area.addEventListener('drop',function(e){
      e.preventDefault(); root.removeAttribute('data-drag');
      showFiles(e.dataTransfer.files);
    });
  });
})();`
