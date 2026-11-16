package swiss

import (
	"encoding/json"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// FahrplanItem is one career entry.
type FahrplanItem struct {
	ID       string   `json:"id"`
	Kind     string   `json:"kind"` // "work" | "education" | "thesis"
	Title    string   `json:"title"`
	Org      string   `json:"org"`
	Short    string   `json:"short"` // chip label; defaults to Org
	Start    string   `json:"start"` // "YYYY-MM"
	End      string   `json:"end"`   // "YYYY-MM"; "" = ongoing
	Location string   `json:"location"`
	Summary  string   `json:"summary"`
	Details  []string `json:"details"`
	Tags     []string `json:"tags"`
	Logo     string   `json:"logo"`
}

type FahrplanProps struct {
	ID        string
	Items     []FahrplanItem
	StartYear int    // first column; default = earliest start year
	Now       string // "YYYY-MM"; ongoing entries end here
	// Labels (German defaults).
	LaneWork, LaneWorkSub, LaneEdu, LaneEduSub string
	KindWork, KindEdu, KindThesis              string
	Today, Present, Current, Hint              string
	Earlier, Later, Entry, ShowDetails, To     string
}

// Fahrplan renders the career timeline as a timetable: one grid column per
// year, square bars in two lanes (work / education) with tracks for
// overlaps, and a chip button per entry stacked above the bars with a leader
// line. Selecting a chip fills the detail panel and dispatches a bubbling
// "fahrplan:select" event ({detail:{id}}); arrow keys step through entries. Under 1024px it turns vertical, newest year on top. Without JS a
// plain list of all entries is shown.
func Fahrplan(p FahrplanProps) g.Node {
	fpDefaults(&p)
	for i := range p.Items {
		if p.Items[i].Short == "" {
			p.Items[i].Short = p.Items[i].Org
		}
		if p.StartYear == 0 || atoiPrefix(p.Items[i].Start) < p.StartYear {
			p.StartYear = atoiPrefix(p.Items[i].Start)
		}
	}
	data, err := json.Marshal(map[string]any{
		"items": p.Items, "start": p.StartYear, "now": p.Now,
		"l": map[string]string{
			"laneWork": p.LaneWork, "laneWorkSub": p.LaneWorkSub, "laneEdu": p.LaneEdu, "laneEduSub": p.LaneEduSub,
			"work": p.KindWork, "education": p.KindEdu, "thesis": p.KindThesis,
			"today": p.Today, "present": p.Present, "current": p.Current,
			"earlier": p.Earlier, "later": p.Later, "entry": p.Entry, "show": p.ShowDetails, "to": p.To,
		},
	})
	if err != nil {
		return g.Raw("<!-- swiss: " + err.Error() + " -->")
	}
	fallback := make([]g.Node, len(p.Items))
	for i, it := range p.Items {
		end := it.End
		if end == "" {
			end = p.Present
		}
		fallback[i] = h.Li(h.B(g.Text(it.Title)), g.Text(" · "+it.Org+" · "+it.Start+" – "+end))
	}
	legend := func(kind, label string) g.Node {
		return h.Span(h.I(g.Attr("data-kind", kind), h.Style("background:var(--k)")), g.Text(label))
	}
	return h.Div(
		g.Attr("data-component", "fahrplan"), h.ID(p.ID),
		h.Div(g.Attr("data-slot", "legend"),
			legend("work", p.KindWork),
			legend("education", p.KindEdu),
			legend("thesis", p.KindThesis),
			h.Span(g.Attr("data-slot", "hint"), g.Text(p.Hint)),
		),
		h.Div(g.Attr("data-slot", "plot"),
			h.Ol(g.Attr("data-slot", "fallback"), g.Group(fallback)),
		),
		h.Div(g.Attr("data-slot", "detail"), g.Attr("aria-live", "polite")),
		h.Script(h.Type("application/json"), g.Attr("data-slot", "data"), g.Raw(strings.ReplaceAll(string(data), "</", `<\/`))),
		h.Script(g.Raw(fahrplanJS)),
	)
}

func atoiPrefix(s string) int {
	n := 0
	for i := 0; i < len(s) && i < 4; i++ {
		if s[i] < '0' || s[i] > '9' {
			break
		}
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func fpDefaults(p *FahrplanProps) {
	set := func(s *string, v string) {
		if *s == "" {
			*s = v
		}
	}
	set(&p.ID, "fahrplan")
	set(&p.LaneWork, "Arbeit")
	set(&p.LaneWorkSub, "Work")
	set(&p.LaneEdu, "Ausbildung")
	set(&p.LaneEduSub, "Education")
	set(&p.KindWork, "Arbeit / Work")
	set(&p.KindEdu, "Ausbildung / Education")
	set(&p.KindThesis, "Abschlussarbeit / Thesis")
	set(&p.Today, "HEUTE")
	set(&p.Present, "heute")
	set(&p.Current, "● Aktuell")
	set(&p.Hint, "▸ Eintrag wählen für Details · ←/→ blättern")
	set(&p.Earlier, "◀ Früher")
	set(&p.Later, "Später ▶")
	set(&p.Entry, "Eintrag / Entry")
	set(&p.ShowDetails, "Details anzeigen")
	set(&p.To, "bis")
}

const fahrplanJS = `(function(){
  var root=document.currentScript.parentNode, D=JSON.parse(root.querySelector('[data-slot=data]').textContent), L=D.l;
  var plot=root.querySelector('[data-slot=plot]'), detail=root.querySelector('[data-slot=detail]');
  var mq=matchMedia('(max-width:1024px)');
  var nowY=+D.now.slice(0,4), START=D.start, YEARS=nowY-START+1, TOTAL=YEARS*12;
  function mi(s){var a=s.split('-');return (+a[0]-START)*12+((+a[1]||1)-1);}
  function fmt(s){var a=s.split('-');return a[1]?a[1]+'.'+a[0]:a[0];}
  function esc(s){return String(s==null?'':s).replace(/[&<>"]/g,function(c){return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c];});}
  var ev=D.items.map(function(t){var o={};for(var k in t) o[k]=t[k];
    o.s=mi(t.start); o.e=(t.end?mi(t.end):mi(D.now))+1; o.ongoing=!t.end; o.chip=t.short+' · '+t.start.slice(0,4); return o;});
  var order=ev.slice().sort(function(a,b){return a.s-b.s||a.e-b.e;});
  var LANES=[{key:'work',de:L.laneWork,en:L.laneWorkSub,kinds:['work']},{key:'edu',de:L.laneEdu,en:L.laneEduSub,kinds:['education','thesis']}];
  LANES.forEach(function(l){
    var tracks=[];
    ev.filter(function(x){return l.kinds.indexOf(x.kind)>=0;}).sort(function(a,b){return a.s-b.s;}).forEach(function(x){
      var t=-1; for(var i=0;i<tracks.length;i++){if(tracks[i]<=x.s){t=i;break;}}
      if(t<0){t=tracks.length;tracks.push(0);}
      tracks[t]=x.e; x.track=t; x.lane=l.key;
    });
    l.tracks=Math.max(1,tracks.length);
  });
  var selected=ev.filter(function(x){return x.ongoing;})[0]||order[order.length-1];
  var ctx=document.createElement('canvas').getContext('2d');
  function textW(t){ctx.font='700 11px "JetBrains Mono", monospace';return ctx.measureText(t).width;}
  function chipW(x){return Math.ceil(4+4+22+6+textW(x.chip)+6+12+6);}
  function initials(x){return ((x.org||x.title).match(/\b[A-ZÄÖÜ]/g)||['?']).slice(0,2).join('');}
  function logo(x,el){return x.logo?'<span data-el="'+el+'"><img src="'+esc(x.logo)+'" alt="" loading="lazy"></span>':'<span data-el="'+el+'" data-mono>'+esc(initials(x))+'</span>';}
  function range(x){return fmt(x.start)+' '+L.to+' '+(x.end?fmt(x.end):L.present);}
  function chip(x,l,t,w){
    return '<button type="button" data-el="chip" data-kind="'+x.kind+'" data-id="'+esc(x.id)+'" style="left:'+l+'px;top:'+t+'px;width:'+w+'px" aria-pressed="'+(x===selected)+'" aria-label="'+esc(x.title+', '+x.org+', '+range(x)+'. '+L.show)+'">'+logo(x,'logo')+'<span data-el="chip-text">'+esc(x.chip)+'</span><span data-el="chip-mark" aria-hidden="true"></span></button>';
  }
  function bar(x,style){return '<span data-el="bar" data-kind="'+x.kind+'" data-id="'+esc(x.id)+'"'+(x.ongoing?' data-running':'')+(x===selected?' data-sel':'')+' style="'+style+'" aria-hidden="true"></span>';}
  function line(x,d){return '<path data-id="'+esc(x.id)+'"'+(x===selected?' data-sel':'')+' d="'+d+'"/>';}
  function cols(){return 'grid-template-columns:repeat('+YEARS+',minmax(0,1fr))';}

  function renderH(){
    var W=plot.clientWidth, px=function(m){return m/TOTAL*W;}, CH=30, RG=6, TH=18, HEAD=30;
    var out='<div data-el="axis" style="'+cols()+'">';
    for(var i=0;i<YEARS;i++) out+='<span>'+(START+i)+'</span>';
    out+='</div><div data-el="body" style="background:repeating-linear-gradient(to right,var(--rule-subtle) 0 1px,transparent 1px calc(100%/'+YEARS+'))">';
    LANES.forEach(function(l){
      var E=ev.filter(function(x){return x.lane===l.key;}).sort(function(a,b){return a.s-b.s;}), rowRight=[];
      E.forEach(function(x){
        var w=chipW(x), cx=Math.max(0,Math.min(px(x.s),W-w)), r=-1;
        for(var i=0;i<rowRight.length;i++){if(rowRight[i]+8<=cx){r=i;break;}}
        if(r<0){r=rowRight.length;rowRight.push(0);}
        rowRight[r]=cx+w; x.cx=cx; x.cw=w; x.row=r;
      });
      var R=rowRight.length, barsTop=HEAD+R*(CH+RG)+14, laneH=barsTop+l.tracks*TH+14, chips='', bars='', lines='';
      E.forEach(function(x){
        var ct=barsTop-14-(x.row+1)*(CH+RG)+RG, bx=px(x.s), bw=Math.max(px(x.e-x.s),10), by=barsTop+x.track*TH;
        chips+=chip(x,x.cx,ct,x.cw);
        bars+=bar(x,'left:'+bx+'px;width:'+bw+'px;top:'+by+'px');
        var lx=Math.round(bx)+3.5; lines+=line(x,'M'+lx+' '+(ct+CH)+' V'+by);
      });
      out+='<div data-el="lane" style="height:'+laneH+'px"><span data-el="lane-label"><b>'+esc(l.de)+'</b>'+esc(l.en)+'</span><svg data-el="lines" aria-hidden="true">'+lines+'</svg>'+bars+chips+'</div>';
    });
    var pct=(mi(D.now)+.5)/TOTAL*100;
    out+='<div data-el="now" style="left:'+pct+'%"><span'+(pct>90?' style="left:auto;right:5px"':'')+'>'+esc(L.today)+'</span></div></div>';
    plot.innerHTML=out;
  }
  function renderV(){
    var W=plot.clientWidth, YH=96, H=YEARS*YH, y=function(m){return (TOTAL-m)/TOTAL*H;}, AX=40, TW=16, LG=10, CH=30;
    var x0=AX, laneX={};
    LANES.forEach(function(l){laneX[l.key]=x0; x0+=l.tracks*TW+LG;});
    var chipX=x0+14, chipMax=W-chipX, prev=-Infinity, placed=[];
    order.slice().sort(function(a,b){return y(a.e)-y(b.e);}).forEach(function(x){
      var cy=Math.max(y(x.e)-4,prev+6); prev=cy+CH; placed.push([x,cy]);
    });
    var bodyH=Math.max(H,prev+8), chips='', bars='', lines='';
    placed.forEach(function(pc){
      var x=pc[0], cy=pc[1], bx=laneX[x.lane]+x.track*TW+3, by=y(x.e), bh=Math.max(y(x.s)-y(x.e),10);
      bars+=bar(x,'left:'+bx+'px;top:'+by+'px;height:'+bh+'px');
      chips+=chip(x,chipX,cy,Math.min(chipMax,chipW(x)));
      var ly=by+5, mid=chipX-7; lines+=line(x,'M'+(bx+10)+' '+ly+' H'+mid+' V'+(cy+CH/2)+' H'+chipX);
    });
    var head='<div data-el="vhead"><span style="width:'+AX+'px"></span>'+LANES.map(function(l){return '<span style="width:'+(l.tracks*TW+LG)+'px">'+esc(l.de.slice(0,4))+'.</span>';}).join('')+'<span style="padding-left:14px">'+esc(L.entry)+'</span></div>';
    var years=''; for(var i=0;i<YEARS;i++) years+='<span data-el="year" style="top:'+(i*YH)+'px">'+(START+YEARS-1-i)+'</span>';
    plot.innerHTML=head+'<div data-el="body" style="height:'+bodyH+'px;--yh:'+YH+'px">'+years+'<svg data-el="lines" aria-hidden="true">'+lines+'</svg>'+bars+chips+'<div data-el="now" style="top:'+y(mi(D.now)+.5)+'px"></div></div>';
  }
  function byId(id){return ev.filter(function(e){return e.id===id;})[0];}
  function announce(){root.dispatchEvent(new CustomEvent('fahrplan:select',{bubbles:true,detail:{id:selected.id}}));}
  function select(x,focus){
    selected=x; render(); renderDetail(); announce();
    if(focus){var c=plot.querySelector('[data-el=chip][data-id="'+CSS.escape(x.id)+'"]'); if(c) c.focus({preventScroll:!mq.matches});}
  }
  function render(){
    root.dataset.orient=mq.matches?'v':'h';
    if(mq.matches) renderV(); else renderH();
    [].forEach.call(plot.querySelectorAll('[data-el=chip]'),function(c){
      var x=byId(c.dataset.id);
      c.onclick=function(){select(x,true);};
      c.onkeydown=function(e){var i=order.indexOf(x);
        if(e.key==='ArrowRight'||e.key==='ArrowDown'){e.preventDefault();select(order[Math.min(order.length-1,i+1)],true);}
        if(e.key==='ArrowLeft'||e.key==='ArrowUp'){e.preventDefault();select(order[Math.max(0,i-1)],true);}};
      function hot(on){[].forEach.call(plot.querySelectorAll('[data-id="'+CSS.escape(x.id)+'"]'),function(el){
        if(el.tagName==='path'){if(on||x===selected) el.setAttribute('data-sel',''); else el.removeAttribute('data-sel');}
        else if(el.dataset.el==='chip'){if(on) el.setAttribute('data-hot',''); else el.removeAttribute('data-hot');}
      });}
      c.onmouseenter=function(){hot(true);}; c.onmouseleave=function(){hot(false);};
    });
    [].forEach.call(plot.querySelectorAll('[data-el=bar]'),function(b){b.onclick=function(){select(byId(b.dataset.id),true);};});
  }
  function renderDetail(){
    var x=selected, i=order.indexOf(x), kind=L[x.kind]||x.kind;
    var det=(x.details||[]).length?'<ul data-el="d-details">'+x.details.map(function(d){return '<li>'+esc(d)+'</li>';}).join('')+'</ul>':'';
    detail.innerHTML='<div data-el="d-main"><div data-component="swiss-kicker"><span>'+esc(kind.split(' / ')[0])+' · '+esc(fmt(x.start))+' — '+esc(x.end?fmt(x.end):L.present)+'</span><span data-slot="aside">'+(x.ongoing?esc(L.current):'')+'</span></div>'+logo(x,'d-logo')+'<div><h3 data-el="d-title">'+esc(x.title)+'</h3><span data-component="swiss-micro" style="color:var(--muted)">'+esc(x.org)+(x.location?' · '+esc(x.location):'')+'</span></div></div>'+
      '<div data-el="d-side">'+(x.summary?'<p style="margin:0">'+esc(x.summary)+'</p>':'')+det+'<div data-component="swiss-pills">'+(x.tags||[]).map(function(t){return '<span data-component="swiss-pill">'+esc(t)+'</span>';}).join('')+'</div>'+
      '<div data-el="d-nav"><span data-component="swiss-micro" style="color:var(--muted)">'+(i+1)+' / '+order.length+'</span><button type="button" data-component="swiss-btn" data-act="prev"'+(i===0?' disabled':'')+'>'+esc(L.earlier)+'</button><button type="button" data-component="swiss-btn" data-act="next"'+(i===order.length-1?' disabled':'')+'>'+esc(L.later)+'</button></div></div>';
    detail.querySelector('[data-act=prev]').onclick=function(){select(order[i-1]);detail.querySelector('[data-act=prev]').focus();};
    detail.querySelector('[data-act=next]').onclick=function(){select(order[i+1]);detail.querySelector('[data-act=next]').focus();};
  }
  render(); renderDetail();
  if(document.readyState==='loading') document.addEventListener('DOMContentLoaded',announce); else setTimeout(announce,0);
  var rt; addEventListener('resize',function(){clearTimeout(rt);rt=setTimeout(render,200);});
  mq.addEventListener('change',render);
  if(document.fonts&&document.fonts.ready) document.fonts.ready.then(render);
})();`
