package cmd

import (
	"os"

	goparse "github.com/Tom5521/xgotext/pkg/go/parse"
	"github.com/Tom5521/xgotext/pkg/po/compiler"
	poparse "github.com/Tom5521/xgotext/pkg/po/parse"
	"github.com/Tom5521/xgotext/pkg/po/types"
)

func join(newParse *goparse.Parser, rawfile *os.File) error {
	baseParse, err := poparse.NewParserFromReader(rawfile, rawfile.Name(), cfg)
	if err != nil {
		return err
	}

	base, _, errs := baseParse.Parse()
	if len(errs) > 0 {
		return errs[0]
	}

	parsed, errs := newParse.Parse()
	if len(errs) > 0 {
		return errs[0]
	}

	merged := types.MergeFiles(base, parsed)

	compiler := compiler.Compiler{
		File:   merged,
		Config: cfg,
	}

	err = rawfile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = rawfile.Seek(0, 0)
	if err != nil {
		return err
	}

	err = compiler.CompileToWriter(rawfile)
	if err != nil {
		return err
	}

	return nil
}
