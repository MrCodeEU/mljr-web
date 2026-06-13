package special

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MapPin struct {
	Lat          float64
	Lng          float64
	AnchorLat    float64 // if set (with AnchorLng), used as the pin's true position instead of Lat/Lng
	AnchorLng    float64
	LineLat      float64 // if set (with LineLng), draw a dashed line from here to the pin (true location vs. jittered marker pos)
	LineLng      float64
	SpreadAngle  float64 // if LineLat/LineLng set, direction (radians) to offset this pin from its anchor by SpreadRadius pixels (stays separated at any zoom, never drifts off-map)
	SpreadRadius float64 // pixel distance for the SpreadAngle offset; 0 means sit exactly at the anchor (no leader line drawn)
	Label        string
	Popup        string // HTML content for popup (optional)
	Icon         string // optional image URL for a logo marker
	Number       int    // optional compact numbered marker
}

type OpenMapProps struct {
	// Center of the map [lat, lng]. Defaults to first pin or [0,0].
	CenterLat float64
	CenterLng float64
	// Zoom level (default 13).
	Zoom int
	// Height of the map container (default "400px").
	Height string
	// ID must be unique per page (default "map").
	ID string
	// TileURL for OpenStreetMap tiles (default: OSM standard).
	// Self-hosted tiles: "https://tiles.yourdomain.com/{z}/{x}/{y}.png"
	TileURL string
	// TileAttrib is the tile attribution string.
	TileAttrib string
}

// OpenMap renders a Leaflet-powered interactive map with clustered markers.
// Requires /static/leaflet.js, /static/leaflet.css, /static/leaflet.markercluster.js,
// /static/MarkerCluster.css and /static/MarkerCluster.Default.css to be served.
// Nearby pins automatically group into a cluster bubble and spiderfy (spread
// with leader lines) on click/max-zoom — isolated pins always sit at their
// true location.
func OpenMap(p OpenMapProps, pins ...MapPin) g.Node {
	if p.Zoom == 0 {
		p.Zoom = 13
	}
	if p.Height == "" {
		p.Height = "400px"
	}
	if p.ID == "" {
		p.ID = "map"
	}
	if p.TileURL == "" {
		p.TileURL = "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
	}
	if p.TileAttrib == "" {
		p.TileAttrib = `&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors`
	}

	// Default center to first pin if not specified
	centerLat := p.CenterLat
	centerLng := p.CenterLng
	if centerLat == 0 && centerLng == 0 && len(pins) > 0 {
		lat, lng := pins[0].Lat, pins[0].Lng
		if pins[0].AnchorLat != 0 || pins[0].AnchorLng != 0 {
			lat, lng = pins[0].AnchorLat, pins[0].AnchorLng
		}
		centerLat, centerLng = lat, lng
	}

	type spreadMember struct {
		iconVar, tooltip string
		angle, radius    float64
	}
	groups := map[[2]float64][]spreadMember{}
	var groupOrder [][2]float64

	// Build pins JS
	var pinParts []string
	for _, pin := range pins {
		lat, lng := pin.Lat, pin.Lng
		if pin.AnchorLat != 0 || pin.AnchorLng != 0 {
			lat, lng = pin.AnchorLat, pin.AnchorLng
		}
		tooltip := ""
		if pin.Popup != "" {
			tooltip = fmt.Sprintf(`.bindTooltip(%s,{direction:'top',offset:[0,-10],opacity:0.96,className:'map-tooltip'})`, jsStr(pin.Popup))
		} else if pin.Label != "" {
			tooltip = fmt.Sprintf(`.bindTooltip(%s,{direction:'top',offset:[0,-10],opacity:0.96,className:'map-tooltip'})`, jsStr(pin.Label))
		}
		iconVar := "pin"
		if pin.Number > 0 {
			iconVar = fmt.Sprintf(`numberIcon(%d)`, pin.Number)
		} else if pin.Icon != "" {
			iconVar = fmt.Sprintf(`logoIcon(%s)`, jsStr(pin.Icon))
		}
		if pin.LineLat != 0 || pin.LineLng != 0 {
			key := [2]float64{pin.LineLat, pin.LineLng}
			if _, ok := groups[key]; !ok {
				groupOrder = append(groupOrder, key)
			}
			groups[key] = append(groups[key], spreadMember{iconVar: iconVar, tooltip: tooltip, angle: pin.SpreadAngle, radius: pin.SpreadRadius})
		} else {
			pinParts = append(pinParts, fmt.Sprintf(
				`cluster.addLayer(L.marker([%v,%v],{icon:%s})%s);`,
				lat, lng, iconVar, tooltip,
			))
		}
	}

	// Duplicate-location groups: zoomed out, show a single count badge in the
	// normal cluster (so it merges naturally with nearby cluster bubbles);
	// zoomed in past the spread threshold, swap to individually placed pins on
	// fixed pixel-radius spokes around a shared anchor dot.
	for _, key := range groupOrder {
		members := groups[key]
		var memberJS strings.Builder
		for _, mb := range members {
			fmt.Fprintf(&memberJS,
				`members.push({mk:L.marker(a,{icon:%s})%s,ln:L.polyline([a,a],{color:'var(--ink,#1a1a1a)',weight:2,opacity:.6,dashArray:'4 4'}),angle:%v,radius:%v});`,
				mb.iconVar, mb.tooltip, mb.angle, mb.radius,
			)
		}
		pinParts = append(pinParts, fmt.Sprintf(`(function(){
  var a=L.latLng(%v,%v);
  var rep=L.marker(a,{icon:numberIcon(%d)});
  cluster.addLayer(rep);
  var dot=L.circleMarker(a,{radius:4,color:'var(--ink,#1a1a1a)',weight:2,fillColor:'var(--accent,#ff5941)',fillOpacity:0,opacity:0}).addTo(map);
  var members=[];
  %s
  members.forEach(function(m){ m.ln.addTo(map); m.mk.addTo(map); m.mk.setOpacity(0); m.ln.setStyle({opacity:0}); });
  var expanded=null;
  function reposition(){
    members.forEach(function(m){
      var p=map.latLngToLayerPoint(a).add([Math.cos(m.angle)*m.radius,Math.sin(m.angle)*m.radius]);
      var b=map.layerPointToLatLng(p);
      m.mk.setLatLng(b);
      m.ln.setLatLngs([a,b]);
    });
    dot.bringToFront();
  }
  function setInteractive(layer,on){
    var el=layer._icon||layer._path;
    if(el) el.style.pointerEvents=on?'':'none';
  }
  function update(){
    var want=map.getZoom()>=12;
    if(want===expanded){ if(want) reposition(); return; }
    expanded=want;
    rep.setOpacity(want?0:1);
    setInteractive(rep,!want);
    dot.setStyle({opacity:want?1:0,fillOpacity:want?1:0});
    members.forEach(function(m){
      m.mk.setOpacity(want?1:0);
      setInteractive(m.mk,want);
      m.ln.setStyle({opacity:want?.6:0});
    });
    if(want) reposition();
  }
  map.on('zoomend',update);
  update();
})();`, key[0], key[1], len(members), memberJS.String()))
	}

	// Default Leaflet markers reference PNGs relative to the CSS file, which
	// are not shipped. Use a self-contained neo-brutalist divIcon pin instead.
	// Cluster bubbles get a matching neo-brutalist treatment via iconCreateFunction.
	script := fmt.Sprintf(`(function(){
  if(typeof L==='undefined'){ console.warn('Leaflet not loaded'); return; }
  var m=document.getElementById('%s');
  if(!m||m._leaflet_id) return;
  var map=L.map('%s',{maxZoom:19}).setView([%v,%v],%d);
  L.tileLayer(%s,{attribution:%s,maxZoom:19}).addTo(map);
  var pin=L.divIcon({className:'',iconSize:[22,22],iconAnchor:[11,22],popupAnchor:[0,-24],
    html:'<div style="width:22px;height:22px;border-radius:50%% 50%% 50%% 0;transform:rotate(-45deg);background:var(--accent,#ff5941);border:3px solid var(--ink,#1a1a1a);box-shadow:2px 2px 0 rgba(0,0,0,.35);box-sizing:border-box"></div>'});
  function logoIcon(src){
    var html='<div style="width:32px;height:32px;border:3px solid var(--ink,#1a1a1a);background:var(--bg,#fff);box-shadow:2px 2px 0 rgba(0,0,0,.35);display:grid;place-items:center;overflow:hidden;box-sizing:border-box"><img src="'+src+'" alt="" style="width:100%%;height:100%%;object-fit:contain;background:#fff"></div>';
    return L.divIcon({className:'',iconSize:[32,32],iconAnchor:[16,32],popupAnchor:[0,-34],html:html});
  }
  function numberIcon(n){
    var html='<div style="width:24px;height:24px;border-radius:50%%;border:3px solid var(--ink,#1a1a1a);background:var(--accent,#ff5941);color:var(--accent-ink,#fff);box-shadow:2px 2px 0 rgba(0,0,0,.35);display:grid;place-items:center;font:900 12px/1 var(--font-mono,monospace);box-sizing:border-box">'+n+'</div>';
    return L.divIcon({className:'',iconSize:[24,24],iconAnchor:[12,12],popupAnchor:[0,-14],html:html});
  }
  var cluster=L.markerClusterGroup({
    maxClusterRadius:28,
    spiderfyOnMaxZoom:true,
    zoomToBoundsOnClick:false,
    showCoverageOnHover:false,
    spiderLegPolylineOptions:{color:'var(--ink,#1a1a1a)',weight:2,opacity:.75,dashArray:'4 4'},
    iconCreateFunction:function(c){
      var n=c.getChildCount();
      var html='<div style="width:34px;height:34px;border-radius:50%%;border:3px solid var(--ink,#1a1a1a);background:var(--accent,#ff5941);color:var(--accent-ink,#fff);box-shadow:2px 2px 0 rgba(0,0,0,.35);display:grid;place-items:center;font:900 13px/1 var(--font-mono,monospace);box-sizing:border-box">'+n+'</div>';
      return L.divIcon({className:'',iconSize:[34,34],html:html});
    }
  });
  cluster.on('clusterclick',function(e){ e.layer.spiderfy(); });
  %s
  cluster.addTo(map);
  setTimeout(function(){ map.invalidateSize(); }, 0);
})();`,
		p.ID, p.ID, centerLat, centerLng, p.Zoom,
		jsStr(p.TileURL), jsStr(p.TileAttrib),
		strings.Join(pinParts, "\n  "),
	)

	return h.Div(
		g.Attr("data-component", "open-map"),
		h.Link(h.Rel("stylesheet"), h.Href("/static/leaflet.css")),
		h.Link(h.Rel("stylesheet"), h.Href("/static/MarkerCluster.css")),
		h.Link(h.Rel("stylesheet"), h.Href("/static/MarkerCluster.Default.css")),
		h.Script(h.Src("/static/leaflet.js")),
		h.Script(h.Src("/static/leaflet.markercluster.js")),
		h.Div(
			h.ID(p.ID),
			g.Attr("data-slot", "map"),
			h.Style(fmt.Sprintf("height:%s;width:100%%;border-radius:var(--radius);border:var(--bw-2) solid var(--ink);z-index:0", p.Height)),
		),
		h.Script(g.Raw(script)),
	)
}
