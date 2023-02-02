package daemon_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"opensvc.com/opensvc/cmd"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/daemon/daemon"
	"opensvc.com/opensvc/daemon/routinehelper"
	"opensvc.com/opensvc/testhelper"
)

func TestMain(m *testing.M) {
	testhelper.Main(m, cmd.ExecuteArgs)
}

func setup(t *testing.T) testhelper.Env {
	env := testhelper.Setup(t)
	env.InstallFile("../../testdata/cluster.conf", "etc/cluster.conf")
	env.InstallFile("../../testdata/ca-cluster1.conf", "etc/namespaces/system/sec/ca-cluster1.conf")
	env.InstallFile("../../testdata/cert-cluster1.conf", "etc/namespaces/system/sec/cert-cluster1.conf")
	rawconfig.LoadSections()
	return env
}

func RunCmd(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("%s error %s\n%s", cmd, err, b)
	} else {
		t.Logf("%s\n%s", cmd, b)
	}
}

func Trace(t *testing.T) {
	RunCmd(t, "ps", "fax")
	RunCmd(t, "netstat", "-petulan")
	pid := os.Getpid()
	RunCmd(t, "ls", "-l", fmt.Sprintf("/proc/%d/fd", pid))
}

func DaemonPorts(t *testing.T, name string) error {
	t.Logf("Verify daemon ports [%s]", name)
	Trace(t)
	var delay time.Duration
	for _, port := range []string{"1214", "1215"} {
		if err := testhelper.TcpPortAvailable(port); err != nil {
			t.Logf("Verify daemon ports [%s] failed for port %s '%s' wait delay then check again", name, port, err)
			Trace(t)
			delay = 5 * time.Second
		}
	}
	time.Sleep(delay)
	for _, port := range []string{"1214", "1215"} {
		if err := testhelper.TcpPortAvailable(port); err != nil {
			t.Logf("Verify daemon ports [%s] failed for port %s '%s'", name, port, err)
			Trace(t)
			return err
		}
	}
	t.Logf("Verify daemon ports [%s] [done]", name)
	return nil
}

func TestDaemon(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipped for non root user")
	}
	require.NoError(t, DaemonPorts(t, fmt.Sprintf("-> %s", t.Name())))
	if t.Failed() {
		t.Fatal("-> TestDaemon DaemonPorts")
	}
	var main *daemon.T
	setup(t)

	t.Log("New")
	main = daemon.New(
		daemon.WithRoutineTracer(routinehelper.NewTracer()),
	)
	require.NotNil(t, main)
	require.False(t, main.Enabled(), "The daemon should not be Enabled after New")
	require.False(t, main.Running(), "The daemon should not be Running after New")
	require.Equalf(t, 0, main.TraceRDump().Count, "found %#v", main.TraceRDump())

	t.Log("Start")
	require.NoError(t, main.Start(context.Background()))
	require.True(t, main.Enabled(), "The daemon should be Enabled after Start")
	require.True(t, main.Running(), "The daemon should be Running after Start")

	t.Log("Restart")
	require.NoError(t, main.Restart(context.Background()))
	require.True(t, main.Enabled(), "The daemon should be Enabled after Restart")
	require.True(t, main.Running(), "The daemon should be Running after Restart")

	t.Log("Stop")
	require.NoError(t, main.Stop())
	require.False(t, main.Enabled(), "The daemon should not be Enabled after Stop")
	require.False(t, main.Running(), "The daemon should not be Running after Stop")
	require.Equalf(t, 0, main.TraceRDump().Count, "Daemon routines should be stopped, found %#v", main.TraceRDump())

	t.Log("Stop")
	require.NoError(t, main.Stop())
	require.False(t, main.Enabled(), "The daemon should not be Enabled after Stop")
	require.False(t, main.Running(), "The daemon should not be Running after Stop")

	t.Log("Restart")
	require.NoError(t, main.Restart(context.Background()))
	require.True(t, main.Enabled(), "The daemon should be Enabled after Restart")
	require.True(t, main.Running(), "The daemon should be Running after Restart")

	t.Log("Restart")
	require.NoError(t, main.Restart(context.Background()))
	require.True(t, main.Enabled(), "The daemon should be Enabled after Restart")
	require.True(t, main.Running(), "The daemon should be Running after Restart")

	t.Log("Stop")
	require.NoError(t, main.Stop())
	require.False(t, main.Enabled(), "The daemon should not be Enabled after Stop")
	require.False(t, main.Running(), "The daemon should not be Running after Stop")

	main.Wait()
	main.Wait() // verify we don't block on calling WaitDone() multiple times
	require.Equalf(t, 0, main.TraceRDump().Count, "Daemon routines should be stopped, found %#v", main.TraceRDump())

	t.Log("RunDaemon")
	main, err := daemon.RunDaemon()
	require.NotNil(t, main)
	require.NoError(t, err)
	require.True(t, main.Enabled(), "The daemon should be Enabled after RunDaemon")
	require.True(t, main.Running(), "The daemon should be Running after RunDaemon")

	t.Log("Stop")
	require.NoError(t, main.Stop())
	require.False(t, main.Enabled(), "The daemon should not be Enabled after Stop")
	require.False(t, main.Running(), "The daemon should not be Running after Stop")
	require.Equalf(t, 0, main.TraceRDump().Count, "Daemon routines should be stopped, found %#v", main.TraceRDump())

	require.NoError(t, DaemonPorts(t, fmt.Sprintf("<- %s", t.Name())))
}
