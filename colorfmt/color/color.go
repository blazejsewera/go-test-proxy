/*
Package color is based on a part of color.go file
from github.com/fatih/color library, v1.19.0,
which is licensed under MIT license.

Source file is available at: https://github.com/fatih/color/blob/v1.19.0/color.go
Lines copied: 76 -- 170, with changes.

Copyright (c) 2013 Fatih Arslan.
License is available at: https://github.com/fatih/color/blob/v1.19.0/LICENSE.md
*/
package color

type StyleAttribute int
type ColorAttribute int
type BgColorAttribute int

const Escape = "\x1b"

const (
	Reset StyleAttribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

const (
	ResetBold StyleAttribute = iota + 22
	ResetItalic
	ResetUnderline
	ResetBlinking
	_
	ResetReversed
	ResetConcealed
	ResetCrossedOut
)

var MapResetAttributes = map[StyleAttribute]StyleAttribute{
	Bold:         ResetBold,
	Faint:        ResetBold,
	Italic:       ResetItalic,
	Underline:    ResetUnderline,
	BlinkSlow:    ResetBlinking,
	BlinkRapid:   ResetBlinking,
	ReverseVideo: ResetReversed,
	Concealed:    ResetConcealed,
	CrossedOut:   ResetCrossedOut,
}

// Foreground text colors
const (
	FgBlack ColorAttribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite

	// used internally for 256 and 24-bit coloring
	foreground
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack ColorAttribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack BgColorAttribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack BgColorAttribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)
