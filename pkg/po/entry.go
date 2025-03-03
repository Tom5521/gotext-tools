package po

import (
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
)

// Location represents the location of a translation string in the source code.
type Location struct {
	Line int    // The line number of the translation.
	File string // The file name where the translation is located.
}

type PluralEntries []PluralEntry

func (p PluralEntries) Sort() PluralEntries {
	var entries PluralEntries
	copy(entries, p)
	slices.SortFunc(entries, func(a, b PluralEntry) int {
		return a.ID - b.ID
	})

	return entries
}

type PluralEntry struct {
	ID  int
	Str string
}

// Entry represents a translatable string, including its context, plural forms,
// and source code locations.
type Entry struct {
	Flags             []string
	Comments          []string
	ExtractedComments []string
	Previous          []string
	ID                string // The original string to be translated.
	Context           string // The context in which the string is used (optional).
	Plural            string // The plural form of the string (optional).
	Plurals           PluralEntries
	Str               string
	Locations         []Location // A list of source code locations for the string.
}

func (e Entry) Hash() uint {
	var b strings.Builder

	if e.Context != "" {
		b.WriteString(e.Context)
		b.WriteByte('4') // EOT byte.
	}
	b.WriteString(e.ID)

	return util.PJWHash(b.String())
}

func (e Entry) String() string {
	return util.Format(e)
}
