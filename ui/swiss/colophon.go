package swiss

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// NavItem is one section in the ruler and the register.
type NavItem struct {
	ID    string // section id (without #)
	Num   string
	Label string // German label
	Sub   string // English sublabel
}

type ColophonProps struct {
	BrandHref string
	Brand     string // bold part, e.g. "mljr.eu"
	BrandSub  string // hidden on mobile, e.g. "— Michael Reinegger · Linz (AT)"
	// Items feed the ruler (desktop) and the register dialog (mobile).
	// Empty = no ruler/register (sub pages).
	Items   []NavItem
	Actions []g.Node // right-hand controls (language, mode, …)
	// PageID is the element made inert while the register is open.
	PageID        string
	NavLabel      string // aria-label of the ruler, default "Abschnitte"
	RegisterLabel string // default "Register"
	CloseLabel    string // default "Schliessen"
	// SkipLabel is the first link in the page, jumping to #main; empty
	// renders no skip link.
	SkipLabel string
}

// Colophon is the sticky top bar: brand, section ruler ("Lineal") with
// per-section fill, actions, a register button on mobile and a scroll
// progress rule. Render Register(p) once at the end of <body>.
func Colophon(p ColophonProps) g.Node {
	defaults(&p)
	seg := func(it NavItem) g.Node {
		txt := func() g.Node {
			return h.Span(g.Attr("data-slot", "txt"),
				h.Span(g.Attr("data-slot", "num"), g.Text(it.Num)),
				h.Span(g.Attr("data-slot", "label"), g.Text(it.Label)),
			)
		}
		return h.Li(h.A(h.Href("#"+it.ID), g.Attr("data-nav", it.ID), g.Attr("aria-label", it.Num+" "+it.Label),
			txt(),
			h.Span(g.Attr("data-slot", "fill"), g.Attr("aria-hidden", "true"), txt()),
		))
	}
	segs := make([]g.Node, len(p.Items))
	for i, it := range p.Items {
		segs[i] = seg(it)
	}
	hasNav := len(p.Items) > 0
	return h.Header(
		g.Attr("data-component", "colophon"),
		g.If(p.SkipLabel != "", h.A(g.Attr("data-slot", "skip"), h.Href("#main"), g.Text(p.SkipLabel))),
		h.A(g.Attr("data-slot", "brand"), h.Href(p.BrandHref),
			h.B(g.Text(p.Brand)),
			g.If(p.BrandSub != "", h.Span(g.Text(p.BrandSub))),
		),
		g.If(hasNav, h.Nav(g.Attr("data-component", "ruler"), g.Attr("aria-label", p.NavLabel), h.Ol(segs...))),
		h.Div(g.Attr("data-slot", "actions"),
			g.Group(p.Actions),
			g.If(hasNav, Btn(false, []g.Node{
				g.Attr("data-slot", "register-btn"), h.ID("swiss-register-btn"),
				g.Attr("aria-haspopup", "dialog"), g.Attr("aria-controls", "swiss-register"),
			}, g.Text(p.RegisterLabel+" ☰"))),
		),
		h.Span(g.Attr("data-slot", "progress"), g.Attr("aria-hidden", "true")),
		h.Script(g.Raw(colophonJS)),
	)
}

// Register is the full-screen index dialog opened from the colophon on mobile.
func Register(p ColophonProps) g.Node {
	defaults(&p)
	if len(p.Items) == 0 {
		return nil
	}
	items := make([]g.Node, len(p.Items))
	for i, it := range p.Items {
		items[i] = h.Li(h.A(h.Href("#"+it.ID), g.Attr("data-nav", it.ID),
			h.Span(g.Attr("data-slot", "num"), g.Text(it.Num)),
			h.Span(g.Attr("data-slot", "label"), g.Text(it.Label),
				g.If(it.Sub != "", h.Small(g.Attr("lang", "en"), g.Text(it.Sub))),
			),
		))
	}
	return h.Div(
		g.Attr("data-component", "register"), h.ID("swiss-register"),
		h.Role("dialog"), g.Attr("aria-modal", "true"), g.Attr("aria-label", p.RegisterLabel),
		g.Attr("data-page", p.PageID), g.Attr("hidden", ""),
		h.Div(g.Attr("data-slot", "head"),
			Micro(g.Text(p.RegisterLabel+" / Index")),
			Btn(false, []g.Node{h.ID("swiss-register-close")}, g.Text(p.CloseLabel+" ×")),
		),
		h.Ol(items...),
		h.Script(g.Raw(registerJS)),
	)
}

func defaults(p *ColophonProps) {
	if p.NavLabel == "" {
		p.NavLabel = "Abschnitte"
	}
	if p.RegisterLabel == "" {
		p.RegisterLabel = "Register"
	}
	if p.CloseLabel == "" {
		p.CloseLabel = "Schliessen"
	}
	if p.BrandHref == "" {
		p.BrandHref = "/"
	}
}

// colophonJS drives the progress rule and the ruler fill. Each section owns a
// scroll range [b[i], b[i+1]): b[0]=0 and the last boundary is max scroll, so
// the ruler starts empty and ends full. Tail sections whose top can never
// reach the header share the last stretch in proportion to their heights.
const colophonJS = `(function(){
  var hdr=document.currentScript.parentNode, prog=hdr.querySelector('[data-slot=progress]');
  var root=document.documentElement;
  function sync(){root.style.setProperty('--swiss-hdr',hdr.offsetHeight+'px');}
  var links=[].slice.call(document.querySelectorAll('a[data-nav]'));
  var ids=[]; links.forEach(function(a){if(ids.indexOf(a.dataset.nav)<0) ids.push(a.dataset.nav);});
  var ticking=false;
  function update(){
    ticking=false;
    var max=Math.max(1,root.scrollHeight-innerHeight);
    prog.style.transform='scaleX('+Math.min(1,scrollY/max)+')';
    var secs=ids.map(function(id){return document.getElementById(id);}).filter(Boolean);
    if(!secs.length) return;
    var hh=hdr.offsetHeight;
    var tops=secs.map(function(s,i){return i===0?0:s.getBoundingClientRect().top+scrollY-hh;});
    var b=tops.slice(); b.push(max);
    var k=-1; for(var i=0;i<tops.length;i++){if(tops[i]>=max){k=i;break;}}
    if(k>0){
      var from=tops[k-1], hs=secs.slice(k-1).map(function(s){return s.offsetHeight;}), sum=hs.reduce(function(a,v){return a+v;},0), acc=from;
      hs.forEach(function(h,j){b[k-1+j]=acc; acc+=(max-from)*h/sum;});
    }
    var act=0; for(var j=0;j<secs.length;j++) if(scrollY>=b[j]) act=j;
    var within=Math.min(1,Math.max(0,(scrollY-b[act])/Math.max(1,b[act+1]-b[act])));
    links.forEach(function(a){
      var i=ids.indexOf(a.dataset.nav);
      if(i===act) a.setAttribute('aria-current','location'); else a.removeAttribute('aria-current');
      a.style.setProperty('--p',i<act?1:i===act?within:0);
    });
  }
  function req(){if(!ticking){ticking=true;requestAnimationFrame(update);}}
  addEventListener('scroll',req,{passive:true});
  addEventListener('resize',function(){sync();req();});
  addEventListener('load',function(){sync();req();});
  sync(); req();
})();`

const registerJS = `(function(){
  var reg=document.currentScript.parentNode, btn=document.getElementById('swiss-register-btn'), close=document.getElementById('swiss-register-close');
  var page=document.getElementById(reg.dataset.page);
  if(!btn) return;
  function open(){reg.hidden=false;if(page) page.inert=true;close.focus();}
  function shut(back){reg.hidden=true;if(page) page.inert=false;if(back) btn.focus();}
  btn.addEventListener('click',open);
  close.addEventListener('click',function(){shut(true);});
  reg.addEventListener('click',function(e){if(e.target.closest('a')) shut(false);});
  addEventListener('keydown',function(e){if(e.key==='Escape'&&!reg.hidden) shut(true);});
})();`
