//go:build showcase

package special

import (
	"mljr-web/ui/registry"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func init() {
	registry.Register(&registry.Component{
		Slug: "open-map", Name: "OpenMap", Category: "special",
		Summary: "Leaflet-powered interactive map using OpenStreetMap tiles. Self-hosted JS (148KB). Drop pins with popups.",
		Code: `special.OpenMap(special.OpenMapProps{
    CenterLat: 52.52, CenterLng: 13.405,
    Zoom: 12, Height: "400px",
},
    special.MapPin{Lat: 52.52, Lng: 13.405, Label: "Berlin HQ"},
)`,
		Render: func(p map[string]string) g.Node {
			return h.Div(
				h.Style("display:flex;flex-direction:column;gap:var(--sp-3)"),
				OpenMap(OpenMapProps{
					CenterLat: 52.52,
					CenterLng: 13.405,
					Zoom:      12,
					Height:    "360px",
					ID:        "showcase-map",
				},
					MapPin{Lat: 52.52, Lng: 13.405, Label: "Berlin", Popup: "<strong>Berlin HQ</strong><br>Mitte District"},
					MapPin{Lat: 52.501, Lng: 13.435, Label: "Office 2"},
				),
				h.P(h.Style("font-size:var(--t-xs);color:var(--muted);margin:0"),
					g.Text("Leaflet 1.9 self-hosted at /static/leaflet.js (148KB). OpenStreetMap tiles — GDPR-friendly alternative: use self-hosted tiles.")),
			)
		},
	})
}
