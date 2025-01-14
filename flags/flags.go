package flags

import "github.com/spf13/pflag"

var (
	Input          string
	Output         string
	ProjectVersion string
	Language       string
	Nplurals       uint

	Exclude []string
	Verbose bool
)

func init() {
	pflag.StringVar(&Input, "input", ".", "Input file or directory path")
	pflag.StringVar(&Output, "output", "default.pot", "Output POT file path")
	pflag.StringVar(
		&ProjectVersion,
		"project-version",
		"",
		"Project version to include in the POT file",
	)
	pflag.StringVar(&Language, "lang", "en", "Language code to include in the POT file")
	pflag.UintVar(
		&Nplurals,
		"nplurals",
		2,
		"Specify the number of plurals of the language in question (default 2)",
	)
	pflag.StringSliceVar(
		&Exclude,
		"exclude",
		nil,
		"Files or directories to exclude from processing (can be specified multiple times)",
	)
	pflag.BoolVar(&Verbose, "verbose", false, "Enable verbose output")
}
