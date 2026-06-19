package pages

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

var slugDisallowed = regexp.MustCompile(`[^a-z0-9-]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugDisallowed.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "group"
	}
	return s
}

// randomSuffix returns a short hex suffix so group slugs stay unique
// without a retry loop on collision.
func randomSuffix() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// randomToken returns a long random hex string for invite links — long
// enough that guessing one isn't feasible.
func randomToken() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
