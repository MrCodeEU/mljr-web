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
		Slug: "ds-anim-scroll", Name: "Scroll Progress", Category: "animation",
		PreviewHeight: "500px",
		Summary: "Motion.scroll() links animation progress to page scroll position. Zero event listeners, zero JS overhead.",
		Code: `// Scroll-linked animation — progress 0→1 as element scrolls through viewport
Motion.scroll(
    Motion.animate('#progress-bar', { scaleX: [0, 1] }, { duration: 1 }),
    { source: document.documentElement }  // whole page scroll
)

// Scroll-triggered with inView
Motion.inView('.section', ({ target }) => {
    return Motion.animate(target,
        { opacity: [0,1], y: [30,0] },
        { duration: 0.5 }
    )
})

// Scroll speed — reads scrollY to drive value
window.addEventListener('scroll', () => {
    const pct = scrollY / (document.body.scrollHeight - innerHeight)
    Motion.animate('#indicator', { scaleX: pct }, { duration: 0 })
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),

				// Sticky scroll progress bar
				h.Div(
					h.Style("position:sticky;top:0;z-index:10;background:var(--bg);padding:var(--sp-3);border-bottom:var(--bw-1) solid var(--line)"),
					h.P(h.Style("font-size:var(--t-xs);opacity:.5;font-weight:700;text-transform:uppercase;letter-spacing:.06em;margin:0 0 var(--sp-2)"), g.Text("Scroll progress bar (sticky)")),
					h.Div(
						h.Style("height:6px;background:var(--line);border-radius:3px;overflow:hidden"),
						h.Div(h.ID("scroll-prog"), h.Style("height:100%;background:var(--accent);width:0%;transform-origin:left;border-radius:3px")),
					),
				),

				// Scrollable content with inView sections
				h.Div(
					h.ID("scroll-content"),
					h.Style("height:380px;overflow-y:auto;display:flex;flex-direction:column;gap:var(--sp-4);padding:var(--sp-4);border:var(--bw-1) solid var(--line);border-radius:var(--radius)"),
					func() g.Node {
						sections := []struct{ color, text string }{
							{"var(--primary)", "Scroll-driven animations use zero event listeners."},
							{"var(--accent)", "Motion.scroll() links progress directly to scroll position."},
							{"var(--success)", "Each section fades and rises as it enters the viewport."},
							{"var(--warning)", "Performance: animations run on the compositor thread."},
							{"var(--primary)", "Works with any scrollable container — not just the page."},
						}
						nodes := make([]g.Node, len(sections))
						for i, s := range sections {
							nodes[i] = h.Div(
								h.Class("scroll-section"),
								h.Style("padding:var(--sp-5);border-left:4px solid "+s.color+";background:var(--surface-2);border-radius:0 var(--radius) var(--radius) 0;opacity:0"),
								h.P(h.Style("margin:0;font-size:var(--t-sm)"), g.Text(s.text)),
							)
						}
						return g.Group(append(
							[]g.Node{h.Div(h.Style("height:40px;display:flex;align-items:center;justify-content:center;opacity:.4;font-size:var(--t-sm)"), g.Text("↓ scroll down"))},
							append(nodes, h.Div(h.Style("height:40px")))...,
						))
					}(),
				),

				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var cont=document.getElementById('scroll-content');
  var bar=document.getElementById('scroll-prog');

  // Scroll progress bar
  cont.addEventListener('scroll',function(){
    var pct=cont.scrollTop/(cont.scrollHeight-cont.clientHeight);
    Motion.animate(bar,{width:(pct*100)+'%'},{duration:0});
  });

  // inView reveals inside scroll container
  document.querySelectorAll('.scroll-section').forEach(function(el,i){
    var anims=[
      function(t){ Motion.animate(t,{opacity:[0,1],x:[-30,0]},{duration:0.5}); },
      function(t){ Motion.animate(t,{opacity:[0,1],y:[20,0]},{duration:0.5}); },
      function(t){ Motion.animate(t,{opacity:[0,1],scale:[0.9,1]},{duration:0.5,easing:'ease-out'}); },
      function(t){ Motion.animate(t,{opacity:[0,1],filter:['blur(8px)','blur(0)']},{duration:0.5}); },
      function(t){ Motion.animate(t,{opacity:[0,1],x:[30,0]},{duration:0.5}); },
    ];
    Motion.inView(el,function(){ anims[i%anims.length](el); },{root:cont,margin:'0px 0px -20px 0px'});
  });
})()`)),
			)
		},
	})
}

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-morph", Name: "Loading Morphs", Category: "animation",
		Summary: "Motion.timeline() morphs skeleton placeholders into real content. Staged entrance with stagger.",
		Code: `// Skeleton → content morph via Motion.timeline
// 1. Hide real content initially
// 2. Show skeleton placeholders
// 3. On data ready: fade skeletons out, stagger content in

Motion.timeline([
    // Fade out skeletons
    ['.skeleton', { opacity:[1,0], scale:[1,0.95] }, { duration:0.2 }],
    // Stagger in real items
    ['.item', { opacity:[0,1], y:[12,0] },
      { duration:0.3, delay: Motion.stagger(0.06), at:'+0.05' }],
])`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

				// Card morph demo
				primitive.Card(primitive.CardProps{},
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0 0 var(--sp-3)"),
						g.Text("Skeleton → content morph")),
					h.Div(h.ID("morph-container"), h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
						// Skeletons
						func() g.Node {
							skels := make([]g.Node, 4)
							for i := range skels {
								skels[i] = h.Div(h.Class("morph-skel"),
									h.Style("height:52px;background:var(--surface-2);border-radius:var(--radius);animation:mljr-pulse 1.5s ease infinite"))
							}
							return g.Group(skels)
						}(),
						// Real items (hidden)
						func() g.Node {
							items := []struct{ name, role, color string }{
								{"Alex Chen", "Senior Engineer", "var(--primary)"},
								{"Jordan Lee", "Product Designer", "var(--accent)"},
								{"Sam Park", "Dev Ops Lead", "var(--success)"},
								{"Morgan Wu", "Data Scientist", "var(--warning)"},
							}
							nodes := make([]g.Node, len(items))
							for i, item := range items {
								nodes[i] = h.Div(
									h.Class("morph-item"),
									h.Style("display:flex;align-items:center;gap:var(--sp-3);opacity:0;padding:var(--sp-2);border-radius:var(--radius);background:var(--surface-2)"),
									h.Div(h.Style("width:36px;height:36px;border-radius:50%;background:"+item.color+";flex-shrink:0")),
									h.Div(
										h.Strong(h.Style("display:block;font-size:var(--t-sm)"), g.Text(item.name)),
										h.Span(h.Style("font-size:var(--t-xs);color:var(--muted)"), g.Text(item.role)),
									),
								)
							}
							return g.Group(nodes)
						}(),
					),
					h.Div(h.Style("display:flex;gap:var(--sp-3);margin-top:var(--sp-4)"),
						primitive.Button(primitive.ButtonProps{Variant: token.Primary, Size: token.SizeSM},
							h.ID("morph-load"),
							g.Text("Load data"),
						),
						primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
							h.ID("morph-reset"),
							g.Text("Reset"),
						),
					),
				),

				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var skels=document.querySelectorAll('.morph-skel');
  var items=document.querySelectorAll('.morph-item');
  var loaded=false;

  document.getElementById('morph-load').addEventListener('click',function(){
    if(loaded) return;
    loaded=true;
    Motion.timeline([
      [skels,{opacity:[1,0],scale:[1,0.95]},{duration:0.25}],
      [items,{opacity:[0,1],y:[12,0]},{duration:0.35,delay:Motion.stagger(0.08),at:'+0.05'}]
    ]).finished.then(function(){
      skels.forEach(function(s){ s.style.display='none'; });
    });
  });

  document.getElementById('morph-reset').addEventListener('click',function(){
    loaded=false;
    skels.forEach(function(s){ s.style.display=''; Motion.animate(s,{opacity:1,scale:1},{duration:0}); });
    items.forEach(function(el){ Motion.animate(el,{opacity:0,y:0},{duration:0}); });
  });
})()`)),
			)
		},
	})
}
