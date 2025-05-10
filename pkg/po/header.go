package po

import (
	"errors"
	"fmt"
	"mime"
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/gotext-tools/v2/internal/slices"
	"github.com/Tom5521/gotext-tools/v2/internal/util"
)

// HeaderField represents a single key-value pair in a header.
type HeaderField struct {
	Key   string // The name of the header field.
	Value string // The value associated with the header field.
}

// Header represents a collection of header fields.
type Header struct {
	Template bool
	Fields   []HeaderField // A slice storing all registered header fields.
}

func NewHeader(options ...HeaderOption) Header {
	h := HeaderConfig{}
	for _, opt := range options {
		opt(&h)
	}

	return h.ToHeader()
}

var (
	headerRegex = regexp.MustCompile(`(.*)\s*:\s*(.*)`)
	exprRegex   = regexp.MustCompile(
		`(?: *(\S+?) *= *(.+?) *; *)`,
	)
)

func parseAdvHeaderField(header string) map[string]string {
	keyValueMap := make(map[string]string)
	header = strings.TrimSpace(header)
	if !strings.HasSuffix(header, ";") {
		header += ";"
	}
	matches := exprRegex.FindAllStringSubmatch(header, -1)
	for _, match := range matches {
		key := match[len(match)-2]
		value := match[len(match)-1]

		keyValueMap[key] = value
	}

	return keyValueMap
}

func (h Header) String() string {
	return util.Format(h)
}

func (h HeaderConfig) String() string {
	return util.Format(h)
}

func (h Header) ToConfig() HeaderConfig {
	mimeVersion := h.Load("MIME-Version")
	if mimeVersion == "" {
		mimeVersion = "1.0"
	}

	var (
		charset   = "UTF-8"
		mediatype = "text/plain"
	)

	mtype, params, err := mime.ParseMediaType(h.Load("Content-Type"))
	if err == nil {
		mediatype = mtype
		chset := params["charset"]
		if util.SupportedCharsets[chset] {
			charset = chset
		}
	}

	var nplurals uint = 2
	var plural string
	{
		pluralForms := h.Load("Plural-Forms")
		matches := parseAdvHeaderField(pluralForms)
		if npluralsStr, ok := matches["nplurals"]; ok {
			np, err := strconv.ParseUint(npluralsStr, 10, strconv.IntSize)
			if err == nil {
				nplurals = uint(np)
			}
		}
		if pluralStr, ok := matches["plural"]; ok {
			plural = pluralStr
		}
	}

	return HeaderConfig{
		Template:                h.Template,
		ProjectIDVersion:        h.Load("Project-Id-Version"),
		ReportMsgidBugsTo:       h.Load("Report-Msgid-Bugs-To"),
		POTCreationDate:         h.Load("POT-Creation-Date"),
		PORevisionDate:          h.Load("PO-Revision-Date"),
		LastTranslator:          h.Load("Last-Translator"),
		LanguageTeam:            h.Load("Language-Team"),
		Language:                h.Load("Language"),
		MimeVersion:             mimeVersion,
		MediaType:               mediatype,
		Charset:                 charset,
		ContentTransferEncoding: "8bit",
		Plural:                  plural,
		Nplurals:                nplurals,
		XGenerator:              h.Load("XGenerator"),
	}
}

func (h Header) Nplurals() (nplurals uint) {
	nplurals = 2
	value := h.Load("Plural-Forms")

	if np, ok := parseAdvHeaderField(value)["nplurals"]; ok {
		n, err := strconv.ParseUint(np, 10, 64)
		if err != nil {
			return
		}

		nplurals = uint(n)
	}

	return
}

func (h Header) ToEntry() Entry {
	var b strings.Builder

	for _, field := range h.Fields {
		if field.Value != " " {
			field.Value = " " + field.Value
		}
		fmt.Fprintf(&b, "\n%s:%s\n", field.Key, field.Value)
	}
	entry := Entry{Str: b.String()}
	if h.Template {
		entry.markAsFuzzy()
	}

	return entry
}

type HeaderConfig struct {
	Template                bool
	ProjectIDVersion        string
	ReportMsgidBugsTo       string
	POTCreationDate         string
	PORevisionDate          string
	LastTranslator          string
	LanguageTeam            string
	Language                string
	MimeVersion             string
	MediaType               string
	Charset                 string
	ContentTransferEncoding string
	Nplurals                uint
	Plural                  string
	XGenerator              string
}

func (h HeaderConfig) Validate() []error {
	var errs []error

	if h.Plural != "" || h.Nplurals != 0 {
		if h.Plural == "" {
			errs = append(errs, errors.New("plural not specified"))
		}
		if h.Nplurals == 0 {
			errs = append(errs, errors.New("nplurals can't be zero"))
		}
	}

	if h.MediaType != "text/plain" && h.MediaType != "" {
		errs = append(errs, fmt.Errorf("media type (%s) must be text/plain", h.MediaType))
	}
	if h.ContentTransferEncoding != "8bit" && h.ContentTransferEncoding != "" {
		errs = append(
			errs,
			fmt.Errorf("content-transfer-encoding(%s) must be 8bit", h.ContentTransferEncoding),
		)
	}
	if !util.SupportedCharsets[h.Charset] && h.Charset != "" {
		errs = append(errs, fmt.Errorf("%q isn't a supported charset", h.Charset))
	}

	return errs
}

func (cfg HeaderConfig) ToHeader() (h Header) {
	h.Template = cfg.Template
	h.sSet("Project-Id-Version", cfg.ProjectIDVersion)
	h.sSet("Report-Msgid-Bugs-To", cfg.ReportMsgidBugsTo)
	if cfg.Template {
		h.sSet("POT-Creation-Date", cfg.POTCreationDate)
	}
	h.sSet("PO-Revision-Date", cfg.PORevisionDate)
	h.sSet("Last-Translator", cfg.LastTranslator)
	h.sSet("Language-Team", cfg.LanguageTeam)
	h.sSet("Language", cfg.Language)
	h.sSet("MIME-Version", cfg.MimeVersion)
	if cfg.MediaType != "" && cfg.Charset != "" {
		h.sSet(
			"Content-Type",
			mime.FormatMediaType(cfg.MediaType, map[string]string{"charset": cfg.Charset}),
		)
	}
	h.sSet("Content-Transfer-Encoding", cfg.ContentTransferEncoding)
	if cfg.Nplurals != 0 && cfg.Plural != "" && !cfg.Template {
		h.sSet(
			"Plural-Forms",
			fmt.Sprintf("nplurals=%d; plural=%s;", cfg.Nplurals, cfg.Plural),
		)
	}
	h.sSet("X-Generator", cfg.XGenerator)

	return
}

func HeaderConfigFromOptions(options ...HeaderOption) HeaderConfig {
	var h HeaderConfig
	for _, opt := range options {
		opt(&h)
	}

	return h
}

var defaultHeaderConfig = HeaderConfig{
	Language:                " ",
	MimeVersion:             "1.0",
	MediaType:               "text/plain",
	Charset:                 "UTF-8",
	ContentTransferEncoding: "8bit",
	Nplurals:                2,
	Plural:                  "(n != 1)",
}

func DefaultHeaderConfig(opts ...HeaderOption) HeaderConfig {
	h := defaultHeaderConfig
	for _, ho := range opts {
		ho(&h)
	}
	return h
}

var defaultTemplateHeaderConfig = HeaderConfig{
	Template:                true,
	ProjectIDVersion:        "PACKAGE VERSION",
	ReportMsgidBugsTo:       " ",
	PORevisionDate:          "YEAR-MO-DA HO:MI+ZONE",
	LastTranslator:          "FULL NAME <EMAIL@ADDRESS>",
	LanguageTeam:            "LANGUAGE <LL@li.org>",
	Language:                " ",
	MimeVersion:             "1.0",
	MediaType:               "text/plain",
	Charset:                 "CHARSET",
	ContentTransferEncoding: "8bit",
}

func DefaultTemplateHeaderConfig(opts ...HeaderOption) HeaderConfig {
	h := defaultTemplateHeaderConfig
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (e Entries) Header() (h Header) {
	i := e.Index("", "")
	if i == -1 {
		return
	}
	entry := e[i]

	h.Template = entry.IsFuzzy()
	header := entry.Str
	lines := strings.Split(header, "\n")
	for _, line := range lines {
		if !headerRegex.MatchString(line) {
			continue
		}
		matches := headerRegex.FindStringSubmatch(line)
		h.Fields = append(h.Fields,
			HeaderField{
				Key:   matches[1],
				Value: matches[2],
			},
		)
	}
	return
}

// Register adds a new header field to the Header object if the key does not already exist.
// Parameters:
//   - key: The name of the header field to register.
//   - d: Optional variadic arguments representing the value(s) to associate with the key.
//     If provided, they are concatenated into a single string using fmt.Sprint.
func (h *Header) Register(key string, d ...string) {
	// Check if the key already exists in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i != -1 {
		// Key already exists; do nothing.
		return
	}

	var values []any
	for _, b := range d {
		values = append(values, b)
	}

	// Append a new HeaderField to the Values slice.
	h.Fields = append(h.Fields,
		HeaderField{
			Key:   key,
			Value: fmt.Sprint(values...), // Concatenate variadic arguments into a single string.
		},
	)
}

// Load retrieves the value associated with a given key from the Header object.
// Parameters:
// - key: The name of the header field to retrieve.
// Returns:
// - The value associated with the key if found; otherwise, an empty string ("").
func (h *Header) Load(key string) string {
	// Search for the key in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i >= 0 {
		// Key found; return its value.
		return h.Fields[i].Value
	}
	// Key not found; return an empty string.
	return ""
}

// Set updates the value of an existing header field or adds a new field if the key does not exist.
// Parameters:
// - key: The name of the header field to update or add.
// - value: The new value to associate with the key.
func (h *Header) Set(key, value string) {
	// Search for the key in the Values slice.
	i := slices.IndexFunc(h.Fields, func(f HeaderField) bool {
		return f.Key == key
	})
	if i >= 0 {
		// Key found; update its value.
		h.Fields[i].Value = value
		return
	}
	// Key not found; append a new HeaderField to the Values slice.
	h.Fields = append(h.Fields,
		HeaderField{
			Key:   key,
			Value: value,
		},
	)
}

func (h *Header) sSet(key, value string) {
	if value == "" || key == "" {
		return
	}
	h.Set(key, value)
}
