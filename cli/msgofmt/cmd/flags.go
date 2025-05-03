package cmd

import (
	"log"

	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
)

var (
	directory   string
	output      string
	endianness  string
	force       bool
	noHashTable bool
	verbose     bool
)

func init() {
	flags := root.Flags()
	flags.StringVarP(&directory, "directory", "D", "",
		`add DIRECTORY to list for input files search`)
	flags.StringVarP(&output, "output-file", "o", "messages.mo",
		`write output to specified file
If output file is -, output is written to standard output.`)
	flags.StringVar(&endianness, "endianness", "native",
		`write out 32-bit numbers in the given byte order
(big or little, default depends on platform)`,
	)
	flags.BoolVarP(&force, "force", "f", false, "Overwrites generated files if they already exist")
	flags.BoolVar(&noHashTable, "no-hash", false, "binary file will not include the hash table")
	flags.BoolVar(&verbose, "verbose", false, "")
}

var compilerCfg = compile.DefaultMoConfig()

func initCfg() {
	compilerCfg.Endianness = func() compile.Endianness {
		switch endianness {
		case "big":
			return compile.BigEndian
		case "little":
			return compile.LittleEndian
		default:
			return compile.NativeEndian
		}
	}()
	compilerCfg.Force = force
	compilerCfg.HashTable = !noHashTable
	compilerCfg.Logger = log.Default()
	compilerCfg.Verbose = verbose
}
