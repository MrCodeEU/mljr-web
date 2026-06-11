package form

import (
	"fmt"
	"time"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CalendarPickerProps struct {
	// Name is the hidden input name for form submission.
	Name string
	// Value is the pre-selected date (ISO 8601: "2026-06-10").
	Value string
	// Min/Max restrict selectable dates (ISO 8601).
	Min string
	Max string
	// Signal is the Datastar signal prefix (default "_cal"). Must be unique per page.
	Signal string
	// Placeholder shown in the trigger button when no date is selected.
	Placeholder string
}

// CalendarPicker renders a custom styled date picker with a month/year grid.
// No external library. JS handles navigation; selected value goes in a hidden input.
func CalendarPicker(p CalendarPickerProps) g.Node {
	if p.Signal == "" {
		p.Signal = "_cal"
	}
	if p.Placeholder == "" {
		p.Placeholder = "Pick a date…"
	}

	sig := p.Signal
	now := time.Now()

	// Initial display month/year
	initYear := now.Year()
	initMonth := int(now.Month())
	if p.Value != "" {
		if t, err := time.Parse("2006-01-02", p.Value); err == nil {
			initYear = t.Year()
			initMonth = int(t.Month())
		}
	}

	script := fmt.Sprintf(`(function(){
  var root=document.querySelector('[data-cal-sig="%s"]');
  if(!root||root._calInit) return;
  root._calInit=true;

  var input=root.querySelector('input[type="hidden"]');
  var trigger=root.querySelector('[data-slot="trigger"]');
  var dropdown=root.querySelector('[data-slot="dropdown"]');
  var grid=root.querySelector('[data-slot="grid"]');
  var monthLabel=root.querySelector('[data-slot="month-label"]');
  var prevBtn=root.querySelector('[data-slot="prev"]');
  var nextBtn=root.querySelector('[data-slot="next"]');

  var year=%d, month=%d; // 1-indexed month
  var selected=input.value||'';
  var min=%s, max=%s;

  var monthNames=['January','February','March','April','May','June','July','August','September','October','November','December'];
  var dayNames=['Su','Mo','Tu','We','Th','Fr','Sa'];

  function fmt2(n){ return n<10?'0'+n:''+n; }
  function dateStr(y,m,d){ return y+'-'+fmt2(m)+'-'+fmt2(d); }

  function render(){
    monthLabel.textContent=monthNames[month-1]+' '+year;
    var firstDay=new Date(year,month-1,1).getDay();
    var daysInMonth=new Date(year,month,0).getDate();
    var html='<div style="display:grid;grid-template-columns:repeat(7,1fr);gap:2px;margin-bottom:4px">';
    dayNames.forEach(function(d){ html+='<div style="text-align:center;font-size:11px;font-weight:700;opacity:.5;padding:4px">'+d+'</div>'; });
    html+='</div><div style="display:grid;grid-template-columns:repeat(7,1fr);gap:2px">';
    for(var i=0;i<firstDay;i++) html+='<div></div>';
    for(var d=1;d<=daysInMonth;d++){
      var ds=dateStr(year,month,d);
      var isSelected=ds===selected;
      var isToday=ds===dateStr(new Date().getFullYear(),new Date().getMonth()+1,new Date().getDate());
      var disabled=(min&&ds<min)||(max&&ds>max);
      var style='text-align:center;padding:6px 2px;border-radius:4px;font-size:13px;font-weight:600;cursor:'+(disabled?'not-allowed':'pointer')+';';
      if(isSelected) style+='background:var(--ink);color:var(--surface);';
      else if(isToday) style+='outline:2px solid var(--accent);';
      else if(!disabled) style+='color:var(--ink);';
      else style+='opacity:.35;';
      html+='<div data-d="'+ds+'" style="'+style+'"'+(disabled?' data-disabled':'')+'>'+d+'</div>';
    }
    html+='</div>';
    grid.innerHTML=html;
    grid.querySelectorAll('[data-d]').forEach(function(el){
      if(!el.hasAttribute('data-disabled')) el.addEventListener('click',function(){ select(el.dataset.d); });
    });
  }

  function select(ds){
    selected=ds;
    input.value=ds;
    var parts=ds.split('-');
    trigger.querySelector('[data-slot="display"]').textContent=parts[2]+'. '+monthNames[parseInt(parts[1])-1]+' '+parts[0];
    dropdown.style.display='none';
    render();
    input.dispatchEvent(new Event('change',{bubbles:true}));
  }

  function toggle(){ dropdown.style.display=dropdown.style.display==='none'?'block':'none'; }

  trigger.addEventListener('click',function(e){ e.stopPropagation(); toggle(); });
  trigger.addEventListener('keydown',function(e){ if(e.key==='Enter'||e.key===' '){ e.preventDefault(); e.stopPropagation(); toggle(); } });
  document.addEventListener('click',function(){ dropdown.style.display='none'; });
  dropdown.addEventListener('click',function(e){ e.stopPropagation(); });
  prevBtn.addEventListener('click',function(){ month--; if(month<1){month=12;year--;} render(); });
  nextBtn.addEventListener('click',function(){ month++; if(month>12){month=1;year++;} render(); });

  render();
  if(selected){ var p=selected.split('-'); trigger.querySelector('[data-slot="display"]').textContent=p[2]+'. '+monthNames[parseInt(p[1])-1]+' '+p[0]; }
})();`,
		sig, initYear, initMonth,
		func() string {
			if p.Min != "" {
				return fmt.Sprintf(`'%s'`, p.Min)
			}
			return "null"
		}(),
		func() string {
			if p.Max != "" {
				return fmt.Sprintf(`'%s'`, p.Max)
			}
			return "null"
		}(),
	)

	displayText := p.Placeholder
	if p.Value != "" {
		if t, err := time.Parse("2006-01-02", p.Value); err == nil {
			displayText = t.Format("02. January 2006")
		}
	}

	return h.Div(
		g.Attr("data-component", "calendar-picker"),
		g.Attr("data-cal-sig", sig),
		h.Input(h.Type("hidden"), h.Name(p.Name), h.Value(p.Value)),
		// Trigger button
		h.Div(
			g.Attr("data-slot", "trigger"),
			g.Attr("tabindex", "0"),
			h.Style("cursor:pointer"),
			g.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>`),
			h.Span(g.Attr("data-slot", "display"), g.Text(displayText)),
		),
		// Dropdown calendar
		h.Div(
			g.Attr("data-slot", "dropdown"),
			h.Style("display:none;position:absolute;z-index:60;background:var(--surface);border:var(--bw-2) solid var(--ink);border-radius:var(--radius);box-shadow:var(--shadow-lg);padding:var(--sp-3);min-width:260px"),
			// Header: prev / month-year / next
			h.Div(
				h.Style("display:flex;align-items:center;justify-content:space-between;margin-bottom:var(--sp-3)"),
				h.Button(h.Type("button"), g.Attr("data-slot", "prev"), h.Style("background:none;border:var(--bw-1) solid var(--line);border-radius:4px;padding:4px 8px;cursor:pointer"), g.Text("‹")),
				h.Span(g.Attr("data-slot", "month-label"), h.Style("font-weight:800;font-size:var(--t-sm)")),
				h.Button(h.Type("button"), g.Attr("data-slot", "next"), h.Style("background:none;border:var(--bw-1) solid var(--line);border-radius:4px;padding:4px 8px;cursor:pointer"), g.Text("›")),
			),
			h.Div(g.Attr("data-slot", "grid")),
		),
		h.Script(g.Raw(script)),
	)
}
