//go:build unix

package exec

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func waitFor(t *testing.T, e *Cmd, timeout time.Duration) error {
	t.Helper()
	select {
	case err := <-e.Wait():
		return err
	case <-time.After(timeout):
		t.Fatal("the process did not report a result in time")
		return nil
	}
}

func TestKillStopsRunningProcess(t *testing.T) {
	e := newCmd("sleep", "60")
	require.NoError(t, e.Start())

	start := time.Now()
	assert.NoError(t, e.Kill(syscall.SIGTERM))
	assert.Less(t, time.Since(start), gracePeriod)
}

func TestKillEscalatesToSigkill(t *testing.T) {
	old := gracePeriod
	gracePeriod = 200 * time.Millisecond
	t.Cleanup(func() { gracePeriod = old })

	// The shell ignores SIGTERM, so only SIGKILL stops it.
	e := newCmd("sh", "-c", "trap '' TERM; while :; do sleep 0.05; done")
	require.NoError(t, e.Start())
	time.Sleep(100 * time.Millisecond) // let the trap take effect

	done := make(chan error, 1)
	go func() { done <- e.Kill(syscall.SIGTERM) }()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Kill did not escalate to SIGKILL")
	}
}

func TestExitCodeIsReported(t *testing.T) {
	e := newCmd("sh", "-c", "exit 3")
	require.NoError(t, e.Start())

	err := waitFor(t, e, 5*time.Second)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exit status 3")
}

func TestCleanExitReportsNil(t *testing.T) {
	e := newCmd("sh", "-c", "exit 0")
	require.NoError(t, e.Start())
	assert.NoError(t, waitFor(t, e, 5*time.Second))
}

func TestKillAfterNaturalExit(t *testing.T) {
	e := newCmd("sh", "-c", "exit 0")
	require.NoError(t, e.Start())
	require.NoError(t, waitFor(t, e, 5*time.Second))

	done := make(chan error, 1)
	go func() { done <- e.Kill(syscall.SIGTERM) }()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Kill blocked on a process that was already gone")
	}
}

func TestKillBeforeStart(t *testing.T) {
	e := newCmd("sleep", "60")
	assert.ErrorIs(t, e.Kill(syscall.SIGTERM), ErrNotStarted)
}

func TestDoubleStart(t *testing.T) {
	e := newCmd("sleep", "60")
	require.NoError(t, e.Start())
	t.Cleanup(func() { _ = e.Kill(syscall.SIGKILL) })

	assert.ErrorIs(t, e.Start(), ErrAlreadyStarted)
}

func TestStartFailure(t *testing.T) {
	e := newCmd("this-command-does-not-exist-djdiscord")
	assert.Error(t, e.Start())
}

func TestKillStopsWholeProcessGroup(t *testing.T) {
	e := newCmd("sh", "-c", "sleep 60 & wait")
	require.NoError(t, e.Start())
	pgid := e.cmd.Process.Pid
	time.Sleep(100 * time.Millisecond)

	require.NoError(t, e.Kill(syscall.SIGTERM))

	// Signal 0 checks that the group exists without sending anything.
	assert.ErrorIs(t, syscall.Kill(-pgid, 0), syscall.ESRCH)
}
