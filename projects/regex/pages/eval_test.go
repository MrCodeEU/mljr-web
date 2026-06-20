package pages

import (
	"strings"
	"testing"
)

func TestEvalRegexRejectsOversizedInput(t *testing.T) {
	result := EvalRegex(EvalInput{
		Pattern: ".",
		Input:   strings.Repeat("a", maxInputBytes+1),
	})
	if result.Err == "" {
		t.Fatal("expected oversized input error")
	}
	if !strings.Contains(result.Err, "input is too large") {
		t.Fatalf("unexpected error: %q", result.Err)
	}
}

func TestEvalRegexLimitsMatchCount(t *testing.T) {
	result := EvalRegex(EvalInput{
		Pattern: "a",
		Input:   strings.Repeat("a", maxMatches+1),
	})
	if result.Err == "" {
		t.Fatal("expected match limit error")
	}
	if result.MatchCount != maxMatches {
		t.Fatalf("expected %d retained matches, got %d", maxMatches, result.MatchCount)
	}
}
