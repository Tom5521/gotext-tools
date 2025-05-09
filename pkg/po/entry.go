package po

import (
	"errors"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
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

func (e *Entry) markAsObsolete() { e.Obsolete = true }

func (e *Entry) markAsFuzzy() {
	if !e.IsFuzzy() {
		e.Flags = append(e.Flags, "fuzzy")
	}
}

func (e Entry) IsHeader() bool {
	return e.ID == "" && e.Context == ""
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
		plurals := e.Plurals
		if !plurals.IsSorted() {
			plurals = plurals.Sort()
		}
		for _, plural := range plurals {
			msgstrs = append(msgstrs, plural.Str)
		}
		str = strings.Join(msgstrs, "\x00")
	}

	return str
}

// Returns the unified msgid, msgid_plural and context according
// to MO format (CTXT \x04 ID (\x00 PLURAL)...).
func (e Entry) UnifiedID() string {
	id := e.ID
	if e.HasContext() {
		id = e.Context + "\x04" + id
	}
	if e.IsPlural() && e.Plural != "" {
		id += "\x00" + e.Plural
	}

	return id
}

// Returns the complete hash of the input ID's, including the msgctxt, msgid and msgid_plural.
// WARNING: This is NOT compatible with the original gettext.
func (e Entry) FullHash() uint32 {
	return util.PJWHash(e.UnifiedID())
}

func (e Entry) Hash() uint32 {
	id := e.ID
	if e.HasContext() {
		id = e.Context + "\x04" + id
	}
	// IDK, why tf if I remove this code, magically works the MO compilation.
	// TODO: Understand this thing.
	// if e.IsPlural() && e.Plural != "" {
	// 	id += "\x00" + e.Plural
	// }

	return util.PJWHash(id)
}

func (e Entry) Equal(x Entry) bool {
	ok1 := e.UnifiedID() == x.UnifiedID() && e.UnifiedStr() == x.UnifiedStr()
	ok2 := slices.CompareFunc(e.Flags, x.Flags, strings.Compare) == 0
	ok3 := e.Obsolete && x.Obsolete
	return ok1 && ok2 && ok3
}

func (e Entry) IsPlural() bool {
	return e.Plural != "" || len(e.Plurals) > 0
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
