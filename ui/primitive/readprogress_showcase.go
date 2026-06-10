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
		Slug: "read-progress", Name: "Read Progress", Category: "primitive",
		Summary: "Thin progress bar fixed at top of page that fills as the user scrolls. Tracks window scroll or a custom container.",
		Code: `// Add once to the page body:
primitive.ReadProgress(primitive.ReadProgressProps{
    Height: "3px",
    Color:  "var(--accent)",
})`,
		Render: func(p map[string]string) g.Node {
			paragraphs := make([]g.Node, 12)
			for i := range paragraphs {
				paragraphs[i] = h.P(
					h.Style("font-size:var(--t-sm);margin-bottom:var(--sp-4)"),
					g.Text(fmt.Sprintf("Paragraph %d — Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco.", i+1)),
				)
			}
			// Scoped to a container
			return h.Div(
				h.Style("position:relative;border:var(--bw-1) solid var(--line);border-radius:var(--radius);overflow:hidden"),
				// Progress bar scoped to container
				h.Div(h.ID("rp-demo"), h.Style("position:sticky;top:0;left:0;right:0;height:3px;background:var(--accent);width:0;z-index:10;transition:width 0.1s linear")),
				h.Div(
					h.ID("rp-container"),
					h.Style("height:400px;overflow-y:auto;padding:var(--sp-5)"),
					h.H3(h.Style("font-weight:800;margin-bottom:var(--sp-4)"), g.Text("Scroll this article…")),
					g.Group(paragraphs),
				),
				h.Script(g.Raw(`(function(){
  var bar=document.getElementById('rp-demo');
  var sc=document.getElementById('rp-container');
  if(!bar||!sc) return;
  sc.addEventListener('scroll',function(){
    var pct=sc.scrollTop/(sc.scrollHeight-sc.clientHeight)||0;
    bar.style.width=(Math.min(1,pct)*100)+'%';
  },{passive:true});
})();`)),
			)
		},
	})
}
