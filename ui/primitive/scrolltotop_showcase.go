//go:build showcase

package primitive

import (
	"fmt"
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "scroll-to-top", Name: "Scroll To Top", Category: "primitive",
		Summary: "Fixed button that appears after scrolling past a threshold. Smooth-scrolls to top on click. Zero dependencies.",
		Code: `primitive.ScrollToTop(primitive.ScrollToTopProps{
    Threshold: 300,
    Position:  "bottom-right",
})`,
		Render: func(p map[string]string) g.Node {
			// Simulate a tall scrollable area
			rows := make([]g.Node, 20)
			for i := range rows {
				rows[i] = h.Div(
					h.Style("padding:var(--sp-4);border-bottom:var(--bw-1) solid var(--line);font-size:var(--t-sm)"),
					g.Text(fmt.Sprintf("Paragraph %d — scroll down to see the button appear at the bottom right.", i+1)),
				)
			}
			return h.Div(
				h.Style("position:relative;height:500px;overflow-y:auto;border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
				h.Div(g.Group(rows)),
				// Inline scroll-to-top scoped to container
				h.Button(
					g.Attr("data-component", "scroll-to-top"),
					h.ID("stt-demo"),
					h.Type("button"),
					h.Style("position:sticky;bottom:var(--sp-4);left:100%;margin-right:var(--sp-3);opacity:0;pointer-events:none;transition:opacity 0.2s"),
					g.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="18 15 12 9 6 15"/></svg>`),
				),
				h.Script(g.Raw(`(function(){
  var c=document.querySelector('[style*="height:500px"]');
  var b=document.getElementById('stt-demo');
  if(!c||!b) return;
  c.addEventListener('scroll',function(){
    b.style.opacity=c.scrollTop>100?'1':'0';
    b.style.pointerEvents=c.scrollTop>100?'auto':'none';
  },{passive:true});
  b.addEventListener('click',function(){ c.scrollTo({top:0,behavior:'smooth'}); });
})();`)),
				h.P(h.Style("padding:var(--sp-3);font-size:var(--t-xs);color:var(--muted)"),
					g.Text("In production: add primitive.ScrollToTop() to page body — button tracks window scroll."),
				),
			)
		},
	})
}
