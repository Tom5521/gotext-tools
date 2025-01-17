package cmd

import "github.com/spf13/pflag"

var (
	// Input.
	filesFrom string
	directory string
	// Output.
	defaultDomain   string
	output          string
	outputDir       string
	forcePo         bool
	noLocation      bool
	addLocation     string
	omitHeader      bool
	copyrightHolder string
	packageName     string
	packageVersion  string
	msgstrPrefix    string
	msgstrSuffix    string
	lang            string
	// Operation Mode.
	joinExisting bool
	excludeFile  []string
	addComments  string
	nplurals     uint
	// Other.
	extractAll bool
)

func init() {
	pflag.StringVarP(&lang, "lang", "l", "en", "Language code to include in the POT file")
	pflag.Uint("nplurals", 2, "Specify the number of plurals forms of the language in question")
	pflag.StringVarP(&filesFrom, "files-from", "f", "", "get list of input files from FILE")
	pflag.StringVarP(
		&directory,
		"directory",
		"D",
		"",
		"add DIRECTORY to list for input files search\nIf input file is -, standard input is read.",
	)
	pflag.StringVarP(
		&defaultDomain,
		"default-domain",
		"d",
		"default",
		"use NAME.pot for output (instead of default.pot)",
	)
	pflag.StringVarP(&output, "output", "o", "", "write output to specified file")
	pflag.StringVarP(
		&outputDir,
		"output-dir",
		"p",
		"",
		"output files will be placed in directory DIR\nIf output file is -, output is written to standard output.",
	)
	pflag.BoolVar(
		&forcePo,
		"force-po",
		false,
		"Always write an output file even if no message is defined.",
	)
	pflag.BoolVarP(
		&noLocation,
		"no-location",
		"n",
		false,
		"Do not write ‘#: filename:line’ lines. Note that using this option makes it harder for technically skilled translators to understand each message’s context.",
	)

	pflag.StringVar(
		&addLocation,
		"add-location",
		"full",
		`Generate ‘#: filename:line’ lines (default).

The optional type can be either ‘full’, ‘file’, or ‘never’. If it is not given or ‘full’, it generates the lines with both file name and line number. If it is ‘file’, the line number part is omitted. If it is ‘never’, it completely suppresses the lines (same as --no-location). `,
	)
	pflag.BoolVar(
		&omitHeader,
		"omit-header",
		false,
		`Don’t write header with ‘msgid ""’ entry. Note: Using this option may lead to an error in subsequent operations if the output contains non-ASCII characters.`,
	)
	pflag.StringVar(
		&copyrightHolder,
		"copyright-holder",
		"",
		`Set the copyright holder in the output. string should specify the copyright holder of the surrounding package. (Note that the msgstr strings, extracted from the package’s sources, belong to the copyright holder of the package.) Translators are expected to transfer or disclaim the copyright for their translations, so that package maintainers can distribute them without legal risk.

If string is empty, the copyright holder field is omitted entirely from the output files, leaving no explicit indication of copyright ownership. This implies that translators should take appropriate steps to ensure the distribution is legally permissible, such as disclaiming their copyright.`,
	)
	pflag.StringVar(
		&packageName,
		"package-name",
		"",
		"Set the package name in the header of the output.",
	)
	pflag.StringVar(
		&packageVersion,
		"package-version",
		"",
		"Set the package version in the header of the output. This option has an effect only if the ‘--package-name’ option is also used.",
	)
	pflag.StringVarP(
		&msgstrPrefix,
		"msgstr-prefix",
		"m",
		"",
		`Use string (or "" if not specified) as prefix for msgstr values.`,
	)
	pflag.StringVarP(
		&msgstrSuffix,
		"msgstr-suffix",
		"M",
		"",
		`Use string (or "" if not specified) as suffix for msgstr values.`,
	)
}
