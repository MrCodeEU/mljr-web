package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ReadProgressProps struct {
	// Height of the bar (default "3px").
	Height string
	// Color: CSS color (default "var(--accent)").
	Color string
	// ZIndex (default 100).
	ZIndex int
	// Top offsets the fixed bar from the viewport top. Use this when the page
	// has a sticky header.
	Top string
	// Target is a CSS selector for the scrollable element (default reads window scroll).
	// Leave empty to track window scroll.
	Target string
}

// ReadProgress renders a thin progress bar fixed at the top of the page
// that fills as the user scrolls.
func ReadProgress(p ReadProgressProps) g.Node {
	if p.Height == "" {
		p.Height = "3px"
	}
	if p.Color == "" {
		p.Color = "var(--accent)"
	}
	if p.ZIndex == 0 {
		p.ZIndex = 100
	}
	if p.Top == "" {
		p.Top = "0"
	}

	var targetExpr string
	if p.Target != "" {
		targetExpr = fmt.Sprintf(`var sc=document.querySelector('%s')||window;`, p.Target)
	} else {
		targetExpr = "var sc=window;"
	}

	var scrollExpr string
	if p.Target != "" {
		scrollExpr = `var el=document.querySelector('` + p.Target + `');if(!el)return;var pct=el.scrollTop/(el.scrollHeight-el.clientHeight)||0;`
	} else {
		scrollExpr = `var pct=(window.scrollY||document.documentElement.scrollTop)/((document.documentElement.scrollHeight-document.documentElement.clientHeight)||1);`
	}

	script := fmt.Sprintf(`(function(){
  var bar=document.getElementById('rp-bar');
  if(!bar) return;
  %s
  function update(){ %s bar.style.width=(Math.min(1,pct)*100)+'%%'; }
  sc.addEventListener('scroll',update,{passive:true});
  update();
})();`, targetExpr, scrollExpr)

	return g.Group{
		h.Div(
			h.ID("rp-bar"),
			g.Attr("data-component", "read-progress"),
			h.Style(fmt.Sprintf(
				"position:fixed;top:%s;left:0;height:%s;background:%s;z-index:%d;width:0;transition:width 0.1s linear",
				p.Top, p.Height, p.Color, p.ZIndex,
			)),
			g.Attr("role", "progressbar"),
			g.Attr("aria-label", "Reading progress"),
			g.Attr("aria-valuemin", "0"),
			g.Attr("aria-valuemax", "100"),
		),
		h.Script(g.Raw(script)),
	}
}
