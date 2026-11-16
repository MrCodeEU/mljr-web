package pages

import (
	"fmt"
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// tabsJS drives every [data-tabs] root on the page: a click on a
// [data-tab="n"] button shows the matching [data-tab-page="n"] and hides
// the rest; a [data-tab-go="n"] pager link does the same and scrolls the
// root back into view if its top has left the screen. Arrow keys move
// between the buttons. Include once per page.
const tabsJS = `(function(){
  function show(root,b){
    var n=b.getAttribute('data-tab');
    root.dispatchEvent(new CustomEvent('tabs:select',{bubbles:true,detail:{id:n}}));
    root.querySelectorAll('[data-tab]').forEach(function(t){t.setAttribute('aria-pressed',String(t.getAttribute('data-tab')===n));});
    root.querySelectorAll('[data-tab-page]').forEach(function(s){
      s.hidden=s.getAttribute('data-tab-page')!==n;
      // embedded iframes run only while their page is visible
      s.querySelectorAll('iframe[data-src]').forEach(function(f){
        if(s.hidden){ if(f.getAttribute('src')) f.removeAttribute('src'); }
        else if(!f.getAttribute('src')) f.setAttribute('src',f.getAttribute('data-src'));
      });
    });
  }
  document.addEventListener('click',function(e){
    var go=e.target.closest('[data-tab-go]');
    if(go){
      var r=go.closest('[data-tabs]'), t=r&&r.querySelector('[data-tab="'+go.getAttribute('data-tab-go')+'"]');
      if(t){ show(r,t); if(r.getBoundingClientRect().top<0) r.scrollIntoView({behavior:'smooth',block:'start'}); }
      return;
    }
    var b=e.target.closest('[data-tab]'); if(!b) return;
    var root=b.closest('[data-tabs]'); if(root) show(root,b);
  });
  document.addEventListener('keydown',function(e){
    if(e.key!=='ArrowRight'&&e.key!=='ArrowLeft') return;
    var b=e.target.closest&&e.target.closest('[data-tab]'); if(!b) return;
    var root=b.closest('[data-tabs]'), all=[].slice.call(root.querySelectorAll('[data-tab]'));
    var i=all.indexOf(b)+(e.key==='ArrowRight'?1:-1); if(i<0||i>=all.length) return;
    e.preventDefault(); all[i].focus(); show(root,all[i]);
  });
})();`

// tabPager renders the previous/next bar under a tabbed page. Titles name
// the neighbours so the reader sees there is more, and what.
func tabPager(lang string, i, n int, title func(int) string, unit string) g.Node {
	var prev, next g.Node
	if i > 0 {
		prev = h.Button(h.Type("button"), h.Class("pager-btn pager-prev"), g.Attr("data-tab-go", strconv.Itoa(i)),
			h.Span(h.Class("pager-dir"), g.Text("◀ "+tr(lang, "Vorherige", "Previous"))),
			h.Span(h.Class("pager-t"), g.Text(title(i-1))),
		)
	}
	if i < n-1 {
		next = h.Button(h.Type("button"), h.Class("pager-btn pager-next"), g.Attr("data-tab-go", strconv.Itoa(i+2)),
			h.Span(h.Class("pager-dir"), g.Text(tr(lang, "Nächste", "Next")+" ▶")),
			h.Span(h.Class("pager-t"), g.Text(title(i+1))),
		)
	} else {
		next = h.Button(h.Type("button"), h.Class("pager-btn pager-next"), g.Attr("data-tab-go", "1"),
			h.Span(h.Class("pager-dir"), g.Text(tr(lang, "Von vorn", "Start over")+" ↺")),
			h.Span(h.Class("pager-t"), g.Text(title(0))),
		)
	}
	return h.Nav(h.Class("pager"), g.Attr("aria-label", unit),
		g.If(prev != nil, prev), g.If(prev == nil, h.Span()),
		h.Span(h.Class("pager-pos"), g.Text(fmt.Sprintf("%s %d / %d", unit, i+1, n))),
		next,
	)
}

// tabLead is the label cell that opens a tab strip.
func tabLead(label string, n int) g.Node {
	return h.Span(h.Class("tab-lead"), h.Span(g.Text(label)), h.B(g.Text(fmt.Sprintf("%d ▸", n))))
}

const tabsCSS = `
.tab-lead { display: flex; flex-direction: column; justify-content: space-between; gap: .3rem; padding: .7rem var(--swiss-pad); background: var(--board); color: #f3f1ea; font-family: var(--font-mono); font-size: .66rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; }
.tab-lead b { font-size: 1rem; letter-spacing: 0; color: var(--act); }
.pager { display: grid; grid-template-columns: 1fr auto 1fr; align-items: stretch; border-top: 3px solid var(--ink); }
.pager-btn { all: unset; box-sizing: border-box; cursor: pointer; display: flex; flex-direction: column; gap: .25rem; padding: .8rem var(--swiss-pad); min-width: 0; }
.pager-next { text-align: right; align-items: flex-end; box-shadow: inset 1px 0 0 var(--ink); }
.pager-prev { box-shadow: inset -1px 0 0 var(--ink); }
.pager-dir { font-family: var(--font-mono); font-size: .68rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; color: var(--main-text); }
.pager-t { font-weight: 800; font-size: 1.05rem; letter-spacing: -0.02em; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 100%; }
.pager-dir::before { content: ""; display: inline-block; width: .6rem; height: .6rem; margin-right: .4rem; background: var(--act); box-shadow: 0 0 0 1px var(--act-fg); vertical-align: 0; }
.pager-next .pager-dir::before { display: none; }
.pager-next .pager-dir::after { content: ""; display: inline-block; width: .6rem; height: .6rem; margin-left: .4rem; background: var(--act); box-shadow: 0 0 0 1px var(--act-fg); }
.pager-btn:hover { background: var(--act); color: var(--act-fg); }
.pager-btn:hover .pager-dir { color: inherit; }
.pager-btn:focus-visible { outline: 3px solid var(--accent); outline-offset: -3px; }
.pager-pos { align-self: center; padding: 0 1rem; font-family: var(--font-mono); font-size: .7rem; font-weight: 700; text-transform: uppercase; letter-spacing: .08em; color: var(--muted); }
@media (max-width: 640px) { .pager-pos { display: none; } .pager { grid-template-columns: 1fr 1fr; } .pager-t { font-size: .9rem; } }
`
