package cmd

var (
	// Input.
	filesFrom string
	directory string
	exclude   []string
	// Output.
	defaultDomain  string
	output         string
	outputDir      string
	forcePo        bool
	noLocation     bool
	addLocation    string
	omitHeader     bool
	packageName    string
	packageVersion string
	msgstrPrefix   string
	msgstrSuffix   string
	lang           string
	// Operation Mode.
	joinExisting bool
	excludeFile  string
	nplurals     uint
	// Other.
	extractAll bool
	verbose    bool
)

func init() {
	flag := root.Flags()

	flag.BoolVar(&verbose, "verbose", false, "increase verbosity level")
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
		"default",
		"use NAME.pot for output (instead of default.pot)",
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
		"Do not write ‘#: filename:line’ lines. Note that using this option makes it harder for technically skilled translators to understand each message’s context.",
	)

	flag.StringVar(
		&addLocation,
		"add-location",
		"full",
		`Generate ‘#: filename:line’ lines (default).

The optional type can be either ‘full’, ‘file’, or ‘never’. If it is not given or ‘full’, it generates the lines with both file name and line number. If it is ‘file’, the line number part is omitted. If it is ‘never’, it completely suppresses the lines (same as --no-location). `,
	)
	flag.BoolVar(
		&omitHeader,
		"omit-header",
		false,
		`Don’t write header with ‘msgid ""’ entry. Note: Using this option may lead to an error in subsequent operations if the output contains non-ASCII characters.`,
	)
	flag.StringVar(
		&packageName,
		"package-name",
		"",
		"Set the package name in the header of the output.",
	)
	flag.StringVar(
		&packageVersion,
		"package-version",
		"",
		"Set the package version in the header of the output. This option has an effect only if the ‘--package-name’ option is also used.",
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
