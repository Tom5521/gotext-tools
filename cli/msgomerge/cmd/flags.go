package cmd

var (
	directory  string
	update     bool
	outputFile string
	// TODO: Finish this.
	// backup           string
	// suffix           string
	noFuzzyMatching bool
	previous        bool
	lang            string
	forcePo         bool
	noLocation      bool
	addLocation     string
)

func init() {
	flags := root.Flags()

	flags.StringVarP(
		&directory,
		"directory",
		"D",
		"",
		`add DIRECTORY to list for input files search`,
	)
	flags.BoolVarP(&update, "update", "U", false, `update def.po,
do nothing if def.po already up to date`)
	flags.StringVarP(&outputFile, "output-file", "o", "-", `write output to specified file
The results are written to standard output if no output file is specified
or if it is -.`)
	flags.BoolVarP(&noFuzzyMatching, "no-fuzzy-matching", "N", false, `do not use fuzzy matching`)
	flags.BoolVar(&previous, "previous", false, "keep previous msgids of translated messages")
	flags.StringVar(&lang, "lang", "en", "set 'Language' field in the header entry")
	flags.BoolVar(&forcePo, "force-po", false, "write PO file even if empty")
	flags.BoolVar(&noLocation, "no-location", false, "suppress '#: filename:line' lines")
	flags.StringVarP(&addLocation, "add-location", "n", "full", `
Generate ‘#: filename:line’ lines (default).

The optional type can be either ‘full’, ‘file’, or ‘never’. 
If it is not given or ‘full’, it generates the lines with both
file name and line number. If it is ‘file’, the line number part is omitted.
If it is ‘never’, it completely suppresses the lines (same as --no-location).`)
}
