// Package fonts embeds the self-hosted web fonts used by the swiss theme
// (Inter + JetBrains Mono, variable, latin + latin-ext subsets from fontsource).
// Projects serve them under /fonts/ via web.MountFonts.
package fonts

import "embed"

//go:embed *.woff2
var FS embed.FS
