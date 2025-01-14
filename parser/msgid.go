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

	for _, location := range t.Locations {
		builder.WriteString(fmt.Sprintf("#: %s:%d\n", location.File, location.Line))
	}

	if t.Context != "" {
		text := fmt.Sprintf(`msgctxt "%s"`+"\n", t.Context)
		if strings.Contains(t.Context, "\n") {
			text = fmt.Sprintf("msgctxt %s\n", formatMultiline(t.Context))
		}
		builder.WriteString(text)
	}

	text := fmt.Sprintf(`msgid "%s"`+"\n", t.ID)
	if strings.Contains(t.ID, "\n") {
		text = fmt.Sprintf("msgid %s\n", formatMultiline(t.ID))
	}
	builder.WriteString(text)

	if t.Plural != "" {
		text = fmt.Sprintf(`msgid_plural "%s"`+"\n", t.Plural)
		if strings.Contains(t.Plural, "\n") {
			text = fmt.Sprintf("msgid_plural %s\n", formatMultiline(t.Plural))
		}
		builder.WriteString(text)
		for i := range flags.Nplurals {
			builder.WriteString(fmt.Sprintf(`msgstr[%d] ""`+"\n", i))
		}
	} else {
		builder.WriteString(`msgstr ""` + "\n")
	}

	return builder.String()
}

func formatMultiline(str string) (n string) {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		n += fmt.Sprintf(`"%s"`, line)
		if i != len(lines)-1 {
			n += "\n"
		}
	}

	return
}
