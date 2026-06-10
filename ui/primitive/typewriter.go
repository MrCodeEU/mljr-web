package primitive

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TypewriterProps struct {
	// Lines is the sequence of text to type out, cycling on completion.
	Lines []string
	// Speed is ms per character typed (default 60).
	Speed int
	// DeleteSpeed is ms per character deleted (default 30).
	DeleteSpeed int
	// Pause is ms to wait at end of a line before deleting (default 1800).
	Pause int
	// NoCursor disables the blinking cursor (cursor shown by default).
	NoCursor bool
	// ID must be unique per page (default "tw").
	ID string
}

// Typewriter animates a type-and-delete loop across multiple phrases.
// No dependencies — pure JS setInterval.
func Typewriter(p TypewriterProps) g.Node {
	if p.Speed == 0 {
		p.Speed = 60
	}
	if p.DeleteSpeed == 0 {
		p.DeleteSpeed = 30
	}
	if p.Pause == 0 {
		p.Pause = 1800
	}
	if p.ID == "" {
		p.ID = "tw"
	}
	if len(p.Lines) == 0 {
		p.Lines = []string{"Hello, World!"}
	}

	escaped := make([]string, len(p.Lines))
	for i, l := range p.Lines {
		escaped[i] = "'" + strings.ReplaceAll(strings.ReplaceAll(l, "\\", "\\\\"), "'", "\\'") + "'"
	}
	linesJS := "[" + strings.Join(escaped, ",") + "]"

	script := fmt.Sprintf(`(function(){
  var el=document.getElementById('%s');
  if(!el) return;
  var lines=%s,li=0,ci=0,deleting=false;
  function tick(){
    var line=lines[li];
    if(!deleting){
      ci++;
      el.textContent=line.slice(0,ci);
      if(ci>=line.length){ deleting=true; setTimeout(tick,%d); return; }
      setTimeout(tick,%d);
    } else {
      ci--;
      el.textContent=line.slice(0,ci);
      if(ci<=0){ deleting=false; li=(li+1)%%lines.length; setTimeout(tick,300); return; }
      setTimeout(tick,%d);
    }
  }
  setTimeout(tick,%d);
})();`, p.ID, linesJS, p.Pause, p.Speed, p.DeleteSpeed, p.Speed)

	return g.Group{
		h.Span(
			h.ID(p.ID),
			g.Attr("data-component", "typewriter"),
			g.If(!p.NoCursor, g.Attr("data-cursor", "true")),
		),
		h.Script(g.Raw(script)),
	}
}
