package po

import (
	"errors"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

// Entry represents a single translatable message in a PO file.
//
// It may include context, plural forms, translation string(s),
// comments, flags, and source code locations.
type Entry struct {
	// Metadata and comments.

	Flags             []string // List of flags (e.g., "fuzzy").
	Comments          []string // Translator comments.
	ExtractedComments []string // Automatically extracted comments.
	Previous          []string // Previous msgid or msgstr lines.

	// Main fields.

	Obsolete  bool          // Indicates whether the entry is obsolete.
	ID        string        // The original string to be translated.
	Context   string        // Context of the string, if any.
	Plural    string        // The plural form of the original string, if applicable.
	Plurals   PluralEntries // List of plural translations.
	Str       string        // Translated string (singular).
	Locations Locations     // List of source code references.
}

// markAsObsolete marks the entry as obsolete.
func (e *Entry) markAsObsolete() { e.Obsolete = true }

// markAsFuzzy adds the "fuzzy" flag to the entry if not already present.
func (e *Entry) markAsFuzzy() {
	if !e.IsFuzzy() {
		e.Flags = append(e.Flags, "fuzzy")
	}
}

// IsHeader reports whether the entry is a header (i.e., both ID and context are empty).
func (e Entry) IsHeader() bool {
	return e.ID == "" && e.Context == ""
}

// Validate checks the entry for internal inconsistencies.
// It returns an error if the entry is both plural and singular.
func (e Entry) Validate() error {
	if e.Str != "" && e.IsPlural() && len(e.Plurals) > 0 {
		return errors.New("the entry can't be plural and singular at the same time")
	}
	return nil
}

// UnifiedStr returns the translation string formatted for MO files.
//
// For plural entries, it joins all plural forms using '\x00'.
// For singular entries, it returns the Str field.
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

// UnifiedID returns the unique identifier for the entry formatted for MO files.
//
// This includes the context, msgid, and plural (if any),
// separated by '\x04' and '\x00' as per gettext MO format.
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

// FullHash returns a hash based on the unified ID including context and plural.
//
// WARNING: This is intended for internal use and is not compatible with gettext's hash.
func (e Entry) FullHash() uint32 {
	return util.PJWHash(e.UnifiedID())
}

// Hash returns a hash based on context and ID.
//
// This is useful for identifying entries with or without plural forms.
func (e Entry) Hash() uint32 {
	id := e.ID
	if e.HasContext() {
		id = e.Context + "\x04" + id
	}
	return util.PJWHash(id)
}

// Equal reports whether two entries are semantically equivalent.
//
// It compares ID, context, translation string, flags, and obsolete status.
func (e Entry) Equal(x Entry) bool {
	ok1 := e.UnifiedID() == x.UnifiedID() && e.UnifiedStr() == x.UnifiedStr()
	ok2 := slices.CompareFunc(e.Flags, x.Flags, strings.Compare) == 0
	ok3 := e.Obsolete && x.Obsolete
	return ok1 && ok2 && ok3
}

// IsPlural reports whether the entry is plural or contains plural forms.
func (e Entry) IsPlural() bool {
	return e.Plural != "" || len(e.Plurals) > 0
}

// HasContext reports whether the entry has a non-empty context.
func (e Entry) HasContext() bool {
	return e.Context != ""
}

// IsFuzzy reports whether the entry is marked as "fuzzy".
func (e Entry) IsFuzzy() bool {
	return slices.Contains(e.Flags, "fuzzy")
}

// String returns a formatted representation of the entry for debugging or display.
func (e Entry) String() string {
	return util.Format(e)
}

