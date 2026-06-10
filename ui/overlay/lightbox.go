package overlay

import (
	"fmt"
	"mljr-web/ui/primitive"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type LightboxImage struct {
	Src     string
	Thumb   string // thumbnail src; defaults to Src
	Alt     string
	Caption string
}

type LightboxProps struct {
	ID         string          // unique per page (default "lb")
	Images     []LightboxImage
	Columns    int    // thumbnail grid columns (default 3)
	ThumbSize  string // CSS size for thumbnails (default "120px")
}

// Lightbox renders a thumbnail grid that opens full-screen on click.
// Click backdrop or × to close. ←/→ keys navigate between images.
func Lightbox(p LightboxProps) g.Node {
	if p.ID == "" {
		p.ID = "lb"
	}
	if p.Columns == 0 {
		p.Columns = 3
	}
	if p.ThumbSize == "" {
		p.ThumbSize = "120px"
	}

	sigOpen := "_" + p.ID + "Open"
	sigIdx := "_" + p.ID + "Idx"

	thumbs := make([]g.Node, len(p.Images))
	for i, img := range p.Images {
		thumb := img.Thumb
		if thumb == "" {
			thumb = img.Src
		}
		thumbs[i] = h.Button(
			h.Type("button"),
			h.Style(fmt.Sprintf("padding:0;border:var(--bw-2) solid var(--line);border-radius:var(--radius);overflow:hidden;cursor:zoom-in;width:%s;height:%s;flex-shrink:0", p.ThumbSize, p.ThumbSize)),
			g.Attr("data-on:click", fmt.Sprintf("$%s=true;$%s=%d", sigOpen, sigIdx, i)),
			h.Img(
				h.Src(thumb),
				h.Alt(img.Alt),
				h.Style("width:100%;height:100%;object-fit:cover;display:block"),
				g.Attr("loading", "lazy"),
			),
		)
	}

	return h.Div(
		g.Attr("data-component", "lightbox"),
		g.Attr("data-signals", fmt.Sprintf(`{"%s":false,"%s":0}`, sigOpen, sigIdx)),

		// Thumbnail grid
		h.Div(
			h.Style(fmt.Sprintf("display:flex;flex-wrap:wrap;gap:var(--sp-2)")),
			g.Group(thumbs),
		),

		// Fullscreen overlay
		h.Div(
			g.Attr("data-slot", "overlay"),
			g.Attr("data-show", "$"+sigOpen),
			h.Style("display:none"),
			h.Div(g.Attr("data-slot", "backdrop"), g.Attr("data-on:click", "$"+sigOpen+"=false")),
			h.Div(
				g.Attr("data-slot", "viewer"),
				// Image rendered via JS (avoids data-attr complexity with arrays)
				h.Img(h.ID(p.ID+"-img"), h.Alt(""), h.Style("max-width:90vw;max-height:80vh;object-fit:contain;display:block")),
				h.Div(
					g.Attr("data-slot", "caption"),
					h.ID(p.ID+"-cap"),
				),
				h.Div(
					g.Attr("data-slot", "controls"),
					primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeIcon, Attrs: []g.Node{
						g.Attr("data-on:click", "$"+sigIdx+"=Math.max(0,$"+sigIdx+"-1)"),
						g.Attr("aria-label", "Previous"),
					}}, g.Raw("←")),
					primitive.Button(primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeIcon, Attrs: []g.Node{
						g.Attr("data-on:click", "$"+sigOpen+"=false"),
						g.Attr("aria-label", "Close"),
					}}, g.Raw("✕")),
					primitive.Button(primitive.ButtonProps{Variant: token.Outline, Size: token.SizeIcon, Attrs: []g.Node{
						g.Attr("data-on:click", fmt.Sprintf("$%s=Math.min(%d,$%s+1)", sigIdx, len(p.Images)-1, sigIdx)),
						g.Attr("aria-label", "Next"),
					}}, g.Raw("→")),
				),
			),
			// JS: sync image src to index signal + keyboard nav
			h.Script(g.Raw(lightboxScript(p))),
		),
	)
}

func lightboxScript(p LightboxProps) string {
	srcs := "["
	caps := "["
	for i, img := range p.Images {
		if i > 0 {
			srcs += ","
			caps += ","
		}
		srcs += fmt.Sprintf("%q", img.Src)
		caps += fmt.Sprintf("%q", img.Caption)
	}
	srcs += "]"
	caps += "]"

	sigOpen := "_" + p.ID + "Open"
	sigIdx := "_" + p.ID + "Idx"

	return fmt.Sprintf(`(function(){
  var srcs=%s; var caps=%s;
  var imgEl=document.getElementById('%s-img');
  var capEl=document.getElementById('%s-cap');
  var open='%s'; var idx='%s';
  function update(){
    var i=window.__ds&&window.__ds.store?window.__ds.store[idx]:0;
    if(imgEl){ imgEl.src=srcs[i]||''; imgEl.alt=caps[i]||''; }
    if(capEl){ capEl.textContent=caps[i]||''; }
  }
  // Poll for signal changes (simple approach)
  setInterval(update, 50);
  document.addEventListener('keydown',function(e){
    if(window.__ds&&window.__ds.store&&window.__ds.store[open]){
      if(e.key==='ArrowLeft'&&window.__ds.store[idx]>0) window.__ds.store[idx]--;
      if(e.key==='ArrowRight'&&window.__ds.store[idx]<srcs.length-1) window.__ds.store[idx]++;
      if(e.key==='Escape') window.__ds.store[open]=false;
    }
  });
})();`, srcs, caps, p.ID, p.ID, sigOpen, sigIdx)
}
