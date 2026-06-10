package pages

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/ui"
	uidata "mljr-web/ui/data"
	"mljr-web/ui/icon"
	"mljr-web/ui/layout"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"
)

func projectsSection(projects []hpdata.Project) g.Node {
	total := len(projects)
	pages := (total + perPage - 1) / perPage

	// Build all project cards — each set of perPage wrapped in a page div
	var pageNodes []g.Node
	for p := 0; p < pages; p++ {
		start := p * perPage
		end := start + perPage
		if end > total {
			end = total
		}
		slice := projects[start:end]

		var cols []g.Node
		tones := []token.Tone{token.ToneNone, token.ToneSky, token.ToneLime, token.ToneViolet, token.TonePink, token.ToneMint}
		for i, proj := range slice {
			cols = append(cols, projectCard(proj, tones[i%len(tones)]))
		}

		pageNodes = append(pageNodes,
			h.Div(
				ui.Show(fmt.Sprintf("$pgPage === %d", p)),
				layout.Grid(layout.GridProps{}, g.Group(cols)),
			),
		)
	}

	return h.Section(
		h.ID("projects"),
		h.Style("padding:var(--sp-12) 0"),
		uidata.PaginationSignals("pg", perPage),
		layout.Container(layout.ContainerProps{},
			sectionHeader("Projects", fmt.Sprintf("%d projects", total), token.ToneLime),
			// top pagination
			h.Div(h.Style("margin-bottom:var(--sp-5);display:flex;justify-content:center"),
				uidata.Pagination(uidata.PaginationProps{ID: "pg", Total: total, PerPage: perPage}),
			),
			h.Div(h.ID("projects-pages"), g.Group(pageNodes)),
			// bottom pagination
			h.Div(h.Style("margin-top:var(--sp-6);display:flex;justify-content:center"),
				uidata.Pagination(uidata.PaginationProps{ID: "pg", Total: total, PerPage: perPage}),
			),
		),
		// Animate page transitions via MutationObserver watching display changes
		h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var container=document.getElementById('projects-pages');
  if(!container) return;
  var obs=new MutationObserver(function(muts){
    muts.forEach(function(m){
      if(m.type==='attributes'&&m.attributeName==='style'){
        var el=m.target;
        if(el.style.display!=='none'){
          Motion.animate(el,{opacity:[0,1],y:[16,0]},{duration:0.35,easing:[0.25,0.46,0.45,0.94]});
        }
      }
    });
  });
  Array.from(container.children).forEach(function(el){
    obs.observe(el,{attributes:true,attributeFilter:['style']});
  });
})();`)),
	)
}

func projectCard(p hpdata.Project, tone token.Tone) g.Node {
	imgs := p.LocalImages()

	var carouselNode g.Node
	if len(imgs) > 1 {
		carouselNode = uidata.Carousel(uidata.CarouselProps{
			ID:     "c" + slugify(p.Name),
			Images: imgs,
			Alt:    p.Name,
		})
	} else if len(imgs) == 1 {
		carouselNode = h.Img(
			h.Src(imgs[0]),
			h.Alt(p.Name),
			h.Style("width:100%;aspect-ratio:16/9;object-fit:cover;border-bottom:var(--border-w) solid var(--line);display:block"),
		)
	}

	var linkNodes []g.Node
	if p.URL != "" {
		linkNodes = append(linkNodes,
			h.A(h.Href(p.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
					icon.Icon("simple-icons:github"),
					g.Text("GitHub"),
				),
			),
		)
	}
	for _, lnk := range p.Links {
		if lnk.URL == "" {
			continue
		}
		linkNodes = append(linkNodes,
			h.A(h.Href(lnk.URL), g.Attr("target", "_blank"), g.Attr("rel", "noopener"),
				primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
					icon.Icon("lucide:arrow-up-right"),
					g.Text(lnk.Name),
				),
			),
		)
	}

	topicNodes := make([]g.Node, 0, 4)
	for _, t := range p.Topics {
		if len(topicNodes) >= 4 {
			break
		}
		topicNodes = append(topicNodes, primitive.Tag(
			primitive.TagProps{Icon: hpdata.TechIcon(t)},
			g.Text(t),
		))
	}

	return layout.Col(layout.ColProps{Span: 4},
		primitive.Card(primitive.CardProps{Tone: tone},
			g.If(carouselNode != nil, carouselNode),
			h.Div(h.Style("display:flex;justify-content:space-between;align-items:flex-start;gap:var(--sp-2)"),
				primitive.Heading(primitive.HeadingProps{Level: 3},
					g.Text(friendlyName(p.Name))),
				g.If(p.Stars > 0,
					primitive.Tag(primitive.TagProps{Tone: token.ToneYellow},
						g.Text(fmt.Sprintf("★ %d", p.Stars))),
				),
			),
			h.P(h.Style("margin:0;font-size:var(--t-sm)"), g.Text(truncate(p.Desc, 110))),
			h.Div(h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2)"),
				primitive.Tag(primitive.TagProps{Tone: token.ToneAccent}, g.Text(p.Language)),
				g.Group(topicNodes),
			),
			g.If(len(linkNodes) > 0,
				h.Div(h.Style("display:flex;gap:var(--sp-1);flex-wrap:wrap;margin-top:auto"), g.Group(linkNodes)),
			),
		),
	)
}
