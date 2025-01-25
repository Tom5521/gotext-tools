package parse

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Tom5521/xgotext/pkg/po/config"
)

type File struct {
	pos     int
	literal string
	Name    string
	Nodes   []Node

	containsHeader bool
	header         map[string]string
}

func (f File) Header() map[string]string {
	if f.header == nil {
		f.processHeader()
	}
	return f.header
}

var headerRegex = regexp.MustCompile(`(.*)\s*:\s*(.*)`)

func (f *File) processHeader() {
	f.header = make(map[string]string)
	for i, node := range f.Nodes {
		n, ok := node.(Msgid)
		if !ok {
			continue
		}
		if n.ID != "" {
			continue
		}

		f.containsHeader = true

		msgstr := f.Nodes[i+1].(Msgstr)
		lines := strings.Split(msgstr.Str, "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			matches := headerRegex.FindStringSubmatch(line)
			f.header[matches[1]] = matches[2]
		}

		break
	}
}

var npluralsRegex = regexp.MustCompile(`nplurals=(\d*)`)

func (f *File) Config() config.Config {
	if f.header == nil {
		f.processHeader()
	}

	if !f.containsHeader {
		return config.Default()
	}

	get := func(key string) string {
		value, ok := f.header[key]
		if !ok {
			return ""
		}
		return value
	}

	cfg := config.Config{
		PackageVersion:   get("Project-Id-Version"),
		MsgidBugsAddress: get("Report-Msgid-Bugs-To"),
		PotCreationDate:  get("POT-Creation-Date"),
		Language:         get("Language"),
	}

	pluralForms := get("Plural-Forms")
	if npluralsRegex.MatchString(pluralForms) {
		matches := npluralsRegex.FindStringSubmatch(pluralForms)
		nplurals, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return cfg
		}

		cfg.Nplurals = uint(nplurals)
	}

	return cfg
}

// func (f File) Translations() []entry.Translation {}

func (f File) Pos() int        { return f.pos }
func (f File) Literal() string { return f.literal }
