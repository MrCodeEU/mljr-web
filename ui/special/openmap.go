package special

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MapPin struct {
	Lat       float64
	Lng       float64
	AnchorLat float64 // optional true location for offset markers/leader lines
	AnchorLng float64
	Label     string
	Popup     string // HTML content for popup (optional)
	Icon      string // optional image URL for a logo marker
	Number    int    // optional compact numbered marker
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

// OpenMap renders a Leaflet-powered interactive map.
// Requires /static/leaflet.js and /static/leaflet.css to be served.
// Uses OpenStreetMap tiles by default — self-host tiles for full privacy.
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
		centerLat = pins[0].Lat
		centerLng = pins[0].Lng
	}

	// Build pins JS
	var pinParts []string
	baseZoom := p.Zoom
	for _, pin := range pins {
		popup := ""
		if pin.Popup != "" {
			popup = fmt.Sprintf(`.bindPopup(%s)`, jsStr(pin.Popup))
		} else if pin.Label != "" {
			popup = fmt.Sprintf(`.bindPopup(%s)`, jsStr(pin.Label))
		}
		iconVar := "pin"
		if pin.Number > 0 {
			iconVar = fmt.Sprintf(`numberIcon(%d)`, pin.Number)
		} else if pin.Icon != "" {
			iconVar = fmt.Sprintf(`logoIcon(%s)`, jsStr(pin.Icon))
		}
		if pin.AnchorLat != 0 || pin.AnchorLng != 0 {
			pinParts = append(pinParts, fmt.Sprintf(
				`(function(){
    var anchor=[%v,%v], delta=[%v,%v], baseZoom=%d;
    var line=L.polyline([anchor,anchor],{color:'var(--ink,#1a1a1a)',weight:2,opacity:.75,dashArray:'4 4',interactive:false}).addTo(map);
    L.circleMarker(anchor,{radius:4,color:'var(--ink,#1a1a1a)',weight:2,fillColor:'var(--accent,#ff5941)',fillOpacity:1,interactive:false}).addTo(map);
    var marker=L.marker(anchor,{icon:%s}).addTo(map)%s;
    function markerPos(){
      var factor=Math.pow(2,baseZoom-map.getZoom());
      factor=Math.max(.45,Math.min(5.5,factor));
      return [anchor[0]+delta[0]*factor,anchor[1]+delta[1]*factor];
    }
    function updateMarkerOffset(){
      var pos=markerPos();
      marker.setLatLng(pos);
      line.setLatLngs([anchor,pos]);
    }
    updateMarkerOffset();
    map.on('zoomend',updateMarkerOffset);
  })();`,
				pin.AnchorLat, pin.AnchorLng, pin.Lat-pin.AnchorLat, pin.Lng-pin.AnchorLng, baseZoom, iconVar, popup,
			))
			continue
		}
		pinParts = append(pinParts, fmt.Sprintf(
			`L.marker([%v,%v],{icon:%s}).addTo(map)%s;`,
			pin.Lat, pin.Lng, iconVar, popup,
		))
	}

	// Default Leaflet markers reference PNGs relative to the CSS file, which
	// are not shipped. Use a self-contained neo-brutalist divIcon pin instead.
	script := fmt.Sprintf(`(function(){
  if(typeof L==='undefined'){ console.warn('Leaflet not loaded'); return; }
  var m=document.getElementById('%s');
  if(!m||m._leaflet_id) return;
  var map=L.map('%s').setView([%v,%v],%d);
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
  %s
  setTimeout(function(){ map.invalidateSize(); }, 0);
})();`,
		p.ID, p.ID, centerLat, centerLng, p.Zoom,
		jsStr(p.TileURL), jsStr(p.TileAttrib),
		strings.Join(pinParts, "\n  "),
	)

	return h.Div(
		g.Attr("data-component", "open-map"),
		h.Link(h.Rel("stylesheet"), h.Href("/static/leaflet.css")),
		h.Script(h.Src("/static/leaflet.js")),
		h.Div(
			h.ID(p.ID),
			g.Attr("data-slot", "map"),
			h.Style(fmt.Sprintf("height:%s;width:100%%;border-radius:var(--radius);border:var(--bw-2) solid var(--ink);z-index:0", p.Height)),
		),
		h.Script(g.Raw(script)),
	)
}
