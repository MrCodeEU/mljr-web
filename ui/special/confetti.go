package special

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ConfettiProps struct {
	// ParticleCount (default 120).
	ParticleCount int
	// Duration in ms (default 3000).
	Duration int
	// Colors: CSS color strings (defaults to accent palette).
	Colors []string
	// Trigger: "click" | "load" | "manual" (default "click").
	Trigger string
	// ButtonLabel shown when Trigger="click" (default "🎉 Confetti").
	ButtonLabel string
}

// Confetti renders a Canvas confetti burst animation.
// Trigger=click wraps a button; Trigger=load fires on page load.
// Call window._confetti() programmatically when Trigger=manual.
func Confetti(p ConfettiProps) g.Node {
	if p.ParticleCount == 0 {
		p.ParticleCount = 120
	}
	if p.Duration == 0 {
		p.Duration = 3000
	}
	if p.Trigger == "" {
		p.Trigger = "click"
	}
	if p.ButtonLabel == "" {
		p.ButtonLabel = "🎉 Celebrate"
	}
	if len(p.Colors) == 0 {
		p.Colors = []string{"#f59e0b", "#ef4444", "#10b981", "#6366f1", "#ec4899", "#f97316", "#14b8a6"}
	}

	colorsJS := "["
	for i, c := range p.Colors {
		if i > 0 {
			colorsJS += ","
		}
		colorsJS += fmt.Sprintf(`'%s'`, c)
	}
	colorsJS += "]"

	script := fmt.Sprintf(`(function(){
  var colors=%s;
  var count=%d;
  var dur=%d;

  function launch(){
    var canvas=document.createElement('canvas');
    canvas.style.cssText='position:fixed;top:0;left:0;width:100vw;height:100vh;pointer-events:none;z-index:9999';
    document.body.appendChild(canvas);
    var ctx=canvas.getContext('2d');
    canvas.width=window.innerWidth;
    canvas.height=window.innerHeight;

    var particles=Array.from({length:count},function(){
      return {
        x:Math.random()*canvas.width,
        y:Math.random()*canvas.height*0.3-canvas.height*0.15,
        vx:(Math.random()-0.5)*12,
        vy:Math.random()*-18-4,
        w:Math.random()*10+5,
        h:Math.random()*6+3,
        angle:Math.random()*360,
        spin:(Math.random()-0.5)*12,
        color:colors[Math.floor(Math.random()*colors.length)],
        opacity:1,
        gravity:0.4+Math.random()*0.3
      };
    });

    var start=performance.now();
    function step(now){
      var elapsed=now-start;
      if(elapsed>dur){ canvas.remove(); return; }
      ctx.clearRect(0,0,canvas.width,canvas.height);
      var progress=elapsed/dur;
      particles.forEach(function(p){
        p.x+=p.vx;
        p.y+=p.vy;
        p.vy+=p.gravity;
        p.vx*=0.99;
        p.angle+=p.spin;
        p.opacity=Math.max(0,1-progress*1.2);
        ctx.save();
        ctx.translate(p.x,p.y);
        ctx.rotate(p.angle*Math.PI/180);
        ctx.globalAlpha=p.opacity;
        ctx.fillStyle=p.color;
        ctx.fillRect(-p.w/2,-p.h/2,p.w,p.h);
        ctx.restore();
      });
      requestAnimationFrame(step);
    }
    requestAnimationFrame(step);
  }

  window._confetti=launch;
  %s
})();`,
		colorsJS, p.ParticleCount, p.Duration,
		func() string {
			switch p.Trigger {
			case "load":
				return "window.addEventListener('load',launch);"
			case "click":
				return `document.querySelectorAll('[data-confetti-trigger]').forEach(function(b){ b.addEventListener('click',launch); });`
			default:
				return ""
			}
		}(),
	)

	var trigger g.Node
	if p.Trigger == "click" {
		trigger = h.Button(
			h.Type("button"),
			g.Attr("data-component", "button"),
			g.Attr("data-variant", "primary"),
			g.Attr("data-confetti-trigger", ""),
			g.Text(p.ButtonLabel),
		)
	}

	return g.Group{
		trigger,
		h.Script(g.Raw(script)),
	}
}
