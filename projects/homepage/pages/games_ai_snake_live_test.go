package pages

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAISnakeLiveRendersStreamHookup(t *testing.T) {
	for _, lang := range []string{"en", "de"} {
		var buf bytes.Buffer
		if err := AISnakeLive(lang, AnalyticsConfig{}).Render(&buf); err != nil {
			t.Fatalf("%s render: %v", lang, err)
		}
		body := buf.String()
		for _, want := range []string{AISnakeStreamPath, `id="ai-snake-live-canvas"`, "/static/games/ai-snake-live/renderer.js", "snakeLiveRender"} {
			if !strings.Contains(body, want) {
				t.Errorf("%s page missing %q", lang, want)
			}
		}
		if strings.Contains(body, "games.ai_snake_live.") {
			t.Errorf("%s page shows a raw i18n key (missing translation)", lang)
		}
	}
}

func TestPlaygroundListsAISnakeLive(t *testing.T) {
	for _, p := range PlaygroundURLPaths() {
		if p == "/games/ai-snake-live" {
			return
		}
	}
	t.Fatal("playground does not list /games/ai-snake-live")
}

// The page calls window.snakeLiveRender and draws into #ai-snake-live-canvas; the
// renderer script must define/use exactly those names or the board stays blank.
func TestAISnakeLiveRendererMatchesPage(t *testing.T) {
	js, err := os.ReadFile("../assets/static/games/ai-snake-live/renderer.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"window.snakeLiveRender", "ai-snake-live-canvas"} {
		if !strings.Contains(string(js), want) {
			t.Errorf("renderer.js does not contain %q", want)
		}
	}
}
