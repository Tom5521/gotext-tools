package compiler

import (
	"fmt"
	"io"
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

func (c PoCompiler) writeHeader(w io.Writer) (err error) {
	if c.Config.OmitHeader {
		return nil
	}
	header := c.File.Header()

	if c.Config.HeaderComments {
		copyright := fmt.Sprintf(copyrightFormat, c.Config.CopyrightHolder, c.Config.PackageName)
		if c.Config.ForeignUser {
			copyright = foreignCopyrightFormat
		}

		_, err = fmt.Fprintf(w, headerFormat, c.Config.Title, copyright)
	} else {
		_, err = fmt.Fprint(w, headerEntry)
	}

	if err != nil {
		return
	}

	if c.Config.HeaderFields {
		for i, field := range header.Fields {
			_, err = fmt.Fprintf(w, headerFieldFormat, field.Key, field.Value)
			if err != nil {
				return
			}

			if i != len(header.Fields) {
				_, err = fmt.Fprint(w, "\n")
				if err != nil {
					return
				}
			}
		}
	}

	if _, err = fmt.Fprintln(w); err != nil {
		return err
	}

	return
}

func (c PoCompiler) writeEntry(w io.Writer, t po.Entry) (err error) {
	nplurals := c.File.Header().Nplurals()

	// Helper function to append formatted lines to the builder.
	fprintfln := func(format string, args ...any) error {
		var comment string
		if c.Config.CommentFuzzy && slices.Contains(t.Flags, "fuzzy") {
			comment = "# "
		}
		_, err = fmt.Fprintf(w, comment+format+"\n", args...)
		return err
	}

	id := formatString(t.ID)
	context := formatString(t.Context)
	plural := formatString(t.Plural)

	// Helper function to handle repeated error checking.
	write := func(format string, args ...any) error {
		if err = fprintfln(format, args...); err != nil {
			return err
		}
		return nil
	}

	for _, comment := range t.Comments {
		if err = write("# %s", comment); err != nil {
			return err
		}
	}
	for _, xcomment := range t.ExtractedComments {
		if err = write("#. %s", xcomment); err != nil {
			return err
		}
	}
	// Add location comments if not suppressed by the configuration.
	if !c.Config.NoLocation && c.Config.AddLocation != PoLocationModeNever {
		switch c.Config.AddLocation {
		case PoLocationModeFull:
			for _, location := range t.Locations {
				if err = write("#: %s:%d", location.File, location.Line); err != nil {
					return err
				}
			}
		case PoLocationModeFile:
			for _, location := range t.Locations {
				if err = write("#: %s", location.File); err != nil {
					return err
				}
			}
		}
	}

	for _, flag := range t.Flags {
		if err = write("#, %s", flag); err != nil {
			return err
		}
	}

	for _, previous := range t.Previous {
		if err = write("#| %s", previous); err != nil {
			return err
		}
	}

	// Add context if available.
	if t.Context != "" {
		if err = write("msgctxt %s", context); err != nil {
			return err
		}
	}

	// Add singular form.
	if err = write("msgid %s", id); err != nil {
		return err
	}

	// Add plural forms if present.
	if t.Plural != "" {
		if err = write("msgid_plural %s", plural); err != nil {
			return err
		}
		if len(t.Plurals) == 0 {
			for i := uint(0); i < nplurals; i++ {
				if err = write(`msgstr[%d] %s`, i, formatPrefixAndSuffix(t.ID, c.Config)); err != nil {
					return err
				}
			}
		} else {
			for _, pe := range t.Plurals {
				if err := write("msgstr[%d] %s", pe.ID, formatPrefixAndSuffix(pe.Str, c.Config)); err != nil {
					return err
				}
			}
		}
	} else {
		// Add empty msgstr for singular strings.
		if err = write(`msgstr %s`, formatPrefixAndSuffix(t.Str, c.Config)); err != nil {
			return err
		}
	}

	if _, err = fmt.Fprintln(w); err != nil {
		return err
	}

	return
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
