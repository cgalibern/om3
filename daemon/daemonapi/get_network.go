package daemonapi

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/opensvc/om3/core/network"
	"github.com/opensvc/om3/core/object"
	"github.com/opensvc/om3/daemon/api"
)

// GetNetworks returns network status list.
func (a *DaemonAPI) GetNetworks(ctx echo.Context, params api.GetNetworksParams) error {
	var items api.NetworkItems
	n, err := object.NewNode(object.WithVolatile(true))
	if err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "Failed to allocate a new object.Node", fmt.Sprint(err))
	}
	cips := GetClusterIPs()
	get := func() network.StatusList {
		if params.Name != nil {
			return network.ShowNetworksByName(n, *params.Name, cips)
		} else {
			return network.ShowNetworks(n, cips)
		}
	}
	for _, stat := range get() {
		item := api.Network{
			Name:    stat.Name,
			Type:    stat.Type,
			Network: stat.Network,
			Free:    stat.Free,
			Used:    stat.Used,
			Size:    stat.Size,
		}
		if len(stat.Errors) > 0 {
			l := append([]string{}, stat.Errors...)
			item.Errors = &l
		}
		items = append(items, item)
	}

	return ctx.JSON(http.StatusOK, api.NetworkList{Kind: "NetworkList", Items: items})
}
