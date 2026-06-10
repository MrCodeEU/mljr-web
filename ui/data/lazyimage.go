package data

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type LazyImageProps struct {
	Src         string
	Alt         string
	Width       string // CSS width
	Height      string // CSS height (also sets aspect ratio container)
	ObjectFit   string // CSS object-fit (default "cover")
	Placeholder string // CSS background color while loading (default "var(--surface-2)")
	Rounded     bool
}

// LazyImage renders an image that loads only when it enters the viewport.
// Shows a skeleton placeholder until loaded.
// Uses native loading="lazy" + IntersectionObserver fade-in.
func LazyImage(p LazyImageProps) g.Node {
	if p.ObjectFit == "" {
		p.ObjectFit = "cover"
	}
	if p.Placeholder == "" {
		p.Placeholder = "var(--surface-2)"
	}

	style := "display:block;width:" + p.Width + ";height:" + p.Height +
		";object-fit:" + p.ObjectFit + ";background:" + p.Placeholder +
		";transition:opacity .3s ease;opacity:0"
	if p.Rounded {
		style += ";border-radius:var(--radius)"
	}

	return g.Group{
		h.Img(
			g.Attr("data-component", "lazy-image"),
			g.Attr("data-src", p.Src),
			h.Alt(p.Alt),
			h.Style(style),
			g.Attr("loading", "lazy"),
		),
		h.Script(g.Raw(`(function(){
  var imgs=document.querySelectorAll('[data-component="lazy-image"]');
  if('IntersectionObserver' in window){
    var io=new IntersectionObserver(function(entries){
      entries.forEach(function(e){
        if(e.isIntersecting){
          var img=e.target;
          img.src=img.dataset.src;
          img.onload=function(){img.style.opacity='1';};
          io.unobserve(img);
        }
      });
    },{rootMargin:'200px'});
    imgs.forEach(function(img){ io.observe(img); });
  } else {
    imgs.forEach(function(img){ img.src=img.dataset.src; img.style.opacity='1'; });
  }
})()`)),
	}
}
