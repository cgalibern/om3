package commands

import (
	"context"
	"fmt"

	"github.com/opensvc/om3/core/actioncontext"
	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/nodeaction"
	"github.com/opensvc/om3/core/nodeselector"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/core/output"
	"github.com/opensvc/om3/core/rawconfig"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/util/hostname"
)

type (
	CmdNodeEval struct {
		OptsGlobal
		OptsLock
		Keywords     []string
		Impersonate  string
		NodeSelector string
	}
)

func (t *CmdNodeEval) Run() error {
	if t.Local {
		return t.doNodeAction()
	}
	c, err := client.New()
	if err != nil {
		return err
	}

	nodenames := []string{hostname.Hostname()}
	if t.NodeSelector != "" {
		sel := nodeselector.New(t.NodeSelector)
		if l, err := sel.Expand(); err != nil {
			return err
		} else {
			nodenames = l
		}
	}
	l := make(api.KeywordItems, 0)
	for _, nodename := range nodenames {
		params := api.GetNodeConfigGetParams{}
		params.Kw = &t.Keywords
		v := true
		params.Evaluate = &v
		if t.Impersonate != "" {
			params.Impersonate = &t.Impersonate
		}
		response, err := c.GetNodeConfigGetWithResponse(context.Background(), nodename, &params)
		if err != nil {
			return err
		}
		switch {
		case response.JSON200 != nil:
			l = append(l, response.JSON200.Items...)
		case response.JSON400 != nil:
			return fmt.Errorf("%s: %s", nodename, *response.JSON400)
		case response.JSON401 != nil:
			return fmt.Errorf("%s: %s", nodename, *response.JSON401)
		case response.JSON403 != nil:
			return fmt.Errorf("%s: %s", nodename, *response.JSON403)
		case response.JSON500 != nil:
			return fmt.Errorf("%s: %s", nodename, *response.JSON500)
		default:
			return fmt.Errorf("%s: unexpected response: %s", nodename, response.Status())
		}
	}
	defaultOutput := "tab=NODE:meta.node,KEYWORD:meta.keyword,VALUE:data.value,EVALUATED_AS:meta.evaluated_as"
	output.Renderer{
		DefaultOutput: defaultOutput,
		Output:        t.Output,
		Color:         t.Color,
		Data:          l,
		Colorize:      rawconfig.Colorize,
	}.Print()

	return nil
}

func (t *CmdNodeEval) doNodeAction() error {
	return nodeaction.New(
		nodeaction.WithLocal(t.Local),
		nodeaction.WithFormat(t.Output),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithLocalFunc(func() (interface{}, error) {
			n, err := object.NewNode()
			if err != nil {
				return nil, err
			}
			ctx := context.Background()
			ctx = actioncontext.WithLockDisabled(ctx, t.Disable)
			ctx = actioncontext.WithLockTimeout(ctx, t.Timeout)
			for _, s := range t.Keywords {
				return n.EvalAs(ctx, s, t.Impersonate)
			}
			return nil, nil
		}),
	).Do()
}
