package special

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TourStep struct {
	// Target is a CSS selector for the element to highlight.
	Target string
	// Title shown in the tooltip.
	Title string
	// Body is the step description.
	Body string
	// Placement: "top" | "bottom" | "left" | "right" (default "bottom").
	Placement string
}

type TourProps struct {
	// Signal prefix (default "_tour").
	Signal string
	// AutoStart launches the tour on page load.
	AutoStart bool
}

// Tour renders an onboarding tour overlay with step-by-step highlights.
// Uses a floating tooltip + semi-transparent spotlight overlay.
// No external library — pure JS + Datastar signal for step tracking.
func Tour(p TourProps, steps ...TourStep) g.Node {
	if p.Signal == "" {
		p.Signal = "_tour"
	}
	if len(steps) == 0 {
		return g.Text("")
	}

	sig := p.Signal

	// Build JS steps array
	stepParts := make([]string, len(steps))
	for i, s := range steps {
		placement := s.Placement
		if placement == "" {
			placement = "bottom"
		}
		stepParts[i] = fmt.Sprintf(`{target:%s,title:%s,body:%s,placement:%s}`,
			jsStr(s.Target), jsStr(s.Title), jsStr(s.Body), jsStr(placement))
	}
	stepsJS := "[" + strings.Join(stepParts, ",") + "]"

	script := fmt.Sprintf(`(function(){
  var steps=%s;
  var overlay=document.getElementById('%sOverlay');
  var tooltip=document.getElementById('%sTip');
  var tipTitle=document.getElementById('%sTipTitle');
  var tipBody=document.getElementById('%sTipBody');
  var tipPrev=document.getElementById('%sTipPrev');
  var tipNext=document.getElementById('%sTipNext');
  var tipCount=document.getElementById('%sTipCount');
  if(!overlay) return;

  var cur=0;

  function show(i){
    cur=i;
    var step=steps[i];
    var target=document.querySelector(step.target);
    if(!target){ next(); return; }

    var r=target.getBoundingClientRect();
    var pad=8;
    // spotlight
    overlay.style.display='block';
    overlay.style.setProperty('--sx',r.left-pad+'px');
    overlay.style.setProperty('--sy',r.top+window.scrollY-pad+'px');
    overlay.style.setProperty('--sw',r.width+pad*2+'px');
    overlay.style.setProperty('--sh',r.height+pad*2+'px');

    // tooltip
    tipTitle.textContent=step.title;
    tipBody.textContent=step.body;
    tipCount.textContent=(i+1)+' / '+steps.length;
    tipPrev.disabled=i===0;
    tipNext.textContent=i===steps.length-1?'Finish':'Next';

    // position tooltip
    var tipH=tooltip.offsetHeight||100;
    var tipW=tooltip.offsetWidth||240;
    var top,left;
    if(step.placement==='bottom'){ top=r.bottom+window.scrollY+pad+8; left=r.left+r.width/2-tipW/2; }
    else if(step.placement==='top'){ top=r.top+window.scrollY-tipH-pad-8; left=r.left+r.width/2-tipW/2; }
    else if(step.placement==='right'){ top=r.top+window.scrollY+r.height/2-tipH/2; left=r.right+pad+8; }
    else { top=r.top+window.scrollY+r.height/2-tipH/2; left=r.left-tipW-pad-8; }
    tooltip.style.top=top+'px';
    tooltip.style.left=Math.max(8,left)+'px';
    tooltip.style.display='block';

    target.scrollIntoView({block:'nearest',behavior:'smooth'});
  }

  function next(){ if(cur<steps.length-1){ show(cur+1); } else { close(); } }
  function prev(){ if(cur>0) show(cur-1); }
  function close(){ overlay.style.display='none'; tooltip.style.display='none'; }

  tipNext.addEventListener('click',next);
  tipPrev.addEventListener('click',prev);
  document.getElementById('%sTipClose').addEventListener('click',close);
  overlay.addEventListener('click',close);

  window['%sStart']=function(){ show(0); };
  document.querySelectorAll('[data-tour-start]').forEach(function(b){ b.addEventListener('click',function(){ show(0); }); });
  %s
})();`,
		stepsJS, sig, sig, sig, sig, sig, sig, sig, sig, sig,
		func() string {
			if p.AutoStart {
				return "window.addEventListener('load',function(){ window['" + sig + "Start'](); });"
			}
			return ""
		}(),
	)

	return g.Group{
		// Spotlight overlay
		h.Div(
			h.ID(sig+"Overlay"),
			h.Style("display:none;position:fixed;inset:0;z-index:1000;background:rgba(0,0,0,0.55);--sx:0;--sy:0;--sw:0;--sh:0;"+
				"clip-path:polygon(0% 0%,0% 100%,var(--sx) 100%,var(--sx) var(--sy),calc(var(--sx) + var(--sw)) var(--sy),calc(var(--sx) + var(--sw)) calc(var(--sy) + var(--sh)),var(--sx) calc(var(--sy) + var(--sh)),var(--sx) 100%,100% 100%,100% 0%)"),
		),
		// Tooltip
		h.Div(
			h.ID(sig+"Tip"),
			g.Attr("data-component", "tour-tooltip"),
			h.Style("display:none;position:absolute;z-index:1001;width:280px;max-width:90vw"),
			h.Div(
				h.Style("background:var(--surface);border:var(--bw-2) solid var(--ink);border-radius:var(--radius);box-shadow:var(--shadow-lg);padding:var(--sp-4)"),
				h.Div(
					h.Style("display:flex;align-items:flex-start;justify-content:space-between;gap:var(--sp-2);margin-bottom:var(--sp-2)"),
					h.Strong(h.ID(sig+"TipTitle"), h.Style("font-size:var(--t-base);font-weight:800")),
					h.Button(h.ID(sig+"TipClose"), h.Type("button"), h.Style("background:none;border:none;cursor:pointer;font-size:1.2rem;color:var(--muted);padding:0;line-height:1"), g.Text("×")),
				),
				h.P(h.ID(sig+"TipBody"), h.Style("font-size:var(--t-sm);color:var(--muted);margin:0 0 var(--sp-3)")),
				h.Div(
					h.Style("display:flex;align-items:center;justify-content:space-between;gap:var(--sp-2)"),
					h.Span(h.ID(sig+"TipCount"), h.Style("font-size:var(--t-xs);color:var(--muted);font-weight:700")),
					h.Div(
						h.Style("display:flex;gap:var(--sp-2)"),
						h.Button(h.ID(sig+"TipPrev"), h.Type("button"), g.Attr("data-component", "button"), g.Attr("data-variant", "outline"), g.Attr("data-size", "sm"), g.Text("Back")),
						h.Button(h.ID(sig+"TipNext"), h.Type("button"), g.Attr("data-component", "button"), g.Attr("data-variant", "primary"), g.Attr("data-size", "sm"), g.Text("Next")),
					),
				),
			),
		),
		h.Script(g.Raw(script)),
	}
}
