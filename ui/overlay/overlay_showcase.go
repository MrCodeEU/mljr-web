//go:build showcase

package overlay

import (
	"fmt"
	"strings"

	"mljr-web/ui/icon"
	"mljr-web/ui/primitive"
	"mljr-web/ui/registry"
	"mljr-web/ui/token"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// renderStr renders a gomponents node to an HTML string.
func renderStr(n g.Node) string {
	var b strings.Builder
	_ = n.Render(&b)
	return b.String()
}

// jsStr escapes a string for embedding inside a single-quoted JS string literal.
func jsStr(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func init() {
	registry.Register(&registry.Component{
		Slug: "modal", Name: "Modal", Category: "overlay",
		Summary: "Scrim + dialog gated by $modalOpen signal. Size: sm|md|lg.",
		Code: `// trigger (anywhere on page)
h.Button(g.Attr("data-on:click", "$modalOpen=true"), g.Text("Open"))

// panel (render once)
overlay.Modal(overlay.ModalProps{Title: "Confirm", Size: token.MD},
    h.P(g.Text("Are you sure?")),
)`,
		Controls: []registry.Control{
			{Name: "size", Type: registry.ControlEnum, Options: []string{"sm", "md", "lg"}, Default: "md"},
		},
		Render: func(p map[string]string) g.Node {
			return h.Div(
				g.Attr("data-signals", "{modalOpen: false}"),
				h.Button(
					g.Attr("data-component", "button"),
					g.Attr("data-variant", "primary"),
					g.Attr("data-on:click", "$modalOpen = true"),
					g.Text("Open"),
				),
				Modal(ModalProps{Title: "Hello", Size: token.Size(p["size"])},
					h.P(g.Text("Modal body text.")),
				),
			)
		},
	})

	registry.Register(&registry.Component{
		Slug: "toast", Name: "Toast", Category: "overlay",
		Summary: "Self-managing notification stack: auto-dismiss timer, progress bar, icons, variant tones, queue up to 6.",
		Code: `// 1. Mount Toaster once per page (e.g. in PageShell body)
overlay.Toaster(overlay.ToasterProps{})

// 2. SSE handler emits toasts into #toaster
sse.PatchElements(web.RenderToString(
    overlay.Toast(overlay.ToastProps{
        ID: "t-"+uuid, Variant: token.ToastSuccess,
        Title: "Saved", AutoDismiss: 5,
    }, g.Text("Changes saved.")),
))`,
		Render: func(p map[string]string) g.Node {
			// data-for does NOT exist in Datastar 1.0.2 — use direct DOM manipulation.
			// _spawnToast(variant) is defined in the inline <script> below.
			// Each button calls it; icons are pre-embedded as SVG strings at render time.

			iconSVG := func(name string) string {
				return jsStr(renderStr(icon.Icon(name)))
			}

			// JS function embedded in <script> — random funny messages per variant.
			script := fmt.Sprintf(`
window._spawnToast = (function() {
  var msgs = {
    info: [
      'Your code compiled on the first try. Screenshot it.',
      'npm audit: 0 vulnerabilities. Universe glitching.',
      'All unit tests passed. Even the flaky ones. Wild.',
      'Your PR has exactly 0 comments. Nobody read it.'
    ],
    success: [
      'Deployment went live. Servers are throwing a party.',
      'PR merged. Tech debt +0. Achievement unlocked.',
      'Migration succeeded. Tables are still there, shockingly.',
      'Docker build cached. You have been blessed today.'
    ],
    warning: [
      'Bundle size grew 2 KB. Shareholders have been notified.',
      '3 unresolved TODOs found. Classic move.',
      'console.log found in production. Police alerted.',
      'You are on page 2 of Google results. Rest in peace.'
    ],
    danger: [
      'Production is down. Twitter already knows. They know.',
      '404: Developer not found. Please provide coffee.',
      'The div is not centered. Design team dispatched.',
      'Someone tried DROP TABLE users again. We good though.'
    ]
  };
  var icons = {
    info:    '%s',
    success: '%s',
    warning: '%s',
    danger:  '%s'
  };
  var titles = {info:'Info', success:'Done!', warning:'Heads up', danger:'Error'};
  return function(variant) {
    var r = document.querySelector('[data-component="toaster"]');
    if (!r) return;
    var items = r.querySelectorAll('[data-component="toast"]');
    if (items.length >= 6) items[0].remove();
    var arr = msgs[variant] || msgs.info;
    var msg = arr[Math.floor(Math.random() * arr.length)];
    var t = document.createElement('div');
    t.setAttribute('data-component', 'toast');
    t.setAttribute('data-variant', variant);
    t.setAttribute('data-autodismiss', '');
    t.style.setProperty('--toast-duration', '5s');
    t.innerHTML =
      '<div data-slot="icon">' + icons[variant] + '</div>' +
      '<div data-slot="content">' +
        '<div data-slot="title">' + titles[variant] + '</div>' +
        '<div style="font-size:var(--t-xs);opacity:.8;margin-top:.15em">' + msg + '</div>' +
      '</div>' +
      '<div data-slot="progress"></div>' +
      '<button data-slot="close" onclick="this.closest(\'[data-component=toast]\').remove()">×</button>';
    t.addEventListener('animationend', function(e) {
      if (e.animationName === 'mljr-toast-out') t.remove();
    });
    r.appendChild(t);
  };
})();`,
				iconSVG("lucide:info"),
				iconSVG("lucide:circle-check"),
				iconSVG("lucide:alert-triangle"),
				iconSVG("lucide:circle-x"),
			)

			toastBtn := func(variant, label string, v token.Variant, icn string) g.Node {
				return primitive.Button(
					primitive.ButtonProps{Variant: v, Size: token.SizeSM},
					g.Attr("data-on:click", fmt.Sprintf("_spawnToast('%s')", variant)),
					icon.Icon(icn),
					g.Text(label),
				)
			}

			return h.Div(
				h.Script(g.Raw(script)),
				h.Style("min-height:360px;position:relative"),

				// spawn buttons
				h.Div(
					h.Style("display:flex;flex-wrap:wrap;gap:var(--sp-2);margin-bottom:var(--sp-4)"),
					toastBtn("info", "Info", token.Outline, "lucide:info"),
					toastBtn("success", "Success", token.Outline, "lucide:circle-check"),
					toastBtn("warning", "Warning", token.Outline, "lucide:alert-triangle"),
					toastBtn("danger", "Error", token.Outline, "lucide:circle-x"),
					primitive.Button(
						primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
						g.Attr("data-on:click",
							`_spawnToast('info');_spawnToast('success');_spawnToast('warning')`),
						g.Text("Spawn 3"),
					),
					primitive.Button(
						primitive.ButtonProps{Variant: token.Ghost, Size: token.SizeSM},
						g.Attr("data-on:click",
							`document.querySelectorAll('[data-component="toast"]').forEach(t=>t.remove())`),
						g.Text("Clear all"),
					),
				),

				h.P(
					h.Style("font-size:var(--t-xs);opacity:.5;margin-bottom:var(--sp-3)"),
					g.Text("Click to spawn · auto-dismiss 5 s · stack max 6"),
				),

				// Toaster — positioned absolute within the showcase container
				h.Div(
					h.ID("demo-toaster"),
					g.Attr("data-component", "toaster"),
					h.Style("position:absolute"),
				),
			)
		},
	})
}
