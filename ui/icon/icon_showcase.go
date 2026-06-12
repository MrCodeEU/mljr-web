//go:build showcase

package icon

import (
	"mljr-web/ui/registry"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "icon", Name: "Icon", Category: "primitive",
		PreviewHeight: "640px",
		Summary:       "Iconify-generated SVGs rendered inline. All registered icons shown below — searchable and copy-on-click.",
		Code: `// import "mljr-web/ui/icon"
icon.Icon("lucide:check")         // inline SVG, currentColor
icon.Icon("simple-icons:github")  // brand icon
icon.Icon("lucide:star", icon.Props{Size: "2rem"})

// Check availability at runtime
icon.Has("lucide:foo") // bool
icon.All()             // []string — all registered names`,
		Render: func(p map[string]string) g.Node {
			all := All()

			// Group by set prefix
			sets := map[string][]string{}
			setOrder := []string{}
			for _, name := range all {
				parts := strings.SplitN(name, ":", 2)
				set := parts[0]
				if _, seen := sets[set]; !seen {
					setOrder = append(setOrder, set)
				}
				sets[set] = append(sets[set], name)
			}

			sections := make([]g.Node, 0, len(setOrder))
			for _, set := range setOrder {
				names := sets[set]
				cells := make([]g.Node, len(names))
				for i, name := range names {
					short := strings.SplitN(name, ":", 2)[1]
					cells[i] = h.Div(
						g.Attr("data-slot", "ic-cell"),
						g.Attr("data-name", name),
						h.Style("display:flex;flex-direction:column;align-items:center;gap:4px;padding:10px 6px;border:1px solid transparent;border-radius:var(--radius);cursor:pointer;transition:border-color .15s,background .15s"),
						h.Title(name),
						Icon(name, Props{Size: "1.5rem"}),
						h.Span(
							h.Style("font-size:9px;text-align:center;color:var(--muted);max-width:72px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;line-height:1.2"),
							g.Text(short),
						),
					)
				}
				sections = append(sections,
					h.Div(
						g.Attr("data-slot", "ic-section"),
						g.Attr("data-set", set),
						h.H3(
							h.Style("font-size:var(--t-xs);text-transform:uppercase;letter-spacing:.08em;font-weight:700;opacity:.45;margin:var(--sp-4) 0 var(--sp-2)"),
							g.Text(set),
						),
						h.Div(
							g.Attr("data-slot", "ic-grid"),
							h.Style("display:flex;flex-wrap:wrap;gap:4px"),
							g.Group(cells),
						),
					),
				)
			}

			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-2)"),

				// Search + count bar
				h.Div(
					h.Style("display:flex;align-items:center;gap:var(--sp-3);position:sticky;top:0;background:var(--bg);padding:var(--sp-3) 0;z-index:1"),
					h.Input(
						h.ID("ic-search"),
						h.Type("search"),
						g.Attr("data-component", "input"),
						h.Placeholder("Search icons…"),
						h.Style("flex:1"),
					),
					h.Span(
						h.ID("ic-count"),
						h.Style("font-size:var(--t-xs);opacity:.5;white-space:nowrap"),
						g.Text(strings.Repeat("", 0)+itoa(len(all))+" icons"),
					),
				),

				// Copy feedback toast
				h.Div(
					h.ID("ic-toast"),
					h.Style("display:none;position:fixed;bottom:var(--sp-5);left:50%;transform:translateX(-50%);background:var(--fg);color:var(--bg);padding:var(--sp-2) var(--sp-4);border-radius:var(--radius);font-size:var(--t-sm);font-weight:700;z-index:100;white-space:nowrap"),
				),

				// Icon grid sections
				h.Div(h.ID("ic-sections"), g.Group(sections)),

				h.Script(g.Raw(`(function(){
  var search=document.getElementById('ic-search');
  var countEl=document.getElementById('ic-count');
  var toast=document.getElementById('ic-toast');
  var toastTimer;

  // Copy on click
  document.querySelectorAll('[data-slot="ic-cell"]').forEach(function(cell){
    cell.addEventListener('mouseenter',function(){
      cell.style.borderColor='var(--line)';
      cell.style.background='var(--surface-2)';
    });
    cell.addEventListener('mouseleave',function(){
      cell.style.borderColor='transparent';
      cell.style.background='';
    });
    cell.addEventListener('click',function(){
      var name=cell.dataset.name;
      navigator.clipboard.writeText('icon.Icon("'+name+'")').catch(function(){});
      clearTimeout(toastTimer);
      toast.textContent='Copied: icon.Icon("'+name+'")';
      toast.style.display='block';
      toastTimer=setTimeout(function(){toast.style.display='none';},2000);
    });
  });

  // Search filter
  function filter(){
    var q=search.value.toLowerCase().trim();
    var visible=0;
    document.querySelectorAll('[data-slot="ic-cell"]').forEach(function(cell){
      var match=!q||cell.dataset.name.includes(q);
      cell.style.display=match?'':'none';
      if(match) visible++;
    });
    // Hide section headers with no visible children
    document.querySelectorAll('[data-slot="ic-section"]').forEach(function(sec){
      var any=Array.from(sec.querySelectorAll('[data-slot="ic-cell"]')).some(function(c){return c.style.display!=='none';});
      sec.style.display=any?'':'none';
    });
    countEl.textContent=visible+' icons';
  }
  search.addEventListener('input',filter);
})()`)),
			)
		},
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
