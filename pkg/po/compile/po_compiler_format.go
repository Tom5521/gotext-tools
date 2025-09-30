package compile

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/color"
	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/pkg/po"
)

type entryBuilder struct {
	po.Entry
	Config PoConfig
}

func (eb *entryBuilder) applyStyle(str string, names ...string) string {
	if eb.Config.Highlight == nil {
		return str
	}

	name := strings.Join(names, " ")
	properties, found := eb.Config.Highlight[name]
	if !found {
		if len(names) > 1 {
			return eb.applyStyle(str, slices.Delete(names, 0, 1)...)
		}

		return str
	}
	str = applyColor(properties.BackgroundColor, str)
	str = applyColor(properties.Color, str)
	str = applyWeight(properties.FontWeight, str)
	str = applyStyle(properties.FontStyle, str)
	str = applyDecoration(properties.TextDecoration, str)
	return str
}

func applyDecoration(d HighlightTextDecoration, str string) string {
	if d != TextDecorationUnderline {
		return str
	}
	return color.Underline.Sprint(str)
}

func applyStyle(s HighlightFontStyle, str string) string {
	switch s {
	case FontStyleItalic, FontStyleOblique:
		return color.Italic.Sprint(str)
	}

	return str
}

func applyWeight(w HighlightFontWeight, str string) string {
	if w != FontWeightBold {
		return str
	}
	return color.Bold.Sprint(str)
}

func applyColor(c TermColorer, str string) string {
	if c == nil {
		return str
	}
	return c.Sprint(str)
}

func (eb *entryBuilder) BuildEntry() []byte {
	var buf bytes.Buffer
	buf.WriteString(eb.comment())
	buf.WriteString(eb.msgid())
	buf.WriteString(eb.msgstr())

	commentFuzzy := eb.IsFuzzy() && eb.Config.CommentFuzzy

	if eb.Obsolete || commentFuzzy {
		entry := buf.String()
		buf.Reset()

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
				buf.WriteString(line + "\n")
				continue
			}

			buf.WriteString(prefix + line + "\n")
		}
	}

	buf.WriteByte('\n')
	return buf.Bytes()
}

func (eb *entryBuilder) BuildHeader(header po.Header) []byte {
	var buf bytes.Buffer

	if eb.Config.HeaderComments {
		copyright := fmt.Sprintf(copyrightFormat, eb.Config.CopyrightHolder, eb.Config.PackageName)
		if eb.Config.ForeignUser {
			copyright = foreignCopyrightFormat
		}

		fmt.Fprintf(&buf, headerFormat, eb.Config.Title, copyright)
	} else {
		fmt.Fprint(&buf, headerEntry)
	}

	if eb.Config.HeaderFields {
		for i, field := range header.Fields {
			fmt.Fprintf(&buf, headerFieldFormat, field.Key, field.Value)

			if i != len(header.Fields) {
				fmt.Fprint(&buf, "\n")
			}
		}
	}

	fmt.Fprintln(&buf)

	buf.WriteString(eb.msgid())
	buf.WriteString(eb.msgstr())

	return buf.Bytes()
}

func (eb *entryBuilder) msgid() string {
	var b strings.Builder
	if eb.HasContext() {
		b.WriteString(eb.keyword("msgctxt"))
		b.WriteString(eb.string(eb.Context, "msgid"))
	}
	b.WriteString(eb.keyword("msgid"))
	b.WriteString(eb.string(eb.ID, "msgid"))

	if eb.IsPlural() {
		b.WriteString(eb.keyword("msgid_plural"))
		b.WriteString(eb.string(eb.Plural, "msgid"))
	}

	return b.String()
}

func (eb *entryBuilder) msgstr() string {
	var msgstr strings.Builder
	const format = "msgstr[%d]"
	if eb.IsPlural() {
		if len(eb.Plurals) == 0 {
			id := eb.string(eb.ID, "msgstr")
			for i := 0; i < 2; i++ {
				fmt.Fprint(&msgstr, eb.keyword(fmt.Sprintf(format, i)))
				fmt.Fprint(&msgstr, id)
			}
			return msgstr.String()
		}
		for _, pe := range eb.Plurals {
			fmt.Fprint(&msgstr, eb.keyword(fmt.Sprintf(format, pe.ID)))
			fmt.Fprint(&msgstr,
				eb.string(
					eb.Config.MsgstrPrefix+pe.Str+eb.Config.MsgstrSuffix,
					"msgstr",
				),
			)
		}

		return msgstr.String()
	}

	fmt.Fprint(&msgstr, eb.keyword("msgstr"))
	fmt.Fprint(&msgstr, eb.string(
		eb.Config.MsgstrPrefix+eb.Str+eb.Config.MsgstrSuffix,
		"msgstr",
	))

	return msgstr.String()
}

func (eb *entryBuilder) comment() string {
	var b strings.Builder
	b.WriteString(eb.translatorComment())
	b.WriteString(eb.extractedComment())
	b.WriteString(eb.referenceComment())
	b.WriteString(eb.flagComment())
	b.WriteString(eb.previousComment())

	return eb.applyStyle(b.String(), "comment")
}

func (eb *entryBuilder) translatorComment() string {
	var b strings.Builder
	for _, comment := range eb.Comments {
		fmt.Fprintf(&b, "# %s\n", comment)
	}
	return b.String()
}

func (eb *entryBuilder) extractedComment() string {
	var b strings.Builder
	for _, comment := range eb.ExtractedComments {
		fmt.Fprintf(&b, "#. %s\n", comment)
	}
	return b.String()
}

func (eb *entryBuilder) referenceComment() string {
	if eb.Config.NoLocation || eb.Config.AddLocation == PoLocationModeNever {
		return ""
	}
	var b strings.Builder

	var writeRef func(id int)
	switch eb.Config.AddLocation {
	case PoLocationModeFull:
		writeRef = func(id int) {
			l := eb.Locations[id]
			fmt.Fprintf(&b, "%s:%d\n", l.File, l.Line)
		}
	case PoLocationModeFile:
		writeRef = func(id int) {
			l := eb.Locations[id]
			fmt.Fprintf(&b, "%s\n", l.File)
		}
	}

	for i := range eb.Locations {
		fmt.Fprint(&b, "#: ")
		writeRef(i)
	}

	return b.String()
}

func (eb *entryBuilder) flagComment() string {
	var comments string
	for _, f := range eb.Flags {
		comment := fmt.Sprintf("#, %s\n", f)
		comments += comment
	}

	return comments
}

func (eb *entryBuilder) previousComment() string {
	var b strings.Builder
	for _, p := range eb.Previous {
		fmt.Fprintf(&b, "#| %s\n", p)
	}
	return b.String()
}

func (eb *entryBuilder) string(str string, styles ...string) string {
	var builder strings.Builder
	if eb.Config.WordWrap {
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			if i != len(lines)-1 {
				line += "\n"
			}
			fmt.Fprint(&builder, `"`)
			fmt.Fprint(&builder, eb.text(escapePOString(line),
				slices.Delete(styles, 0, 1)...))
			builder.WriteString(eb.applyStyle(`"`, "string") + "\n")
		}
		return eb.applyStyle(
			builder.String(),
			append(styles, "string")...,
		)
	}

	fmt.Fprint(&builder, `"`)
	fmt.Fprint(&builder, eb.text(escapePOString(str),
		slices.Delete(styles, 0, 1)...,
	))
	builder.WriteString(eb.applyStyle(`"`, "string") + "\n")
	// fmt.Fprintf(&builder, "\"%s\"\n", escapePOString(str))

	return eb.applyStyle(
		builder.String(),
		append(styles, "string")...,
	)
}

func (eb *entryBuilder) text(str string, styles ...string) string {
	return eb.applyStyle(str, append(styles, "text")...)
}

func (eb *entryBuilder) keyword(kw string) string {
	return eb.applyStyle(kw, "keyword") + " "
}
