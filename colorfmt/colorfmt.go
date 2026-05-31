package colorfmt

import (
	"errors"
	"fmt"
	"io"

	"github.com/blazejsewera/go-test-proxy/colorfmt/color"
)

const escape = color.Escape

type Style int
type Color int

const (
	Normal Style = iota
	Bold
	Italic
	Underline
)

const (
	Base Color = iota
	Red
	Green
	Yellow
	Blue
	BrightWhite
)

type ColorFmt struct {
	colorEnabled bool
	stdout       io.Writer
	stderr       io.Writer
}

func New(colorEnabled bool, stdout io.Writer, stderr io.Writer) *ColorFmt {
	c := &ColorFmt{colorEnabled, stdout, stderr}
	if c.colorEnabled {
		_, err1 := fmt.Fprint(c.stdout, beginFormatting(Normal, Base))
		_, err2 := fmt.Fprint(c.stderr, beginFormatting(Normal, Base))
		if err := errors.Join(err1, err2); err != nil {
			panic(err)
		}
	}
	return c
}

func (c *ColorFmt) Cprint(style Style, color Color, s string) {
	c.styledPrint(style, color, s)
}

func (c *ColorFmt) Cprintf(style Style, color Color, format string, v ...any) {
	formatted := fmt.Sprintf(format, v...)
	c.styledPrint(style, color, formatted)
}

func (c *ColorFmt) Cerrprintf(style Style, color Color, format string, v ...any) {
	formatted := fmt.Sprintf(format, v...)
	c.styledStdErrPrint(style, color, formatted)
}

func (c *ColorFmt) styledPrint(style Style, color Color, s string) {
	if !c.colorEnabled {
		_, err := fmt.Fprint(c.stdout, s)
		if err != nil {
			panic(err)
		}
		return
	}

	_, err := fmt.Fprintf(c.stdout, "%s%s%s", beginFormatting(style, color), s, endFormatting())
	if err != nil {
		panic(err)
	}
}

func (c *ColorFmt) styledStdErrPrint(style Style, color Color, s string) {
	if !c.colorEnabled {
		_, err := fmt.Fprint(c.stderr, s)
		if err != nil {
			panic(err)
		}
		return
	}

	_, err := fmt.Fprintf(c.stderr, "%s%s%s", beginFormatting(style, color), s, endFormatting())
	if err != nil {
		panic(err)
	}
}

func beginFormatting(style Style, color Color) string {
	styleAttr := mapStyle(style)
	colorAttr := mapColor(color)
	return fmt.Sprintf("%s[%d;%dm", escape, styleAttr, colorAttr)
}

func endFormatting() string {
	return fmt.Sprintf("%s[%dm", color.Escape, color.Reset)
}

func mapStyle(s Style) color.StyleAttribute {
	switch s {
	case Normal:
		return color.Reset
	case Bold:
		return color.Bold
	case Italic:
		return color.Italic
	case Underline:
		return color.Underline
	default:
		return color.Reset
	}
}

func mapColor(c Color) color.ColorAttribute {
	switch c {
	case Base:
		return color.FgWhite
	case Red:
		return color.FgRed
	case Green:
		return color.FgGreen
	case Yellow:
		return color.FgYellow
	case Blue:
		return color.FgBlue
	case BrightWhite:
		return color.FgHiWhite
	default:
		return color.FgWhite
	}
}
