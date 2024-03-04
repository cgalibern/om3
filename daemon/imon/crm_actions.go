package imon

import (
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/opensvc/om3/core/env"
	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/daemon/msgbus"
	"github.com/opensvc/om3/util/command"
	"github.com/opensvc/om3/util/pubsub"
)

var (
	cmdPath string

	maxRunners = 25

	// runners chan limit number of // crmActions to maxRunners
	runners = make(chan struct{}, maxRunners)

	// crmAction can be used to define alternate crmAction for tests
	crmAction func(title string, cmdArgs ...string) error
)

func init() {
	var err error
	cmdPath, err = os.Executable()
	if err != nil {
		cmdPath = "/bin/false"
	}
}

// SetCmdPathForTest set the opensvc command path for tests
func SetCmdPathForTest(s string) {
	// TODO use another method to create dedicated side effects for tests
	cmdPath = s
}

func (t *Manager) orchestrateAfterAction(state, newState instance.MonitorState) {
	select {
	case <-t.ctx.Done():
		return
	default:
	}
	t.cmdC <- cmdOrchestrate{state: state, newState: newState}
}

func (t *Manager) crmBoot() error {
	return t.crmAction("boot", t.path.String(), "boot", "--local")
}

func (t *Manager) crmDelete() error {
	return t.crmAction("delete", t.path.String(), "delete", "--local")
}

func (t *Manager) crmFreeze() error {
	return t.crmAction("freeze", t.path.String(), "freeze", "--local")
}

func (t *Manager) crmProvisionNonLeader() error {
	return t.crmAction("provision non leader", t.path.String(), "provision", "--local")
}

func (t *Manager) crmProvisionLeader() error {
	return t.crmAction("provision leader", t.path.String(), "provision", "--local", "--leader", "--disable-rollback")
}

func (t *Manager) crmResourceStartStandby(rids []string) error {
	s := strings.Join(rids, ",")
	return t.crmAction("start", t.path.String(), "startstandby", "--local", "--rid", s)
}

func (t *Manager) crmResourceStart(rids []string) error {
	s := strings.Join(rids, ",")
	return t.crmAction("start", t.path.String(), "start", "--local", "--rid", s)
}

func (t *Manager) crmShutdown() error {
	return t.crmAction("shutdown", t.path.String(), "shutdown")
}

func (t *Manager) crmStart() error {
	return t.crmAction("start", t.path.String(), "start", "--local")
}

func (t *Manager) crmStatus() error {
	return t.crmAction("status", t.path.String(), "status", "-r")
}

func (t *Manager) crmStop() error {
	return t.crmAction("stop", t.path.String(), "stop", "--local")
}

func (t *Manager) crmUnfreeze() error {
	return t.crmAction("unfreeze", t.path.String(), "unfreeze", "--local")
}

func (t *Manager) crmUnprovisionNonLeader() error {
	return t.crmAction("unprovision non leader", t.path.String(), "unprovision", "--local")
}

func (t *Manager) crmUnprovisionLeader() error {
	return t.crmAction("unprovision leader", t.path.String(), "unprovision", "--local", "--leader")
}

func (t *Manager) crmAction(title string, cmdArgs ...string) error {
	if crmAction != nil {
		return crmAction(title, cmdArgs...)
	}
	return t.crmDefaultAction(title, cmdArgs...)
}

func (t *Manager) crmDefaultAction(title string, cmdArgs ...string) error {
	runners <- struct{}{}
	defer func() {
		<-runners
	}()
	sid := uuid.New()
	cmd := command.New(
		command.WithName(cmdPath),
		command.WithArgs(cmdArgs),
		command.WithLogger(t.log),
		command.WithVarEnv(
			env.OriginSetenvArg(env.ActionOriginDaemonMonitor),
			env.ActionOrchestrationIDVar+"="+t.state.OrchestrationID.String(),
			"OSVC_SESSION_ID="+sid.String(),
		),
	)
	labels := []pubsub.Label{t.labelLocalhost, t.labelPath, {"origin", "imon"}, {"sid", sid.String()}}
	if title != "" {
		t.loggerWithState().Infof("-> exec %s", append([]string{cmdPath}, cmdArgs...))
	} else {
		t.loggerWithState().Debugf("-> exec %s", append([]string{cmdPath}, cmdArgs...))
	}
	t.pubsubBus.Pub(&msgbus.Exec{Command: cmd.String(), Node: t.localhost, Origin: "imon", Title: title}, labels...)
	startTime := time.Now()
	if err := cmd.Run(); err != nil {
		duration := time.Now().Sub(startTime)
		t.pubsubBus.Pub(&msgbus.ExecFailed{Command: cmd.String(), Duration: duration, ErrS: err.Error(), Node: t.localhost, Origin: "imon", Title: title}, labels...)
		t.loggerWithState().Errorf("<- exec %s: %s", append([]string{cmdPath}, cmdArgs...), err)
		return err
	}
	duration := time.Now().Sub(startTime)
	t.pubsubBus.Pub(&msgbus.ExecSuccess{Command: cmd.String(), Duration: duration, Node: t.localhost, Origin: "imon", Title: title}, labels...)
	if title != "" {
		t.loggerWithState().Infof("<- exec %s", append([]string{cmdPath}, cmdArgs...))
	} else {
		t.loggerWithState().Debugf("<- exec %s", append([]string{cmdPath}, cmdArgs...))
	}
	return nil
}
