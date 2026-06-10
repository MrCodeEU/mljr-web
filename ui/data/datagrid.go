package data

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type DataGridColumn struct {
	Key      string // matches key in row map
	Label    string
	Sortable bool
	Width    string // optional CSS width
}

type DataGridProps struct {
	ID       string // unique per page (default "dg")
	Columns  []DataGridColumn
	Rows     []map[string]string
	Search   bool
	PageSize int // rows per page (0 = no pagination)
}

// DataGrid renders a sortable, filterable, paginated data table.
// All logic runs client-side via an inline script — no server roundtrip needed.
func DataGrid(p DataGridProps) g.Node {
	if p.ID == "" {
		p.ID = "dg"
	}
	if p.PageSize == 0 {
		p.PageSize = 10
	}

	// Header row
	ths := make([]g.Node, len(p.Columns))
	for i, col := range p.Columns {
		thAttrs := []g.Node{
			g.Attr("data-key", col.Key),
			h.Style("text-align:left;padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-2) solid var(--line);font-size:var(--t-xs);text-transform:uppercase;letter-spacing:.06em;white-space:nowrap" +
				func() string {
					if col.Width != "" {
						return ";width:" + col.Width
					}
					return ""
				}()),
		}
		if col.Sortable {
			thAttrs = append(thAttrs,
				g.Attr("data-sortable", ""),
				h.Style("cursor:pointer;user-select:none"),
			)
		}
		ths[i] = h.Th(append(thAttrs, g.Text(col.Label))...)
	}

	// Body rows
	trs := make([]g.Node, len(p.Rows))
	for i, row := range p.Rows {
		tds := make([]g.Node, len(p.Columns))
		for j, col := range p.Columns {
			tds[j] = h.Td(
				g.Attr("data-key", col.Key),
				h.Style("padding:var(--sp-2) var(--sp-3);border-bottom:var(--bw-1) solid var(--line);font-size:var(--t-sm)"),
				g.Text(row[col.Key]),
			)
		}
		trs[i] = h.Tr(tds...)
	}

	toolbar := []g.Node{}
	if p.Search {
		toolbar = append(toolbar,
			h.Input(
				g.Attr("data-slot", "search"),
				g.Attr("data-component", "input"),
				h.Type("search"),
				h.Placeholder("Filter…"),
				h.Style("width:220px"),
			),
		)
	}
	toolbar = append(toolbar,
		h.Span(g.Attr("data-slot", "info"), h.Style("font-size:var(--t-xs);color:var(--muted);margin-left:auto")),
		h.Button(g.Attr("data-slot", "prev"), g.Attr("data-component", "button"), g.Attr("data-variant", "outline"), g.Attr("data-size", "sm"), h.Type("button"), h.Disabled(), g.Text("←")),
		h.Button(g.Attr("data-slot", "next"), g.Attr("data-component", "button"), g.Attr("data-variant", "outline"), g.Attr("data-size", "sm"), h.Type("button"), g.Text("→")),
	)

	return h.Div(
		h.ID(p.ID),
		g.Attr("data-component", "data-grid"),
		h.Div(
			g.Attr("data-slot", "toolbar"),
			h.Style("display:flex;align-items:center;gap:var(--sp-2);margin-bottom:var(--sp-3)"),
			g.Group(toolbar),
		),
		h.Div(
			g.Attr("data-slot", "table-wrap"),
			h.Style("overflow-x:auto;border:var(--bw-2) solid var(--line);border-radius:var(--radius)"),
			h.Table(
				h.Style("width:100%;border-collapse:collapse"),
				h.THead(h.Tr(ths...)),
				h.TBody(trs...),
			),
		),
		h.Script(g.Raw(dataGridScript(p))),
	)
}

func dataGridScript(p DataGridProps) string {
	return fmt.Sprintf(`(function(){
  var root=document.getElementById('%s');
  if(!root) return;
  var tbody=root.querySelector('tbody');
  var allRows=Array.from(tbody.querySelectorAll('tr'));
  var search=root.querySelector('[data-slot="search"]');
  var info=root.querySelector('[data-slot="info"]');
  var prev=root.querySelector('[data-slot="prev"]');
  var next=root.querySelector('[data-slot="next"]');
  var state={q:'',sort:'',dir:'asc',page:1};
  var pageSize=%d;

  allRows.forEach(function(row){
    row._data={};
    Array.from(row.querySelectorAll('td[data-key]')).forEach(function(td){
      row._data[td.dataset.key]=td.textContent.trim();
    });
  });

  function render(){
    var q=state.q.toLowerCase();
    var filtered=allRows.filter(function(row){
      if(!q) return true;
      return Object.values(row._data).some(function(v){return v.toLowerCase().includes(q);});
    });
    if(state.sort){
      filtered.sort(function(a,b){
        var av=a._data[state.sort]||'';
        var bv=b._data[state.sort]||'';
        var n=Number(av),m=Number(bv);
        var cmp=(!isNaN(n)&&!isNaN(m))?(n-m):av.localeCompare(bv,undefined,{numeric:true});
        return state.dir==='asc'?cmp:-cmp;
      });
    }
    var total=filtered.length;
    var start=(state.page-1)*pageSize;
    var visible=new Set(filtered.slice(start,start+pageSize));
    allRows.forEach(function(row){ row.style.display=visible.has(row)?'':'none'; });
    if(info) info.textContent=(total?((start+1)+'-'+Math.min(start+pageSize,total)+' of '+total):'0 results');
    if(prev) prev.disabled=state.page<=1;
    if(next) next.disabled=(start+pageSize)>=total;
  }

  if(search){ search.addEventListener('input',function(){ state.q=this.value; state.page=1; render(); }); }
  if(prev)  prev.addEventListener('click',function(){ state.page--; render(); });
  if(next)  next.addEventListener('click',function(){ state.page++; render(); });

  root.querySelectorAll('th[data-sortable]').forEach(function(th){
    th.addEventListener('click',function(){
      var key=th.dataset.key;
      if(state.sort===key){ state.dir=state.dir==='asc'?'desc':'asc'; }
      else { state.sort=key; state.dir='asc'; }
      root.querySelectorAll('th[data-sortable]').forEach(function(t){ delete t.dataset.sorted; });
      th.dataset.sorted=state.dir;
      render();
    });
  });

  render();
})();`, p.ID, p.PageSize)
}
