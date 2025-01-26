package compiler

import (
	"fmt"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/config"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

// Format generates the PO file representation of the Translation.
// The output is influenced by the provided configuration.
//
// Parameters:
//   - cfg: The `poconfig.Config` object used to control formatting behavior.
//
// Returns:
//   - A string representing the translation in PO file format.
//
// Example:
//
//	translation := Translation{
//	    ID:      "Hello",
//	    Context: "Greeting",
//	    Plural:  "Hellos",
//	    Locations: []Location{
//	        {Line: 10, File: "example.go"},
//	    },
//	}
//	config := poconfig.DefaultConfig()
//	formatted := translation.Format(config)
//	fmt.Println(formatted)
func FormatTranslation(t types.Translation, cfg config.Config) string {
	var builder strings.Builder

	// Helper function to append formatted lines to the builder.
	fprintfln := func(format string, args ...any) {
		fmt.Fprintf(&builder, format+"\n", args...)
	}

	id := formatString(t.ID)
	context := formatString(t.Context)
	plural := formatString(t.Plural)

	// Add location comments if not suppressed by the configuration.
	if !cfg.NoLocation || cfg.AddLocation == "never" {
		switch cfg.AddLocation {
		case "full":
			for _, location := range t.Locations {
				fprintfln("#: %s:%d", location.File, location.Line)
			}
		case "file":
			for _, location := range t.Locations {
				fprintfln("#: %s", location.File)
			}
		}
	}

	// Add context if available.
	if t.Context != "" {
		fprintfln("msgctxt %s", context)
	}

	// Add singular form.
	fprintfln("msgid %s", id)

	// Add plural forms if present.
	if t.Plural != "" {
		fprintfln("msgid_plural %s", plural)
		for i := range cfg.Nplurals {
			if i == 1 {
				fprintfln("msgstr[%d] %s", i, formatPrefixAndSuffix(t.Plural, cfg))
				continue
			}
			fprintfln(`msgstr[%d] %s`, i, formatPrefixAndSuffix(t.ID, cfg))
		}
	} else {
		// Add empty msgstr for singular strings.
		fprintfln(`msgstr %s`, formatPrefixAndSuffix(t.ID, cfg))
	}

	return builder.String()
}

func formatPrefixAndSuffix(id string, cfg config.Config) string {
	text := `""`

	if cfg.Msgstr.Prefix != "" {
		text = formatString(cfg.Msgstr.Prefix + id)
	}
	if cfg.Msgstr.Suffix != "" {
		text = formatString(id + cfg.Msgstr.Suffix)
	}

	return text
}

// formatString applies formatting rules to a string to make it compatible
// with PO file syntax. It escapes special characters and handles multiline strings.
//
// Parameters:
//   - str: The input string.
//
// Returns:
//   - The formatted string.
func formatString(str string) string {
	str = fixSpecialChars(str)
	return formatMultiline(str)
}

// formatMultiline formats a string as a PO-compatible multiline string.
// Line breaks are escaped with `\n`.
//
// Parameters:
//   - str: The input string.
//
// Returns:
//   - A multiline-formatted string.
func formatMultiline(str string) string {
	var builder strings.Builder
	builder.Grow(len(str) * 2)

	builder.WriteRune('"')

	for _, char := range str {
		if char == '\n' {
			builder.WriteString("\\n")
			continue
		}
		builder.WriteRune(char)
	}

	builder.WriteRune('"')

	return builder.String()
}

// fixSpecialChars escapes special characters (`"` and `\`) in a string.
//
// Parameters:
//   - str: The input string.
//
// Returns:
//   - The string with escaped special characters.
func fixSpecialChars(str string) string {
	var builder strings.Builder
	builder.Grow(len(str) * 2)

	for _, char := range str {
		if char == '"' || char == '\\' {
			builder.WriteRune('\\')
		}
		builder.WriteRune(char)
	}

	return builder.String()
}
