package oxcmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/nodeselector"
	"github.com/opensvc/om3/core/output"
	"github.com/opensvc/om3/core/rawconfig"
	"github.com/opensvc/om3/daemon/api"
)

type (
	CmdNodeDrivers struct {
		OptsGlobal
		NodeSelector string
	}
)

func (t *CmdNodeDrivers) Run() error {
	c, err := client.New()
	if err != nil {
		return err
	}

	if t.NodeSelector == "" {
		t.NodeSelector = "*"
	}

	sel := nodeselector.New(t.NodeSelector, nodeselector.WithClient(c))
	nodenames, err := sel.Expand()
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*5)

	l := make(api.DriverItems, 0)
	q := make(chan api.DriverItems)
	errC := make(chan error)
	doneC := make(chan string)
	todo := len(nodenames)

	for _, nodename := range nodenames {
		go func(nodename string) {
			defer func() { doneC <- nodename }()
			response, err := c.GetNodeDriverWithResponse(ctx, nodename)
			if err != nil {
				errC <- err
				return
			}
			switch {
			case response.JSON200 != nil:
				q <- response.JSON200.Items
			case response.JSON400 != nil:
				errC <- fmt.Errorf("%s: %s", nodename, *response.JSON400)
			case response.JSON401 != nil:
				errC <- fmt.Errorf("%s: %s", nodename, *response.JSON401)
			case response.JSON403 != nil:
				errC <- fmt.Errorf("%s: %s", nodename, *response.JSON403)
			case response.JSON500 != nil:
				errC <- fmt.Errorf("%s: %s", nodename, *response.JSON500)
			default:
				errC <- fmt.Errorf("%s: unexpected response: %s", nodename, response.Status())
			}
		}(nodename)
	}

	var (
		errs error
		done int
	)

	for {
		select {
		case err := <-errC:
			errs = errors.Join(errs, err)
		case items := <-q:
			l = append(l, items...)
		case <-doneC:
			done++
			if done == todo {
				goto out
			}
		case <-ctx.Done():
			errs = errors.Join(errs, ctx.Err())
			goto out
		}
	}

out:

	defaultOutput := "tab=NODE:meta.node,NAME:data.name"
	output.Renderer{
		DefaultOutput: defaultOutput,
		Output:        t.Output,
		Color:         t.Color,
		Data:          api.DriverList{Items: l, Kind: "DriverList"},
		Colorize:      rawconfig.Colorize,
	}.Print()

	return errs
}
