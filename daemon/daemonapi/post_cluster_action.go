package daemonapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/hostname"
)

func (a *DaemonApi) PostClusterActionAbort(ctx echo.Context) error {
	return a.PostClusterAction(ctx, node.MonitorGlobalExpectAborted)
}

func (a *DaemonApi) PostClusterActionFreeze(ctx echo.Context) error {
	return a.PostClusterAction(ctx, node.MonitorGlobalExpectFrozen)
}

func (a *DaemonApi) PostClusterActionUnfreeze(ctx echo.Context) error {
	return a.PostClusterAction(ctx, node.MonitorGlobalExpectThawed)
}

func (a *DaemonApi) PostClusterAction(ctx echo.Context, globalExpect node.MonitorGlobalExpect) error {
	var (
		value = node.MonitorUpdate{}
	)
	if mon := node.MonitorData.Get(hostname.Hostname()); mon == nil {
		return JSONProblemf(ctx, http.StatusNotFound, "Not found", "node monitor not found: %s", hostname.Hostname())
	}
	value = node.MonitorUpdate{
		GlobalExpect:             &globalExpect,
		CandidateOrchestrationId: uuid.New(),
	}
	msg := msgbus.SetNodeMonitor{
		Node:  hostname.Hostname(),
		Value: value,
		Err:   make(chan error),
	}
	a.EventBus.Pub(&msg, labelNode, labelApi)
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	var errs error
	for {
		select {
		case <-ticker.C:
			return JSONProblemf(ctx, http.StatusRequestTimeout, "set monitor", "timeout waiting for monitor commit")
		case err := <-msg.Err:
			if err != nil {
				errs = errors.Join(errs, err)
			} else if errs != nil {
				return JSONProblemf(ctx, http.StatusConflict, "set monitor", "%s", errs)
			} else {
				return ctx.JSON(http.StatusOK, api.OrchestrationQueued{
					OrchestrationId: value.CandidateOrchestrationId,
				})
			}
		case <-ctx.Request().Context().Done():
			return JSONProblemf(ctx, http.StatusGone, "set monitor", "")
		}
	}
}