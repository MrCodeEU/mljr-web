package pages

import (
	"encoding/json"
	"slices"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui/swiss"
)

// careerMap sits under the Netzplan: one square logo marker per location.
// Selecting a stop or timeline bar flies there and highlights its marker;
// clicking a marker selects its stop. stopOf maps timeline IDs to the
// Netzplan panel keys the tabs use. Tiles render greyscale.
func careerMap(d hpdata.SiteData, stopOf map[string]string, lang string) g.Node {
	type pin struct {
		Lat  float64  `json:"lat"`
		Lng  float64  `json:"lng"`
		Logo string   `json:"logo"`
		Name string   `json:"name"`
		IDs  []string `json:"ids"`
	}
	var pins []*pin
	byLoc := map[orgLocation]*pin{}
	for _, t := range d.Timeline {
		loc, ok := orgLocationFor(t.ID, t.Organization)
		if !ok {
			continue
		}
		p := byLoc[loc]
		if p == nil {
			p = &pin{Lat: loc.lat, Lng: loc.lng, Logo: t.Logo, Name: t.Organization}
			byLoc[loc] = p
			pins = append(pins, p)
		}
		if k := stopOf[t.ID]; k != "" && !slices.Contains(p.IDs, k) {
			p.IDs = append(p.IDs, k)
		}
	}
	data, err := json.Marshal(pins)
	if err != nil {
		return nil
	}
	return h.Div(h.Class("cmap"),
		h.Div(h.Class("cmap-head"),
			swiss.Kicker(tr(lang, "Karte · Ort der gewählten Station", "Map · location of the selected stop"), h.Span(h.ID("cmap-name"), g.Text("—"))),
		),
		h.Div(h.ID("cmap"), h.Class("cmap-map")),
		h.Script(h.Type("application/json"), h.ID("cmap-data"), g.Raw(strings.ReplaceAll(string(data), "</", `<\/`))),
		h.Script(g.Raw(careerMapJS)),
	)
}

const careerMapJS = `(function(){
  var pins=JSON.parse(document.getElementById('cmap-data').textContent), map, marks=[], want=null;
  function icon(p,on){
    var html='<div class="cmap-pin'+(on?' on':'')+'">'+(p.logo?'<img src="'+p.logo+'" alt="">':'')+'</div>';
    return L.divIcon({className:'',iconSize:on?[40,40]:[28,28],iconAnchor:on?[20,20]:[14,14],html:html});
  }
  function show(id){
    want=id; if(!map) return;
    marks.forEach(function(m){
      var on=m.p.ids.indexOf(id)>=0; m.mk.setIcon(icon(m.p,on)); m.mk.setZIndexOffset(on?1000:0);
      if(on){ map.flyTo([m.p.lat,m.p.lng],Math.max(map.getZoom(),11),{duration:.8}); document.getElementById('cmap-name').textContent=m.p.name; }
    });
  }
  function init(){
    if(map||typeof L==='undefined') return;
    map=L.map('cmap',{scrollWheelZoom:false,zoomControl:true}).setView([48.2,14.3],9);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',{maxZoom:18,attribution:'&copy; OpenStreetMap'}).addTo(map);
    pins.forEach(function(p){
      var mk=L.marker([p.lat,p.lng],{icon:icon(p,false),title:p.name}).addTo(map); marks.push({p:p,mk:mk});
      // marker click selects its (first) stop on the Netzplan
      mk.on('click',function(){ var t=p.ids[0]&&document.querySelector('#nz [data-tab="'+p.ids[0]+'"]'); if(t) t.dispatchEvent(new MouseEvent('click',{bubbles:true})); });
    });
    if(want) show(want);
  }
  document.addEventListener('tabs:select',function(e){ if(e.target.id==='nz') show(e.detail.id); });
  var cur=document.querySelector('#nz [data-tab][aria-pressed="true"]'); if(cur) want=cur.getAttribute('data-tab');
  var el=document.getElementById('cmap');
  // Leaflet (~150 KB) is only fetched once the map comes near the viewport
  new IntersectionObserver(function(es,o){
    if(!es[0].isIntersecting) return; o.disconnect();
    if(typeof L!=='undefined'){ init(); return; }
    var css=document.createElement('link'); css.rel='stylesheet'; css.href='/static/leaflet.css'; document.head.appendChild(css);
    var js=document.createElement('script'); js.src='/static/leaflet.js'; js.onload=init; document.head.appendChild(js);
  },{rootMargin:'400px'}).observe(el);
})();`

const careerMapCSS = `
.cmap { border-top: 3px solid var(--ink); }
.cmap-head { padding: .75rem var(--swiss-pad) .5rem; }
.cmap-pin { cursor: pointer; }
.cmap-map { height: 360px; border-top: 1px solid var(--ink); z-index: 0; }
.cmap-pin { width: 100%; height: 100%; background: #fff; border: 2px solid var(--ink); display: grid; place-items: center; overflow: hidden; box-sizing: border-box; }
.leaflet-container .leaflet-marker-pane .cmap-pin img { width: 100%; height: 100%; padding: 2px; box-sizing: border-box; object-fit: contain; filter: grayscale(1); }
.cmap-pin.on { border: 4px solid var(--accent); }
.leaflet-container .leaflet-marker-pane .cmap-pin.on img { filter: none; }
#cmap-name { color: var(--main-text); }
`
