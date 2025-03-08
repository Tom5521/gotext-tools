package compiler

import (
	"fmt"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/po"
)

const (
	copyrightFormat = `# Copyright (C) %s
# This file is distributed under the same license as the %s package.`
	foreignCopyrightFormat = `# This file is put in the public domain.`
	headerFormat           = `# %s
%s
#
#, fuzzy
msgid ""
msgstr ""
`
	headerFieldFormat = `"%s: %s\n"`
)

func (c PoCompiler) formatHeader() string {
	if c.Config.OmitHeader {
		return ""
	}
	var b strings.Builder
	header := c.File.Header()

	copyright := fmt.Sprintf(copyrightFormat, c.Config.CopyrightHolder, c.Config.PackageName)
	if c.Config.ForeignUser {
		copyright = foreignCopyrightFormat
	}

	fmt.Fprintf(&b, headerFormat, c.Config.Title, copyright)

	for i, field := range header.Fields {
		fmt.Fprintf(&b, headerFieldFormat, field.Key, field.Value)
		if i != len(header.Fields) {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

func (c PoCompiler) formatEntry(t po.Entry) string {
	var builder strings.Builder
	nplurals := c.File.Header().Nplurals()

	// Helper function to append formatted lines to the builder.
	fprintfln := func(format string, args ...any) {
		var prefix string
		if t.Comment {
			prefix = "# "
		}
		fmt.Fprintf(&builder, prefix+format+"\n", args...)
	}

	id := formatString(t.ID)
	context := formatString(t.Context)
	plural := formatString(t.Plural)

	for _, comment := range t.Comments {
		fprintfln("# %s", comment)
	}
	for _, xcomment := range t.ExtractedComments {
		fprintfln("#. %s", xcomment)
	}
	// Add location comments if not suppressed by the configuration.
	if !c.Config.NoLocation && c.Config.AddLocation != PoLocationModeNever {
		switch c.Config.AddLocation {
		case PoLocationModeFull:
			for _, location := range t.Locations {
				fprintfln("#: %s:%d", location.File, location.Line)
			}
		case PoLocationModeFile:
			for _, location := range t.Locations {
				fprintfln("#: %s", location.File)
			}
		}
	}

	for _, flag := range t.Flags {
		fprintfln("#, %s", flag)
	}

	for _, previous := range t.Previous {
		fprintfln("#| %s", previous)
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

		for i := range util.ROverNumber(nplurals) {
			if i == 1 {
				fprintfln("msgstr[%d] %s", i, formatPrefixAndSuffix(t.Plural, c.Config))
				continue
			}
			fprintfln(`msgstr[%d] %s`, i, formatPrefixAndSuffix(t.ID, c.Config))
		}
	} else {
		// Add empty msgstr for singular strings.
		fprintfln(`msgstr %s`, formatPrefixAndSuffix(t.ID, c.Config))
	}

	return builder.String()
}

func formatPrefixAndSuffix(id string, cfg PoConfig) string {
	text := `""`

	if cfg.MsgstrPrefix != "" {
		text = formatString(cfg.MsgstrPrefix + id)
	}
	if cfg.MsgstrSuffix != "" {
		text = formatString(id + cfg.MsgstrSuffix)
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
