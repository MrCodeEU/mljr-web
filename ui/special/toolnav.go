package special

type toolLink struct {
	ID    string
	Label string
	URL   string
}

// liveTools are the standalone tool sites, linked from each other by
// ToolShell.
var liveTools = []toolLink{
	{ID: "regex", Label: "Regex Lab", URL: "https://regex.mljr.eu"},
	{ID: "cron", Label: "Cron Explorer", URL: "https://cron.mljr.eu"},
	{ID: "codec", Label: "Codec", URL: "https://codec.mljr.eu"},
}
