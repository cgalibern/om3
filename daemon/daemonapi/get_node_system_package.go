package daemonapi

import (
	"errors"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/opensvc/om3/core/clusternode"
	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/daemon/api"
)

func (a *DaemonAPI) GetNodeSystemPackage(ctx echo.Context, nodename api.InPathNodeName) error {
	if a.localhost == nodename {
		return a.getLocalNodeSystemPackage(ctx)
	} else if !clusternode.Has(nodename) {
		return JSONProblemf(ctx, http.StatusBadRequest, "Invalid parameters", "%s is not a cluster node", nodename)
	} else {
		return a.getPeerNodeSystemPackage(ctx, nodename)
	}
}

func (a *DaemonAPI) getPeerNodeSystemPackage(ctx echo.Context, nodename string) error {
	if data := node.StatusData.Get(nodename); data == nil {
		return JSONProblemf(ctx, http.StatusNotFound, "node status data not found", "%s", nodename)
	}
	c, err := newProxyClient(ctx, nodename)
	if err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "New client", "%s: %s", nodename, err)
	} else if !clusternode.Has(nodename) {
		return JSONProblemf(ctx, http.StatusBadRequest, "Invalid nodename", "field 'nodename' with value '%s' is not a cluster node", nodename)
	}
	if resp, err := c.GetNodeSystemPackageWithResponse(ctx.Request().Context(), nodename); err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "Request peer", "%s: %s", nodename, err)
	} else if len(resp.Body) > 0 {
		return ctx.JSONBlob(resp.StatusCode(), resp.Body)
	}
	return nil
}

func (a *DaemonAPI) getLocalNodeSystemPackage(ctx echo.Context) error {
	n, err := object.NewNode()
	if err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "New node", "%s", err)
	}
	data, err := n.LoadPkg()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return JSONProblemf(ctx, http.StatusNotFound, "Load package cache", "waiting for cached value: %s", err)
		} else {
			return JSONProblemf(ctx, http.StatusInternalServerError, "Load package cache", "%s", err)
		}
	}
	items := make(api.PackageItems, len(data))
	for i := 0; i < len(data); i++ {
		items[i] = api.PackageItem{
			Kind: "PackageItem",
			Data: api.Package{
				Version:     data[i].Version,
				Name:        data[i].Name,
				Arch:        data[i].Arch,
				Type:        data[i].Type,
				InstalledAt: data[i].InstalledAt,
				Sig:         data[i].Sig,
			},
			Meta: api.NodeMeta{
				Node: a.localhost,
			},
		}
	}

	return ctx.JSON(http.StatusOK, api.PackageList{Kind: "PackageList", Items: items})
}