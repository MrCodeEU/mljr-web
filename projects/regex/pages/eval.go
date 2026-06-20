package pages

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/dlclark/regexp2/v2"
)

const (
	maxPatternBytes     = 16 << 10
	maxInputBytes       = 256 << 10
	maxReplacementBytes = 64 << 10
	maxMatches          = 10000
	maxGroupsPerMatch   = 1000
)

type EvalInput struct {
	Pattern string
	FlagI   bool
	FlagM   bool
	FlagS   bool
	Input   string
	Replace string
}

type MatchInfo struct {
	Index  int
	Value  string
	Start  int
	End    int
	Groups []string
}

type EvalResult struct {
	Err            string
	MatchCount     int
	Highlighted    string // HTML-safe: plain text escaped, matches in <mark>
	Matches        []MatchInfo
	Replaced       string
	ReplaceApplied bool
	InputRaw       string
	PatternEmpty   bool
}

func EvalRegex(inp EvalInput) EvalResult {
	if len(inp.Pattern) > maxPatternBytes {
		return limitError(inp, fmt.Sprintf("pattern is too large; limit is %d KiB", maxPatternBytes/1024))
	}
	if len(inp.Input) > maxInputBytes {
		return limitError(inp, fmt.Sprintf("input is too large; limit is %d KiB", maxInputBytes/1024))
	}
	if len(inp.Replace) > maxReplacementBytes {
		return limitError(inp, fmt.Sprintf("replacement is too large; limit is %d KiB", maxReplacementBytes/1024))
	}
	if inp.Pattern == "" {
		return EvalResult{
			Highlighted:  html.EscapeString(inp.Input),
			InputRaw:     inp.Input,
			PatternEmpty: true,
		}
	}

	var flags regexp2.RegexOptions
	if inp.FlagI {
		flags |= regexp2.IgnoreCase
	}
	if inp.FlagM {
		flags |= regexp2.Multiline
	}
	if inp.FlagS {
		flags |= regexp2.Singleline
	}

	re, err := regexp2.Compile(inp.Pattern, flags)
	if err != nil {
		return EvalResult{
			Err:         err.Error(),
			Highlighted: html.EscapeString(inp.Input),
			InputRaw:    inp.Input,
		}
	}
	re.MatchTimeout = 500 * time.Millisecond

	var spans [][2]int
	var matches []MatchInfo
	idx := 0

	m, err := re.FindStringMatch(inp.Input)
	if err != nil {
		return EvalResult{
			Err:         "match error: " + err.Error(),
			Highlighted: html.EscapeString(inp.Input),
			InputRaw:    inp.Input,
		}
	}
	for m != nil {
		if idx >= maxMatches {
			return EvalResult{
				Err:         fmt.Sprintf("too many matches; limit is %d", maxMatches),
				MatchCount:  len(matches),
				Highlighted: buildHighlight(inp.Input, spans),
				Matches:     matches,
				InputRaw:    inp.Input,
			}
		}
		bstart, blen := m.ByteRange()
		end := bstart + blen
		spans = append(spans, [2]int{bstart, end})

		groups := m.Groups()
		if len(groups) > maxGroupsPerMatch+1 {
			return EvalResult{
				Err:         fmt.Sprintf("too many capture groups; limit is %d", maxGroupsPerMatch),
				MatchCount:  len(matches),
				Highlighted: buildHighlight(inp.Input, spans),
				Matches:     matches,
				InputRaw:    inp.Input,
			}
		}
		var groupStrs []string
		for _, g := range groups[1:] { // groups[0] = full match
			groupStrs = append(groupStrs, g.String())
		}

		matches = append(matches, MatchInfo{
			Index:  idx,
			Value:  m.String(),
			Start:  bstart,
			End:    end,
			Groups: groupStrs,
		})
		idx++

		m, err = re.FindNextMatch(m)
		if err != nil {
			return EvalResult{
				Err:         "match error: " + err.Error(),
				Highlighted: html.EscapeString(inp.Input),
				InputRaw:    inp.Input,
			}
		}
	}

	highlighted := buildHighlight(inp.Input, spans)

	replaced := inp.Input
	if inp.Replace != "" {
		if r, rerr := re.Replace(inp.Input, inp.Replace, -1, maxMatches); rerr == nil {
			replaced = r
		}
	}

	return EvalResult{
		MatchCount:     len(matches),
		Highlighted:    highlighted,
		Matches:        matches,
		Replaced:       replaced,
		ReplaceApplied: inp.Replace != "",
		InputRaw:       inp.Input,
	}
}

func limitError(inp EvalInput, msg string) EvalResult {
	return EvalResult{
		Err:         msg,
		Highlighted: html.EscapeString(inp.Input),
		InputRaw:    inp.Input,
	}
}

func buildHighlight(input string, spans [][2]int) string {
	if len(spans) == 0 {
		return html.EscapeString(input)
	}
	var sb strings.Builder
	pos := 0
	for _, sp := range spans {
		if sp[0] < pos {
			continue
		}
		if pos < sp[0] {
			sb.WriteString(html.EscapeString(input[pos:sp[0]]))
		}
		sb.WriteString(`<mark class="rx-match">`)
		sb.WriteString(html.EscapeString(input[sp[0]:sp[1]]))
		sb.WriteString(`</mark>`)
		pos = sp[1]
	}
	if pos < len(input) {
		sb.WriteString(html.EscapeString(input[pos:]))
	}
	return sb.String()
}
