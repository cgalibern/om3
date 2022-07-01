package commands

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/nodeaction"
	"opensvc.com/opensvc/core/object"
)

type (
	// CmdNodeComplianceFix is the cobra flag set of the sysreport command.
	CmdNodeComplianceFix struct {
		OptsGlobal
		object.OptsNodeComplianceFix
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *CmdNodeComplianceFix) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *CmdNodeComplianceFix) cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fix",
		Short: "run compliance fixes",
		Run: func(_ *cobra.Command, _ []string) {
			t.run()
		},
	}
}

func (t *CmdNodeComplianceFix) run() {
	nodeaction.New(
		nodeaction.WithLocal(t.Local),
		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithFormat(t.Format),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithRemoteAction("compliance fix"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"format":    t.Format,
			"force":     t.Force,
			"module":    t.Module,
			"moduleset": t.Moduleset,
			"attach":    t.Attach,
		}),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			return object.NewNode().ComplianceFix(t.OptsNodeComplianceFix)
		}),
	).Do()
}
