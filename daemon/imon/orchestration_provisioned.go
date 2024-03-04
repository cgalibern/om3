package imon

import (
	"sort"

	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/provisioned"
	"github.com/opensvc/om3/core/status"
)

func (t *Manager) orchestrateProvisioned() {
	switch t.state.State {
	case instance.MonitorStateIdle,
		instance.MonitorStateStopFailed,
		instance.MonitorStateUnprovisionFailed:
		t.provisionedFromIdle()
	case instance.MonitorStateWaitLeader:
		t.provisionedFromWaitLeader()
	case instance.MonitorStateProvisionFailed:
		t.provisionedFromProvisionFailed()
	}
}

func (t *Manager) provisionedFromProvisionFailed() {
	if t.provisionedClearIfReached() {
		return
	}
}

func (t *Manager) provisionedFromIdle() {
	if t.provisionedClearIfReached() {
		return
	}
	if t.isProvisioningLeader() {
		t.doAction(t.crmProvisionLeader, instance.MonitorStateProvisioning, instance.MonitorStateIdle, instance.MonitorStateProvisionFailed)
		return
	} else {
		t.transitionTo(instance.MonitorStateWaitLeader)
	}
}

func (t *Manager) provisionedFromWaitLeader() {
	if t.provisionedClearIfReached() {
		t.transitionTo(instance.MonitorStateIdle)
		return
	}
	if !t.hasLeaderProvisioned() {
		return
	}
	t.doAction(t.crmProvisionNonLeader, instance.MonitorStateProvisioning, instance.MonitorStateIdle, instance.MonitorStateProvisionFailed)
	return
}

func (t *Manager) provisionedClearIfReached() bool {
	reached := func(msg string) bool {
		t.log.Infof(msg)
		t.doneAndIdle()
		t.state.LocalExpect = instance.MonitorLocalExpectNone
		t.updateIfChange()
		return true
	}
	if t.instStatus[t.localhost].Provisioned.IsOneOf(provisioned.True, provisioned.NotApplicable) {
		return reached("provisioned orchestration: instance is provisioned -> set reached, clear local expect")
	}
	if t.instStatus[t.localhost].Avail == status.NotApplicable {
		return reached("provisioned orchestration: instance availability is n/a -> set reached, clear local expect")
	}
	return false
}

func (t *Manager) leaders() []string {
	l := make([]string, 0)
	for node, instMon := range t.instMonitor {
		if instMon.IsLeader {
			l = append(l, node)
		}
	}
	if t.state.IsLeader {
		l = append(l, t.localhost)
	}
	return l
}

// provisioningLeader returns one of all leaders.
// Select the first in alphalexical order.
func (t *Manager) provisioningLeader() string {
	leaders := t.leaders()
	switch len(leaders) {
	case 0:
		return ""
	case 1:
		return leaders[0]
	default:
		sort.StringSlice(leaders).Sort()
		return leaders[0]
	}
}

func (t *Manager) isProvisioningLeader() bool {
	if t.provisioningLeader() == t.localhost {
		return true
	}
	return false
}

func (t *Manager) hasLeaderProvisioned() bool {
	leader := t.provisioningLeader()
	if leaderInstanceStatus, ok := t.instStatus[leader]; !ok {
		return false
	} else if leaderInstanceStatus.Provisioned.IsOneOf(provisioned.True, provisioned.NotApplicable) {
		return true
	}
	return false
}
