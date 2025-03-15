package po

import (
	"slices"

	"github.com/Tom5521/xgotext/internal/util"
)

// Entry represents a translatable string, including its context, plural forms,
// and source code locations.
type Entry struct {
	// Comments.

	Flags             []string
	Comments          []string
	ExtractedComments []string
	Previous          []string

	// Fields.

	ID        string // The original string to be translated.
	Context   string // The context in which the string is used (optional).
	Plural    string // The plural form of the string (optional).
	Plurals   PluralEntries
	Str       string
	Locations Locations // A list of source code locations for the string.
}

func (e Entry) Equal(x Entry) bool {
	return util.Equal(e, x)
}

func (e Entry) IsPlural() bool {
	return e.Plural != ""
}

func (e Entry) HasContext() bool {
	return e.Context != ""
}

func (e Entry) IsFuzzy() bool {
	return slices.Contains(e.Flags, "fuzzy")
}

func (e Entry) String() string {
	return util.Format(e)
}
