package special

import (
	"fmt"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// LogoScatterProps configures the mljr.eu logo scatter animation.
type LogoScatterProps struct {
	// ID is the SVG element id and gradient id prefix. Required.
	ID string
	// SVGStyle is the full inline style for the SVG element.
	// Defaults to centered absolute position at Size.
	SVGStyle string
	// Size sets CSS width/height when SVGStyle is empty (default "280px").
	Size string
	// InitialOpacity is the assembled-state SVG opacity (default 0.65).
	InitialOpacity float64
	// Mode is "loop" (auto scatter/assemble cycle, default) or "scroll"
	// (IntersectionObserver-triggered).
	Mode string
	// TriggerID is the element to observe in scroll mode.
	TriggerID string
	// WithBackground wraps the SVG in a full-page absolute positioned bg div.
	WithBackground bool
	// BackgroundStyle overrides the wrapper style when WithBackground is true.
	BackgroundStyle string
	// PieceCopies controls how many copies of each SVG path are animated.
	// Defaults to 2, which means original + one clone.
	PieceCopies int
	// WrapInLoad wraps the script in window.addEventListener('load', ...).
	WrapInLoad bool
}

// LogoScatter renders the mljr.eu logo scatter animation.
// Motion v10 (window.Motion) must be loaded before this script executes.
func LogoScatter(p LogoScatterProps) g.Node {
	if p.ID == "" {
		p.ID = "logo-svg"
	}
	if p.Size == "" {
		p.Size = "280px"
	}
	if p.InitialOpacity == 0 {
		p.InitialOpacity = 0.65
	}
	if p.SVGStyle == "" {
		p.SVGStyle = fmt.Sprintf(
			"position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);overflow:visible;width:%s;height:%s",
			p.Size, p.Size,
		)
	}
	if p.PieceCopies <= 0 {
		p.PieceCopies = 2
	}

	svgHTML := buildLogoSVG(p.ID, p.SVGStyle)
	script := buildScatterScript(p)
	if p.WrapInLoad {
		script = "window.addEventListener('load',function(){\n" + script + "\n});"
	}

	if p.WithBackground {
		bgStyle := p.BackgroundStyle
		if bgStyle == "" {
			bgStyle = "position:absolute;inset:0;pointer-events:none;overflow:hidden;z-index:0"
		}
		return g.Group{
			h.Div(
				h.Style(bgStyle),
				h.ID(p.ID+"-bg"),
				g.Raw(svgHTML),
			),
			h.Script(g.Raw(script)),
		}
	}
	return g.Group{
		g.Raw(svgHTML),
		h.Script(g.Raw(script)),
	}
}

func buildLogoSVG(id, style string) string {
	r := strings.NewReplacer("__ID__", id, "__STYLE__", style, "__P__", id+"-")
	return r.Replace(logoSVGTpl)
}

// logoSVGTpl is the mljr.eu logo SVG template.
// __ID__ → SVG element id, __STYLE__ → inline style, __P__ → gradient id prefix.
const logoSVGTpl = `<svg id="__ID__"
   style="__STYLE__"
   viewBox="0 0 2666.6667 2666.6667" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(2582.0449,-2891.0076,-2891.0076,-2582.0449,11.75041,1476.6624)" id="__P__lg4"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(2582.0449,-2891.0076,-2891.0076,-2582.0449,196.49147,1641.6602)" id="__P__lg8"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(2582.0449,-2891.0076,-2891.0076,-2582.0449,-19.747679,1448.5305)" id="__P__lg12"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(2582.0449,-2891.0076,-2891.0076,-2582.0449,165.89384,1614.3325)" id="__P__lg16"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(2582.0449,-2891.0076,-2891.0076,-2582.0449,-63.247246,1409.6798)" id="__P__lg20"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(4629.0181,-5182.9175,-5182.9175,-4629.0181,-215.09644,2602.1555)" id="__P__lg24"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(4629.0181,-5182.9175,-5182.9175,-4629.0181,-134.03297,2532.7686)" id="__P__lg28"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(4629.0181,-5182.9175,-5182.9175,-4629.0181,-289.07684,2632.8066)" id="__P__lg32"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(1603.6569,-1795.5474,-1795.5474,-1603.6569,904.22754,1593.01)" id="__P__lg36"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
    <linearGradient x1="0" y1="0" x2="1" y2="0" gradientUnits="userSpaceOnUse" gradientTransform="matrix(4248.605,-4756.9849,-4756.9849,-4248.605,53.146076,2777.4116)" id="__P__lg40"><stop style="stop-color:#f28d1d" offset="0"/><stop style="stop-color:#a40054" offset="0.63"/><stop style="stop-color:#a40054" offset="1"/></linearGradient>
  </defs>
  <g id="__P__layer">
    <path d="M 450.428,1162.83 V 755.809 c 0,-12.912 10.474,-23.375 23.372,-23.375 c 12.912,0 23.378,10.463 23.378,23.375 v 407.021 c 0,12.917 -10.466,23.374 -23.378,23.374 c -12.898,0 -23.372,-10.457 -23.372,-23.374" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg4);stroke:none"/>
    <path d="m 450.428,1425.996 v -189.662 c 0,-12.912 10.474,-23.375 23.372,-23.375 c 12.912,0 23.378,10.463 23.378,23.375 v 189.662 c 0,12.92 -10.466,23.369 -23.378,23.369 c -12.898,0 -23.372,-10.449 -23.372,-23.369" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg8);stroke:none"/>
    <path d="M 526.082,937.601 V 684.822 c 0,-12.899 10.477,-23.375 23.379,-23.375 c 12.918,0 23.374,10.476 23.374,23.375 v 252.779 c 0,12.898 -10.456,23.368 -23.374,23.368 c -12.902,0 -23.379,-10.47 -23.379,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg12);stroke:none"/>
    <path d="m 526.082,1355.016 v -340.298 c 0,-12.906 10.477,-23.374 23.379,-23.374 c 12.918,0 23.374,10.468 23.374,23.374 v 340.298 c 0,12.911 -10.456,23.367 -23.374,23.367 c -12.902,0 -23.379,-10.456 -23.379,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg16);stroke:none"/>
    <path d="m 601.757,664.046 v -50.219 c 0,-12.898 10.456,-23.375 23.371,-23.375 c 12.913,0 23.372,10.477 23.372,23.375 v 50.219 c 0,12.898 -10.459,23.368 -23.372,23.368 c -12.915,0 -23.371,-10.47 -23.371,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg20);stroke:none"/>
    <path d="M 1350.146,937.601 V 709.308 c 0,-12.898 10.455,-23.375 23.371,-23.375 c 12.911,0 23.37,10.477 23.37,23.375 v 228.293 c 0,12.898 -10.459,23.368 -23.37,23.368 c -12.916,0 -23.371,-10.47 -23.371,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg24);stroke:none"/>
    <path d="M 1431.209,868.214 V 639.922 c 0,-12.899 10.456,-23.375 23.371,-23.375 c 12.912,0 23.37,10.476 23.37,23.375 v 228.292 c 0,12.899 -10.458,23.367 -23.37,23.367 c -12.915,0 -23.371,-10.468 -23.371,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg28);stroke:none"/>
    <path d="M 1276.165,968.252 V 739.959 c 0,-12.899 10.456,-23.375 23.371,-23.375 c 12.912,0 23.37,10.476 23.37,23.375 v 228.293 c 0,12.898 -10.458,23.367 -23.37,23.367 c -12.915,0 -23.371,-10.469 -23.371,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg32);stroke:none"/>
    <path d="m 1431.209,995.871 v -38.13 c 0,-12.899 10.456,-23.375 23.371,-23.375 c 12.912,0 23.37,10.476 23.37,23.375 v 38.13 c 0,12.898 -10.458,23.367 -23.37,23.367 c -12.915,0 -23.371,-10.469 -23.371,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg36);stroke:none"/>
    <path d="m 1502.826,1241.482 v -226.763 c 0,-12.899 10.456,-23.375 23.372,-23.375 c 12.911,0 23.37,10.476 23.37,23.375 v 226.763 c 0,12.898 -10.459,23.368 -23.37,23.368 c -12.916,0 -23.372,-10.47 -23.372,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg40);stroke:none"/>
    <path d="M 601.757,1067.432 V 735.02 c 0,-12.906 10.456,-23.375 23.371,-23.375 c 12.913,0 23.372,10.469 23.372,23.375 v 332.412 c 0,12.92 -10.459,23.382 -23.372,23.382 c -12.915,0 -23.371,-10.462 -23.371,-23.382" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg4);stroke:none"/>
    <path d="m 601.757,1284.035 v -147.336 c 0,-12.911 10.456,-23.381 23.371,-23.381 c 12.913,0 23.372,10.47 23.372,23.381 v 147.336 c 0,12.897 -10.459,23.368 -23.372,23.368 c -12.915,0 -23.371,-10.471 -23.371,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg8);stroke:none"/>
    <path d="M 677.411,1213.033 V 989.2 c 0,-12.918 10.46,-23.381 23.375,-23.381 c 12.908,0 23.368,10.463 23.368,23.381 v 223.833 c 0,12.913 -10.46,23.383 -23.368,23.383 c -12.915,0 -23.375,-10.47 -23.375,-23.383" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg12);stroke:none"/>
    <path d="M 753.079,1142.047 V 918.2 c 0,-12.905 10.467,-23.368 23.382,-23.368 c 12.898,0 23.367,10.463 23.367,23.368 v 223.847 c 0,12.911 -10.469,23.38 -23.367,23.38 c -12.915,0 -23.382,-10.469 -23.382,-23.38" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg16);stroke:none"/>
    <path d="M 828.744,1071.066 V 847.211 c 0,-12.898 10.463,-23.374 23.371,-23.374 c 12.895,0 23.385,10.476 23.385,23.374 v 223.855 c 0,12.918 -10.49,23.381 -23.385,23.381 c -12.908,0 -23.371,-10.463 -23.371,-23.381" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg20);stroke:none"/>
    <path d="m 1276.165,1399.759 v -223.834 c 0,-12.918 10.459,-23.381 23.375,-23.381 c 12.908,0 23.366,10.463 23.366,23.381 v 223.834 c 0,12.912 -10.458,23.381 -23.366,23.381 c -12.916,0 -23.375,-10.469 -23.375,-23.381" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg24);stroke:none"/>
    <path d="m 1350.141,1355.002 v -223.847 c 0,-12.905 10.467,-23.368 23.382,-23.368 c 12.898,0 23.368,10.463 23.368,23.368 v 223.847 c 0,12.912 -10.47,23.381 -23.368,23.381 c -12.915,0 -23.382,-10.469 -23.382,-23.381" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg28);stroke:none"/>
    <path d="m 1431.202,1307.086 v -223.853 c 0,-12.898 10.462,-23.375 23.37,-23.375 c 12.895,0 23.386,10.477 23.386,23.375 v 223.853 c 0,12.919 -10.491,23.381 -23.386,23.381 c -12.908,0 -23.37,-10.462 -23.37,-23.381" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg32);stroke:none"/>
    <path d="M 1502.822,797.85 V 574.003 c 0,-12.905 10.467,-23.368 23.382,-23.368 c 12.898,0 23.367,10.463 23.367,23.368 v 223.847 c 0,12.911 -10.469,23.381 -23.367,23.381 c -12.915,0 -23.382,-10.47 -23.382,-23.381" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg36);stroke:none"/>
    <path d="M 1207.029,883.912 V 657.867 c 0,-12.912 10.449,-23.374 23.388,-23.374 c 12.898,0 23.36,10.462 23.36,23.374 v 226.045 c 0,12.906 -10.462,23.376 -23.36,23.376 c -12.939,0 -23.388,-10.47 -23.388,-23.376" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg40);stroke:none"/>
    <path d="m 1207.029,1425.996 v -461.66 c 0,-12.898 10.449,-23.374 23.388,-23.374 c 12.898,0 23.36,10.476 23.36,23.374 v 461.66 c 0,12.92 -10.462,23.369 -23.36,23.369 c -12.939,0 -23.388,-10.449 -23.388,-23.369" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg4);stroke:none"/>
    <path d="m 1131.367,710.42 v -25.598 c 0,-12.899 10.463,-23.375 23.362,-23.375 c 12.925,0 23.388,10.476 23.388,23.375 v 25.598 c 0,12.906 -10.463,23.375 -23.388,23.375 c -12.899,0 -23.362,-10.469 -23.362,-23.375" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg8);stroke:none"/>
    <path d="M 1131.367,1114.176 V 777.938 c 0,-12.925 10.463,-23.374 23.362,-23.374 c 12.925,0 23.388,10.449 23.388,23.374 v 336.238 c 0,12.897 -10.463,23.367 -23.388,23.367 c -12.899,0 -23.362,-10.47 -23.362,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg12);stroke:none"/>
    <path d="m 1131.367,1355.016 v -169.854 c 0,-12.918 10.463,-23.367 23.362,-23.367 c 12.925,0 23.388,10.449 23.388,23.367 v 169.854 c 0,12.911 -10.463,23.367 -23.388,23.367 c -12.899,0 -23.362,-10.456 -23.362,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg16);stroke:none"/>
    <path d="M 1055.733,894.301 V 613.827 c 0,-12.898 10.449,-23.375 23.362,-23.375 c 12.884,0 23.361,10.477 23.361,23.375 v 280.474 c 0,12.912 -10.477,23.375 -23.361,23.375 c -12.913,0 -23.362,-10.463 -23.362,-23.375" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg20);stroke:none"/>
    <path d="M 1055.733,1284.035 V 965.262 c 0,-12.913 10.449,-23.375 23.362,-23.375 c 12.884,0 23.361,10.462 23.361,23.375 v 318.773 c 0,12.897 -10.477,23.368 -23.361,23.368 c -12.913,0 -23.362,-10.471 -23.362,-23.368" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg24);stroke:none"/>
    <path d="M 980.046,1212.871 V 989.01 c 0,-12.905 10.476,-23.375 23.387,-23.375 c 12.899,0 23.375,10.47 23.375,23.375 v 223.861 c 0,12.912 -10.476,23.367 -23.375,23.367 c -12.911,0 -23.387,-10.455 -23.387,-23.367" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg28);stroke:none"/>
    <path d="M 904.398,1141.876 V 918.029 c 0,-12.918 10.463,-23.374 23.375,-23.374 c 12.884,0 23.374,10.456 23.374,23.374 v 223.847 c 0,12.912 -10.49,23.382 -23.374,23.382 c -12.912,0 -23.375,-10.47 -23.375,-23.382" transform="matrix(1.3333333,0,0,-1.3333333,0,2666.6667)" style="fill:url(#__P__lg32);stroke:none"/>
  </g>
</svg>`

func buildScatterScript(p LogoScatterProps) string {
	isScroll := p.Mode == "scroll"
	bgID := p.ID + "-bg"
	initOp := p.InitialOpacity
	copies := p.PieceCopies

	var b strings.Builder

	b.WriteString("(function(){\n")
	fmt.Fprintf(&b, "  var svg=document.getElementById('%s');\n", p.ID)
	if isScroll {
		fmt.Fprintf(&b, "  var bg=document.getElementById('%s');\n", bgID)
		b.WriteString("  if(!svg||!bg||typeof Motion==='undefined') return;\n")
	} else {
		b.WriteString("  if(!svg||typeof Motion==='undefined') return;\n")
	}

	fmt.Fprintf(&b, "  var pieceCopies=%d;\n", copies)
	b.WriteString(`  Array.from(svg.querySelectorAll('path')).forEach(function(p){
    for(var copy=1;copy<pieceCopies;copy++){
      var cl=p.cloneNode(true); cl.removeAttribute('id');
      p.parentNode.appendChild(cl);
    }
  });
  function finite(n,fallback){return Number.isFinite(n)?n:fallback;}
  var vbW=finite(svg.viewBox.baseVal.width,0);
  var pieces=Array.from(svg.querySelectorAll('path')).map(function(path){
    var gEl=document.createElementNS('http://www.w3.org/2000/svg','g');
    gEl.style.transformBox='fill-box';
    gEl.style.transformOrigin='center';
    path.parentNode.insertBefore(gEl,path);
    gEl.appendChild(path);
    return {g:gEl,fr:(Math.random()-0.5)*55,fs:1.3+Math.random()*1.3,
            fx:0,fy:0,sx:0,sy:0,sr:(Math.random()-0.5)*7,
            dur:1,sd:4+Math.random()*5,sdl:Math.random()*2};
  });
  function pieceReady(p){
    return Number.isFinite(p.fx)&&Number.isFinite(p.fy)&&Number.isFinite(p.sx)&&Number.isFinite(p.sy)&&
           Number.isFinite(p.fr)&&Number.isFinite(p.fs)&&Number.isFinite(p.sr)&&Number.isFinite(p.dur)&&
           Number.isFinite(p.sd)&&Number.isFinite(p.sdl);
  }
  function coordsReady(){return pieces.length>0&&pieces.every(pieceReady);}
`)

	b.WriteString("  function recalc(){\n")
	b.WriteString("    if(!Number.isFinite(vbW)||vbW<=0) return false;\n")
	if isScroll {
		b.WriteString("    var svgR=svg.getBoundingClientRect();\n")
		b.WriteString("    if(!svgR||!Number.isFinite(svgR.width)||!Number.isFinite(svgR.height)||svgR.width<1||svgR.height<1) return false;\n")
		b.WriteString("    var svgScale=finite(vbW/svgR.width,1);\n")
		b.WriteString("    var cx=svgR.left+window.scrollX+svgR.width/2;\n")
		b.WriteString("    var cy=svgR.top+window.scrollY+svgR.height/2;\n")
		b.WriteString("    var dw=document.documentElement.scrollWidth;\n")
		b.WriteString("    var dh=document.documentElement.scrollHeight;\n")
	} else {
		b.WriteString("    var svgR=svg.getBoundingClientRect();\n")
		b.WriteString("    if(!svgR||!Number.isFinite(svgR.width)||!Number.isFinite(svgR.height)||svgR.width<1||svgR.height<1) return false;\n")
		b.WriteString("    var svgScale=finite(vbW/svgR.width,1);\n")
		b.WriteString("    var cx=svgR.left+svgR.width/2;\n")
		b.WriteString("    var cy=svgR.top+svgR.height/2;\n")
		b.WriteString("    var dw=window.innerWidth;\n")
		b.WriteString("    var dh=window.innerHeight;\n")
	}
	b.WriteString(`    cx=finite(cx,0); cy=finite(cy,0);
    dw=Math.max(1,finite(dw,window.innerWidth||1));
    dh=Math.max(1,finite(dh,window.innerHeight||1));
    pieces.forEach(function(p){
      var roll=Math.random();
      var tx=roll<0.4?dw*0.04+Math.random()*dw*0.22
             :roll<0.6?dw*0.24+Math.random()*dw*0.52
             :          dw*0.72+Math.random()*dw*0.22;
      var ty=dh*0.05+Math.random()*dh*0.90;
      p.fx=finite((tx-cx)*svgScale,0);
      p.fy=finite((ty-cy)*svgScale,0);
      p.sx=finite((Math.random()-0.5)*60*svgScale,0);
      p.sy=finite((Math.random()-0.5)*40*svgScale,0);
      p.dur=finite(Math.max(0.7,Math.min(3.0,Math.hypot(tx-cx,ty-cy)/600)),1);
    });
    return coordsReady();
  }
`)

	b.WriteString("  var swAnims=[],scAnims=[],swayTimers=[],scattered=false")
	if !isScroll {
		b.WriteString(",loopTimer=null")
	}
	b.WriteString(";\n")
	b.WriteString("  function tf(x,y,r,s){x=finite(x,0);y=finite(y,0);r=finite(r,0);s=finite(s,1);return 'translateX('+x+'px) translateY('+y+'px) rotate('+r+'deg) scale('+s+')'}\n")
	b.WriteString("  function stopAll(){\n")
	if !isScroll {
		b.WriteString("    clearTimeout(loopTimer);\n")
	}
	b.WriteString(`    swAnims.forEach(function(a){try{a.stop();}catch(e){}});swAnims=[];
    scAnims.forEach(function(a){try{a.stop();}catch(e){}});scAnims=[];
    swayTimers.forEach(function(t){clearTimeout(t);});swayTimers=[];
  }
  function sway(p){
    if(!pieceReady(p)) return;
    swAnims.push(Motion.animate(p.g,
      {transform:[tf(p.fx,p.fy,p.fr,p.fs),tf(p.fx+p.sx,p.fy+p.sy,p.fr+p.sr,p.fs)]},
      {duration:p.sd,easing:'ease-in-out',repeat:Infinity,direction:'alternate',delay:p.sdl}
    ));
  }
`)

	if isScroll {
		b.WriteString(`  function scatter(){
    if(scattered) return;
    if(!coordsReady()&&!recalc()) return;
    if(!coordsReady()) return;
    scattered=true; stopAll();
`)
	} else {
		b.WriteString(`  function scatter(){
    if(!coordsReady()&&!recalc()) return;
    if(!coordsReady()) return;
    stopAll(); scattered=true;
`)
	}
	b.WriteString(`    Motion.animate(svg,{opacity:1},{duration:0.5});
    pieces.forEach(function(p){
      scAnims.push(Motion.animate(p.g,
        {transform:tf(p.fx,p.fy,p.fr,p.fs)},
        {duration:p.dur,easing:'cubic-bezier(0.65,0,0.35,1)'}
      ));
      swayTimers.push(setTimeout(function(){sway(p);},p.dur*1000+100));
    });
`)
	if !isScroll {
		b.WriteString("    loopTimer=setTimeout(function(){assemble();},7000);\n")
	}
	b.WriteString("  }\n")

	if isScroll {
		b.WriteString(`  function assemble(){
    if(!scattered) return; scattered=false; stopAll();
`)
	} else {
		b.WriteString(`  function assemble(){
    stopAll(); scattered=false;
`)
	}
	fmt.Fprintf(&b, "    Motion.animate(svg,{opacity:%v},{duration:0.7});\n", initOp)
	b.WriteString(`    pieces.forEach(function(p){
      scAnims.push(Motion.animate(p.g,
        {transform:tf(0,0,0,1)},
        {duration:1.1,easing:'cubic-bezier(0.65,0,0.35,1)'}
      ));
    });
`)
	if !isScroll {
		b.WriteString("    loopTimer=setTimeout(function(){recalc();scatter();},3500);\n")
	}
	b.WriteString("  }\n")

	b.WriteString(`  var resizeTimer;
  window.addEventListener('resize',function(){
    clearTimeout(resizeTimer);
    resizeTimer=setTimeout(function(){
`)
	if isScroll {
		b.WriteString(`      var wasScattered=scattered;
      if(wasScattered){scattered=false;stopAll();}
      if(recalc()&&wasScattered){scatter();}
`)
	} else {
		b.WriteString(`      recalc();
      if(scattered){scatter();}
`)
	}
	b.WriteString("    },250);\n  });\n")

	if isScroll {
		if p.TriggerID != "" {
			fmt.Fprintf(&b, "  var triggerEl=document.getElementById('%s');\n", p.TriggerID)
			b.WriteString(`  if(triggerEl){
    new IntersectionObserver(function(entries){
      if(entries[0].isIntersecting){assemble();}
      else if(window.scrollY>80){scatter();}
    },{threshold:0.9}).observe(triggerEl);
  }
`)
		}
		b.WriteString("  recalc();\n")
	} else {
		fmt.Fprintf(&b, "  Motion.animate(svg,{opacity:%v},{duration:0});\n", initOp)
		b.WriteString("  recalc();\n")
		b.WriteString("  setTimeout(function(){scatter();},700);\n")
	}

	b.WriteString("})();")
	return b.String()
}
