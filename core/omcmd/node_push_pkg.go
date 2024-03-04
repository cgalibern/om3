package commands

import (
	"context"
	"fmt"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/nodeaction"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/util/xsession"
)

type (
	CmdNodePushPkg struct {
		OptsGlobal
		NodeSelector string
	}
)

func (t *CmdNodePushPkg) Run() error {
	return nodeaction.New(
		nodeaction.WithLocal(t.Local),
		nodeaction.WithRemoteNodes(t.NodeSelector),
		nodeaction.WithFormat(t.Output),
		nodeaction.WithColor(t.Color),
		nodeaction.WithServer(t.Server),
		nodeaction.WithLocalFunc(func() (interface{}, error) {
			n, err := object.NewNode()
			if err != nil {
				return nil, err
			}
			l, err := n.PushPkg()
			if err != nil {
				return nil, err
			}
			fmt.Printf("Pushed %d installed packages information.", len(l))
			return nil, nil
		}),
		nodeaction.WithRemoteFunc(func(ctx context.Context, nodename string) (interface{}, error) {
			c, err := client.New(client.WithURL(t.Server))
			if err != nil {
				return nil, err
			}
			params := api.PostNodeActionPushPkgParams{}
			{
				sid := xsession.ID
				params.RequesterSid = &sid
			}
			response, err := c.PostNodeActionPushPkgWithResponse(ctx, nodename, &params)
			if err != nil {
				return nil, err
			}
			switch {
			case response.JSON200 != nil:
				return *response.JSON200, nil
			case response.JSON401 != nil:
				return nil, fmt.Errorf("node %s: %s", nodename, *response.JSON401)
			case response.JSON403 != nil:
				return nil, fmt.Errorf("node %s: %s", nodename, *response.JSON403)
			case response.JSON500 != nil:
				return nil, fmt.Errorf("node %s: %s", nodename, *response.JSON500)
			default:
				return nil, fmt.Errorf("node %s: unexpected response: %s", nodename, response.Status())
			}
		}),
	).Do()
}
