package primitive

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type WordRotateProps struct {
	// Words is the list of words to cycle through. Required.
	Words []string
	// Interval is the ms between word changes (default 2000).
	Interval int
	// ID must be unique per page (default "wr").
	ID string
	// Class is extra CSS class for the visible word span.
	Class string
}

// WordRotate cycles through a list of words with a fade+slide animation.
// Uses requestAnimationFrame-free setInterval + CSS transitions.
func WordRotate(p WordRotateProps) g.Node {
	if p.Interval == 0 {
		p.Interval = 2000
	}
	if p.ID == "" {
		p.ID = "wr"
	}
	if len(p.Words) == 0 {
		return g.Text("")
	}

	// Escape words for JS array
	escaped := make([]string, len(p.Words))
	for i, w := range p.Words {
		escaped[i] = "'" + strings.ReplaceAll(w, "'", "\\'") + "'"
	}
	wordsJS := "[" + strings.Join(escaped, ",") + "]"

	script := fmt.Sprintf(`(function(){
  var el=document.getElementById('%s');
  if(!el) return;
  var words=%s, i=0;
  setInterval(function(){
    el.style.opacity='0';
    el.style.transform='translateY(-8px)';
    setTimeout(function(){
      i=(i+1)%%words.length;
      el.textContent=words[i];
      el.style.opacity='1';
      el.style.transform='translateY(0)';
    },200);
  },%d);
})();`, p.ID, wordsJS, p.Interval)

	style := "display:inline-block;transition:opacity 0.2s ease,transform 0.2s ease"
	if p.Class != "" {
		style += ";" + p.Class
	}

	return g.Group{
		h.Span(
			h.ID(p.ID),
			g.Attr("data-component", "word-rotate"),
			h.Style(style),
			g.Text(p.Words[0]),
		),
		h.Script(g.Raw(script)),
	}
}
