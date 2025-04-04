package po

import (
	"errors"
	"slices"
	"strings"

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

	Obsolete  bool
	ID        string // The original string to be translated.
	Context   string // The context in which the string is used (optional).
	Plural    string // The plural form of the string (optional).
	Plurals   PluralEntries
	Str       string
	Locations Locations // A list of source code locations for the string.
}

// Check for possible errors and inconsistencies in the entry.
func (e Entry) Validate() error {
	if e.Str != "" && e.IsPlural() && len(e.Plurals) > 0 {
		return errors.New("the entry cant be plural and singular at the same time")
	}
	return nil
}

// Returns the msgstr (STR) or plurals
// (PluralStr1 \x00 PluralStr2 \x00 PluralStr3...)
// unified according to the MO format.
func (e Entry) UnifiedStr() string {
	str := e.Str
	if e.IsPlural() {
		var msgstrs []string
		plurals := e.Plurals.Sort()
		for _, plural := range plurals {
			msgstrs = append(msgstrs, plural.Str)
		}
		str = strings.Join(msgstrs, "\x00")
	}

	return str
}

// Returns the unified msgid, msgid_plural and context according
// to MO format (CTXT \x04 ID \x00 PLURAL).
func (e Entry) UnifiedID() string {
	id := e.ID
	if e.HasContext() {
		id = e.Context + "\x04" + id
	}
	if e.Plural != "" {
		id += "\x00" + e.Plural
	}

	return id
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
