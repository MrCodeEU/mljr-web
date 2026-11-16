package special

import (
	"regexp"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// LogoGridProps configures LogoGrid, the swiss-theme take on LogoScatter:
// a flat two-colour logo that sits assembled in a slot and, on scroll,
// snaps its bars onto the column lines of the page grid ("Raster-Snap").
type LogoGridProps struct {
	SlotID  string // element the assembled logo is centred on (required)
	SheetID string // element whose content box defines the column grid (required)
	HeroID  string // bars land below 55% of this element's height (required)
	Copies  int    // pieces per bar on desktop; mobile always uses the originals. Default 2
}

var logoPathRe = regexp.MustCompile(`<path d="([^"]+)" transform="([^"]+)"`)

// rSplitX: paths starting at or right of this x (pre-transform) belong to the R.
const rSplitX = 1050

// LogoGrid renders the logo layer. Place it as the first child of a
// position:relative wrapper that spans the whole page; it fills that wrapper.
// Sections marked data-noscatter are kept free of bars. Static under
// prefers-reduced-motion.
func LogoGrid(p LogoGridProps) g.Node {
	if p.Copies < 1 {
		p.Copies = 2
	}
	var b strings.Builder
	b.WriteString(`<svg data-slot="logo" viewBox="0 0 2666.6667 2666.6667" xmlns="http://www.w3.org/2000/svg" focusable="false">`)
	for c := 0; c < p.Copies; c++ {
		if c > 0 {
			b.WriteString(`<g data-copy="1">`)
		}
		for _, m := range logoPathRe.FindAllStringSubmatch(logoSVGTpl, -1) {
			b.WriteString(`<g data-piece><path d="` + m[1] + `" transform="` + m[2] + `"`)
			if startX(m[1]) >= rSplitX {
				b.WriteString(` data-r="1"`)
			}
			b.WriteString(`/></g>`)
		}
		if c > 0 {
			b.WriteString(`</g>`)
		}
	}
	b.WriteString(`</svg>`)

	r := strings.NewReplacer("__SLOT__", p.SlotID, "__SHEET__", p.SheetID, "__HERO__", p.HeroID)
	return h.Div(
		g.Attr("data-component", "logo-grid"),
		g.Attr("aria-hidden", "true"),
		g.Raw(b.String()),
		h.Script(g.Raw(r.Replace(logoGridJS))),
	)
}

func startX(d string) float64 {
	f := strings.FieldsFunc(strings.TrimLeft(d, "mM "), func(r rune) bool { return r == ',' || r == ' ' })
	if len(f) == 0 {
		return 0
	}
	x, _ := strconv.ParseFloat(f[0], 64)
	return x
}

const logoGridJS = `(function(){
  var layer=document.currentScript.parentNode;
  if(document.readyState==='loading') document.addEventListener('DOMContentLoaded',init); else init();
  function init(){
  var svg=layer.querySelector('svg');
  var slot=document.getElementById('__SLOT__'), sheet=document.getElementById('__SHEET__'), hero=document.getElementById('__HERO__');
  if(!slot||!sheet||!hero||!svg.animate) return;
  var reduced=matchMedia('(prefers-reduced-motion: reduce)').matches, mq=matchMedia('(max-width:1024px)');
  var VB=2666.6667, MIN_K=0.15, MARGIN=1, OP_A=1, OP_S=0.3, scattered=false, anims=[];
  var pieces=[].map.call(svg.querySelectorAll('g[data-piece]'),function(el){
    el.style.transformBox='fill-box'; el.style.transformOrigin='center';
    return {g:el,copy:!!el.closest('[data-copy]')};
  });
  function live(){return mq.matches?pieces.filter(function(pc){return !pc.copy;}):pieces;}
  function tf(x,y,r,s){return 'translate('+x+'px,'+y+'px) rotate('+r+'deg) scale('+s+')';}
  function stopAll(keep){anims.forEach(function(a){if(keep){try{a.commitStyles();}catch(e){}}a.cancel();});anims=[];}
  function shuffle(a){for(var i=a.length-1;i>0;i--){var j=Math.random()*(i+1)|0,t=a[i];a[i]=a[j];a[j]=t;}return a;}
  function rect(el){var r=el.getBoundingClientRect();return {l:r.left+scrollX,t:r.top+scrollY,w:r.width,h:r.height};}
  function place(){
    var c=rect(slot), p=rect(layer);
    var size=Math.min(Math.max(120,Math.min(c.w,c.h)*1.7),mq.matches?200:520);
    if(mq.matches) size=Math.min(size,Math.max(c.w,c.h)*1.3);
    svg.style.width=size+'px'; svg.style.height=size+'px';
    svg.style.left=(c.l-p.l+c.w/2-size/2)+'px'; svg.style.top=(c.t-p.t+c.h/2-size/2)+'px';
  }
  function measure(){
    var s=rect(svg), p=rect(layer), k=s.w/VB;
    live().forEach(function(pc){var b=pc.g.getBBox();pc.k=k;pc.cx=s.l-p.l+(b.x+b.width/2)*k;pc.cy=s.t-p.t+(b.y+b.height/2)*k;});
  }
  function targets(){
    var sr=rect(sheet), p=rect(layer), cs=getComputedStyle(sheet);
    var left=sr.l-p.l+parseFloat(cs.paddingLeft), w=sr.w-parseFloat(cs.paddingLeft)-parseFloat(cs.paddingRight);
    var cols=mq.matches?4:12, lines=[];
    for(var i=MARGIN;i<=cols-MARGIN;i++) lines.push(left+i*w/cols);
    var hr=rect(hero), top=hr.t-p.t+hr.h*0.55, bottom=layer.scrollHeight-120;
    var keep=[].map.call(document.querySelectorAll('[data-noscatter]'),function(el){var r=rect(el);return [r.t-p.t-90,r.t-p.t+r.h+90];});
    function freeY(){for(var n=0;n<40;n++){var y=Math.round((top+Math.random()*(bottom-top))/24)*24;if(!keep.some(function(z){return y>z[0]&&y<z[1];})) return y;}return top;}
    shuffle(live().slice()).forEach(function(pc,j){
      var horiz=Math.random()<0.18, tx=lines[j%lines.length], ty=freeY();
      pc.t={x:(tx-pc.cx)/pc.k,y:(ty-pc.cy)/pc.k,r:horiz?90:0,s:Math.max(1,MIN_K/pc.k)};
      pc.dur=Math.max(700,Math.min(2200,Math.hypot(tx-pc.cx,ty-pc.cy)*0.9)); pc.delay=j*14;
    });
  }
  function scatter(instant){
    if(reduced) return;
    stopAll(true); scattered=true; targets();
    anims.push(svg.animate([{opacity:OP_S}],{duration:instant?0:500,fill:'forwards'}));
    live().forEach(function(pc){anims.push(pc.g.animate([{transform:tf(pc.t.x,pc.t.y,pc.t.r,pc.t.s)}],{duration:instant?0:pc.dur,delay:instant?0:pc.delay,easing:'cubic-bezier(.7,0,.2,1)',fill:'forwards'}));});
  }
  function assemble(instant){
    stopAll(true); scattered=false;
    anims.push(svg.animate([{opacity:OP_A}],{duration:instant?0:700,fill:'forwards'}));
    live().forEach(function(pc,j){anims.push(pc.g.animate([{transform:tf(0,0,0,1)}],{duration:instant?0:1100,delay:instant?0:j*6,easing:'cubic-bezier(.65,0,.35,1)',fill:'forwards'}));});
  }
  function relayout(){
    var was=scattered; stopAll(false); pieces.forEach(function(pc){pc.g.style.transform='';});
    place(); measure();
    if(was) scatter(true); else assemble(true);
  }
  relayout();
  new IntersectionObserver(function(es){if(es[0].isIntersecting) assemble(false); else if(scrollY>80) scatter(false);},{threshold:.6}).observe(slot);
  var rt; addEventListener('resize',function(){clearTimeout(rt);rt=setTimeout(relayout,250);});
  if(document.fonts&&document.fonts.ready) document.fonts.ready.then(relayout);
  addEventListener('load',relayout);
  }
})();`
