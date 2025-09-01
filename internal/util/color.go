package util

import (
	"fmt"

	"github.com/xo/terminfo"
)

type Color uint8

const (
	Default Color = 39

	Normal       = Default
	Bold   Color = iota + 1
	_            // Skip 2
	Italic
	Underline
)

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

func (c Color) Sprint(a ...any) string {
	str := fmt.Sprint(a...)
	lvl, err := terminfo.ColorLevelFromEnv()
	if err != nil {
		return str
	}
	if lvl <= terminfo.ColorLevelNone {
		return str
	}

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, str)
}

func (c Color) Sprintf(format string, args ...any) string {
	return c.Sprint(fmt.Sprintf(format, args...))
}
