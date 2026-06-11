package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AudioPlayerProps struct {
	Src    string
	Title  string
	Artist string
	Cover  string // optional cover image URL
	// Signal prefix (default "_ap"). Must be unique per page.
	Signal string
}

// AudioPlayer renders a custom-styled HTML5 audio player with Canvas waveform animation.
func AudioPlayer(p AudioPlayerProps) g.Node {
	sig := p.Signal
	if sig == "" {
		sig = "_ap"
	}
	aid := sig + "El"
	cid := sig + "Cv"

	script := fmt.Sprintf(`(function(){
  var audio=document.getElementById('%s');
  var canvas=document.getElementById('%s');
  var root=audio.closest('[data-component="audio-player"]');
  if(!audio||!canvas) return;

  // Waveform bars animation
  var ctx=canvas.getContext('2d');
  var bars=40;
  var playing=false;
  var frame;
  var heights=Array.from({length:bars},function(_,i){ return 0.15+0.6*Math.abs(Math.sin(i*0.4)); });

  function draw(){
    canvas.width=canvas.offsetWidth*window.devicePixelRatio||300;
    canvas.height=canvas.offsetHeight*window.devicePixelRatio||48;
    ctx.clearRect(0,0,canvas.width,canvas.height);
    // Canvas fillStyle cannot resolve CSS custom properties — read them off the root element.
    var cs=getComputedStyle(root);
    var accent=cs.getPropertyValue('--accent').trim()||'#6366f1';
    var lineCol=cs.getPropertyValue('--line').trim()||'#d4d4d8';
    var w=canvas.width,h=canvas.height,bw=w/bars*0.6,gap=w/bars*0.4;
    var pct=audio.duration?audio.currentTime/audio.duration:0;
    for(var i=0;i<bars;i++){
      var x=i*(bw+gap);
      var amp=heights[i];
      if(playing){ amp=heights[i]=0.15+0.65*Math.abs(Math.sin(i*0.4+Date.now()*0.003+i*0.1)); }
      var bh=h*amp;
      ctx.fillStyle=(i/bars)<pct?accent:lineCol;
      ctx.fillRect(x,(h-bh)/2,bw,bh);
    }
    if(playing) frame=requestAnimationFrame(draw);
  }
  draw();

  // Time format
  function fmt(s){ var m=Math.floor(s/60),sec=Math.floor(s%%60); return m+':'+(sec<10?'0':'')+sec; }

  // Event listeners
  audio.addEventListener('play',function(){ playing=true; root.setAttribute('data-playing','true'); draw(); });
  audio.addEventListener('pause',function(){ playing=false; root.setAttribute('data-playing',''); cancelAnimationFrame(frame); draw(); });
  audio.addEventListener('ended',function(){ playing=false; root.setAttribute('data-playing',''); cancelAnimationFrame(frame); draw(); });
  audio.addEventListener('timeupdate',function(){
    document.getElementById('%sCur').textContent=fmt(audio.currentTime);
    if(!playing) draw();
  });
  audio.addEventListener('loadedmetadata',function(){ document.getElementById('%sDur').textContent=fmt(audio.duration); });

  // Controls
  root.querySelector('[data-slot="play-btn"]').addEventListener('click',function(){ audio.paused?audio.play():audio.pause(); });
  var bar=root.querySelector('[data-slot="progress"]');
  bar.addEventListener('input',function(){ audio.currentTime=(bar.value/1000)*(audio.duration||0); draw(); });
  audio.addEventListener('timeupdate',function(){ bar.value=audio.duration?Math.round(audio.currentTime/audio.duration*1000):0; });
  root.querySelector('[data-slot="mute-btn"]').addEventListener('click',function(){ audio.muted=!audio.muted; root.setAttribute('data-muted',audio.muted?'true':''); });
})();`, aid, cid, sig, sig)

	return h.Div(
		g.Attr("data-component", "audio-player"),
		// Hidden native audio element
		g.El("audio", h.ID(aid), g.Attr("data-slot", "audio"), h.Src(p.Src), h.Preload("metadata")),
		// Cover / icon
		g.If(p.Cover != "", h.Img(h.Src(p.Cover), g.Attr("data-slot", "cover"), h.Alt(p.Title))),
		// Track info + canvas waveform
		h.Div(g.Attr("data-slot", "body"),
			h.Div(g.Attr("data-slot", "info"),
				g.If(p.Title != "", h.Div(g.Attr("data-slot", "title"), g.Text(p.Title))),
				g.If(p.Artist != "", h.Div(g.Attr("data-slot", "artist"), g.Text(p.Artist))),
			),
			h.Canvas(h.ID(cid), g.Attr("data-slot", "waveform"), h.Style("width:100%;height:48px;display:block;cursor:pointer")),
		),
		// Controls row
		h.Div(g.Attr("data-slot", "controls"),
			h.Button(h.Type("button"), g.Attr("data-slot", "play-btn"), g.Attr("aria-label", "Play/Pause"),
				g.Raw(`<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>`),
			),
			h.Span(h.ID(sig+"Cur"), g.Attr("data-slot", "time"), g.Text("0:00")),
			h.Input(h.Type("range"), g.Attr("data-slot", "progress"), h.Min("0"), h.Max("1000"), h.Value("0")),
			h.Span(h.ID(sig+"Dur"), g.Attr("data-slot", "time"), g.Text("0:00")),
			h.Button(h.Type("button"), g.Attr("data-slot", "mute-btn"), g.Attr("aria-label", "Mute"),
				g.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07"/></svg>`),
			),
		),
		h.Script(g.Raw(script)),
	)
}
