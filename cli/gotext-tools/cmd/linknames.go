package cmd

import (
	_ "unsafe"

	_ "github.com/Tom5521/gotext-tools/v2/cli/msgofmt/cmd"
	_ "github.com/Tom5521/gotext-tools/v2/cli/msgomerge/cmd"
	_ "github.com/Tom5521/gotext-tools/v2/cli/msgounfmt/cmd"
	_ "github.com/Tom5521/gotext-tools/v2/cli/xgotext/cmd"
	"github.com/spf13/cobra"
)

//go:linkname msgofmt github.com/Tom5521/gotext-tools/v2/cli/msgofmt/cmd.root
//go:linkname msgomerge github.com/Tom5521/gotext-tools/v2/cli/msgomerge/cmd.root
//go:linkname xgotext github.com/Tom5521/gotext-tools/v2/cli/xgotext/cmd.root
//go:linkname msgounfmt github.com/Tom5521/gotext-tools/v2/cli/msgounfmt/cmd.root

var (
	msgofmt   *cobra.Command
	msgomerge *cobra.Command
	msgounfmt *cobra.Command
	xgotext   *cobra.Command
)
