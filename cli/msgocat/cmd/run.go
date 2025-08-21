package cmd

import (
	"fmt"
	"os"

	"github.com/Tom5521/gotext-tools/v2/pkg/po"
	"github.com/Tom5521/gotext-tools/v2/pkg/po/compile"
	"github.com/spf13/cobra"
)

// run is the main execution function for the command.
func run(cmd *cobra.Command, args []string) error {
	args, err := expandArgsWithFiles(args)
	if err != nil {
		return err
	}

	state := processState{
		entriesMap:   make(map[string]entryInfo),
		entriesSlice: make(po.Entries, 0),
	}

	for _, file := range args {
		err = mergePoFiles(file, &state)
		if err != nil {
			return fmt.Errorf("error processing PO file (%s): %w", file, err)
		}
	}

	return compile.PoToWriter(
		state.entriesSlice,
		os.Stdout,

		compile.PoWithWordWrap(true),
		compile.PoWithHighlight(compile.DefaultHighlight),
		compile.PoWithOmitHeader(true),
	)
}
