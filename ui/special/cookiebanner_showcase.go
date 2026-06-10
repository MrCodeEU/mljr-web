//go:build showcase

package special

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "cookie-banner", Name: "Cookie Banner", Category: "special",
		Summary: "GDPR cookie consent banner. Renders via JS — hides instantly if already accepted. localStorage persist. No Datastar needed.",
		Code: `special.CookieBanner(special.CookieBannerProps{
    Message:    "We use cookies to improve your experience.",
    PolicyHref: "/privacy",
    Accept:     "Accept all",
    Decline:    "Decline",
    Storage:    "my-app-cookie-consent",
    Position:   "bottom",
})`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("position:relative;min-height:200px;border:var(--bw-1) dashed var(--line);border-radius:var(--radius);padding:var(--sp-5);overflow:hidden"),
				h.P(h.Style("color:var(--muted);font-size:var(--t-sm);margin-bottom:var(--sp-4)"), g.Text("Cookie banner renders at bottom of page. Clear storage below to reset.")),
				h.Button(
					h.Type("button"),
					g.Attr("data-component", "button"),
					g.Attr("data-variant", "outline"),
					h.Style("font-size:var(--t-sm)"),
					g.Attr("data-on:click", "localStorage.removeItem('demo-cookie-consent');window.location.reload()"),
					g.Text("Reset consent (reload to see banner)"),
				),
				// Render scoped banner for demo
				h.Script(g.Raw(`(function(){
  var key='demo-cookie-consent';
  if(localStorage.getItem(key)) return;
  var wrap=document.currentScript.parentElement;
  var bar=document.createElement('div');
  bar.style.cssText='position:absolute;bottom:0;left:0;right:0;background:var(--surface);border-top:var(--bw-2) solid var(--ink);padding:var(--sp-3) var(--sp-4);display:flex;align-items:center;gap:var(--sp-3);flex-wrap:wrap';
  bar.innerHTML='<p style="margin:0;flex:1;font-size:var(--t-sm)">We use cookies to improve your experience. <a href="#" style="color:var(--accent);font-weight:700">Privacy Policy</a></p>';
  var btns=document.createElement('div');
  btns.style.cssText='display:flex;gap:var(--sp-2)';
  var d=document.createElement('button');d.type='button';d.textContent='Decline';
  d.style.cssText='padding:var(--sp-1) var(--sp-3);font-size:var(--t-sm);font-weight:700;background:none;border:var(--bw-2) solid var(--ink);border-radius:var(--radius);cursor:pointer';
  d.onclick=function(){localStorage.setItem(key,'declined');bar.remove();};
  var a=document.createElement('button');a.type='button';a.textContent='Accept all';
  a.style.cssText='padding:var(--sp-1) var(--sp-3);font-size:var(--t-sm);font-weight:700;background:var(--ink);color:var(--surface);border:var(--bw-2) solid var(--ink);border-radius:var(--radius);cursor:pointer';
  a.onclick=function(){localStorage.setItem(key,'accepted');bar.remove();};
  btns.appendChild(d);btns.appendChild(a);bar.appendChild(btns);wrap.appendChild(bar);
})();`)),
			)
		},
	})
}
