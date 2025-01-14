package main

import (
	"os"

	"github.com/Tom5521/xgotext/flags"
	"github.com/Tom5521/xgotext/parser"
	"github.com/gookit/color"
	"github.com/spf13/pflag"
)

func main() {
	help := pflag.Bool("help", false, "Print this message.")
	pflag.Parse()

	if *help {
		pflag.PrintDefaults()
		return
	}

	parser, err := parser.NewParser(flags.Input)
	if err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
	parser.Parse()

	err = os.WriteFile(flags.Output, parser.Compile(), 0o600)
	if err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
}
