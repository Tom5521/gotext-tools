package po

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int
	File string
}

func (l Location) Equal(l2 Location) bool {
	return util.Equal(l, l2)
}

type Locations []Location

func CompareLocation(a, b Location) int {
	line := CompareLocationByLine(a, b)
	file := CompareLocationByFile(a, b)

	if file != 0 {
		return file
	}

	return line
}

func CompareLocationByLine(a, b Location) int {
	return a.Line - b.Line
}

func CompareLocationByFile(a, b Location) int {
	return strings.Compare(filepath.Clean(a.File), filepath.Clean(b.File))
}

func (l Locations) Equal(l2 Locations) bool {
	return util.Equal(l, l2)
}

func (l Locations) Sort() Locations {
	slices.SortFunc(l, CompareLocation)
	return l
}

func (l Locations) SortByLine() Locations {
	slices.SortFunc(l, CompareLocationByLine)

	return l
}

func (l Locations) SortByFile() Locations {
	slices.SortFunc(l, CompareLocationByFile)

	return l
}
