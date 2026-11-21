package swiss

import (
	"encoding/json"
	"strconv"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	"mljr-web/ui/icon"
)

// FlapRow is one line of the departures board.
type FlapRow struct {
	Icon  string // ui/icon name, e.g. "simple-icons:go"
	Code  string // 3 chars; derived from Name when empty
	Name  string
	Gate  string // 3-char group short, e.g. "SPR"
	Value string // last column, up to 5 chars (e.g. a count)
	Group string // full group name (screen-reader table)
}

type FlapBoardProps struct {
	ID       string
	Title    string // "Abflug"
	TitleSub string // "DEPARTURES"
	Rows     []FlapRow
	PerPage  int // default 8
	// IntervalMS between automatic page flips; default 15000.
	IntervalMS int
	// Column headings, page label ("Seite"), control labels, sr table caption.
	ColCode, ColName, ColGate, ColValue  string
	PageLabel, Pause, Resume, Prev, Next string
	Caption, ColGroup                    string
}

// FlapBoard renders an airport split-flap display that pages through rows.
// Tiles flip sequentially through the character drum; icon tiles spin
// through random icons first. The board itself is aria-hidden: a visually
// hidden table carries the content, and pause/prev/next satisfy WCAG 2.2.2.
// Under prefers-reduced-motion nothing flips or auto-advances.
func FlapBoard(p FlapBoardProps) g.Node {
	fbDefaults(&p)
	type jsRow struct {
		Icon, Code, Name, Gate, Lvl string
	}
	rows := make([]jsRow, len(p.Rows))
	icons := map[string]string{}
	trs := make([]g.Node, len(p.Rows))
	for i, r := range p.Rows {
		code := r.Code
		if code == "" {
			code = FlapCode(r.Name)
		}
		rows[i] = jsRow{Icon: r.Icon, Code: code, Name: r.Name, Gate: r.Gate, Lvl: r.Value}
		if svg, ok := icon.SVG(r.Icon); ok {
			icons[r.Icon] = svg
		}
		trs[i] = h.Tr(h.Td(g.Text(r.Group)), h.Td(g.Text(r.Name)), h.Td(g.Text(r.Value)))
	}
	data, err := json.Marshal(map[string]any{
		"rows": rows, "icons": icons, "per": p.PerPage, "interval": p.IntervalMS,
		"cols": [][]any{{"icon", 1, ""}, {"code", 3, p.ColCode}, {"name", 20, p.ColName}, {"gate", 3, p.ColGate}, {"lvl", 5, p.ColValue}},
		"page": p.PageLabel, "pause": p.Pause, "resume": p.Resume,
	})
	if err != nil {
		return g.Raw("<!-- swiss: " + err.Error() + " -->")
	}
	pages := (len(p.Rows) + p.PerPage - 1) / p.PerPage
	return h.Div(
		g.Attr("data-component", "flapboard"), h.ID(p.ID),
		h.Div(g.Attr("data-slot", "board"), g.Attr("aria-hidden", "true"),
			h.Div(g.Attr("data-slot", "head"),
				h.Span(g.Attr("data-slot", "title"), g.Text(p.Title), h.Small(g.Text(p.TitleSub))),
				h.Span(g.Attr("data-slot", "clock"), g.Text("--:--")),
			),
			h.Div(g.Attr("data-slot", "rows")),
		),
		h.Div(g.Attr("data-slot", "controls"),
			Micro(g.Attr("data-slot", "page"), g.Text(p.PageLabel+" 1/"+strconv.Itoa(pages))),
			h.Span(
				Btn(false, []g.Node{g.Attr("data-act", "prev"), g.Attr("aria-label", p.Prev)}, g.Text("◀")),
				Btn(false, []g.Node{g.Attr("data-act", "pause"), g.Attr("aria-pressed", "false")}, g.Text(p.Pause)),
				Btn(false, []g.Node{g.Attr("data-act", "next"), g.Attr("aria-label", p.Next)}, g.Text("▶")),
			),
		),
		h.Table(g.Attr("data-component", "swiss-sr"),
			h.Caption(g.Text(p.Caption)),
			h.THead(h.Tr(h.Th(g.Text(p.ColGroup)), h.Th(g.Text(p.ColName)), h.Th(g.Text(p.ColValue)))),
			h.TBody(trs...),
		),
		h.Script(h.Type("application/json"), g.Attr("data-slot", "data"), g.Raw(strings.ReplaceAll(string(data), "</", `<\/`))),
		h.Script(g.Raw(flapJS)),
	)
}

// FlapCode derives a 3-char code: first letter + following consonants.
func FlapCode(name string) string {
	var u []rune
	for _, r := range strings.ToUpper(name) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '#' || r == '.' || r == '+' {
			u = append(u, r)
		}
	}
	if len(u) <= 3 {
		return string(u)
	}
	out := []rune{u[0]}
	for _, r := range u[1:] {
		if !strings.ContainsRune("AEIOU", r) {
			out = append(out, r)
		}
	}
	out = append(out, u[1:]...)
	return string(out[:3])
}

func fbDefaults(p *FlapBoardProps) {
	set := func(s *string, v string) {
		if *s == "" {
			*s = v
		}
	}
	set(&p.ID, "flapboard")
	set(&p.Title, "Abflug")
	set(&p.TitleSub, "DEPARTURES")
	set(&p.ColCode, "Code")
	set(&p.ColName, "Kompetenz")
	set(&p.ColGate, "Gleis")
	set(&p.ColValue, "Wert")
	set(&p.PageLabel, "Seite")
	set(&p.Pause, "Anhalten")
	set(&p.Resume, "Weiter")
	set(&p.Prev, "Vorherige Seite")
	set(&p.Next, "Nächste Seite")
	set(&p.Caption, "Kompetenzen nach Bereich")
	set(&p.ColGroup, "Bereich")
	if p.PerPage <= 0 {
		p.PerPage = 8
	}
	if p.IntervalMS <= 0 {
		p.IntervalMS = 15000
	}
}

const flapJS = `(function(){
  var root=document.currentScript.parentNode, D=JSON.parse(root.querySelector('[data-slot=data]').textContent);
  var board=root.querySelector('[data-slot=board]'), rowsEl=root.querySelector('[data-slot=rows]');
  var reduced=matchMedia('(prefers-reduced-motion: reduce)').matches, mq=matchMedia('(max-width:1024px)');
  var CHARS=' ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#.+/-■□', COLS=D.cols, ROWS=D.per, SK=D.rows;
  var PAGES=Math.ceil(SK.length/ROWS), ICONS=D.icons, KEYS=Object.keys(ICONS);
  function unit(c){return c[0]==='icon'?1.6:1.2;}
  function esc(s){return String(s).replace(/[&<>"]/g,function(c){return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c];});}
  var html='<div data-el="row" data-head>'+COLS.map(function(c){return '<span data-el="col" data-key="'+c[0]+'" style="width:calc('+(c[1]*unit(c))+'em + '+(c[1]-1)*2+'px)"><span>'+esc(c[2])+'</span></span>';}).join('')+'</div>';
  for(var r=0;r<ROWS;r++) html+='<div data-el="row">'+COLS.map(function(c){return '<span data-el="col" data-key="'+c[0]+'">'+(c[0]==='icon'?'<span data-el="tile" data-icon></span>':new Array(c[1]+1).join('<span data-el="tile"> </span>'))+'</span>';}).join('')+'</div>';
  rowsEl.innerHTML=html;
  var tiles=[], iconTiles=[];
  [].forEach.call(rowsEl.querySelectorAll('[data-el=row]:not([data-head])'),function(row,r){
    iconTiles.push({el:row.querySelector('[data-icon]'),cur:'',target:'',next:0,spins:0,alt:false});
    var c=0; [].forEach.call(row.querySelectorAll('[data-el=tile]:not([data-icon])'),function(el){tiles.push({el:el,r:r,c:c++,cur:' ',target:' ',next:0,alt:false});});
  });
  var page=0, paused=reduced, running=false, timer=null, pageEl=root.querySelector('[data-slot=page]');
  function setPage(p){
    page=(p+PAGES)%PAGES;
    var slice=SK.slice(page*ROWS,page*ROWS+ROWS), now=performance.now(), k=0;
    for(var r=0;r<ROWS;r++){
      var s=slice[r]||{}, it=iconTiles[r];
      it.target=s.Icon||''; it.next=now+r*70; it.spins=reduced||!s.Icon?0:5;
      COLS.forEach(function(c){
        if(c[0]==='icon') return;
        var key={code:'Code',name:'Name',gate:'Gate',lvl:'Lvl'}[c[0]];
        var txt=String(s[key]||'').toUpperCase(); while(txt.length<c[1]) txt+=' '; txt=txt.slice(0,c[1]);
        for(var i=0;i<txt.length;i++){var t=tiles[k++]; t.target=txt[i]; t.next=now+r*70+t.c*12;}
      });
    }
    pageEl.textContent=D.page+' '+(page+1)+'/'+PAGES;
    if(!running){running=true;requestAnimationFrame(step);}
  }
  function step(now){
    var busy=false;
    tiles.forEach(function(t){
      if(t.cur===t.target) return; busy=true;
      if(now<t.next) return;
      var ti=CHARS.indexOf(t.target);
      t.cur=(reduced||ti<0)?t.target:CHARS[(CHARS.indexOf(t.cur)+1)%CHARS.length];
      t.el.textContent=t.cur;
      if(!reduced){t.alt=!t.alt;t.el.dataset.flap=t.alt?'a':'b';}
      t.next=now+26;
    });
    iconTiles.forEach(function(it){
      if(it.cur===it.target&&!it.spins) return; busy=true;
      if(now<it.next) return;
      var key=it.spins>0?KEYS[Math.random()*KEYS.length|0]:it.target;
      if(it.spins>0) it.spins--; else it.cur=it.target;
      it.el.innerHTML=ICONS[key]||'';
      if(!reduced){it.alt=!it.alt;it.el.dataset.flap=it.alt?'a':'b';}
      it.next=now+90;
    });
    if(busy) requestAnimationFrame(step); else running=false;
  }
  // Size tiles so the board fills its box: W = tf*(sum units + 0.7(C-1)) + 2px*(N-C)
  function fit(){
    var cs=getComputedStyle(board), W=board.clientWidth-parseFloat(cs.paddingLeft)-parseFloat(cs.paddingRight);
    var cols=COLS.filter(function(c){return !(mq.matches&&c[0]==='gate');});
    var N=cols.reduce(function(a,c){return a+c[1];},0), C=cols.length, em=cols.reduce(function(a,c){return a+c[1]*unit(c);},0);
    board.style.setProperty('--tf',Math.max(8,Math.min(26,(W-2*(N-C)-2)/(em+0.7*(C-1))))+'px');
  }
  fit();
  var rt; addEventListener('resize',function(){clearTimeout(rt);rt=setTimeout(fit,150);});
  function schedule(){clearInterval(timer);if(!paused) timer=setInterval(function(){setPage(page+1);},D.interval);}
  var pb=root.querySelector('[data-act=pause]');
  root.querySelector('[data-act=next]').onclick=function(){setPage(page+1);schedule();};
  root.querySelector('[data-act=prev]').onclick=function(){setPage(page-1);schedule();};
  pb.onclick=function(){paused=!paused;pb.textContent=paused?D.resume:D.pause;pb.setAttribute('aria-pressed',paused);schedule();};
  if(reduced){pb.textContent=D.resume;pb.setAttribute('aria-pressed','true');}
  new IntersectionObserver(function(es,o){if(es[0].isIntersecting){o.disconnect();setPage(0);schedule();}},{threshold:.3}).observe(board);
  var clk=root.querySelector('[data-slot=clock]');
  function clock(){var d=new Date();clk.textContent=String(d.getHours()).padStart(2,'0')+':'+String(d.getMinutes()).padStart(2,'0');}
  clock(); setInterval(clock,15000);
})();`
