package cmd

import (
	"github.com/Tom5521/xgotext/pkg/parsers"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

var (
	ParserCfg   parsers.Config
	CompilerCfg compiler.Config
	HeadersCfg  types.HeaderConfig
)

var (
	// CLI.

	filesFrom    string
	directory    string
	output       string
	outputDir    string
	joinExisting bool
	excludeFile  string

	// Parser.

	exclude    []string
	extractAll bool

	// Header.

	lang             string
	packageVersion   string
	nplurals         uint
	msgidBugsAddress string

	// Compiler.

	forcePo         bool
	noLocation      bool
	addLocation     string
	omitHeader      bool
	packageName     string
	foreignUser     bool
	title           string
	copyrightHolder string
	msgstrPrefix    string
	msgstrSuffix    string

	// Other.
	defaultDomain string
	// verbose       bool
)

func init() {
	flag := root.Flags()

	flag.StringVar(
		&msgidBugsAddress,
		"msgid-bugs-address",
		"",
		`Set the reporting address for msgid bugs.
This is the email address or URL to which the translators shall report bugs
in the untranslated strings:
    - Strings which are not entire sentences; see the maintainer guidelines
    in Preparing Translatable Strings.
    - Strings which use unclear terms or require additional context to be understood.
    - Strings which make invalid assumptions about notation of date, time or money.
    - Pluralisation problems.
    - Incorrect English spelling.
    - Incorrect formatting. 

It can be your email address, or a mailing list address where translators 
can write to without being subscribed, or the URL of a web page through which
the translators can contact you.

The default value is empty, which means that translators will be clueless!
Don’t forget to specify this option.
`,
	)
	flag.StringVar(&title, "title", "SOME DESCRIPTIVE TITLE", "Set the title of the pot file.")
	flag.StringVar(
		&copyrightHolder,
		"copyright-holder",
		"YEAR THE PACKAGE'S COPYRIGHT HOLDER",
		`Set the copyright holder in the output.
string should be the copyright holder of the surrounding package.
(Note that the msgstr strings, extracted from the package’s sources,
belong to the copyright holder of the package.) Translators are expected to transfer
or disclaim the copyright for their translations, so that package maintainers
can distribute them without legal risk. If string is empty, the output files
are marked as being in the public domain; in this case, the translators are
expected to disclaim their copyright, again so that package maintainers can
distribute them without legal risk.
`,
	)
	flag.BoolVar(
		&foreignUser,
		"foreign-user",
		false,
		`Omit FSF copyright in output. 
This option is equivalent to ‘--copyright-holder=''’. It can be useful
for packages outside the GNU project that want their translations to
be in the public domain.`,
	)
	// flag.BoolVar(&verbose, "verbose", false, "increase verbosity level")
	flag.StringSliceVarP(&exclude, "exclude", "X", nil, "Specifies which files will be omitted.")
	flag.BoolVarP(&extractAll, "extract-all", "a", false, "Extract all strings.")
	flag.StringVarP(
		&excludeFile,
		"exclude-file",
		"x",
		"",
		"Entries from file are not extracted. file should be a PO or POT file.",
	)
	flag.BoolVarP(&joinExisting, "join-existing", "j", false, "Join messages with existing file.")
	flag.StringVarP(&lang, "lang", "l", "en", "Language code to include in the POT file")
	flag.UintVar(
		&nplurals,
		"nplurals",
		2,
		"Specify the number of plurals forms of the language in question",
	)
	flag.StringVarP(&filesFrom, "files-from", "f", "", "get list of input files from FILE")
	flag.StringVarP(
		&directory,
		"directory",
		"D",
		"",
		"add DIRECTORY to list for input files search\nIf input file is -, standard input is read.",
	)
	flag.StringVarP(
		&defaultDomain,
		"default-domain",
		"d",
		"messages",
		"use NAME.pot for output (instead of messages.pot)",
	)
	flag.StringVarP(&output, "output", "o", "", "write output to specified file")
	flag.StringVarP(
		&outputDir,
		"output-dir",
		"p",
		"",
		`output files will be placed in directory DIR
If output file is -, output is written to standard output.`,
	)
	flag.BoolVar(
		&forcePo,
		"force-po",
		false,
		"Always write an output file even if no message is defined.",
	)
	flag.BoolVarP(
		&noLocation,
		"no-location",
		"n",
		false,
		`Do not write ‘#: filename:line’ lines.
Note that using this option makes it harder for technically
skilled translators to understand each message’s context.
`,
	)

	flag.StringVar(
		&addLocation,
		"add-location",
		"full",
		`Generate ‘#: filename:line’ lines (default).

The optional type can be either ‘full’, ‘file’, or ‘never’. 
If it is not given or ‘full’, it generates the lines with both file
name and line number. If it is ‘file’, the line number part is omitted. 
If it is ‘never’, it completely suppresses the lines (same as --no-location).
`,
	)
	flag.BoolVar(
		&omitHeader,
		"omit-header",
		false,
		`Don’t write header with ‘msgid ""’ entry. 
Note: Using this option may lead to an error in subsequent
operations if the output contains non-ASCII characters.
`,
	)
	flag.StringVar(
		&packageName,
		"package-name",
		"PACKAGE",
		"Set the package name in the header of the output.",
	)
	flag.StringVar(
		&packageVersion,
		"package-version",
		"PACKAGE VERSION",
		`Set the package version in the header of the output. 
This option has an effect only if the ‘--package-name’ option is also used.`,
	)
	flag.StringVarP(
		&msgstrPrefix,
		"msgstr-prefix",
		"m",
		"",
		`Use string (or "" if not specified) as prefix for msgstr values.`,
	)
	flag.StringVarP(
		&msgstrSuffix,
		"msgstr-suffix",
		"M",
		"",
		`Use string (or "" if not specified) as suffix for msgstr values.`,
	)
}

func initConfig() {
	ParserCfg = parsers.Config{
		Language:   lang,
		Exclude:    exclude,
		ExtractAll: extractAll,
	}
	CompilerCfg = compiler.Config{
		ForcePo:         forcePo,
		OmitHeader:      omitHeader,
		PackageName:     packageName,
		CopyrightHolder: copyrightHolder,
		ForeignUser:     foreignUser,
		Title:           title,
		NoLocation:      noLocation,
		AddLocation:     addLocation,
		MsgstrPrefix:    msgstrPrefix,
		MsgstrSuffix:    msgstrSuffix,
	}
	HeadersCfg = types.HeaderConfig{
		Nplurals:          nplurals,
		ProjectIDVersion:  packageVersion,
		ReportMsgidBugsTo: msgidBugsAddress,
		Language:          lang,
	}
}
