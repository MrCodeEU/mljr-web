package layout

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TOCItem struct {
	ID    string // heading element id (href="#id")
	Label string
	Level int // 1=h1, 2=h2, 3=h3
}

type TOCProps struct {
	// Title shown above the list (default "Contents").
	Title string
	// Sticky makes the TOC position:sticky.
	Sticky bool
	// ContentSelector is the CSS selector to observe for headings (default "main").
	ContentSelector string
}

// TableOfContents renders a navigation list of heading links with scroll-spy.
// Items can be passed explicitly, or pass none to auto-generate via JS from headings.
func TableOfContents(p TOCProps, items ...TOCItem) g.Node {
	if p.Title == "" {
		p.Title = "Contents"
	}
	if p.ContentSelector == "" {
		p.ContentSelector = "main"
	}

	style := "display:flex;flex-direction:column;gap:var(--sp-1)"
	if p.Sticky {
		style += ";position:sticky;top:var(--sp-6)"
	}

	var itemNodes []g.Node
	for _, it := range items {
		depth := it.Level + 1 // Level 1=h2 → depth 2, matching heading tag depth in CSS
		itemNodes = append(itemNodes, h.Li(
			h.A(
				h.Href("#"+it.ID),
				g.Attr("data-slot", "item"),
				g.Attr("data-depth", fmt.Sprintf("%d", depth)),
				g.Attr("data-toc-href", it.ID),
				g.Text(it.Label),
			),
		))
	}

	// JS: IntersectionObserver scroll-spy + auto-populate if no items passed
	script := `(function(){
  document.querySelectorAll('[data-component="toc"]:not([data-toc-init])').forEach(function(toc){
    toc.setAttribute('data-toc-init','1');
    var sel=toc.dataset.content||'main';
    var list=toc.querySelector('[data-slot="list"]');

    // Auto-populate from DOM headings if list is empty
    if(!list.children.length){
      var container=document.querySelector(sel)||document;
      var hs=Array.from(container.querySelectorAll('h2,h3,h4'));
      hs.forEach(function(heading,i){
        if(!heading.id) heading.id='toc-h-'+i;
        var li=document.createElement('li');
        var a=document.createElement('a');
        a.href='#'+heading.id;
        a.setAttribute('data-slot','item');
        a.setAttribute('data-depth',heading.tagName.charAt(1));
        a.setAttribute('data-toc-href',heading.id);
        a.textContent=heading.textContent;
        li.appendChild(a);
        list.appendChild(li);
      });
    }

    // Scroll-spy
    var links=Array.from(toc.querySelectorAll('[data-toc-href]'));
    if(!links.length) return;
    var targets=links.map(function(l){ return document.getElementById(l.dataset.tocHref); }).filter(Boolean);
    var observer=new IntersectionObserver(function(entries){
      entries.forEach(function(e){
        if(e.isIntersecting){
          links.forEach(function(l){ l.removeAttribute('data-state'); });
          var active=toc.querySelector('[data-toc-href="'+e.target.id+'"]');
          if(active) active.setAttribute('data-state','active');
        }
      });
    },{rootMargin:'-10% 0px -80% 0px'});
    targets.forEach(function(t){ observer.observe(t); });
  });
})();`

	return h.Nav(
		g.Attr("data-component", "toc"),
		g.Attr("data-content", p.ContentSelector),
		h.Style(style),
		h.Div(g.Attr("data-slot", "title"), g.Text(p.Title)),
		h.Ul(g.Attr("data-slot", "list"), h.Style("list-style:none;margin:0;padding:0"), g.Group(itemNodes)),
		h.Script(g.Raw(script)),
	)
}
