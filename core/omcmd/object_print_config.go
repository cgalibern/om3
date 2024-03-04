package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/clientcontext"
	"github.com/opensvc/om3/core/naming"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/core/objectselector"
	"github.com/opensvc/om3/core/output"
	"github.com/opensvc/om3/core/rawconfig"
	"github.com/opensvc/om3/daemon/api"
)

type (
	CmdObjectPrintConfig struct {
		OptsGlobal
		Eval        bool
		Impersonate string
	}
)

type result map[string]rawconfig.T

func (t *CmdObjectPrintConfig) extract(selector string) (result, error) {
	data := make(result)
	c, err := client.New(client.WithURL(t.Server))
	if err != nil {
		return data, err
	}
	paths, err := objectselector.New(
		selector,
		objectselector.WithClient(c),
	).Expand()
	if err != nil {
		return data, err
	}
	for _, p := range paths {
		if d, err := t.extractOne(p, c); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
		} else {
			data[p.String()] = d
		}
	}
	return data, nil
}

func (t *CmdObjectPrintConfig) extractOne(p naming.Path, c *client.T) (rawconfig.T, error) {
	if data, err := t.extractFromDaemon(p, c); err == nil {
		return data, nil
	} else if clientcontext.IsSet() {
		return rawconfig.T{}, err
	} else if p.Exists() {
		return t.extractLocal(p)
	} else {
		return rawconfig.T{}, fmt.Errorf("%w, and no local instance to read from", err)
	}
}

func (t *CmdObjectPrintConfig) extractLocal(p naming.Path) (rawconfig.T, error) {
	obj, err := object.NewConfigurer(p)
	if err != nil {
		return rawconfig.T{}, err
	}
	if t.Eval {
		if t.Impersonate != "" {
			return obj.EvalConfigAs(t.Impersonate)
		}
		return obj.EvalConfig()
	}
	return obj.PrintConfig()
}

func (t *CmdObjectPrintConfig) extractFromDaemon(p naming.Path, c *client.T) (rawconfig.T, error) {
	var nodenames []string
	var errs error
	resp, err := c.GetObjectWithResponse(context.Background(), p.Namespace, p.Kind, p.Name)
	if err != nil {
		return rawconfig.T{}, err
	}
	switch {
	case resp.JSON200 != nil:
		if len(resp.JSON200.Data.Scope) == 0 {
			return rawconfig.T{}, nil
		} else {
			nodenames = resp.JSON200.Data.Scope
		}
	default:
		return rawconfig.T{}, fmt.Errorf("unexpected GetObject response: %s", resp.Status())
	}
	params := api.GetObjectConfigParams{
		Evaluate:    &t.Eval,
		Impersonate: &t.Impersonate,
	}
	for _, nodename := range nodenames {
		scopeClient, err := client.New(client.WithURL(nodename))
		if err != nil {
			return rawconfig.T{}, err
		}
		resp, err := scopeClient.GetObjectConfigWithResponse(context.Background(), p.Namespace, p.Kind, p.Name, &params)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		} else if resp.StatusCode() != http.StatusOK {
			errs = errors.Join(errs, fmt.Errorf("get object config: %s", resp.Status()))
			continue
		}
		data := rawconfig.T{}
		if b, err := json.Marshal(resp.JSON200.Data); err != nil {
			errs = errors.Join(errs, err)
			continue
		} else if err := json.Unmarshal(b, &data); err != nil {
			errs = errors.Join(errs, err)
			continue
		} else {
			return data, nil
		}
	}
	return rawconfig.T{}, errs
}

func (t *CmdObjectPrintConfig) Run(selector, kind string) error {
	var (
		data result
		err  error
	)
	mergedSelector := mergeSelector(selector, t.ObjectSelector, kind, "")
	if data, err = t.extract(mergedSelector); err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("no match")
	}
	var render func() string
	if _, err := naming.ParsePath(selector); err == nil {
		// single object selection
		render = func() string {
			d, _ := data[selector]
			return d.Render()
		}
		output.Renderer{
			Output:        t.Output,
			Color:         t.Color,
			Data:          data[selector].Data,
			HumanRenderer: render,
			Colorize:      rawconfig.Colorize,
		}.Print()
	} else {
		render = func() string {
			s := ""
			for p, d := range data {
				s += "#\n"
				s += "# path: " + p + "\n"
				s += "#\n"
				s += strings.Repeat("#", 78) + "\n"
				s += d.Render()
			}
			return s
		}
		output.Renderer{
			Output:        t.Output,
			Color:         t.Color,
			Data:          data,
			HumanRenderer: render,
			Colorize:      rawconfig.Colorize,
		}.Print()
	}
	return nil
}
