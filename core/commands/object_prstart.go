package commands

import (
	"context"

	"github.com/spf13/cobra"
	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/flag"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/objectaction"
	"opensvc.com/opensvc/core/path"
)

type (
	// CmdObjectPRStart is the cobra flag set of the start command.
	CmdObjectPRStart struct {
		OptsGlobal
		OptsLock
		OptsResourceSelector
		OptDryRun
		OptTo
		OptForce
	}
)

// Init configures a cobra command and adds it to the parent command.
func (t *CmdObjectPRStart) Init(kind string, parent *cobra.Command, selector *string) {
	cmd := t.cmd(kind, selector)
	parent.AddCommand(cmd)
	flag.Install(cmd, t)
}

func (t *CmdObjectPRStart) cmd(kind string, selector *string) *cobra.Command {
	return &cobra.Command{
		Use:     "prstart",
		Short:   "preempt devices exclusive write access reservation",
		Aliases: []string{"prov"},
		Run: func(cmd *cobra.Command, args []string) {
			t.run(selector, kind)
		},
	}
}

func (t *CmdObjectPRStart) run(selector *string, kind string) {
	mergedSelector := mergeSelector(*selector, t.ObjectSelector, kind, "")
	objectaction.New(
		objectaction.WithObjectSelector(mergedSelector),
		objectaction.WithRID(t.RID),
		objectaction.WithTag(t.Tag),
		objectaction.WithSubset(t.Subset),
		objectaction.WithLocal(true),
		objectaction.WithFormat(t.Format),
		objectaction.WithColor(t.Color),
		objectaction.WithRemoteNodes(t.NodeSelector),
		objectaction.WithRemoteAction("prstart"),
		objectaction.WithDigest(),
		objectaction.WithLocalRun(func(p path.T) (any, error) {
			o, err := object.NewActor(p)
			if err != nil {
				return nil, err
			}
			ctx := context.Background()
			ctx = actioncontext.WithLockDisabled(ctx, t.Disable)
			ctx = actioncontext.WithLockTimeout(ctx, t.Timeout)
			ctx = actioncontext.WithRID(ctx, t.RID)
			ctx = actioncontext.WithTag(ctx, t.Tag)
			ctx = actioncontext.WithSubset(ctx, t.Subset)
			ctx = actioncontext.WithTo(ctx, t.To)
			ctx = actioncontext.WithForce(ctx, t.Force)
			ctx = actioncontext.WithDryRun(ctx, t.DryRun)
			return nil, o.PRStart(ctx)
		}),
	).Do()
}