package cmd

import (
	"fmt"
	"os"

	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	poparse "github.com/Tom5521/xgotext/pkg/po/parse"
)

func join(newParse *goparse.Parser, rawfile *os.File) error {
	baseParse, err := poparse.NewParserFromReader(
		rawfile,
		rawfile.Name(),
		poparse.WithConfig(PoParserCfg),
	)
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

	po.MergeFiles(false, base, parsed)

	compiler := compiler.New(base, compiler.WithConfig(CompilerCfg))

	err = rawfile.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating file %s: %w", rawfile.Name(), err)
	}

	_, err = rawfile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("error seeking file(%s) offset: %w", rawfile.Name(), err)
	}

	err = compiler.ToWriter(rawfile)
	if err != nil {
		return err
	}

	return nil
}
