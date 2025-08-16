package po

import (
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
func (e Entry) Validate() []error {
	var errs []error
	if e.Str != "" && e.IsPlural() && len(e.Plurals) > 0 {
		errs = append(errs, &InvalidEntryError{
			ID:     e.UnifiedID(),
			Reason: ErrBadPluralEntry,
		},
		)
	}
	return errs
}

// UnifiedStr returns the translation string formatted for MO files.
//
// For plural entries, it joins all plural forms using '\x00'.
// For singular entries, it returns the Str field.
func (e Entry) UnifiedStr() string {
	var builder strings.Builder
	builder.WriteString(e.Str)
	if len(e.Plurals) > 0 {
		plurals := slices.Clone(e.Plurals)
		if !plurals.IsSorted() {
			plurals = plurals.Sort()
		}
		for i, pe := range plurals {
			builder.WriteString(pe.Str)
			if i != len(plurals)-1 {
				builder.WriteByte(0)
			}
		}
	}

	return builder.String()
}

// UnifiedID returns the unique identifier for the entry formatted for MO files.
//
// This includes the context, msgid, and plural (if any),
// separated by '\x04' and '\x00' as per gettext MO format.
func (e Entry) UnifiedID() string {
	var builder strings.Builder

	if e.HasContext() {
		builder.WriteString(e.Context)
		builder.WriteByte(4)
	}

	builder.WriteString(e.ID)

	if e.Plural != "" {
		builder.WriteByte(0)
		builder.WriteString(e.Plural)
	}

	return builder.String()
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
	var builder strings.Builder
	if e.HasContext() {
		builder.WriteString(e.Context)
		builder.WriteByte(4)
	}
	builder.WriteString(e.ID)
	return util.PJWHash(builder.String())
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

func (e Entry) HasComments() bool {
	return len(e.Flags) > 0 || len(e.Comments) > 0 ||
		len(e.ExtractedComments) > 0 || len(e.Previous) > 0 ||
		len(e.Locations) > 0
}

// String returns a formatted representation of the entry for debugging or display.
func (e Entry) String() string {
	return util.Format(e)
}
