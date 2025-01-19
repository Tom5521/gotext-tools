package compiler

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/Tom5521/xgotext/internal/util"
	"github.com/Tom5521/xgotext/pkg/poconfig"
	"github.com/Tom5521/xgotext/pkg/poentry"
)

//go:embed header.pot
var baseHeader string

type Compiler struct {
	Translations []poentry.Translation
	Config       poconfig.Config
}

func (c Compiler) CompileToWriter(w io.Writer) error {
	_, err := fmt.Fprintf(
		w,
		baseHeader,
		c.Config.PackageVersion,
		c.Config.Language,
		c.Config.Nplurals,
	)
	if err != nil {
		return err
	}

	translations := util.CleanDuplicates(c.Translations)
	for _, t := range translations {
		_, err = fmt.Fprintln(w, t.Format(c.Config.Nplurals))
		if err != nil {
			return err
		}

	}
	return nil
}

func (c Compiler) CompileToFile(f string) error {
	flags := os.O_RDWR
	if c.Config.ForcePo {
		flags = flags | os.O_CREATE
	}
	file, err := os.OpenFile(f, flags, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	if c.Config.ForcePo {
		err = file.Truncate(0)
		if err != nil {
			return err
		}
	}

	return c.CompileToWriter(file)
}

func (c Compiler) CompileToString() string {
	var b strings.Builder

	c.CompileToWriter(&b)

	return b.String()
}

func (c Compiler) CompileToBytes() []byte {
	var b bytes.Buffer

	c.CompileToWriter(&b)

	return b.Bytes()
}

var ErrNotImplementedYet = errors.New("not implemented yet (sorry)")

// TODO: Implement domains.
func (c Compiler) CompileToDir(d string) error    { return ErrNotImplementedYet }
func (c Compiler) CompileToDirFS(dfs fs.FS) error { return ErrNotImplementedYet }
