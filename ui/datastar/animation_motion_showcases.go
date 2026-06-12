//go:build showcase

package datastar

import (
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// ---------------------------------------------------------------------------
// ds-anim-spring — spring physics with configurable stiffness / damping
// ---------------------------------------------------------------------------

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-spring", Name: "Spring Physics", Category: "animation",
		Summary: "Motion.spring() produces physics-based easing. Tune stiffness, damping, and mass for wildly different feels.",
		Code: `// Motion.spring() returns an easing generator — pass to easing option
// No duration needed; spring settles at its natural time.
Motion.animate(el, { y: 150 }, {
    easing: Motion.spring({ stiffness: 300, damping: 8, mass: 1 }),
})

// Configs by feel:
// Stiff  → stiffness:800, damping:25   (fast, minimal overshoot)
// Bouncy → stiffness:200, damping:5    (oscillates several times)
// Slow   → stiffness:100, damping:15   (heavy, lazy)
// Snap   → stiffness:600, damping:40   (crisp, one bounce)`,
		Render: func(p map[string]string) g.Node {
			type ball struct {
				id  string
				lbl string
				col string
			}
			balls := []ball{
				{"sp-stiff", "Stiff", "var(--primary)"},
				{"sp-bouncy", "Bouncy", "var(--accent)"},
				{"sp-slow", "Slow", "var(--success)"},
				{"sp-snap", "Snap", "var(--warning)"},
			}
			ballNodes := make([]g.Node, len(balls))
			for i, b := range balls {
				ballNodes[i] = h.Div(
					h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-2)"),
					h.Div(
						h.ID(b.id),
						h.Style("width:52px;height:52px;border-radius:50%;background:"+b.col+";border:var(--bw-2) solid var(--line)"),
					),
					h.Span(h.Style("font-size:var(--t-xs);font-weight:700;text-transform:uppercase;letter-spacing:.05em"), g.Text(b.lbl)),
				)
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),
				h.Div(
					h.Style("display:flex;justify-content:space-around;align-items:flex-start;flex-wrap:wrap;gap:var(--sp-4);padding:var(--sp-5) var(--sp-3) 0"),
					g.Group(ballNodes),
				),
				h.Div(
					h.Style("display:flex;gap:var(--sp-3);justify-content:center;flex-wrap:wrap"),
					h.Button(
						h.ID("sp-drop"),
						g.Attr("data-component", "button"),
						g.Attr("data-variant", string(token.Primary)),
						g.Attr("data-size", string(token.SizeMD)),
						h.Type("button"),
						g.Text("Drop"),
					),
					h.Button(
						h.ID("sp-reset"),
						g.Attr("data-component", "button"),
						g.Attr("data-variant", string(token.Outline)),
						g.Attr("data-size", string(token.SizeMD)),
						h.Type("button"),
						g.Text("Reset"),
					),
				),
				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var configs={
    'sp-stiff': {stiffness:800,damping:25,mass:1},
    'sp-bouncy':{stiffness:200,damping:5, mass:1},
    'sp-slow':  {stiffness:100,damping:15,mass:2},
    'sp-snap':  {stiffness:600,damping:40,mass:0.5}
  };
  var anims=[];
  document.getElementById('sp-drop').addEventListener('click',function(){
    anims.forEach(function(a){try{a.stop();}catch(e){}});
    anims=[];
    Object.keys(configs).forEach(function(id){
      var el=document.getElementById(id);
      if(!el) return;
      // Spring to target: let it overshoot and settle naturally (no explicit duration)
      anims.push(Motion.animate(el,{y:150},{easing:Motion.spring(configs[id])}));
    });
  });
  document.getElementById('sp-reset').addEventListener('click',function(){
    anims.forEach(function(a){try{a.stop();}catch(e){}});
    anims=[];
    Object.keys(configs).forEach(function(id){
      var el=document.getElementById(id);
      if(el) Motion.animate(el,{y:0},{duration:0.3,easing:'ease-in'});
    });
  });
})()`)),
			)
		},
	})
}

// ---------------------------------------------------------------------------
// ds-anim-inview — viewport-triggered reveals via Motion.inView
// ---------------------------------------------------------------------------

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-inview", Name: "inView Reveals", Category: "animation",
		Summary: "Motion.inView() triggers animations when elements enter a scrollable viewport — zero scroll event listeners.",
		Code: `// Motion.inView(element, callback, options)
// callback fires when element enters the root viewport.
// Supports custom root for scrollable containers.

Motion.inView(el, ({ target }) => {
    Motion.animate(target, { opacity:[0,1], x:[-40,0] }, { duration:0.5 })
})

// With custom root (scrollable container):
Motion.inView(el, callback, {
    root: scrollContainer,
    margin: '-10% 0px',
})`,
		Render: func(p map[string]string) g.Node {
			type card struct {
				title string
				body  string
				color string
			}
			cards := []card{
				{"Slides in from left", "x:[-40,0], opacity:[0,1]", "var(--primary)"},
				{"Rises from below", "y:[30,0], opacity:[0,1]", "var(--accent)"},
				{"Zooms in", "scale:[0.8,1], opacity:[0,1]", "var(--success)"},
				{"Blurs in", "filter:blur(12px)→blur(0)", "var(--warning)"},
				{"Flips up", "rotateX(60deg)→0, opacity", "var(--primary)"},
			}
			cardNodes := make([]g.Node, len(cards))
			for i, c := range cards {
				cardNodes[i] = h.Div(
					g.Attr("data-slot", "iv-item"),
					h.Style("padding:var(--sp-4);border-left:4px solid "+c.color+";background:var(--surface-2);border-radius:0 var(--radius) var(--radius) 0;opacity:0"),
					h.P(h.Style("font-weight:700;margin:0 0 var(--sp-1)"), g.Text(c.title)),
					h.Code(h.Style("font-size:var(--t-xs);opacity:.6"), g.Text(c.body)),
				)
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
				h.P(
					h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0"),
					g.Text("↓ Scroll the container below"),
				),
				h.Div(
					h.ID("iv-scroll"),
					h.Style("height:280px;overflow-y:auto;display:flex;flex-direction:column;gap:var(--sp-5);padding:var(--sp-4);border:var(--bw-1) solid var(--line);border-radius:var(--radius);scroll-behavior:smooth"),
					h.Div(h.Style("height:60px;display:flex;align-items:center;justify-content:center;opacity:.4;font-size:var(--t-sm)"), g.Text("↓ keep scrolling")),
					g.Group(cardNodes),
					h.Div(h.Style("height:40px")),
				),
				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var root=document.getElementById('iv-scroll');
  var fns=[
    function(el){ Motion.animate(el,{opacity:[0,1],x:[-40,0]},{duration:0.5}); },
    function(el){ Motion.animate(el,{opacity:[0,1],y:[30,0]},{duration:0.5}); },
    function(el){ Motion.animate(el,{opacity:[0,1],scale:[0.8,1]},{duration:0.5,easing:'ease-out'}); },
    function(el){ Motion.animate(el,{opacity:[0,1],filter:['blur(12px)','blur(0px)']},{duration:0.6}); },
    function(el){ Motion.animate(el,{opacity:[0,1],transform:['perspective(400px) rotateX(60deg)','perspective(400px) rotateX(0deg)']},{duration:0.5,easing:'ease-out'}); }
  ];
  document.querySelectorAll('[data-slot="iv-item"]').forEach(function(el,i){
    Motion.inView(el,function(){fns[i%fns.length](el);},{root:root,margin:'0px 0px -20px 0px'});
  });
})()`)),
			)
		},
	})
}

// ---------------------------------------------------------------------------
// ds-anim-text — three distinct text reveal effects
// ---------------------------------------------------------------------------

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-text", Name: "Text Reveal", Category: "animation",
		Summary: "Three distinct character-level reveal effects: blur rise, scale scatter, and wave rotation.",
		Code: `// Split text into spans, then stagger with different keyframes per effect.

// Effect 1 — Blur rise
Motion.animate(chars, { opacity:[0,1], y:[24,0], filter:['blur(8px)','blur(0px)'] },
    { delay: Motion.stagger(0.03), duration:0.5, easing:[0.22,1,0.36,1] })

// Effect 2 — Scale scatter (random x offsets per char)
Motion.animate(chars[i], { opacity:[0,1], x:[rand,0], scale:[0,1] },
    { delay: i*0.04, duration:0.4 })

// Effect 3 — Wave rotation
Motion.animate(chars, { opacity:[0,1], rotateZ:[-25,0], y:[12,0] },
    { delay: Motion.stagger(0.025, {from:'center'}), duration:0.4 })`,
		Render: func(p map[string]string) g.Node {
			type effect struct {
				id     string
				phrase string
				label  string
				color  string
			}
			effects := []effect{
				{"tr-blur", "Motion makes the web alive.", "Blur rise — y + blur stagger", "var(--primary)"},
				{"tr-scatter", "No framework. No bundle.", "Scale scatter — random x offsets", "var(--accent)"},
				{"tr-wave", "Fast, native, beautiful.", "Wave rotation — from center", "var(--success)"},
			}
			rows := make([]g.Node, len(effects))
			for i, e := range effects {
				rows[i] = h.Div(
					h.Style("padding:var(--sp-4);border:var(--bw-1) solid var(--line);border-radius:var(--radius);border-top:3px solid "+e.color+";display:flex;flex-direction:column;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0;color:"+e.color), g.Text(e.label)),
					h.Div(
						h.ID(e.id),
						h.Style("font-size:var(--t-xl);font-weight:900;font-family:var(--font-display);letter-spacing:-.02em;min-height:2em;display:flex;align-items:center;flex-wrap:wrap"),
						g.Attr("data-phrase", e.phrase),
					),
					h.Button(
						g.Attr("data-slot", "tr-play"),
						g.Attr("data-component", "button"),
						g.Attr("data-variant", string(token.Outline)),
						g.Attr("data-size", string(token.SizeSM)),
						h.Type("button"),
						g.Attr("data-anim", e.id),
						g.Text("Reveal"),
					),
				)
			}
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-4)"),
				g.Group(rows),
				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;

  function buildChars(stage){
    var phrase=stage.dataset.phrase||'';
    stage.innerHTML='';
    return Array.from(phrase).map(function(ch){
      var s=document.createElement('span');
      s.textContent=ch===' '?' ':ch;
      s.style.display='inline-block';
      s.style.opacity='0';
      stage.appendChild(s);
      return s;
    });
  }

  var effects={
    'tr-blur': function(chars){
      Motion.animate(chars,
        {opacity:[0,1],y:[24,0],filter:['blur(8px)','blur(0px)']},
        {delay:Motion.stagger(0.03),duration:0.5,easing:[0.22,1,0.36,1]}
      );
    },
    'tr-scatter': function(chars){
      chars.forEach(function(c,i){
        var rx=(Math.random()-0.5)*60;
        Motion.animate(c,
          {opacity:[0,1],x:[rx,0],scale:[0,1]},
          {delay:i*0.04,duration:0.45,easing:'ease-out'}
        );
      });
    },
    'tr-wave': function(chars){
      Motion.animate(chars,
        {opacity:[0,1],rotateZ:[-25,0],y:[12,0]},
        {delay:Motion.stagger(0.025,{from:'center'}),duration:0.4,easing:'ease-out'}
      );
    }
  };

  document.querySelectorAll('[data-slot="tr-play"]').forEach(function(btn){
    var effectId=btn.dataset.anim;
    var stage=document.getElementById(effectId);
    var chars=buildChars(stage);
    btn.addEventListener('click',function(){
      chars=buildChars(stage);
      var fn=effects[effectId];
      if(fn) fn(chars);
    });
  });
})()`)),
			)
		},
	})
}

// ---------------------------------------------------------------------------
// ds-anim-gesture — hover and focus micro-interactions
// ---------------------------------------------------------------------------

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-anim-gesture", Name: "Gesture & Hover", Category: "animation",
		Summary: "Micro-interactions on hover, focus, and click. Motion handles enter/leave without CSS keyframes.",
		Code: `// Hover enter / leave
el.addEventListener('mouseenter', () =>
    Motion.animate(el, { scale:1.04, y:-2 }, { duration:0.2, easing:'ease-out' }))
el.addEventListener('mouseleave', () =>
    Motion.animate(el, { scale:1, y:0 }, { duration:0.3, easing:'ease-out' }))

// Magnetic pull — track mouse relative to element center
el.addEventListener('mousemove', e => {
    var rect = el.getBoundingClientRect()
    var dx = (e.clientX - rect.left - rect.width/2) * 0.3
    var dy = (e.clientY - rect.top  - rect.height/2) * 0.3
    Motion.animate(el, { x:dx, y:dy }, { duration:0.3, easing:'ease-out' })
})

// Click ripple
el.addEventListener('click', e => {
    var dot = document.createElement('span')
    dot.style.cssText = 'position:absolute;width:8px;height:8px;border-radius:50%;background:var(--primary);pointer-events:none'
    dot.style.left = (e.offsetX - 4) + 'px'
    dot.style.top  = (e.offsetY - 4) + 'px'
    el.appendChild(dot)
    Motion.animate(dot, { scale:[0,10], opacity:[0.6,0] }, { duration:0.5 })
        .then(() => dot.remove())
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-6)"),

				// Hover lift cards
				h.Div(
					h.Style("display:grid;grid-template-columns:repeat(3,1fr);gap:var(--sp-3)"),
					func() g.Node {
						colors := []string{"#00ADD8", "var(--accent)", "var(--success)"}
						labels := []string{"Hover me", "No really", "Try it!"}
						nodes := make([]g.Node, 3)
						for i := range nodes {
							nodes[i] = h.Div(
								g.Attr("data-slot", "gest-card"),
								h.Style("padding:var(--sp-4);border:var(--bw-2) solid var(--line);border-radius:var(--radius);background:var(--surface-2);cursor:pointer;text-align:center;font-weight:700;border-top:4px solid "+colors[i]),
								g.Text(labels[i]),
							)
						}
						return g.Group(nodes)
					}(),
				),

				// Magnetic button — use a div wrapper so ID lands on the element we animate
				h.Div(
					h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0"), g.Text("Magnetic button")),
					h.Button(
						h.ID("magnet-btn"),
						g.Attr("data-component", "button"),
						g.Attr("data-variant", string(token.Primary)),
						g.Attr("data-size", string(token.SizeMD)),
						h.Type("button"),
						g.Text("Pull me"),
					),
				),

				// Click ripple — overflow:hidden clips ripple to button bounds
				h.Div(
					h.Style("display:flex;flex-direction:column;align-items:center;gap:var(--sp-3)"),
					h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin:0"), g.Text("Click ripple")),
					h.Div(
						h.ID("ripple-btn"),
						h.Style("padding:var(--sp-3) var(--sp-6);border:var(--bw-2) solid var(--line);border-radius:var(--radius);cursor:pointer;position:relative;overflow:hidden;background:var(--surface-2);font-weight:700;text-align:center;user-select:none"),
						g.Text("Click anywhere"),
					),
				),

				h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;

  // Hover lift
  document.querySelectorAll('[data-slot="gest-card"]').forEach(function(el){
    el.addEventListener('mouseenter',function(){
      Motion.animate(el,{y:-6,scale:1.03},{duration:0.25,easing:'ease-out'});
    });
    el.addEventListener('mouseleave',function(){
      Motion.animate(el,{y:0,scale:1},{duration:0.3,easing:'ease-out'});
    });
    el.addEventListener('mousedown',function(){
      Motion.animate(el,{scale:0.97},{duration:0.1});
    });
    el.addEventListener('mouseup',function(){
      Motion.animate(el,{scale:1.03},{duration:0.15});
    });
  });

  // Magnetic button — ID is on the <button> element directly
  var mag=document.getElementById('magnet-btn');
  if(mag){
    mag.addEventListener('mousemove',function(e){
      var r=mag.getBoundingClientRect();
      var dx=(e.clientX-r.left-r.width/2)*0.3;
      var dy=(e.clientY-r.top -r.height/2)*0.3;
      Motion.animate(mag,{x:dx,y:dy},{duration:0.3,easing:'ease-out'});
    });
    mag.addEventListener('mouseleave',function(){
      Motion.animate(mag,{x:0,y:0},{duration:0.4,easing:'ease-out'});
    });
  }

  // Click ripple — offsetX/offsetY give position relative to element
  var rippleEl=document.getElementById('ripple-btn');
  if(rippleEl){
    rippleEl.addEventListener('click',function(e){
      var dot=document.createElement('span');
      dot.style.cssText='position:absolute;width:8px;height:8px;border-radius:50%;background:var(--primary);pointer-events:none;transform-origin:center center';
      dot.style.left=(e.offsetX-4)+'px';
      dot.style.top =(e.offsetY-4)+'px';
      rippleEl.appendChild(dot);
      Motion.animate(dot,{scale:[0,10],opacity:[0.6,0]},{duration:0.5,easing:'ease-out'}).then(function(){dot.remove();});
    });
  }
})()`)),
			)
		},
	})
}
