package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CountdownProps struct {
	// Target is an ISO 8601 datetime string (e.g. "2026-12-31T23:59:59").
	Target string
	// ID must be unique per page when multiple countdowns appear.
	ID      string
	Compact bool // show only total seconds remaining vs dd:hh:mm:ss units
	// OnExpire is a JS expression to run when the countdown reaches zero.
	OnExpire string
}

// Countdown renders a live countdown timer to a target datetime.
// Updates every second via JS setInterval. No server required.
func Countdown(p CountdownProps) g.Node {
	if p.ID == "" {
		p.ID = "countdown"
	}

	var display g.Node
	if p.Compact {
		display = h.Span(h.ID(p.ID+"-val"), g.Attr("data-component", "countdown"), g.Text("–"))
	} else {
		display = h.Div(
			g.Attr("data-component", "countdown"),
			h.Style("display:flex;gap:var(--sp-3);align-items:flex-end"),
			unit(p.ID+"-d", "Days"),
			sep(),
			unit(p.ID+"-h", "Hours"),
			sep(),
			unit(p.ID+"-m", "Min"),
			sep(),
			unit(p.ID+"-s", "Sec"),
		)
	}

	onExpire := p.OnExpire
	if onExpire == "" {
		onExpire = ""
	}

	script := fmt.Sprintf(`(function(){
  var target=new Date('%s').getTime();
  var compact=%v;
  var onExpire=function(){%s};

  function pad(n){return n<10?'0'+n:String(n);}
  function tick(){
    var now=Date.now();
    var diff=Math.max(0,target-now);
    if(diff===0) onExpire();
    var s=Math.floor(diff/1000)%%60;
    var m=Math.floor(diff/60000)%%60;
    var h=Math.floor(diff/3600000)%%24;
    var d=Math.floor(diff/86400000);
    if(compact){
      var el=document.getElementById('%s-val');
      if(el) el.textContent=Math.floor(diff/1000)+'s';
    } else {
      var set=function(id,v){ var e=document.getElementById(id); if(e) e.textContent=pad(v); };
      set('%s-d',d); set('%s-h',h); set('%s-m',m); set('%s-s',s);
    }
  }
  tick();
  setInterval(tick,1000);
})();`,
		p.Target, p.Compact, onExpire,
		p.ID, p.ID, p.ID, p.ID, p.ID)

	return g.Group{
		display,
		h.Script(g.Raw(script)),
	}
}

func unit(id, label string) g.Node {
	return h.Div(
		h.Style("display:flex;flex-direction:column;align-items:center;min-width:3ch"),
		h.Span(
			h.ID(id),
			h.Style("font-size:var(--t-2xl);font-weight:900;font-family:var(--font-display);letter-spacing:-.02em;line-height:1"),
			g.Text("–"),
		),
		h.Span(
			h.Style("font-size:9px;text-transform:uppercase;letter-spacing:.08em;opacity:.5;margin-top:2px"),
			g.Text(label),
		),
	)
}

func sep() g.Node {
	return h.Span(h.Style("font-size:var(--t-xl);font-weight:900;opacity:.3;padding-bottom:var(--sp-3)"), g.Text(":"))
}
