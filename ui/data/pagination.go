package data

import (
	"fmt"
	"math"

	"mljr-web/ui"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PaginationProps struct {
	ID      string
	Total   int
	PerPage int
	Attrs   []g.Node
}

func PaginationSignals(id string, perPage int) g.Node {
	if id == "" {
		id = "pg"
	}
	sig := id + "Page"
	return h.Div(ui.Signals(fmt.Sprintf(`{%s:0}`, sig)))
}

func Pagination(p PaginationProps) g.Node {
	id := p.ID
	if id == "" {
		id = "pg"
	}
	perPage := p.PerPage
	if perPage <= 0 {
		perPage = 6
	}
	sig := id + "Page"

	pages := int(math.Ceil(float64(p.Total) / float64(perPage)))
	if pages < 1 {
		pages = 1
	}
	maxPage := pages - 1

	prevExpr := fmt.Sprintf("$%s = Math.max(0,$%s-1)", sig, sig)
	nextExpr := fmt.Sprintf("$%s = Math.min(%d,$%s+1)", sig, maxPage, sig)

	// prev: disabled when on first page
	prevDisabled := fmt.Sprintf(`{"data-state": $%s === 0 ? "disabled" : ""}`, sig)
	// next: disabled when on last page
	nextDisabled := fmt.Sprintf(`{"data-state": $%s === %d ? "disabled" : ""}`, sig, maxPage)

	btns := make([]g.Node, 0, pages+2)
	btns = append(btns,
		h.Button(
			g.Attr("data-slot", "prev"),
			g.Attr("data-attr", prevDisabled),
			ui.On("click", fmt.Sprintf("if($%s>0){%s}", sig, prevExpr)),
			g.Text("←"),
		),
	)
	for i := range pages {
		idx := i
		// active state driven by Datastar — data-attr sets data-state="active" when current page
		activeAttr := fmt.Sprintf(`{"data-state": $%s === %d ? "active" : ""}`, sig, idx)
		btns = append(btns,
			h.Button(
				g.Attr("data-slot", "btn"),
				g.Attr("data-attr", activeAttr),
				ui.On("click", fmt.Sprintf("$%s = %d", sig, idx)),
				g.Text(fmt.Sprintf("%d", idx+1)),
			),
		)
	}
	btns = append(btns,
		h.Button(
			g.Attr("data-slot", "next"),
			g.Attr("data-attr", nextDisabled),
			ui.On("click", fmt.Sprintf("if($%s<%d){%s}", sig, maxPage, nextExpr)),
			g.Text("→"),
		),
	)

	return h.Div(
		g.Attr("data-component", "pagination"),
		g.Group(p.Attrs),
		g.Group(btns),
	)
}

// PageAnimation selects the entrance animation PaginatedPages plays when a
// page becomes visible.
type PageAnimation string

const (
	PageAnimSlideUp   PageAnimation = "slide-up"
	PageAnimSlideLeft PageAnimation = "slide-left"
	PageAnimFade      PageAnimation = "fade"
	PageAnimScale     PageAnimation = "scale"
	PageAnimFlip      PageAnimation = "flip"
	PageAnimNone      PageAnimation = "none"
)

type PaginatedPagesProps struct {
	// ID is the signal prefix shared with Pagination/PaginationSignals
	// (default "pg"): page i is shown while $<ID>Page === i.
	ID string
	// Animation played when a page becomes visible (default PageAnimSlideUp).
	Animation PageAnimation
	Attrs     []g.Node
}

// PaginatedPages wraps one node per page and shows only the active one,
// driven by the same $<ID>Page signal Pagination writes. On every
// hidden→visible transition the entrance animation replays exactly once:
// the observer reacts only to style mutations whose old value was
// display:none, so the CSS animation itself can never re-trigger it.
func PaginatedPages(p PaginatedPagesProps, pages ...g.Node) g.Node {
	id := p.ID
	if id == "" {
		id = "pg"
	}
	anim := p.Animation
	if anim == "" {
		anim = PageAnimSlideUp
	}
	sig := id + "Page"
	containerID := "pp-" + id

	wrapped := make([]g.Node, len(pages))
	for i, page := range pages {
		wrapped[i] = h.Div(
			g.Attr("data-slot", "page"),
			ui.Show(fmt.Sprintf("$%s === %d", sig, i)),
			page,
		)
	}

	var script g.Node
	if anim != PageAnimNone {
		script = h.Script(g.Raw(fmt.Sprintf(`(function(){
  var c=document.getElementById('%s');
  if(!c)return;
  var obs=new MutationObserver(function(muts){
    muts.forEach(function(m){
      var el=m.target;
      var was=(m.oldValue||'');
      var wasHidden=was.indexOf('display:none')>-1||was.indexOf('display: none')>-1;
      if(wasHidden&&el.style.display!=='none'){
        el.removeAttribute('data-anim');
        void el.offsetWidth;
        el.setAttribute('data-anim','');
      }
    });
  });
  Array.prototype.forEach.call(c.children,function(el){
    obs.observe(el,{attributes:true,attributeFilter:['style'],attributeOldValue:true});
  });
})();`, containerID)))
	}

	return h.Div(
		h.ID(containerID),
		g.Attr("data-component", "paginated-pages"),
		g.Attr("data-animation", string(anim)),
		g.Group(p.Attrs),
		g.Group(wrapped),
		script,
	)
}
