package exec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

// gracePeriod is how long Kill waits before it sends SIGKILL.
var gracePeriod = 20 * time.Second

var (
	ErrNotStarted     = errors.New("process is not started")
	ErrAlreadyStarted = errors.New("process is already started")
)

type Cmd struct {
	cmd        *exec.Cmd
	waitStream chan error
	stopping   atomic.Bool
	started    atomic.Bool
}

func New(version string) *Cmd {
	e := newCmd("java", "-Dnogui=true", "-jar", fmt.Sprintf("JMusicBot-%s.jar", version))
	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr
	return e
}

func newCmd(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	setSysProcAttr(cmd)

	return &Cmd{
		cmd: cmd,
		// The buffer lets the watcher goroutine deliver its result without a
		// reader.
		waitStream: make(chan error, 1),
	}
}

// Start runs JMusicBot in the background and returns once the process is up.
func (e *Cmd) Start() error {
	if !e.started.CompareAndSwap(false, true) {
		return ErrAlreadyStarted
	}

	startStream := make(chan error, 1)
	go func() {
		// setPdeathsig sets a parent-death signal on the child on linux. The
		// kernel sends it when the thread that created the child exits, which
		// can happen before this process exits. The lock holds that thread for
		// the full life of the child. See https://go.dev/issue/27505.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(e.waitStream)

		if err := e.cmd.Start(); err != nil {
			startStream <- err
			return
		}
		startStream <- nil

		err := e.cmd.Wait()
		if e.stopping.Load() {
			e.waitStream <- nil
			return
		}
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
			// ExitCode returns -1 when a signal stopped the process, so the
			// message uses the full ProcessState.
			e.waitStream <- fmt.Errorf("JMusicBot exited with %s", exitErr.ProcessState)
			return
		}
		e.waitStream <- err
	}()

	return <-startStream
}

// Wait returns the channel that delivers the result of JMusicBot. The channel
// delivers nil on a clean exit or after Kill, then closes.
func (e *Cmd) Wait() <-chan error {
	return e.waitStream
}

// Kill sends sig to the process group of JMusicBot and waits for the exit. It
// sends SIGKILL to the same group after gracePeriod.
func (e *Cmd) Kill(sig syscall.Signal) error {
	if !e.started.Load() || e.cmd.Process == nil {
		return ErrNotStarted
	}
	e.stopping.Store(true)

	// Setpgid makes the child's PGID equal to its PID.
	pgid := e.cmd.Process.Pid
	if err := signalGroup(pgid, sig); err != nil {
		return err
	}

	timer := time.NewTimer(gracePeriod)
	defer timer.Stop()

	select {
	case err := <-e.waitStream:
		return err
	case <-timer.C:
		if err := signalGroup(pgid, syscall.SIGKILL); err != nil {
			return err
		}
		return <-e.waitStream
	}
}
