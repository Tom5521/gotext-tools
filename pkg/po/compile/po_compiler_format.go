package compile

import (
	"fmt"
	"io"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
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

func (c PoCompiler) writeHeader(w io.Writer) {
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

func (c PoCompiler) fprintfln(w io.Writer, e po.Entry, format string, args ...any) {
	var prefix string
	if !strings.HasPrefix(format, "#") {
		if c.Config.CommentFuzzy && e.IsFuzzy() {
			prefix = "# "
		}
		if e.Obsolete {
			prefixRune := '~'
			if c.Config.UseCustomObsoletePrefix {
				prefixRune = c.Config.CustomObsoletePrefixRune
			}
			prefix = string([]rune{'#', prefixRune, ' '})
		}
	}
	str := fmt.Sprintf(prefix+format, args...)

	fmt.Fprintln(w, str)
}

func (c PoCompiler) writeComment(w io.Writer, e po.Entry) {
	write := func(format string, args ...any) {
		c.fprintfln(w, e, format, args...)
	}

	for _, comment := range e.Comments {
		write("# %s", comment)
	}
	for _, xcomment := range e.ExtractedComments {
		write("#. %s", xcomment)
	}
	// Add location comments if not suppressed by the configuration.
	if !c.Config.NoLocation && c.Config.AddLocation != PoLocationModeNever {
		switch c.Config.AddLocation {
		case PoLocationModeFull:
			for _, location := range e.Locations {
				write("#: %s:%d", location.File, location.Line)
			}
		case PoLocationModeFile:
			for _, location := range e.Locations {
				write("#: %s", location.File)
			}
		}
	}

	for _, flag := range e.Flags {
		write("#, %s", flag)
	}

	for _, previous := range e.Previous {
		write("#| %s", previous)
	}
}

// TODO: Rename this to something like manageWordWrapping
func (c PoCompiler) formatMultiline(str string) string {
	var builder strings.Builder

	if c.Config.WordWrap {
		c.processWordWrap(&builder, str)
	} else {
		fmt.Fprintf(&builder, "%q", str)
	}

	return builder.String()
}

func (c PoCompiler) processWordWrap(builder *strings.Builder, str string) {
	lines := strings.Split(str, "\n")

	// TODO: Explain why the hell this thing writes **that** to the builder.
	if len(lines) > 1 {
		builder.WriteString("\"\"\n")
	}
	for i, line := range lines {
		if line == "" {
			continue
		}
		isLastLine := i == len(lines)-1
		if !isLastLine {
			line += "\n"
		}
		fmt.Fprintf(builder, "%q", line)
		if !isLastLine {
			builder.WriteByte('\n')
		}
	}
}

func (c PoCompiler) formatMsgstr(i string) string {
	return c.formatMultiline(c.formatPrefixAndSuffix(i))
}

func (c PoCompiler) formatMsgid(i string) string {
	return c.formatMultiline(i)
}

func (c PoCompiler) writeEntry(w io.Writer, e po.Entry) {
	// Helper function to append formatted lines to the builder.
	write := func(format string, args ...any) {
		c.fprintfln(w, e, format, args...)
	}

	c.writeComment(w, e)

	id := c.formatMsgid(e.ID)
	context := c.formatMsgid(e.Context)
	plural := c.formatMsgid(e.Plural)

	// Add context if available.
	if e.HasContext() {
		write("msgctxt %s", context)
	}

	// Add singular form.
	write("msgid %s", id)

	// Add plural forms if present.
	if e.IsPlural() {
		write("msgid_plural %s", plural)

		if len(e.Plurals) == 0 {
			for i := uint(0); i < c.nplurals; i++ {
				write(`msgstr[%d] %s`, i, c.formatMsgstr(e.ID))
			}
		} else {
			for _, pe := range e.Plurals {
				write("msgstr[%d] %s", pe.ID, c.formatMsgstr(pe.Str))
			}
		}
	} else {
		// Add empty msgstr for singular strings.
		write(`msgstr %s`, c.formatMsgstr(e.Str))
	}

	fmt.Fprintln(w)
}

func (c PoCompiler) formatPrefixAndSuffix(id string) string {
	return c.Config.MsgstrPrefix + id + c.Config.MsgstrSuffix
}
