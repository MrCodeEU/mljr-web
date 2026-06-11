package pages

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestLegalPagesRenderOwnerAndPrivacyText(t *testing.T) {
	for name, node := range map[string]any{
		"impressum":   Impressum(AnalyticsConfig{}),
		"datenschutz": Datenschutz(AnalyticsConfig{}),
	} {
		renderable, ok := node.(interface {
			Render(w io.Writer) error
		})
		if !ok {
			t.Fatalf("%s is not renderable", name)
		}
		var buf bytes.Buffer
		if err := renderable.Render(&buf); err != nil {
			t.Fatalf("%s render: %v", name, err)
		}
		body := buf.String()
		if !strings.Contains(body, "Michael Reinegger") {
			t.Fatalf("%s page does not contain owner name", name)
		}
		if !strings.Contains(body, "hello@mljr.eu") {
			t.Fatalf("%s page does not contain contact email", name)
		}
	}
}

func TestAnalyticsHeadRequiresScriptAndWebsiteID(t *testing.T) {
	if got := AnalyticsHead(AnalyticsConfig{UmamiScriptSrc: "/umami/script.js"}); got != nil {
		t.Fatalf("AnalyticsHead with missing website ID returned %d nodes", len(got))
	}
	if got := AnalyticsHead(AnalyticsConfig{UmamiWebsiteID: "site-id"}); got != nil {
		t.Fatalf("AnalyticsHead with missing script returned %d nodes", len(got))
	}
	if got := AnalyticsHead(AnalyticsConfig{UmamiScriptSrc: "/umami/script.js", UmamiWebsiteID: "site-id"}); len(got) != 1 {
		t.Fatalf("AnalyticsHead returned %d nodes, want 1", len(got))
	}
}
