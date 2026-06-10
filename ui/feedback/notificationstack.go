package feedback

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type NotificationStackProps struct {
	// Position: "top-right" (default) | "top-left" | "bottom-right" | "bottom-left"
	Position string
	// Max is the maximum number of notifications visible at once (default 5).
	Max int
}

// NotificationStack renders a fixed-position container for persistent notifications.
// Add notifications at runtime via: window._pushNotification({title, body, variant, duration})
// variant: "info" | "success" | "warning" | "error"
// duration: ms until auto-dismiss (0 = manual dismiss only)
func NotificationStack(p NotificationStackProps) g.Node {
	if p.Position == "" {
		p.Position = "top-right"
	}
	if p.Max == 0 {
		p.Max = 5
	}

	return h.Div(
		g.Attr("data-component", "notification-stack"),
		g.Attr("data-position", p.Position),
		h.Div(
			g.Attr("data-slot", "container"),
			h.ID("notification-stack"),
		),
		h.Script(g.Raw(notificationStackScript(p.Max))),
	)
}

func notificationStackScript(max int) string {
	return `(function(){
  var container=document.getElementById('notification-stack');
  if(!container) return;
  var max=` + fmt.Sprintf("%d", max) + `;
  var count=0;

  window._pushNotification=function(opts){
    opts=opts||{};
    var title=opts.title||'Notification';
    var body=opts.body||'';
    var variant=opts.variant||'info';
    var duration=typeof opts.duration==='number'?opts.duration:5000;

    // Trim to max
    while(container.children.length>=max){
      container.firstChild.remove(); count--;
    }

    var id='notif-'+(++count);
    var n=document.createElement('div');
    n.id=id;
    n.setAttribute('data-component','notification');
    n.setAttribute('data-variant',variant);
    n.style.cssText='animation:mljr-toast-in .25s ease';
    n.innerHTML='<div data-slot="body">'+
      '<strong data-slot="title">'+escHtml(title)+'</strong>'+
      (body?'<p data-slot="message">'+escHtml(body)+'</p>':'')+
      '</div>'+
      '<button data-slot="dismiss" type="button" aria-label="Dismiss" onclick="this.closest(\'[data-component=notification]\').remove()">✕</button>';
    container.appendChild(n);

    if(duration>0){
      setTimeout(function(){
        if(n.parentNode){
          n.style.animation='mljr-toast-out .2s ease forwards';
          setTimeout(function(){ n.remove(); },200);
        }
      },duration);
    }
  };

  function escHtml(s){ return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;'); }
})();`
}
