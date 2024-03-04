package imon

import (
	"time"

	"github.com/rs/zerolog"

	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/naming"
	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/core/provisioned"
	"github.com/opensvc/om3/core/status"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/command"
	"github.com/opensvc/om3/util/toc"
)

type (
	todoMap map[string]bool
)

func (t todoMap) Add(rid string) {
	t[rid] = true
}

func (t todoMap) Del(rid string) {
	delete(t, rid)
}

func (t todoMap) IsEmpty() bool {
	return len(t) == 0
}

func newTodoMap() todoMap {
	m := make(todoMap)
	return m
}

func (t *Manager) orchestrateResourceRestart() {
	todoRestart := newTodoMap()
	todoStandby := newTodoMap()

	pubMonitorAction := func(rid string) {
		t.pubsubBus.Pub(
			&msgbus.InstanceMonitorAction{
				Path:   t.path,
				Node:   t.localhost,
				Action: t.instConfig.MonitorAction,
				RID:    rid,
			},
			t.labelPath,
			t.labelLocalhost)
	}

	// doPreMonitorAction executes a user-defined command before imon
	// runs the MonitorAction. This command can detect a situation where
	// the MonitorAction can not succeed, and decide to do another action.
	doPreMonitorAction := func() error {
		if t.instConfig.PreMonitorAction == "" {
			return nil
		}
		t.log.Infof("execute pre monitor action: %s", t.instConfig.PreMonitorAction)
		cmdArgs, err := command.CmdArgsFromString(t.instConfig.PreMonitorAction)
		if err != nil {
			return err
		}
		if len(cmdArgs) == 0 {
			return nil
		}
		cmd := command.New(
			command.WithName(cmdArgs[0]),
			command.WithVarArgs(cmdArgs[1:]...),
			command.WithLogger(t.log),
			command.WithStdoutLogLevel(zerolog.InfoLevel),
			command.WithStderrLogLevel(zerolog.ErrorLevel),
			command.WithTimeout(60*time.Second),
		)
		return cmd.Run()
	}

	doMonitorAction := func(rid string) {
		switch t.instConfig.MonitorAction {
		case instance.MonitorActionCrash:
		case instance.MonitorActionFreezeStop:
		case instance.MonitorActionReboot:
		case instance.MonitorActionSwitch:
		case instance.MonitorActionNone:
			t.log.Errorf("skip monitor action: not configured")
			return
		default:
			t.log.Errorf("skip monitor action: not supported: %s", t.instConfig.MonitorAction)
			return
		}

		if err := doPreMonitorAction(); err != nil {
			t.log.Errorf("pre monitor action: %s", err)
		}

		t.log.Infof("do %s monitor action", t.instConfig.MonitorAction)
		pubMonitorAction(rid)

		switch t.instConfig.MonitorAction {
		case instance.MonitorActionCrash:
			if err := toc.Crash(); err != nil {
				t.log.Errorf("monitor action: %s", err)
			}
		case instance.MonitorActionFreezeStop:
			t.doFreezeStop()
			t.doStop()
		case instance.MonitorActionReboot:
			if err := toc.Reboot(); err != nil {
				t.log.Errorf("monitor action: %s", err)
			}
		case instance.MonitorActionSwitch:
			t.createPendingWithDuration(stopDuration)
			t.doAction(t.crmStop, instance.MonitorStateStopping, instance.MonitorStateStartFailed, instance.MonitorStateStopFailed)
		}
	}

	resetTimer := func(rid string, rmon *instance.ResourceMonitor) {
		todoRestart.Del(rid)
		todoStandby.Del(rid)
		if rmon.Restart.Timer != nil {
			t.log.Infof("resource %s is up, reset delayed restart", rid)
			t.change = rmon.StopRestartTimer()
			t.state.Resources.Set(rid, *rmon)
		}
	}

	resetRemaining := func(rid string, rcfg *instance.ResourceConfig, rmon *instance.ResourceMonitor) {
		if rmon.Restart.Remaining != rcfg.Restart {
			t.log.Infof("resource %s is up, reset restart count to the max (%d -> %d)", rid, rmon.Restart.Remaining, rcfg.Restart)
			t.state.MonitorActionExecutedAt = time.Time{}
			rmon.Restart.Remaining = rcfg.Restart
			t.state.Resources.Set(rid, *rmon)
			t.change = true
		}
	}

	resetRemainingAndTimer := func(rid string, rcfg *instance.ResourceConfig, rmon *instance.ResourceMonitor) {
		resetRemaining(rid, rcfg, rmon)
		resetTimer(rid, rmon)
	}

	resetTimers := func() {
		for rid, rmon := range t.state.Resources {
			resetTimer(rid, &rmon)
		}
	}

	planFor := func(rid string, resStatus status.T, started bool) {
		rcfg := t.instConfig.Resources.Get(rid)
		rmon := t.state.Resources.Get(rid)
		switch {
		case rcfg == nil:
			return
		case rmon == nil:
			return
		case rcfg.IsDisabled:
			t.log.Debugf("resource %s restart skip: disable=%v", rid, rcfg.IsDisabled)
			resetRemainingAndTimer(rid, rcfg, rmon)
		case resStatus.Is(status.NotApplicable, status.Undef):
			t.log.Debugf("resource %s restart skip: status=%s", rid, resStatus)
			resetRemainingAndTimer(rid, rcfg, rmon)
		case resStatus.Is(status.Up, status.StandbyUp):
			t.log.Debugf("resource %s restart skip: status=%s", rid, resStatus)
			resetRemainingAndTimer(rid, rcfg, rmon)
		case rmon.Restart.Timer != nil:
			t.log.Debugf("resource %s restart skip: already has a delay timer", rid)
		case !t.state.MonitorActionExecutedAt.IsZero():
			t.log.Debugf("resource %s restart skip: already ran the monitor action", rid)
		case started:
			t.log.Infof("resource %s status %s, restart remaining %d out of %d", rid, resStatus, rmon.Restart.Remaining, rcfg.Restart)
			if rmon.Restart.Remaining == 0 {
				t.state.MonitorActionExecutedAt = time.Now()
				t.change = true
				doMonitorAction(rid)
			} else {
				todoRestart.Add(rid)
			}
		case rcfg.IsStandby:
			t.log.Infof("resource %s status %s, standby restart remaining %d out of %d", rid, resStatus, rmon.Restart.Remaining, rcfg.Restart)
			todoStandby.Add(rid)
		default:
			t.log.Debugf("resource %s restart skip: instance not started", rid)
			resetTimer(rid, rmon)
		}
	}

	getRidsAndDelay := func(todo todoMap) ([]string, time.Duration) {
		var maxDelay time.Duration
		rids := make([]string, 0)
		now := time.Now()
		for rid := range todo {
			rcfg := t.instConfig.Resources.Get(rid)
			if rcfg == nil {
				continue
			}
			rmon := t.state.Resources.Get(rid)
			if rmon == nil {
				continue
			}
			if rcfg.RestartDelay != nil {
				notBefore := rmon.Restart.LastAt.Add(*rcfg.RestartDelay)
				if now.Before(notBefore) {
					delay := notBefore.Sub(now)
					if delay > maxDelay {
						maxDelay = delay
					}
				}
			}
			rids = append(rids, rid)
		}
		return rids, maxDelay
	}

	doRestart := func() {
		rids, delay := getRidsAndDelay(todoRestart)
		if len(rids) == 0 {
			return
		}
		timer := time.AfterFunc(delay, func() {
			now := time.Now()
			for _, rid := range rids {
				rmon := t.state.Resources.Get(rid)
				if rmon == nil {
					continue
				}
				rmon.Restart.LastAt = now
				rmon.Restart.Timer = nil
				t.state.Resources.Set(rid, *rmon)
				t.change = true
			}
			action := func() error {
				return t.crmResourceStart(rids)
			}
			t.doTransitionAction(action, instance.MonitorStateStarting, instance.MonitorStateIdle, instance.MonitorStateStartFailed)
		})
		for _, rid := range rids {
			rmon := t.state.Resources.Get(rid)
			if rmon == nil {
				continue
			}
			rmon.DecRestartRemaining()
			rmon.Restart.Timer = timer
			t.state.Resources.Set(rid, *rmon)
			t.change = true
		}
	}

	doStandby := func() {
		rids, delay := getRidsAndDelay(todoStandby)
		if len(rids) == 0 {
			return
		}
		timer := time.AfterFunc(delay, func() {
			now := time.Now()
			for _, rid := range rids {
				rmon := t.state.Resources.Get(rid)
				if rmon == nil {
					continue
				}
				rmon.Restart.LastAt = now
				rmon.Restart.Timer = nil
				t.state.Resources.Set(rid, *rmon)
				t.change = true
			}
			action := func() error {
				return t.crmResourceStartStandby(rids)
			}
			t.doTransitionAction(action, instance.MonitorStateStarting, instance.MonitorStateIdle, instance.MonitorStateStartFailed)
		})
		for _, rid := range rids {
			rmon := t.state.Resources.Get(rid)
			if rmon == nil {
				continue
			}
			rmon.DecRestartRemaining()
			rmon.Restart.Timer = timer
			t.state.Resources.Set(rid, *rmon)
			t.change = true
		}
	}

	// discard the cluster object
	if t.path.String() == "cluster" {
		return
	}

	// discard all execpt svc and vol
	switch t.path.Kind {
	case naming.KindSvc, naming.KindVol:
	default:
		return
	}

	// discard if the instance status does not exist
	if _, ok := t.instStatus[t.localhost]; !ok {
		resetTimers()
		return
	}

	// don't run on frozen nodes
	if t.nodeStatus[t.localhost].IsFrozen() {
		resetTimers()
		return
	}

	// don't run when the node is not idle
	if t.nodeMonitor[t.localhost].State != node.MonitorStateIdle {
		resetTimers()
		return
	}

	// don't run on frozen instances
	if t.instStatus[t.localhost].IsFrozen() {
		resetTimers()
		return
	}

	// discard not provisioned
	if instanceStatus := t.instStatus[t.localhost]; instanceStatus.Provisioned.IsOneOf(provisioned.False, provisioned.Mixed, provisioned.Undef) {
		t.log.Debugf("skip restart: provisioned=%s", instanceStatus.Provisioned)
		resetTimers()
		return
	}

	// discard if the instance has no monitor data
	instMonitor, ok := t.GetInstanceMonitor(t.localhost)
	if !ok {
		t.log.Debugf("skip restart: no instance monitor")
		resetTimers()
		return
	}

	// discard if the instance is not idle nor start failed
	switch instMonitor.State {
	case instance.MonitorStateIdle, instance.MonitorStateStartFailed:
		// pass
	default:
		t.log.Debugf("skip restart: state=%s", instMonitor.State)
		return
	}

	started := instMonitor.LocalExpect == instance.MonitorLocalExpectStarted

	for rid, rstat := range t.instStatus[t.localhost].Resources {
		planFor(rid, rstat.Status, started)
	}
	doStandby()
	doRestart()
}
