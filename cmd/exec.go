package cmd

import (
	"github.com/spf13/cobra"
	"github.com/opensvc/om3/util/capexec"
)

var (
	xo capexec.T

	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Execute a command with cappings and limits",
		Run: func(_ *cobra.Command, args []string) {
			xo.Exec(args)
		},
	}
)

func init() {
	root.AddCommand(execCmd)
	flags := execCmd.Flags()
	xo.FlagSet(flags)
}
