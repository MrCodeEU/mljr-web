package primitive

import (
	"fmt"

	"mljr-web/ui/icon"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ScrollToTopProps struct {
	// Threshold is pixels scrolled before button appears (default 300).
	Threshold int
	// Position: "bottom-right" (default) | "bottom-left"
	Position string
}

// ScrollToTop renders a fixed button that appears after scrolling past a threshold.
// Clicking smoothly scrolls back to the top.
func ScrollToTop(p ScrollToTopProps) g.Node {
	if p.Threshold == 0 {
		p.Threshold = 300
	}
	if p.Position == "" {
		p.Position = "bottom-right"
	}

	posStyle := "bottom:var(--sp-6);right:var(--sp-6)"
	if p.Position == "bottom-left" {
		posStyle = "bottom:var(--sp-6);left:var(--sp-6)"
	}

	script := fmt.Sprintf(`(function(){
  var btn=document.getElementById('stt-btn');
  if(!btn) return;
  window.addEventListener('scroll',function(){
    btn.style.opacity=window.scrollY>%d?'1':'0';
    btn.style.pointerEvents=window.scrollY>%d?'auto':'none';
  },{passive:true});
  btn.addEventListener('click',function(){ window.scrollTo({top:0,behavior:'smooth'}); });
})();`, p.Threshold, p.Threshold)

	return g.Group{
		h.Button(
			h.ID("stt-btn"),
			g.Attr("data-component", "scroll-to-top"),
			h.Type("button"),
			h.Style(posStyle+";position:fixed;z-index:90;opacity:0;pointer-events:none;transition:opacity 0.2s"),
			h.Aria("label", "Scroll to top"),
			icon.Icon("lucide:arrow-up"),
		),
		h.Script(g.Raw(script)),
	}
}
