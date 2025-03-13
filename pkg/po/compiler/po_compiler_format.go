package compiler

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po"
)

const (
	copyrightFormat = `# Copyright (C) %s
# This file is distributed under the same license as the %s package.`
	foreignCopyrightFormat = `# This file is put in the public domain.`
	headerFormat           = `# %s
%s
#
` + headerEntry
	headerEntry = `msgid ""
msgstr ""
`
	headerFieldFormat = `"%s: %s\n"`
)

func (c PoCompiler) writeHeader(w *bufio.Writer) {
	if c.Config.OmitHeader {
		return
	}

	if c.Config.HeaderComments {
		copyright := fmt.Sprintf(copyrightFormat, c.Config.CopyrightHolder, c.Config.PackageName)
		if c.Config.ForeignUser {
			copyright = foreignCopyrightFormat
		}

		fmt.Fprintf(w, headerFormat, c.Config.Title, copyright)
	} else {
		fmt.Fprint(w, headerEntry)
	}

	if c.Config.HeaderFields {
		for i, field := range c.header.Fields {
			fmt.Fprintf(w, headerFieldFormat, field.Key, field.Value)

			if i != len(c.header.Fields) {
				fmt.Fprint(w, "\n")
			}
		}
	}

	fmt.Fprintln(w)
}

func (c PoCompiler) writeEntry(w *bufio.Writer, t po.Entry) {
	// Helper function to append formatted lines to the builder.
	write := func(format string, args ...any) {
		var comment string
		if c.Config.CommentFuzzy && slices.Contains(t.Flags, "fuzzy") {
			comment = "# "
		}
		fmt.Fprintf(w, comment+format+"\n", args...)
	}

	id := formatString(t.ID)
	context := formatString(t.Context)
	plural := formatString(t.Plural)

	for _, comment := range t.Comments {
		write("# %s", comment)
	}
	for _, xcomment := range t.ExtractedComments {
		write("#. %s", xcomment)
	}
	// Add location comments if not suppressed by the configuration.
	if !c.Config.NoLocation && c.Config.AddLocation != PoLocationModeNever {
		switch c.Config.AddLocation {
		case PoLocationModeFull:
			for _, location := range t.Locations {
				write("#: %s:%d", location.File, location.Line)
			}
		case PoLocationModeFile:
			for _, location := range t.Locations {
				write("#: %s", location.File)
			}
		}
	}

	for _, flag := range t.Flags {
		write("#, %s", flag)
	}

	for _, previous := range t.Previous {
		write("#| %s", previous)
	}

	// Add context if available.
	if t.Context != "" {
		write("msgctxt %s", context)
	}

	// Add singular form.
	write("msgid %s", id)

	// Add plural forms if present.
	if t.Plural != "" {
		write("msgid_plural %s", plural)

		if len(t.Plurals) == 0 {
			for i := uint(0); i < c.nplurals; i++ {
				write(`msgstr[%d] %s`, i, formatPrefixAndSuffix(t.ID, c.Config))
			}
		} else {
			for _, pe := range t.Plurals {
				write("msgstr[%d] %s", pe.ID, formatPrefixAndSuffix(pe.Str, c.Config))
			}
		}
	} else {
		// Add empty msgstr for singular strings.
		write(`msgstr %s`, formatPrefixAndSuffix(t.Str, c.Config))
	}

	fmt.Fprintln(w)
}

func formatPrefixAndSuffix(id string, cfg PoConfig) string {
	return formatString(cfg.MsgstrPrefix + id + cfg.MsgstrSuffix)
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
