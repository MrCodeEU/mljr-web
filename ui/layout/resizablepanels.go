package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ResizablePanelsProps struct {
	ID           string // unique per page (default "rp")
	Direction    string // "horizontal" (default) | "vertical"
	InitialSplit int    // initial first-panel percent (default 50)
	Min          int    // minimum percent for each panel (default 15)
}

// ResizablePanels renders two panels with a drag handle between them.
// Pass exactly two children: first panel content and second panel content.
func ResizablePanels(p ResizablePanelsProps, first, second g.Node) g.Node {
	if p.ID == "" {
		p.ID = "rp"
	}
	if p.Direction == "" {
		p.Direction = "horizontal"
	}
	if p.InitialSplit == 0 {
		p.InitialSplit = 50
	}
	if p.Min == 0 {
		p.Min = 15
	}

	isH := p.Direction == "horizontal"
	flexDir := "row"
	if !isH {
		flexDir = "column"
	}
	firstSize := fmt.Sprintf("%d%%", p.InitialSplit)
	secondSize := fmt.Sprintf("%d%%", 100-p.InitialSplit)
	cursor := "col-resize"
	if !isH {
		cursor = "row-resize"
	}

	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "resizable-panels"),
		g.Attr("data-direction", p.Direction),
		h.Style("display:flex;flex-direction:"+flexDir+";width:100%;height:100%;overflow:hidden"),
		h.Div(
			g.Attr("data-slot", "first"),
			h.Style(func() string {
				if isH {
					return "width:" + firstSize + ";overflow:auto;flex-shrink:0"
				}
				return "height:" + firstSize + ";overflow:auto;flex-shrink:0"
			}()),
			first,
		),
		h.Div(
			g.Attr("data-slot", "handle"),
			h.Style("flex-shrink:0;cursor:"+cursor+";background:var(--line);"+
				func() string {
					if isH {
						return "width:4px"
					}
					return "height:4px;width:100%"
				}()+";transition:background .1s"),
		),
		h.Div(
			g.Attr("data-slot", "second"),
			h.Style(func() string {
				if isH {
					return "width:" + secondSize + ";overflow:auto;flex:1;min-width:0"
				}
				return "height:" + secondSize + ";overflow:auto;flex:1;min-height:0"
			}()),
			second,
		),
		h.Script(g.Raw(resizableScript(p))),
	)
}

func resizableScript(p ResizablePanelsProps) string {
	isH := p.Direction == "horizontal"
	return fmt.Sprintf(`(function(){
  var root=document.getElementById('%s');
  if(!root) return;
  var handle=root.querySelector('[data-slot="handle"]');
  var first=root.querySelector('[data-slot="first"]');
  var second=root.querySelector('[data-slot="second"]');
  var isH=%v; var min=%d;
  var dragging=false;

  handle.addEventListener('mousedown',start);
  handle.addEventListener('touchstart',function(e){ start(e.touches[0]); });

  function start(e){
    e.preventDefault();
    dragging=true;
    document.addEventListener('mousemove',move);
    document.addEventListener('touchmove',function(e){ move(e.touches[0]); },{passive:false});
    document.addEventListener('mouseup',stop,{once:true});
    document.addEventListener('touchend',stop,{once:true});
    handle.style.background='var(--accent)';
  }
  function stop(){
    dragging=false;
    document.removeEventListener('mousemove',move);
    handle.style.background='';
  }
  function move(e){
    if(!dragging) return;
    var r=root.getBoundingClientRect();
    var pct;
    if(isH){ pct=Math.max(min,Math.min(100-min,(e.clientX-r.left)/r.width*100)); }
    else    { pct=Math.max(min,Math.min(100-min,(e.clientY-r.top)/r.height*100)); }
    if(isH){ first.style.width=pct+'%%'; second.style.width=(100-pct)+'%%'; }
    else   { first.style.height=pct+'%%'; second.style.height=(100-pct)+'%%'; }
  }
})();`, p.ID, isH, p.Min)
}
