package color

import (
	"fmt"

	"github.com/xo/terminfo"
)

type Color uint8

const (
	Default Color = 39
	Normal        = Default
)

const (
	Reset Color = iota
	Bold
	_
	Italic
	Underline
)

// Regular colors.
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

// Background colors.
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

const (
	ResetFormat = "\x1b[0m"
	ColorFormat = "\x1b[%dm"
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

	return fmt.Sprintf("\x1b[%dm%s"+ResetFormat, c, str)
}

func (c Color) Sprintf(format string, args ...any) string {
	return c.Sprint(fmt.Sprintf(format, args...))
}
