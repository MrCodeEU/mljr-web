package swiss

import (
	"strconv"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GalleryProps struct {
	Images []string
	Alt    string
	// Aspect of the stage, e.g. "4/3" (default). Images are fitted inside it
	// on the paper-tone ground and never scaled up, cropped or stretched.
	Aspect string
	// Label prefixes the counter, default "Abb." (Abbildung).
	Label      string
	Prev, Next string // aria-labels; default German
}

// Gallery shows one image at a time on a flat stage with a caption bar
// (counter + ◀ ▶). Without JS the first image shows. Include GalleryScript
// once per page.
func Gallery(p GalleryProps) g.Node {
	if len(p.Images) == 0 {
		return nil
	}
	if p.Aspect == "" {
		p.Aspect = "4/3"
	}
	if p.Label == "" {
		p.Label = "Abb."
	}
	if p.Prev == "" {
		p.Prev = "Vorheriges Bild"
	}
	if p.Next == "" {
		p.Next = "Nächstes Bild"
	}
	n := len(p.Images)
	imgs := make([]g.Node, n)
	for i, src := range p.Images {
		alt := p.Alt
		if n > 1 {
			alt += " (" + strconv.Itoa(i+1) + "/" + strconv.Itoa(n) + ")"
		}
		imgs[i] = h.Img(h.Src(src), h.Alt(alt), g.Attr("loading", "lazy"), g.Attr("decoding", "async"), g.If(i > 0, g.Attr("hidden", "")))
	}
	return h.Figure(
		g.Attr("data-component", "gallery"),
		g.Attr("data-count", strconv.Itoa(n)),
		h.Div(append([]g.Node{g.Attr("data-slot", "stage"), h.Style("aspect-ratio:" + p.Aspect)}, imgs...)...),
		h.FigCaption(g.Attr("data-slot", "bar"),
			h.Span(g.Attr("data-slot", "counter"), g.Attr("aria-live", "polite"),
				g.Text(p.Label+" "), h.Span(g.Attr("data-slot", "pos"), g.Text("1")), g.Text(" / "+strconv.Itoa(n)),
			),
			g.If(n > 1, h.Span(g.Attr("data-slot", "ticks"), g.Attr("aria-hidden", "true"), g.Group(ticks(n)))),
			g.If(n > 1, h.Span(g.Attr("data-slot", "nav"),
				h.Button(h.Type("button"), g.Attr("data-component", "swiss-btn"), g.Attr("data-dir", "-1"), g.Attr("aria-label", p.Prev), g.Text("◀")),
				h.Button(h.Type("button"), g.Attr("data-component", "swiss-btn"), g.Attr("data-dir", "1"), g.Attr("aria-label", p.Next), g.Text("▶")),
			)),
		),
	)
}

// GalleryScript wires every Gallery on the page via one delegated listener.
func ticks(n int) []g.Node {
	out := make([]g.Node, n)
	for i := range out {
		out[i] = h.I(g.If(i == 0, g.Attr("data-cur", "")))
	}
	return out
}

func GalleryScript() g.Node {
	return h.Script(g.Raw(`(function(){
  document.addEventListener('click',function(e){
    var b=e.target.closest('[data-component=gallery] [data-dir]'); if(!b) return;
    var f=b.closest('[data-component=gallery]'), imgs=f.querySelectorAll('[data-slot=stage] img'), n=imgs.length, cur=0;
    for(var i=0;i<n;i++) if(!imgs[i].hidden) cur=i;
    var nx=(cur+(+b.dataset.dir)+n)%n;
    imgs[cur].hidden=true; imgs[nx].hidden=false;
    f.querySelector('[data-slot=pos]').textContent=nx+1;
    f.querySelectorAll('[data-slot=ticks] i').forEach(function(t,i){t.toggleAttribute('data-cur',i===nx);});
  });
})();`))
}
