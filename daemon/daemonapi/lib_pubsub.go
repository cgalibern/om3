package daemonapi

import (
	"time"

	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/daemon/daemondata"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/plog"
)

func (a *DaemonApi) announceSub(name string) {
	a.EventBus.Pub(&msgbus.ClientSubscribed{Time: time.Now(), Name: name})
}

func (a *DaemonApi) announceUnsub(name string) {
	a.EventBus.Pub(&msgbus.ClientUnsubscribed{Time: time.Now(), Name: name})
}

func (a *DaemonApi) announceNodeState(log *plog.Logger, state node.MonitorState) {
	log.Infof("announce node state %s", state)
	a.EventBus.Pub(&msgbus.SetNodeMonitor{Node: a.localhost, Value: node.MonitorUpdate{State: &state}}, labelApi)
	time.Sleep(2 * daemondata.PropagationInterval())
}