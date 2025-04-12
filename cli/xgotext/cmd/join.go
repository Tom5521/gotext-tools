package cmd

import (
	"os"

	goparse "github.com/Tom5521/gotext-tools/pkg/go/parse"
	"github.com/Tom5521/gotext-tools/pkg/po"
	"github.com/Tom5521/gotext-tools/pkg/po/compile"
	poparse "github.com/Tom5521/gotext-tools/pkg/po/parse"
)

func join(newParse *goparse.Parser, rawfile *os.File) error {
	baseParse, err := poparse.NewPoFromReader(
		rawfile,
		rawfile.Name(),
		poparse.PoWithConfig(PoParserCfg),
	)
	if err != nil {
		return err
	}

	poParsed := baseParse.Parse()
	if len(baseParse.Errors()) > 0 {
		return baseParse.Errors()[0]
	}

	goParsed := newParse.Parse()
	if len(newParse.Errors()) > 0 {
		return newParse.Errors()[0]
	}

	poParsed.Entries = po.Merge(poParsed.Entries, goParsed.Entries)

	compiler := compile.NewPo(poParsed, compile.PoWithConfig(CompilerCfg))

	// Truncate file.
	rawfile, err = os.Create(rawfile.Name())
	if err != nil {
		return err
	}

	err = compiler.ToWriter(rawfile)
	if err != nil {
		return err
	}

	return nil
}
