//go:build unix

package exec

import (
	"errors"
	"os/exec"
	"syscall"
)

func setSysProcAttr(cmd *exec.Cmd) {
	attr := &syscall.SysProcAttr{
		// A new process group takes the child out of the terminal's foreground
		// group, so Ctrl+C reaches this program alone and Kill forwards it.
		Setpgid: true,
	}
	setPdeathsig(attr)
	cmd.SysProcAttr = attr
}

// signalGroup sends sig to the whole process group of pid. It ignores a group
// that no longer exists.
func signalGroup(pid int, sig syscall.Signal) error {
	if err := syscall.Kill(-pid, sig); err != nil && !errors.Is(err, syscall.ESRCH) {
		return err
	}
	return nil
}
