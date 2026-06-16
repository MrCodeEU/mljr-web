package pages

import (
	"html"
	"regexp"
	"strings"
)

// EvalInput holds signals from the client.
type EvalInput struct {
	Pattern string
	FlagI   bool
	FlagM   bool
	FlagS   bool
	Input   string
	Replace string
}

// MatchInfo describes one match and its captured groups.
type MatchInfo struct {
	Index  int
	Value  string
	Start  int
	End    int
	Groups []string
}

// EvalResult is returned by EvalRegex.
type EvalResult struct {
	Err            string
	MatchCount     int
	Highlighted    string // HTML-safe: non-match text escaped, matches in <mark>
	Matches        []MatchInfo
	Replaced       string
	ReplaceApplied bool // true when Replace field was non-empty
	InputRaw       string
	PatternEmpty   bool // true when Pattern was ""
}

// EvalRegex compiles the pattern (with flag prefixes), finds all matches, and
// builds the highlighted HTML. It never panics.
func EvalRegex(inp EvalInput) EvalResult {
	if inp.Pattern == "" {
		return EvalResult{
			Highlighted:  html.EscapeString(inp.Input),
			InputRaw:     inp.Input,
			PatternEmpty: true,
		}
	}

	prefix := ""
	if inp.FlagI {
		prefix += "(?i)"
	}
	if inp.FlagM {
		prefix += "(?m)"
	}
	if inp.FlagS {
		prefix += "(?s)"
	}

	re, err := regexp.Compile(prefix + inp.Pattern)
	if err != nil {
		return EvalResult{
			Err:         err.Error(),
			Highlighted: html.EscapeString(inp.Input),
			InputRaw:    inp.Input,
		}
	}

	allIdx := re.FindAllStringSubmatchIndex(inp.Input, -1)
	highlighted := buildHighlight(inp.Input, allIdx)

	matches := make([]MatchInfo, 0, len(allIdx))
	for i, m := range allIdx {
		val := ""
		if m[0] >= 0 && m[1] >= 0 {
			val = inp.Input[m[0]:m[1]]
		}
		groups := make([]string, 0)
		for j := 2; j+1 < len(m); j += 2 {
			g := ""
			if m[j] >= 0 && m[j+1] >= 0 {
				g = inp.Input[m[j]:m[j+1]]
			}
			groups = append(groups, g)
		}
		matches = append(matches, MatchInfo{
			Index:  i,
			Value:  val,
			Start:  m[0],
			End:    m[1],
			Groups: groups,
		})
	}

	replaced := re.ReplaceAllString(inp.Input, inp.Replace)

	return EvalResult{
		MatchCount:     len(allIdx),
		Highlighted:    highlighted,
		Matches:        matches,
		Replaced:       replaced,
		ReplaceApplied: inp.Replace != "",
		InputRaw:       inp.Input,
	}
}

func buildHighlight(input string, allIdx [][]int) string {
	if len(allIdx) == 0 {
		return html.EscapeString(input)
	}
	var sb strings.Builder
	pos := 0
	for _, m := range allIdx {
		if m[0] < 0 || m[1] < 0 {
			continue
		}
		if pos < m[0] {
			sb.WriteString(html.EscapeString(input[pos:m[0]]))
		}
		sb.WriteString(`<mark class="rx-match">`)
		sb.WriteString(html.EscapeString(input[m[0]:m[1]]))
		sb.WriteString(`</mark>`)
		pos = m[1]
	}
	if pos < len(input) {
		sb.WriteString(html.EscapeString(input[pos:]))
	}
	return sb.String()
}
