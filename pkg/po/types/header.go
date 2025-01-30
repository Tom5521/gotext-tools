package types

import (
	"slices"
	"time"
)

type HeaderField struct {
	Key   string
	Value string
}

type Header struct {
	Values []HeaderField
}

func DefaultHeader() (h Header) {
	h.Register("Project-Id-Version")
	h.Register("Report-Msgid-Bugs-To")
	h.Register("POT-Creation-Date", time.Now().Format(time.DateTime))
	h.Register("PO-Revision-Date")
	h.Register("Last-Translator")
	h.Register("Language-Team")
	h.Register("Language")
	h.Register("MIME-Version", "1.0")
	h.Register("Content-Type", "text/plain; charset=CHARSET")
	h.Register("Content-Transfer-Encoding", "8bit")
	h.Register("Plural-Forms", "nplurals=%d; plural=(n != 1);")

	return h
}

func (h *Header) Register(key string, d ...string) {
	i := slices.IndexFunc(h.Values, func(f HeaderField) bool {
		return f.Key == key
	})

	if i != -1 {
		return
	}

	var v string
	if len(d) > 0 {
		v = d[0]
	}

	h.Values = append(h.Values,
		HeaderField{
			Key:   key,
			Value: v,
		},
	)
}

func (h *Header) Load(key string) string {
	i := slices.IndexFunc(h.Values, func(f HeaderField) bool {
		return f.Key == key
	})

	if i >= 0 {
		return h.Values[i].Value
	}
	return ""
}

func (h *Header) Set(key, value string) {
	i := slices.IndexFunc(h.Values, func(f HeaderField) bool {
		return f.Key == key
	})

	if i >= 0 {
		h.Values[i].Value = value
		return
	}
	h.Values = append(h.Values,
		HeaderField{
			Key:   key,
			Value: value,
		},
	)
}
