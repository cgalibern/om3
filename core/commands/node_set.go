package commands

import (
	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/nodeaction"
	"opensvc.com/opensvc/core/object"
)

type (
	// NodeSet is the cobra flag set of the start command.
	NodeSet struct {
		OptsGlobal
		object.OptsSet
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *NodeSet) Init(parent *cobra.Command) {
	cmd := t.cmd()
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *NodeSet) cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "set a configuration key value",
		Run: func(_ *cobra.Command, _ []string) {
			t.run()
		},
	}
}

func (t *NodeSet) run() {
	nodeaction.New(
		nodeaction.LocalFirst(),
		nodeaction.WithLocal(t.Local),
		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithFormat(t.Format),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithRemoteAction("set"),
		nodeaction.WithRemoteOptions(map[string]interface{}{
			"kw": t.KeywordOps,
		}),
		nodeaction.WithLocalRun(func() (interface{}, error) {
			return nil, object.NewNode().Set(t.OptsSet)
		}),
	).Do()
}
