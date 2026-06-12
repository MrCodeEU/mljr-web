package form

import (
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type OTPInputProps struct {
	Name   string
	Length int    // number of digits, default 6
	Label  string // sr-only label
}

// OTPInput renders a one-time-password input as N single-digit boxes.
// Auto-advances focus, handles paste, backspace navigation.
// Hidden input named p.Name carries the concatenated value.
func OTPInput(p OTPInputProps) g.Node {
	if p.Length == 0 {
		p.Length = 6
	}
	if p.Label == "" {
		p.Label = "One-time password"
	}

	boxes := make([]g.Node, p.Length)
	for i := range boxes {
		boxes[i] = h.Input(
			g.Attr("data-slot", "otp-box"),
			h.Type("text"),
			g.Attr("inputmode", "numeric"),
			g.Attr("maxlength", "1"),
			g.Attr("pattern", "[0-9]"),
			g.Attr("data-component", "input"),
			g.Attr("style", "width:3rem;text-align:center;font-size:var(--t-xl);font-weight:900;font-family:var(--font-display);letter-spacing:0;padding:.5rem"),
			g.Attr("aria-label", fmt.Sprintf("%s digit %d", p.Label, i+1)),
			g.Attr("autocomplete", func() string {
				if i == 0 {
					return "one-time-code"
				}
				return "off"
			}()),
		)
	}

	return h.Div(
		g.Attr("data-component", "otp-input"),
		h.Div(
			g.Attr("data-slot", "boxes"),
			g.Attr("style", "display:flex;gap:var(--sp-2)"),
			g.Group(boxes),
		),
		h.Input(
			h.Type("hidden"),
			h.Name(p.Name),
			g.Attr("data-slot", "value"),
		),
		h.Script(g.Raw(otpScript)),
	)
}

const otpScript = `(function(){
  document.querySelectorAll('[data-component="otp-input"]').forEach(function(root){
    var boxes=Array.from(root.querySelectorAll('[data-slot="otp-box"]'));
    var hidden=root.querySelector('[data-slot="value"]');
    function sync(){ hidden.value=boxes.map(function(b){return b.value;}).join(''); }

    boxes.forEach(function(box,i){
      box.addEventListener('input',function(){
        box.value=box.value.replace(/[^0-9]/g,'').slice(-1);
        sync();
        if(box.value&&i<boxes.length-1) boxes[i+1].focus();
      });
      box.addEventListener('keydown',function(e){
        if(e.key==='Backspace'&&!box.value&&i>0){
          boxes[i-1].focus(); boxes[i-1].value=''; sync();
        } else if(e.key==='ArrowLeft'&&i>0){ e.preventDefault(); boxes[i-1].focus(); }
        else if(e.key==='ArrowRight'&&i<boxes.length-1){ e.preventDefault(); boxes[i+1].focus(); }
      });
      box.addEventListener('paste',function(e){
        e.preventDefault();
        var data=(e.clipboardData||window.clipboardData).getData('text').replace(/[^0-9]/g,'');
        data.slice(0,boxes.length-i).split('').forEach(function(ch,j){
          if(boxes[i+j]) boxes[i+j].value=ch;
        });
        sync();
        var next=Math.min(i+data.length,boxes.length-1);
        boxes[next].focus();
      });
      box.addEventListener('focus',function(){ box.select(); });
    });
  });
})();`
