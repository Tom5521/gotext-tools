package config

import (
	"errors"
)

// Config represents the configuration options for PO (Portable Object) file processing.
// This struct contains various settings that influence the parsing and generation of PO files.
type Config struct {
	DefaultDomain  string   // The default domain name for translations.
	ForcePo        bool     // If true, forces the output to always be in PO format.
	NoLocation     bool     // If true, omits location comments from the output.
	AddLocation    string   // Determines how location comments are added (e.g., "never").
	OmitHeader     bool     // If true, omits the header section from the output.
	PackageName    string   // The name of the package associated with the translations.
	PackageVersion string   // The version of the package associated with the translations.
	Language       string   // The language code for translations (e.g., "en").
	Nplurals       uint     // The number of plural forms for the language.
	Exclude        []string // A list of paths to exclude from processing.
	JoinExisting   bool     // If true, combines new translations with existing ones.
	ExtractAll     bool     // If true, extracts all translatable strings regardless of other settings.
	Msgstr         struct {
		Prefix string // A prefix to prepend to all `msgstr` values.
		Suffix string // A suffix to append to all `msgstr` values.
	}
}

// DefaultConfig returns a default configuration for PO file processing.
// The default configuration sets the following values:
//   - DefaultDomain: "default"
//   - Language: "en"
//   - Nplurals: 2
func DefaultConfig() Config {
	return Config{
		DefaultDomain: "default",
		Language:      "en",
		Nplurals:      2,
	}
}

// Validate checks the Config object for potential inconsistencies or invalid settings.
//
// Returns:
//   - A slice of errors if any issues are found; otherwise, an empty slice.
//
// Validation rules:
//   - If `NoLocation` is true and `AddLocation` is not set to "never", a conflict error is reported.
//   - If `Nplurals` is 0, an error is reported.
//
// Example:
//
//	config := DefaultConfig()
//	config.NoLocation = true
//	config.AddLocation = "always"
//	errs := config.Validate()
//	if len(errs) > 0 {
//	    fmt.Println("Configuration errors:", errs)
//	}
func (c Config) Validate() (errs []error) {
	if c.NoLocation && c.AddLocation != "never" {
		errs = append(errs, errors.New("noLocation and AddLocation are in conflict"))
	}

	if c.Nplurals == 0 {
		errs = append(errs, errors.New("nplurals is equal to 0"))
	}

	return
}
