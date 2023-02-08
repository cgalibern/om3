package commands

import (
	"github.com/pkg/errors"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/clientcontext"
	"github.com/opensvc/om3/core/network"
	"github.com/opensvc/om3/core/object"
)

type (
	CmdNetworkSetup struct {
		OptsGlobal
	}
)

func (t *CmdNetworkSetup) Run() error {
	if t.Local || !clientcontext.IsSet() {
		return t.doLocal()
	} else {
		return t.doDaemon()
	}
}

func (t *CmdNetworkSetup) doLocal() error {
	n, err := object.NewNode()
	if err != nil {
		return err
	}
	return network.Setup(n)
}

func (t *CmdNetworkSetup) doDaemon() error {
	var (
		c   *client.T
		err error
	)
	if c, err = client.New(client.WithURL(t.Server)); err != nil {
		return err
	}
	return errors.Errorf("TODO %v", c)
}
