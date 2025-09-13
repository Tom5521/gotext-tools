package compile

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

type entryBuilder struct {
	po.Entry
	buffer *bytes.Buffer
	Config PoConfig
}

func (eb *entryBuilder) print(a ...any) {
	fmt.Fprint(eb.buffer, a...)
}

func (eb *entryBuilder) printf(format string, args ...any) {
	fmt.Fprintf(eb.buffer, format, args...)
}

func (eb *entryBuilder) println(a ...any) {
	fmt.Fprintln(eb.buffer, a...)
}

func (eb *entryBuilder) BuildEntry() []byte {
	defer eb.buffer.Reset()
	eb.comment()
	eb.msgid()
	eb.msgstr()

	commentFuzzy := eb.IsFuzzy() && eb.Config.CommentFuzzy

	if eb.Obsolete || commentFuzzy {
		entry := eb.buffer.String()
		var fixed string
		prefix := "#"
		if eb.Obsolete {
			if eb.Config.UseCustomObsoletePrefix {
				prefix += string(eb.Config.CustomObsoletePrefixRune)
			} else {
				prefix += "~"
			}
		}
		prefix += " "

		for _, line := range strings.Split(entry, "\n") {
			if strings.HasPrefix(line, "#") || line == "" {
				fixed += line + "\n"
				continue
			}

			fixed += prefix + line + "\n"
		}

		eb.buffer.Reset()
		eb.buffer.WriteString(fixed)
	}
	return eb.buffer.Bytes()
}

func (eb *entryBuilder) BuildHeader(header po.Header) []byte {
	defer eb.buffer.Reset()

	if eb.Config.HeaderComments {
		copyright := fmt.Sprintf(copyrightFormat, eb.Config.CopyrightHolder, eb.Config.PackageName)
		if eb.Config.ForeignUser {
			copyright = foreignCopyrightFormat
		}

		eb.printf(headerFormat, eb.Config.Title, copyright)
	} else {
		eb.print(headerEntry)
	}

	if eb.Config.HeaderFields {
		for i, field := range header.Fields {
			eb.printf(headerFieldFormat, field.Key, field.Value)

			if i != len(header.Fields) {
				eb.print("\n")
			}
		}
	}

	eb.println()

	eb.msgid()
	eb.msgstr()

	return eb.buffer.Bytes()
}

func (eb *entryBuilder) msgid() {
	if eb.HasContext() {
		eb.print("msgctxt ")
		eb.string(eb.Context)
	}
	eb.print("msgid ")
	eb.string(eb.ID)

	if eb.IsPlural() {
		eb.print("msgid_plural ")
		eb.string(eb.Plural)
	}
}

func (eb *entryBuilder) msgstr() {
	const format = "msgstr[%d] "
	if eb.IsPlural() {
		if len(eb.Plurals) == 0 {
			for i := 0; i < 2; i++ {
				eb.printf(format, i)
				eb.string(eb.ID)
			}
			return
		}
		for _, pe := range eb.Plurals {
			eb.printf(format, pe.ID)
			eb.string(
				eb.Config.MsgstrPrefix + pe.Str + eb.Config.MsgstrSuffix,
			)
		}

		return
	}

	eb.print("msgstr ")
	eb.string(
		eb.Config.MsgstrPrefix + eb.Str + eb.Config.MsgstrSuffix,
	)
}

func (eb *entryBuilder) comment() {
	eb.translatorComment()
	eb.extractedComment()
	eb.referenceComment()
	eb.flagComment()
	eb.previousComment()
}

func (eb *entryBuilder) translatorComment() {
	for _, comment := range eb.Comments {
		eb.printf("# %s\n", comment)
	}
}

func (eb *entryBuilder) extractedComment() {
	for _, comment := range eb.ExtractedComments {
		eb.printf("#. %s\n", comment)
	}
}

func (eb *entryBuilder) referenceComment() {
	if eb.Config.NoLocation || eb.Config.AddLocation == PoLocationModeNever {
		return
	}
	var writeRef func(id int)
	switch eb.Config.AddLocation {
	case PoLocationModeFull:
		writeRef = func(id int) {
			l := eb.Locations[id]
			eb.printf("%s:%d\n", l.File, l.Line)
		}
	case PoLocationModeFile:
		writeRef = func(id int) {
			l := eb.Locations[id]
			eb.printf("%s\n", l.File)
		}
	}

	for i := range eb.Locations {
		eb.print("#: ")
		writeRef(i)
	}
}

func (eb *entryBuilder) flagComment() {
	for _, f := range eb.Flags {
		eb.printf("#, %s\n", f)
	}
}

func (eb *entryBuilder) previousComment() {
	for _, p := range eb.Previous {
		eb.printf("#| %s\n", p)
	}
}

func (eb *entryBuilder) string(str string) {
	if eb.Config.WordWrap {
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			if i != len(lines)-1 {
				line += "\n"
			}
			eb.printf("\"%s\"\n", escapePOString(line))
		}
		return
	}
	eb.printf("\"%s\"\n", escapePOString(str))
}
