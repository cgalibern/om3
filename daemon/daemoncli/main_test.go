package daemoncli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func init() {
	WaitRunningDelay = 10 * time.Millisecond
}

func TestDaemonStartThenStop(t *testing.T) {
	require.False(t, Running())
	go func() {
		require.Nil(t, Start())
	}()
	require.Nil(t, WaitRunning())
	require.True(t, Running())
	require.Nil(t, Stop())
	require.False(t, Running())
}

func TestDaemonReStartThenStop(t *testing.T) {
	require.False(t, Running())
	go func() {
		require.Nil(t, ReStart())
	}()
	require.Nil(t, WaitRunning())
	require.True(t, Running())
	require.Nil(t, Stop())
	require.False(t, Running())
}

func TestStop(t *testing.T) {
	require.False(t, Running())
	require.Nil(t, Stop())
	require.False(t, Running())
}
