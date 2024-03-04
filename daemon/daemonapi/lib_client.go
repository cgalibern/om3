package daemonapi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/clusternode"
	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/util/funcopt"
)

func (a *DaemonAPI) proxy(ctx echo.Context, nodename string, fn func(*client.T) (*http.Response, error)) error {
	if data := node.StatusData.Get(nodename); data == nil {
		return JSONProblemf(ctx, http.StatusNotFound, "node status data not found", "%s", nodename)
	}
	c, err := newProxyClient(ctx, nodename)
	if err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "New client", "%s: %s", nodename, err)
	} else if !clusternode.Has(nodename) {
		return JSONProblemf(ctx, http.StatusBadRequest, "Invalid nodename", "field 'nodename' with value '%s' is not a cluster node", nodename)
	}
	if resp, err := fn(c); err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "Request peer", "%s: %s", nodename, err)
	} else {
		return ctx.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
	}
	return nil
}

func newProxyClient(ctx echo.Context, nodename string, opts ...funcopt.O) (*client.T, error) {
	options := []funcopt.O{
		client.WithURL(nodename),
		client.WithAuthorization(ctx.Request().Header.Get("authorization")),
	}
	options = append(options, opts...)
	return client.New(options...)
}
