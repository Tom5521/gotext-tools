package main

import "github.com/spf13/pflag"

var (
	input          string
	output         string
	proyectVersion string
	language       string

	exclude []string
	verbose bool
)

func init() {
	pflag.StringVar(&input, "input", ".", "Input file or directory path")
	pflag.StringVar(&output, "output", "default.pot", "Output POT file path")
	pflag.StringVar(
		&proyectVersion,
		"proyect-version",
		"",
		"Project version to include in the POT file",
	)
	pflag.StringVar(&language, "lang", "", "Language code to include in the POT file")
	pflag.StringSliceVar(
		&exclude,
		"exclude",
		nil,
		"Files or directories to exclude from processing (can be specified multiple times)",
	)
	pflag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
}
