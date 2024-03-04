package imon

import (
	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/status"
	"github.com/opensvc/om3/core/topology"
)

func (t *Manager) orchestrateNone() {
	t.clearStartFailed()
	t.clearBootFailed()
	if t.objStatus.Orchestrate == "ha" {
		t.orchestrateHAStart()
		t.orchestrateHAStop()
	}
}

func (t *Manager) orchestrateHAStop() {
	if t.objStatus.Topology != topology.Flex {
		return
	}
	if v, _ := t.isExtraInstance(); !v {
		return
	}
	t.stop()
}

func (t *Manager) orchestrateHAStart() {
	// we are here because we are ha object with global expect None
	switch t.state.State {
	case instance.MonitorStateReady:
		t.cancelReadyState()
	case instance.MonitorStateStarted:
		// started means the action start has been done. This state is a
		// waiter step to verify if received started like local instance status
		// to transition state: started -> idle
		// It prevents unexpected transition state -> ready
		if t.isLocalStarted() {
			t.log.Infof("instance is now started, enable resource restart")
			t.state.LocalExpect = instance.MonitorLocalExpectStarted
			t.transitionTo(instance.MonitorStateIdle)
		}
		return
	}
	if v, reason := t.isStartable(); !v {
		if t.pendingCancel != nil && t.state.State == instance.MonitorStateReady {
			t.log.Infof("instance is not startable, clear the ready state: %s", reason)
			t.clearPending()
			t.transitionTo(instance.MonitorStateIdle)
		}
		return
	}
	if t.isLocalStarted() {
		return
	}
	t.orchestrateStarted()
}

// clearBootFailed clears the boot failed state when the following conditions are met:
//
// + local avail is Down, StandbyDown, NotApplicable
// + global expect is none
func (t *Manager) clearBootFailed() {
	if t.state.State != instance.MonitorStateBootFailed {
		return
	}
	switch t.instStatus[t.localhost].Avail {
	case status.Down:
	case status.StandbyDown:
	case status.NotApplicable:
	default:
		return
	}
	for _, instanceMonitor := range t.instMonitor {
		switch instanceMonitor.GlobalExpect {
		case instance.MonitorGlobalExpectNone:
		default:
			return
		}
	}
	t.log.Infof("clear instance %s: local instance avail is %s, object avail is %s",
		t.state.State, t.instStatus[t.localhost].Avail, t.objStatus.Avail)
	t.transitionTo(instance.MonitorStateIdle)
}

func (t *Manager) clearStartFailed() {
	if t.state.State != instance.MonitorStateStartFailed {
		return
	}
	if t.objStatus.Avail != status.Up {
		return
	}
	for _, instanceMonitor := range t.instMonitor {
		switch instanceMonitor.GlobalExpect {
		case instance.MonitorGlobalExpectNone:
		default:
			return
		}
	}
	t.log.Infof("clear instance start failed: the object is up")
	t.transitionTo(instance.MonitorStateIdle)
}
