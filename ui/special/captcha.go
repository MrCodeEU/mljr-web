package special

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CaptchaProps struct {
	ChallengeURL string // e.g. "/api/altcha"
	Name         string // hidden input name, default "altcha"
	// Lazy loads the widget script only when the widget nears the viewport
	// and starts the proof-of-work when the form gets focus, instead of both
	// on page load.
	Lazy  bool
	Attrs []g.Node
}

// captchaLazyJS injects the loader once the widget is within 600px of the
// viewport. Inline, so CSP allows it via its hash.
const captchaLazyJS = `(function(){
  var w=document.currentScript.parentNode, done=false;
  function load(){ if(done) return; done=true; var s=document.createElement('script'); s.type='module'; s.src='/static/altcha-loader.js'; document.head.appendChild(s); }
  if(!window.IntersectionObserver) return load();
  var io=new IntersectionObserver(function(es){ if(es.some(function(e){return e.isIntersecting;})){ io.disconnect(); load(); } },{rootMargin:'600px 0px'});
  io.observe(w);
})();`

// Captcha renders the self-hosted altcha widget. The JS is vendored at
// /static/altcha.js (no CDN, no telemetry). The widget performs proof-of-work
// client-side and writes its solution to a hidden input named p.Name inside the
// widget element.
func Captcha(p CaptchaProps) g.Node {
	if p.Name == "" {
		p.Name = "altcha"
	}
	if p.ChallengeURL == "" {
		p.ChallengeURL = "/api/altcha"
	}
	return h.Div(
		g.Attr("data-component", "captcha"),
		g.Group(p.Attrs),
		// altcha-loader imports the vendored ESM bundle and registers the local
		// SHA worker used by altcha-lib-go's legacy SHA-256 challenges.
		g.If(!p.Lazy, h.Script(h.Type("module"), h.Src("/static/altcha-loader.js"))),
		g.If(p.Lazy, h.Script(g.Raw(captchaLazyJS))),
		// altcha-widget: auto="onload" starts proof-of-work immediately,
		// "onfocus" when the surrounding form is focused.
		g.El("altcha-widget",
			g.Attr("challenge", p.ChallengeURL),
			g.Attr("name", p.Name),
			g.If(!p.Lazy, g.Attr("auto", "onload")),
			g.If(p.Lazy, g.Attr("auto", "onfocus")),
			g.Attr("hidefooter", ""),
		),
	)
}

// HoneypotField renders an invisible input that bots fill and humans don't.
// Verify server-side: reject if the field value is non-empty.
func HoneypotField(name string) g.Node {
	if name == "" {
		name = "_hp"
	}
	return h.Input(
		h.Type("text"),
		h.Name(name),
		g.Attr("tabindex", "-1"),
		g.Attr("autocomplete", "off"),
		h.Style("position:absolute;left:-9999px;width:1px;height:1px;overflow:hidden"),
		g.Attr("aria-hidden", "true"),
	)
}
