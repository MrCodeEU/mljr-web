package pages

// Example is a preset pattern/input pair loadable from the dropdown.
type Example struct {
	Key     string
	Name    string
	Pattern string
	Input   string
	Replace string
	FlagI   bool
	FlagM   bool
	FlagS   bool
}

// ExampleGroup is a labeled set of examples shown as an <optgroup>.
type ExampleGroup struct {
	Name  string
	Items []Example
}

var ExampleGroups = []ExampleGroup{
	{
		Name: "Basic",
		Items: []Example{
			{Key: "basic-word", Name: "Word tokens", Pattern: `\w+`, Input: "Hello, world! 42 items."},
			{Key: "basic-digit", Name: "Numbers", Pattern: `\d+`, Input: "Order #1234 costs $56.78 (qty: 3)"},
			{Key: "basic-vowel", Name: "Character class — vowels", Pattern: `[aeiou]`, Input: "Hello world", FlagI: true},
			{Key: "basic-email", Name: "Email addresses", Pattern: `[\w.+-]+@[\w-]+\.[a-z]{2,}`, Input: "Contact foo@bar.com or baz@qux.io", FlagI: true},
		},
	},
	{
		Name: "Groups & Replace",
		Items: []Example{
			{Key: "group-swap", Name: "Swap first ↔ last name", Pattern: `(\w+)\s+(\w+)`, Input: "John Doe\nJane Smith\nAlex Jones", Replace: "$2, $1"},
			{Key: "group-named", Name: "Named groups — ISO date", Pattern: `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})`, Input: "Events: 2024-01-15 and 2024-12-31", Replace: "${day}/${month}/${year}"},
			{Key: "group-noncap", Name: "Non-capturing group", Pattern: `(?:foo|bar)\d+`, Input: "foo123 bar456 baz789 qux000"},
		},
	},
	{
		Name: "Anchors",
		Items: []Example{
			{Key: "anchor-linestart", Name: "Line-start words (multiline)", Pattern: `^\w+`, Input: "foo bar\nbaz qux\nhello world", FlagM: true},
			{Key: "anchor-boundary", Name: "Word boundary — 'cat'", Pattern: `\bcat\b`, Input: "cat concatenate catfish scattered cats"},
		},
	},
	{
		Name: "PCRE: Lookaround",
		Items: []Example{
			{Key: "look-pos-ahead", Name: "Lookahead — word before number", Pattern: `\w+(?=\s\d)`, Input: "foo 123 bar baz 456 qux"},
			{Key: "look-pos-behind", Name: "Lookbehind — price after $", Pattern: `(?<=\$)\d+`, Input: "Prices: $100, $250, €300, £45"},
			{Key: "look-neg-ahead", Name: "Negative lookahead", Pattern: `\b\w+\b(?!\s\d)`, Input: "foo 123 bar baz 456 qux"},
			{Key: "look-neg-behind", Name: "Negative lookbehind — bare decimal", Pattern: `(?<!\d)\.\d+`, Input: "3.14 .5 12.0 .99"},
		},
	},
	{
		Name: "PCRE: Backreferences",
		Items: []Example{
			{Key: "backref-dup", Name: "Duplicate words", Pattern: `\b(\w+)\s+\1\b`, Input: "the the quick brown fox fox jumps over", FlagI: true},
			{Key: "backref-tag", Name: "Balanced HTML tags", Pattern: `<(\w+)[^>]*>.*?</\1>`, Input: "<div>hello</div>\n<p>world</p>\n<span>!</span>", FlagS: true},
		},
	},
	{
		Name: "Practical",
		Items: []Example{
			{Key: "prac-url", Name: "URLs", Pattern: `https?://[\w./%-]+`, Input: "Visit https://example.com or http://foo.bar/path?q=1#anchor"},
			{Key: "prac-ipv4", Name: "IPv4 addresses", Pattern: `\b(?:\d{1,3}\.){3}\d{1,3}\b`, Input: "IPs: 192.168.1.1 and 10.0.0.1 and 999.999.999.999"},
			{Key: "prac-strip-html", Name: "Strip HTML tags", Pattern: `<[^>]+>`, Input: "<h1>Hello</h1>\n<p>World <strong>!</strong></p>", Replace: ""},
			{Key: "prac-camel-snake", Name: "camelCase → snake_case", Pattern: `([a-z])([A-Z])`, Input: "helloWorld myVariableName camelCaseString", Replace: "$1_$2", FlagI: false},
		},
	},
}

// FindExample returns the Example with the given Key.
func FindExample(key string) (Example, bool) {
	for _, grp := range ExampleGroups {
		for _, ex := range grp.Items {
			if ex.Key == key {
				return ex, true
			}
		}
	}
	return Example{}, false
}
