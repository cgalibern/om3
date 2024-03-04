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

func (a *DaemonAPI) PostPeerActionUnfreeze(ctx echo.Context, nodename string, params api.PostPeerActionUnfreezeParams) error {
	if nodename == a.localhost {
		return a.localNodeActionUnfreeze(ctx, params)
	}
	return a.proxy(ctx, nodename, func(c *client.T) (*http.Response, error) {
		return c.PostPeerActionUnfreeze(ctx.Request().Context(), nodename, &params)
	})
}

func (a *DaemonAPI) localNodeActionUnfreeze(ctx echo.Context, params api.PostPeerActionUnfreezeParams) error {
	if v, err := assertGrant(ctx, rbac.GrantRoot); !v {
		return err
	}
	log := LogHandler(ctx, "PostPeerActionUnfreeze")
	var requesterSid uuid.UUID
	args := []string{"node", "unfreeze", "--local"}
	if params.RequesterSid != nil {
		requesterSid = *params.RequesterSid
	}
	if sid, err := a.apiExec(ctx, naming.Path{}, requesterSid, args, log); err != nil {
		return JSONProblemf(ctx, http.StatusInternalServerError, "", "%s", err)
	} else {
		return ctx.JSON(http.StatusOK, api.NodeActionAccepted{SessionID: sid})
	}
}
