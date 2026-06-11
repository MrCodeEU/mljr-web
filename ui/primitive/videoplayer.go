package primitive

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type VideoPlayerProps struct {
	// Src is the video URL (required).
	Src string
	// Poster is the thumbnail shown before play.
	Poster string
	// AspectRatio: CSS aspect-ratio (default "16/9").
	AspectRatio string
	// Autoplay starts muted on load.
	Autoplay bool
	// Loop repeats the video.
	Loop bool
	// Signal is the Datastar signal prefix (default "_vp"). Must be unique per page.
	Signal string
}

// VideoPlayer renders a custom-styled HTML5 video player.
// Native <video> handles codecs; custom controls built with Datastar signals + JS bridge.
func VideoPlayer(p VideoPlayerProps) g.Node {
	if p.AspectRatio == "" {
		p.AspectRatio = "16/9"
	}
	sig := p.Signal
	if sig == "" {
		sig = "_vp"
	}

	// Unique video element id based on signal prefix
	vid := sig + "El"

	// JS bridge: native video state ↔ DOM attributes + control wiring
	bridge := fmt.Sprintf(`(function(){
  var v=document.getElementById('%s');
  if(!v) return;
  v.addEventListener('play',function(){ v.closest('[data-component="video-player"]').setAttribute('data-playing','true'); });
  v.addEventListener('pause',function(){ v.closest('[data-component="video-player"]').setAttribute('data-playing',''); });
  v.addEventListener('ended',function(){ v.closest('[data-component="video-player"]').setAttribute('data-playing',''); });
  v.addEventListener('loadedmetadata',function(){ document.getElementById('%sDur').textContent=fmt(v.duration); });
  v.addEventListener('timeupdate',function(){
    var pct=v.duration?v.currentTime/v.duration:0;
    document.getElementById('%sBar').value=Math.round(pct*1000);
    document.getElementById('%sCur').textContent=fmt(v.currentTime);
  });
  function fmt(s){ var m=Math.floor(s/60),sec=Math.floor(s%%60); return m+':'+(sec<10?'0':'')+sec; }
  %s
  // Controls wiring
  var root=v.closest('[data-component="video-player"]');
  root.querySelector('[data-slot="play-btn"]').addEventListener('click',function(){ v.paused?v.play():v.pause(); });
  var bar=document.getElementById('%sBar');
  bar.addEventListener('input',function(){ v.currentTime=(bar.value/1000)*(v.duration||0); });
  root.querySelector('[data-slot="mute-btn"]').addEventListener('click',function(){ v.muted=!v.muted; root.setAttribute('data-muted',v.muted?'true':''); });
  var vol=root.querySelector('[data-slot="volume"]');
  vol.addEventListener('input',function(){ v.volume=vol.value/100; v.muted=false; });
  root.querySelector('[data-slot="fs-btn"]').addEventListener('click',function(){ (v.requestFullscreen||v.webkitRequestFullscreen).call(v); });
  var spd=root.querySelector('[data-slot="speed"]');
  spd.addEventListener('change',function(){ v.playbackRate=parseFloat(spd.value); });
})();`,
		vid, vid, vid, vid,
		func() string {
			if p.Autoplay {
				return "v.muted=true;v.play().catch(function(){});"
			}
			return ""
		}(),
		vid,
	)

	videoAttrs := []g.Node{
		h.ID(vid),
		g.Attr("data-slot", "video"),
		h.Style(fmt.Sprintf("width:100%%;aspect-ratio:%s;display:block;background:#000;border-radius:var(--radius) var(--radius) 0 0", p.AspectRatio)),
	}
	if p.Poster != "" {
		videoAttrs = append(videoAttrs, h.Poster(p.Poster))
	}
	if p.Loop {
		videoAttrs = append(videoAttrs, g.Attr("loop", ""))
	}
	videoAttrs = append(videoAttrs, h.Src(p.Src))

	return h.Div(
		g.Attr("data-component", "video-player"),
		g.El("video", videoAttrs...),
		h.Div(
			g.Attr("data-slot", "controls"),
			// Play/Pause
			h.Button(h.Type("button"), g.Attr("data-slot", "play-btn"), g.Attr("aria-label", "Play/Pause"),
				g.Raw(`<svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M8 5v14l11-7z"/></svg>`),
			),
			// Current time
			h.Span(h.ID(vid+"Cur"), g.Attr("data-slot", "time"), g.Text("0:00")),
			// Progress bar
			h.Input(h.Type("range"), h.ID(vid+"Bar"), g.Attr("data-slot", "progress"),
				h.Min("0"), h.Max("1000"), h.Value("0"),
			),
			// Duration
			h.Span(h.ID(vid+"Dur"), g.Attr("data-slot", "time"), g.Text("0:00")),
			// Mute
			h.Button(h.Type("button"), g.Attr("data-slot", "mute-btn"), g.Attr("aria-label", "Mute"),
				g.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07"/></svg>`),
			),
			// Volume
			h.Input(h.Type("range"), g.Attr("data-slot", "volume"), g.Attr("aria-label", "Volume"),
				h.Min("0"), h.Max("100"), h.Value("100"),
			),
			// Speed
			g.El("select", g.Attr("data-slot", "speed"), g.Attr("aria-label", "Playback speed"),
				h.Option(g.Attr("value", "0.5"), g.Text("0.5×")),
				h.Option(g.Attr("value", "0.75"), g.Text("0.75×")),
				h.Option(g.Attr("value", "1"), g.Attr("selected", ""), g.Text("1×")),
				h.Option(g.Attr("value", "1.25"), g.Text("1.25×")),
				h.Option(g.Attr("value", "1.5"), g.Text("1.5×")),
				h.Option(g.Attr("value", "2"), g.Text("2×")),
			),
			// Fullscreen
			h.Button(h.Type("button"), g.Attr("data-slot", "fs-btn"), g.Attr("aria-label", "Fullscreen"),
				g.Raw(`<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 3 21 3 21 9"/><polyline points="9 21 3 21 3 15"/><line x1="21" y1="3" x2="14" y2="10"/><line x1="3" y1="21" x2="10" y2="14"/></svg>`),
			),
		),
		h.Script(g.Raw(bridge)),
	)
}
