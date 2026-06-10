package special

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CookieBannerProps struct {
	// Message is the consent notice text.
	Message string
	// PolicyHref links to the privacy policy.
	PolicyHref string
	// Accept is the accept button label (default "Accept all").
	Accept string
	// Decline is the decline button label (empty = no decline button).
	Decline string
	// Storage is the localStorage key (default "mljr-cookie-consent").
	Storage string
	// Position: "bottom" (default) | "top"
	Position string
}

// CookieBanner renders a GDPR cookie consent banner.
// Hides itself immediately if localStorage key is already set.
// No Datastar dependency — pure JS.
func CookieBanner(p CookieBannerProps) g.Node {
	if p.Message == "" {
		p.Message = "We use cookies to improve your experience and analyze site usage."
	}
	if p.Accept == "" {
		p.Accept = "Accept all"
	}
	if p.Storage == "" {
		p.Storage = "mljr-cookie-consent"
	}
	if p.Position == "" {
		p.Position = "bottom"
	}

	posStyle := "bottom:0;left:0;right:0"
	if p.Position == "top" {
		posStyle = "top:0;left:0;right:0"
	}

	declineBtn := ""
	if p.Decline != "" {
		declineBtn = fmt.Sprintf(`var d=document.createElement('button');d.type='button';d.textContent=%s;d.style.cssText='padding:var(--sp-2) var(--sp-4);font-size:var(--t-sm);font-weight:700;background:none;border:var(--bw-2) solid var(--ink);border-radius:var(--radius);cursor:pointer';d.onclick=function(){localStorage.setItem(%s,'declined');banner.remove();};btns.appendChild(d);`,
			jsStr(p.Decline), jsStr(p.Storage))
	}

	policyLink := ""
	if p.PolicyHref != "" {
		policyLink = fmt.Sprintf(` <a href=%s style="color:var(--accent);text-decoration:underline;font-weight:700">Privacy Policy</a>`, jsStr(p.PolicyHref))
	}

	script := fmt.Sprintf(`(function(){
  if(localStorage.getItem(%s)) return;
  var banner=document.createElement('div');
  banner.style.cssText='position:fixed;%s;z-index:200;background:var(--surface);border-top:var(--bw-2) solid var(--ink);padding:var(--sp-4) var(--sp-6);display:flex;align-items:center;gap:var(--sp-4);flex-wrap:wrap;box-shadow:var(--shadow-lg)';
  var msg=document.createElement('p');
  msg.style.cssText='margin:0;flex:1;font-size:var(--t-sm);';
  msg.innerHTML=%s+%s;
  banner.appendChild(msg);
  var btns=document.createElement('div');
  btns.style.cssText='display:flex;gap:var(--sp-2);flex-wrap:wrap';
  %s
  var a=document.createElement('button');
  a.type='button';a.textContent=%s;
  a.style.cssText='padding:var(--sp-2) var(--sp-4);font-size:var(--t-sm);font-weight:700;background:var(--ink);color:var(--surface);border:var(--bw-2) solid var(--ink);border-radius:var(--radius);cursor:pointer';
  a.onclick=function(){localStorage.setItem(%s,'accepted');banner.remove();};
  btns.appendChild(a);
  banner.appendChild(btns);
  document.body.appendChild(banner);
})();`,
		jsStr(p.Storage),
		posStyle,
		jsStr(p.Message),
		policyLink,
		declineBtn,
		jsStr(p.Accept),
		jsStr(p.Storage),
	)

	return h.Script(g.Raw(script))
}
