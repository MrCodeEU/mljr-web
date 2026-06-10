//go:build showcase

package datastar

import (
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-transition", Name: "Page Transitions", Category: "animation",
		Summary: "Simulate route transitions with Motion.timeline — outgoing content exits, incoming content enters.",
		Code: `// Out → In page transition pattern
async function navigate(href) {
    // 1. Animate current page out
    await Motion.animate('#page-content',
        { opacity:[1,0], y:[0,-20] },
        { duration:0.2 }
    ).finished

    // 2. Fetch and swap content (or use Datastar @get)
    await loadNewContent(href)

    // 3. Animate new page in
    Motion.animate('#page-content',
        { opacity:[0,1], y:[20,0] },
        { duration:0.25, easing:'ease-out' }
    )
}

// With Datastar:
data-on:click="
  await Motion.animate('#view',{opacity:0,y:-16},{duration:0.15}).finished;
  @get('/view/'+$tab);
"`,
		Render: func(p map[string]string) g.Node {
			pages := []struct{ id, label, color, content string }{
				{"home", "Home", "var(--primary)", "Welcome to the home view. Motion fades this content out before loading the next page."},
				{"about", "About", "var(--accent)", "The about page. Notice the smooth directional slide transition between views."},
				{"work", "Work", "var(--success)", "Portfolio projects live here. Each transition feels intentional, not jarring."},
			}
			tabs := make([]g.Node, len(pages))
			for i, pg := range pages {
				tabs[i] = primitive.Button(primitive.ButtonProps{
					Variant: token.Outline,
					Attrs:   []g.Node{g.Attr("data-tab-btn", pg.id)},
				}, g.Text(pg.label))
			}
			views := make([]g.Node, len(pages))
			for i, pg := range pages {
				views[i] = h.Div(
					g.Attr("data-tab-view", pg.id),
					h.Style("display:none;padding:var(--sp-6);border-left:4px solid "+pg.color),
					h.Strong(h.Style("font-size:var(--t-lg);font-weight:900"), g.Text(pg.label)),
					h.P(h.Style("color:var(--muted);margin:var(--sp-2) 0 0"), g.Text(pg.content)),
				)
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				h.Div(h.Style("display:flex;gap:var(--sp-2)"), g.Group(tabs)),
				h.Div(
					h.ID("trans-container"),
					h.Style("border:var(--bw-1) solid var(--line);border-radius:var(--radius);min-height:140px;overflow:hidden"),
					g.Group(views),
				),
				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var container=document.getElementById('trans-container');
  var views=Array.from(container.querySelectorAll('[data-tab-view]'));
  var buttons=document.querySelectorAll('[data-tab-btn]');
  var current=views[0];
  views[0].style.display='block';
  buttons[0].setAttribute('data-variant','primary');

  buttons.forEach(function(btn){
    btn.addEventListener('click',function(){
      var target=container.querySelector('[data-tab-view="'+btn.dataset.tabBtn+'"]');
      if(target===current) return;
      var dir=Array.from(views).indexOf(target)>Array.from(views).indexOf(current)?1:-1;
      Motion.animate(current,{opacity:[1,0],x:[0,dir*-30]},{duration:0.18}).finished.then(function(){
        current.style.display='none';
        target.style.display='block';
        Motion.animate(target,{opacity:[0,1],x:[dir*30,0]},{duration:0.22,easing:'ease-out'});
        current=target;
      });
      buttons.forEach(function(b){ b.setAttribute('data-variant','outline'); });
      btn.setAttribute('data-variant','primary');
    });
  });
})()`)),
			)
		},
	})
}
