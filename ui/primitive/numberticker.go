package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NumberTickerProps struct {
	// Value is the target number to animate to.
	Value float64
	// From is the starting value (default 0).
	From float64
	// Duration is the animation duration in ms (default 1200).
	Duration int
	// Decimals controls decimal places displayed (default 0).
	Decimals int
	// Prefix is displayed before the number (e.g. "€", "$").
	Prefix string
	// Suffix is displayed after the number (e.g. "%", "K").
	Suffix string
	// Locale enables toLocaleString formatting (default true).
	Locale bool
	// TriggerOnView starts animation when element enters viewport.
	TriggerOnView bool
	// ID must be unique per page when multiple tickers appear.
	ID string
}

// NumberTicker renders an animated number counter.
// Counts from From → Value using requestAnimationFrame (no Motion needed).
func NumberTicker(p NumberTickerProps) g.Node {
	if p.Duration == 0 {
		p.Duration = 1200
	}
	if p.ID == "" {
		p.ID = "ntick"
	}

	initial := p.Prefix + fmt.Sprintf("%."+fmt.Sprintf("%d", p.Decimals)+"f", p.From) + p.Suffix

	locale := "true"
	if !p.Locale {
		locale = "false"
	}

	script := fmt.Sprintf(`(function(){
  var el=document.getElementById('%s');
  if(!el) return;
  var from=%v,to=%v,dur=%d,dec=%d,pre='%s',suf='%s',loc=%s;
  function format(v){
    var s=v.toFixed(dec);
    if(loc){ var n=Number(s); s=n.toLocaleString(undefined,{minimumFractionDigits:dec,maximumFractionDigits:dec}); }
    return pre+s+suf;
  }
  function run(){
    var start=null;
    function step(ts){
      if(!start) start=ts;
      var pct=Math.min(1,(ts-start)/dur);
      var ease=1-Math.pow(1-pct,3); // ease-out cubic
      el.textContent=format(from+ease*(to-from));
      if(pct<1) requestAnimationFrame(step);
      else el.textContent=format(to);
    }
    requestAnimationFrame(step);
  }
  if(%v && 'IntersectionObserver' in window){
    var io=new IntersectionObserver(function(e){ if(e[0].isIntersecting){ run(); io.disconnect(); } });
    io.observe(el);
  } else { run(); }
})();`,
		p.ID, p.From, p.Value, p.Duration, p.Decimals,
		p.Prefix, p.Suffix, locale, p.TriggerOnView)

	return g.Group{
		h.Span(
			h.ID(p.ID),
			g.Attr("data-component", "number-ticker"),
			g.Text(initial),
		),
		h.Script(g.Raw(script)),
	}
}
