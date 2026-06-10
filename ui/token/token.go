// Package token holds the typed enums consumed by ui components and the
// showcase registry. Mirrors the data-* values handled by ui/css/core.css.
package token

type Variant string

const (
	Primary   Variant = "primary"
	Secondary Variant = "secondary"
	Outline   Variant = "outline"
	Danger    Variant = "danger"
	Ghost     Variant = "ghost"
	VarTone   Variant = "tone"
)

type Size string

const (
	SizeSM   Size = "sm"
	SizeMD   Size = "md"
	SizeLG   Size = "lg"
	SizeIcon Size = "icon"
)

type Tone string

const (
	ToneNone    Tone = ""
	ToneYellow  Tone = "yellow"
	ToneLime    Tone = "lime"
	ToneCyan    Tone = "cyan"
	ToneViolet  Tone = "violet"
	TonePink    Tone = "pink"
	ToneSky     Tone = "sky"
	ToneMint    Tone = "mint"
	ToneBlush   Tone = "blush"
	ToneAccent  Tone = "accent"
	ToneAccent2 Tone = "accent-2"
)

type ToastVariant string

const (
	ToastInfo    ToastVariant = "info"
	ToastSuccess ToastVariant = "success"
	ToastWarning ToastVariant = "warning"
	ToastDanger  ToastVariant = "danger"
)

type Theme string

const (
	ThemeSwissBrut Theme = "swissbrut"
	ThemeInk       Theme = "ink"
)

type Mode string

const (
	ModeLight Mode = "light"
	ModeDark  Mode = "dark"
)
