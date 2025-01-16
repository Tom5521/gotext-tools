package parser

import (
	"fmt"
	"strings"

	"github.com/Tom5521/xgotext/flags"
)

type Location struct {
	Line int
	File string
}

type Translation struct {
	ID        string
	Context   string
	Plural    string
	Locations []Location
}

func (t Translation) String() string {
	var builder strings.Builder

	fprintfln := func(format string, args ...any) {
		fmt.Fprintf(&builder, format+"\n", args...)
	}

	id := formatString(t.ID)
	context := formatString(t.Context)
	plural := formatString(t.Plural)

	for _, location := range t.Locations {
		fprintfln("# %s:%d", location.File, location.Line)
	}

	if t.Context != "" {
		fprintfln("msgctxt %s", context)
	}

	fprintfln("msgid %s", id)

	if t.Plural != "" {
		fprintfln("msgid_plural %s", plural)
		for i := range flags.Nplurals {
			fprintfln(`msgstr[%d] ""`, i)
		}
	} else {
		fprintfln(`msgstr ""`)
	}

	return builder.String()
}

func formatString(str string) string {
	str = fixSpecialChars(str)
	return formatMultiline(str)
}

func formatMultiline(str string) string {
	var builder strings.Builder
	builder.Grow(len(str) * 2)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		fmt.Fprintf(&builder, `"%s"`, line)
		if i != len(lines)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}

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
