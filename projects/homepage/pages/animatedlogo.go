package pages

import (
	"mljr-web/ui/special"

	g "maragu.dev/gomponents"
)

// AnimatedLogoBackground renders one document-space scatter logo behind the homepage.
// It must not size itself from document scroll height; that causes phantom
// page height and throws off scroll progress calculations.
func AnimatedLogoBackground() g.Node {
	return special.LogoScatter(special.LogoScatterProps{
		ID:              "logo-svg-hp",
		SVGStyle:        "position:absolute;top:8vh;left:45%;transform:translateX(-55%);overflow:visible;width:min(540px,120vw);height:min(540px,120vw);opacity:0.26",
		InitialOpacity:  0.26,
		Mode:            "scroll",
		TriggerID:       "hero",
		WithBackground:  true,
		BackgroundStyle: "position:absolute;inset:0;pointer-events:none;overflow:visible;z-index:0",
		PieceCopies:     3,
		WrapInLoad:      true,
	})
}
