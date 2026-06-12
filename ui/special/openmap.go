package special

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type MapPin struct {
	Lat   float64
	Lng   float64
	Label string
	Popup string // HTML content for popup (optional)
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
	for _, pin := range pins {
		popup := ""
		if pin.Popup != "" {
			popup = fmt.Sprintf(`.bindPopup(%s)`, jsStr(pin.Popup))
		} else if pin.Label != "" {
			popup = fmt.Sprintf(`.bindPopup(%s)`, jsStr(pin.Label))
		}
		pinParts = append(pinParts, fmt.Sprintf(
			`L.marker([%v,%v],{icon:pin}).addTo(map)%s;`,
			pin.Lat, pin.Lng, popup,
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
  %s
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
