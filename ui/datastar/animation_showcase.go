//go:build showcase

package datastar

import (
	"mljr-web/ui"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "ds-animation", Name: "Animation (Motion)", Category: "datastar",
		Summary: "Motion v10 (window.Motion, 24 KB) provides animate(), timeline(), stagger(), and inView(). Call from any Datastar expression.",
		Code: `// Motion is available globally as window.Motion in all showcase previews.
// Vendor: /static/motion.min.js (motion v10 UMD build, 24 KB)

// Basic animate: (targets, keyframes, options)
Motion.animate('#el', {opacity:[0,1], y:[8,0]}, {duration:0.3, easing:'ease-out'})

// From Datastar event
data-on:click="Motion.animate(el, {scale:[0.95,1]}, {duration:0.15}); $open=true"

// Stagger — offset each animation start
Motion.animate('.list-item', {opacity:[0,1], x:[-12,0]},
    {duration:0.35, delay: Motion.stagger(0.06)})

// Timeline — sequence with offsets
Motion.timeline([
    ['.a', {opacity:[0,1]}, {duration:0.2}],
    ['.b', {opacity:[0,1]}, {duration:0.2, at:'-0.1'}],
    ['.c', {opacity:[0,1]}, {duration:0.2, at:'-0.1'}],
])

// Scroll-triggered
Motion.inView('.hero', ({target}) => {
    Motion.animate(target, {opacity:[0,1], y:[24,0]}, {duration:0.5})
})

// Number counter
Motion.animate(el, {}, {duration:1.5,
    onUpdate: p => el.textContent = Math.round(p * 9420).toLocaleString()
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				ui.Signals(`{"_aOpen":false,"_aCount":0}`),
				h.Div(h.Style("display:flex;flex-direction:column;gap:var(--sp-5)"),

					// Basic animate
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("animate() — basic")),
						h.Div(h.Style("display:flex;gap:var(--sp-3);flex-wrap:wrap"),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", `Motion.animate('#anim-box',{x:[0,120]},{duration:0.4,easing:'spring(1,80,10)'})`),
								g.Text("Slide →"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", `Motion.animate('#anim-box',{rotate:[0,360]},{duration:0.6,easing:'ease-in-out'})`),
								g.Text("Spin"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
								g.Attr("data-on:click", `Motion.animate('#anim-box',{scale:[1,1.5,1]},{duration:0.4})`),
								g.Text("Pulse"),
							),
							primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
								g.Attr("data-on:click", `Motion.animate('#anim-box',{x:0,rotate:0,scale:1},{duration:0.3})`),
								g.Text("Reset"),
							),
						),
						h.Div(h.Style("margin-top:var(--sp-4);height:80px;display:flex;align-items:center"),
							h.Div(
								h.ID("anim-box"),
								h.Style("width:60px;height:60px;background:var(--accent);border-radius:var(--radius);display:flex;align-items:center;justify-content:center;color:var(--accent-ink);font-weight:900;font-family:var(--font-display)"),
								g.Text("A"),
							),
						),
					),

					// Stagger
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("stagger() — offset each item")),
						primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
							g.Attr("data-on:click", `Motion.animate('.stagger-item',{opacity:[0,1],x:[-16,0]},{duration:0.4,delay:Motion.stagger(0.07)})`),
							g.Text("Animate list"),
						),
						h.Div(h.Style("margin-top:var(--sp-3);display:flex;flex-direction:column;gap:var(--sp-2)"),
							func() g.Node {
								items := []string{"Deploy to production", "Write unit tests", "Fix that one bug", "Delete console.logs", "Open a PR"}
								nodes := make([]g.Node, len(items))
								for i, item := range items {
									nodes[i] = h.Div(
										h.Class("stagger-item"),
										h.Style("padding:var(--sp-2) var(--sp-3);background:var(--surface-2);border-radius:var(--radius);font-size:var(--t-sm);opacity:0"),
										g.Text(item),
									)
								}
								return g.Group(nodes)
							}(),
						),
					),

					// Timeline
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("timeline() — sequenced choreography")),
						primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
							g.Attr("data-on:click", `
Motion.timeline([
  ['.tl-a',{opacity:[0,1],y:[12,0]},{duration:.25}],
  ['.tl-b',{opacity:[0,1],y:[12,0]},{duration:.25,at:'-0.1'}],
  ['.tl-c',{opacity:[0,1],y:[12,0]},{duration:.25,at:'-0.1'}],
  ['.tl-d',{opacity:[0,1],scale:[0.8,1]},{duration:.3,at:'-0.05'}]
]).play()`),
							g.Text("Play sequence"),
						),
						h.Div(h.Style("margin-top:var(--sp-3);display:flex;gap:var(--sp-3);align-items:flex-end"),
							h.Div(h.Class("tl-a"), h.Style("opacity:0;width:50px;height:80px;background:var(--accent);border-radius:var(--radius)")),
							h.Div(h.Class("tl-b"), h.Style("opacity:0;width:50px;height:60px;background:var(--accent-2);border-radius:var(--radius)")),
							h.Div(h.Class("tl-c"), h.Style("opacity:0;width:50px;height:100px;background:var(--success);border-radius:var(--radius)")),
							h.Div(h.Class("tl-d"), h.Style("opacity:0;width:50px;height:40px;background:var(--warning);border-radius:var(--radius)")),
						),
					),

					// Number counter
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("number counter — onUpdate callback")),
						primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeSM},
							g.Attr("data-on:click", `var el=document.getElementById('counter-val');Motion.animate(el,{},{duration:1.5,easing:'ease-out',onUpdate:p=>el.textContent=Math.round(p*9420).toLocaleString()})`),
							g.Text("Count up"),
						),
						h.Div(
							h.ID("counter-val"),
							h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display);margin-top:var(--sp-3)"),
							g.Text("0"),
						),
					),

					// Datastar + Motion button
					primitive.Card(primitive.CardProps{},
						h.P(h.Style("font-size:var(--t-xs);text-transform:uppercase;font-weight:700;letter-spacing:.06em;opacity:.5;margin-bottom:var(--sp-3)"), g.Text("Motion + Datastar — animate then update state")),
						h.Div(h.Style("display:flex;gap:var(--sp-3);align-items:center"),
							h.Button(
								g.Attr("data-component", "button"),
								g.Attr("data-attr", `{"data-variant":$_aOpen?"primary":"outline"}`),
								// Animate the panel, THEN flip the signal
								g.Attr("data-on:click",
									`var panel=document.getElementById('anim-panel');`+
										`if(!$_aOpen){`+
										`panel.style.display='block';`+
										`Motion.animate(panel,{opacity:[0,1],y:[-8,0]},{duration:0.2});`+
										`}else{`+
										`Motion.animate(panel,{opacity:0,y:-8},{duration:0.15}).then(()=>panel.style.display='none');`+
										`}`+
										`$_aOpen=!$_aOpen`),
								g.Attr("data-text", "$_aOpen?'▲ Close':'▼ Open'"),
							),
						),
						h.Div(
							h.ID("anim-panel"),
							h.Style("display:none;margin-top:var(--sp-3);padding:var(--sp-4);background:var(--surface-2);border-radius:var(--radius)"),
							h.P(g.Text("Animated enter/exit. Motion handles opacity+y, Datastar manages "),
								h.Code(g.Text("$_aOpen")),
								g.Text(" state."),
							),
						),
					),
				),
			)
		},
	})
}
