package cmd

import (
	"os"

	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	poparse "github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

func join(newParse *goparse.Parser, rawfile *os.File) error {
	baseParse, err := poparse.NewParserFromReader(rawfile, rawfile.Name())
	if err != nil {
		return err
	}

	base := baseParse.Parse()
	if len(baseParse.Errors()) > 0 {
		return baseParse.Errors()[0]
	}

	parsed := newParse.Parse()
	if len(newParse.Errors()) > 0 {
		return newParse.Errors()[0]
	}

	types.MergeFiles(false, base, parsed)

	compiler := compiler.New(base, compiler.WithConfig(CompilerCfg))

	err = rawfile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = rawfile.Seek(0, 0)
	if err != nil {
		return err
	}

	err = compiler.ToWriter(rawfile)
	if err != nil {
		return err
	}

	return nil
}
