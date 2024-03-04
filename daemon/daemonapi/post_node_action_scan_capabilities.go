package daemonapi

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/naming"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/daemon/rbac"
)

func (a *DaemonAPI) PostNodeActionScanCapabilities(ctx echo.Context, nodename string, params api.PostNodeActionScanCapabilitiesParams) error {
	if nodename == a.localhost {
		return a.localNodeActionScanCapabilities(ctx, params)
	}
	return a.proxy(ctx, nodename, func(c *client.T) (*http.Response, error) {
		return c.PostNodeActionScanCapabilities(ctx.Request().Context(), nodename, &params)
	})
}

func (a *DaemonAPI) localNodeActionScanCapabilities(ctx echo.Context, params api.PostNodeActionScanCapabilitiesParams) error {
	if v, err := assertGrant(ctx, rbac.GrantRoot); !v {
		return err
	}
	log := LogHandler(ctx, "PostNodeActionScanCapabilities")
	var requesterSid uuid.UUID
	args := []string{"node", "scan", "capabilities", "--local"}
	if params.RequesterSid != nil {
		requesterSid = *params.RequesterSid
	}
	if sid, err := a.apiExec(ctx, naming.Path{}, requesterSid, args, log); err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "", "%s", err)
	} else {
		return ctx.JSON(http.StatusOK, api.NodeActionAccepted{SessionID: sid})
	}
}
