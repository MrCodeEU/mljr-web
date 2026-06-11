package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type TiltCardProps struct {
	// MaxTilt is the max rotation in degrees (default 12).
	MaxTilt float64
	// Scale on hover (default 1.03).
	Scale float64
	// Perspective in px (default 800).
	Perspective int
	// Shine adds a reflective gradient overlay.
	Shine bool
}

// TiltCard wraps content in a 3D perspective tilt on pointer hover.
// Uses pointer events + CSS perspective — no Motion required.
func TiltCard(p TiltCardProps, children ...g.Node) g.Node {
	if p.MaxTilt == 0 {
		p.MaxTilt = 12
	}
	if p.Scale == 0 {
		p.Scale = 1.03
	}
	if p.Perspective == 0 {
		p.Perspective = 800
	}

	script := fmt.Sprintf(`(function(){
  document.querySelectorAll('[data-component="tilt-card"]:not([data-tilt-init])').forEach(function(el){
    el.setAttribute('data-tilt-init','1');
    var shine=el.querySelector('[data-slot="shine"]');
    var maxT=%v,scl=%v,persp=%d;
    el.addEventListener('pointermove',function(e){
      var r=el.getBoundingClientRect();
      var x=(e.clientX-r.left)/r.width-0.5;
      var y=(e.clientY-r.top)/r.height-0.5;
      el.style.transform='perspective('+persp+'px) rotateX('+(-y*maxT*2)+'deg) rotateY('+(x*maxT*2)+'deg) scale('+scl+')';
      if(shine){ shine.style.background='radial-gradient(circle at '+(x*100+50)+'%% '+(y*100+50)+'%%,rgba(255,255,255,0.25) 0%%,rgba(255,255,255,0) 70%%)'; }
    });
    el.addEventListener('pointerleave',function(){
      el.style.transform='perspective('+persp+'px) rotateX(0) rotateY(0) scale(1)';
      if(shine) shine.style.background='none';
    });
  });
})();`, p.MaxTilt, p.Scale, p.Perspective)

	nodes := append([]g.Node{
		g.Attr("data-component", "tilt-card"),
	}, children...)
	if p.Shine {
		nodes = append(nodes, h.Div(g.Attr("data-slot", "shine")))
	}
	nodes = append(nodes, h.Script(g.Raw(script)))

	return h.Div(nodes...)
}
