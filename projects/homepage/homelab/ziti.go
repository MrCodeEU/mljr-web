package homelab

import (
	"context"
	"log"
	"sort"
	"time"
)

// controllerRouters names the routers that run next to the Ziti controller
// on the public host. Every other router sits at home and only dials out.
var controllerRouters = map[string]bool{"mljr": true}

// Mesh queries for the ziti-mesh-exporter metrics in VictoriaMetrics. Each
// query aggregates down to router names, so identity, client and service
// labels (personal devices, internal service names) never reach the page.
const (
	qRouterOnline = `max by (router) (ziti_mesh_router_online)`
	qLinkUp       = `max by (source, target) (ziti_mesh_link_up)`
	qLinkLatency  = `max by (source, target) (ziti_mesh_link_latency_seconds)`
	qTerminators  = `count by (router) (ziti_mesh_terminators > 0)`
	qCircuits     = `sum(ziti_mesh_circuits)`
)

// fetchZitiMesh fills the mesh part of s from the OpenZiti exporter
// metrics. Without router data it leaves MeshOK false, so poll keeps the
// previous mesh.
func (p *Poller) fetchZitiMesh(ctx context.Context, s *Snapshot) {
	routers, err := p.promQuerySeries(ctx, qRouterOnline)
	if err != nil {
		log.Printf("homelab: ziti routers: %v", err)
		return
	}
	hosts := map[string]*MeshHost{}
	for _, r := range routers {
		name := r.Metric["router"]
		if name == "" {
			continue
		}
		hosts[name] = &MeshHost{Name: name, Controller: controllerRouters[name], Online: r.Value == 1}
	}
	if len(hosts) == 0 {
		return
	}

	// Links and their latency. A link to a router that is offline counts as
	// down, whatever the other end still reports.
	latency := map[[2]string]time.Duration{}
	if lat, err := p.promQuerySeries(ctx, qLinkLatency); err == nil {
		for _, l := range lat {
			latency[[2]string{l.Metric["source"], l.Metric["target"]}] = time.Duration(l.Value * float64(time.Second))
		}
	}
	if links, err := p.promQuerySeries(ctx, qLinkUp); err != nil {
		log.Printf("homelab: ziti links: %v", err)
	} else {
		for _, l := range links {
			src, dst := hosts[l.Metric["source"]], hosts[l.Metric["target"]]
			if src == nil || dst == nil {
				continue
			}
			link := MeshLink{Source: src.Name, Target: dst.Name, Up: l.Value == 1 && src.Online && dst.Online}
			if link.Up {
				link.Latency = latency[[2]string{src.Name, dst.Name}]
			}
			s.MeshLinks = append(s.MeshLinks, link)
			for _, h := range []*MeshHost{src, dst} {
				h.LinksTotal++
				if link.Up {
					h.LinksUp++
				}
			}
		}
	}

	if terms, err := p.promQuerySeries(ctx, qTerminators); err == nil {
		for _, t := range terms {
			if h := hosts[t.Metric["router"]]; h != nil {
				h.Services = int(t.Value)
			}
		}
	}

	// No series means no open circuit: the exporter only publishes a series
	// per live connection, and the routers above prove it is reachable.
	s.MeshCircuits = 0
	if c, err := p.promQuerySeries(ctx, qCircuits); err == nil && len(c) > 0 {
		s.MeshCircuits = int(c[0].Value)
	}

	for _, h := range hosts {
		s.Mesh = append(s.Mesh, *h)
	}
	sortMesh(s.Mesh, s.MeshLinks)
	s.MeshOK = true
}

// sortMesh orders the controller first, then routers by name, and links by
// their endpoints, so the page stays stable between polls.
func sortMesh(hosts []MeshHost, links []MeshLink) {
	sort.Slice(hosts, func(a, b int) bool {
		if hosts[a].Controller != hosts[b].Controller {
			return hosts[a].Controller
		}
		return hosts[a].Name < hosts[b].Name
	})
	sort.Slice(links, func(a, b int) bool {
		if links[a].Source != links[b].Source {
			return links[a].Source < links[b].Source
		}
		return links[a].Target < links[b].Target
	})
}
