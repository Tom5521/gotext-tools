package compile

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

func escapePOString(s string) string {
	var buf strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\t':
			buf.WriteString(`\t`)
		case '\r':
			buf.WriteString(`\r`)
		default:
			if strconv.IsPrint(r) {
				buf.WriteRune(r)
			} else {
				fmt.Fprintf(&buf, "\\x%02x", r)
			}
		}
	}
	return buf.String()
}

func (c PoCompiler) compileEntries() {
	for _, e := range c.File.Entries {
		c.entry(e)
	}
}

func (c PoCompiler) entry(entry po.Entry) {
	eb := EntryBuilder{entry, c.writer, c.Config}
	eb.Build()
}

type EntryBuilder struct {
	po.Entry
	Builder io.Writer
	Config  PoConfig
}

func (eb EntryBuilder) Build() {
	eb.comment()
	eb.msgid()
	eb.msgstr()
}

func (eb EntryBuilder) header() {}
func (eb EntryBuilder) msgid() {
	if eb.HasContext() {
		eb.keyword("msgctxt")
		eb.string(eb.Context)
	}
	eb.keyword("msgid")
	eb.string(eb.ID)

	if eb.IsPlural() {
		eb.keyword("msgid_plural")
		eb.string(eb.Plural)
	}
}

func (eb EntryBuilder) msgstr() {
	const format = "msgstr[%d]"
	if eb.IsPlural() {
		for _, pe := range eb.Plurals {
			eb.keyword(fmt.Sprintf(format, pe.ID))
			eb.string(pe.Str)
		}
		return
	}

	eb.keyword("msgstr")
	eb.string(eb.Str)
}

func (eb EntryBuilder) fuzzy()        {}
func (eb EntryBuilder) obsolete()     {}
func (eb EntryBuilder) translated()   {}
func (eb EntryBuilder) untranslated() {}

func (eb EntryBuilder) comment() {
	eb.translatorComment()
	eb.extractedComment()
	eb.referenceComment()
	eb.flagComment()
	eb.previousComment()
}

func (eb EntryBuilder) translatorComment() {}
func (eb EntryBuilder) extractedComment()  {}

func (eb EntryBuilder) referenceComment() {
	eb.reference()
}

func (eb EntryBuilder) reference() {}

func (eb EntryBuilder) flagComment() {
	eb.flag()
}

func (eb EntryBuilder) flag() {
	eb.fuzzyFlag()
}

func (eb EntryBuilder) fuzzyFlag() {}

func (eb EntryBuilder) previousComment() {
	eb.previous()
}

func (eb EntryBuilder) previous() {}

func (eb EntryBuilder) keyword(kw string) {
	fmt.Fprint(eb.Builder, kw+" ")
}

func (eb EntryBuilder) string(str string) {
	if eb.Config.WordWrap {
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			if i != len(lines)-1 {
				line += "\n"
			}
			fmt.Fprint(eb.Builder, `"`)
			eb.text(escapePOString(line))
			fmt.Fprintln(eb.Builder, `"`)
		}

		return
	}
	fmt.Fprint(eb.Builder, `"`)
	eb.text(escapePOString(str))
	fmt.Fprintln(eb.Builder, `"`)
}

func (eb EntryBuilder) text(txt string) {
	fmt.Fprint(eb.Builder, txt)
	/*
		 	eb.escapeSequence()
			eb.formatDirective()
			eb.invalidFormatDirective()
	*/
}

func (eb EntryBuilder) escapeSequence()         {}
func (eb EntryBuilder) formatDirective()        {}
func (eb EntryBuilder) invalidFormatDirective() {}
