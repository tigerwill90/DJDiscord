//go:build linux

package exec

import "syscall"

// setPdeathsig makes the kernel kill the child with this program. It covers a
// SIGKILL, a panic or a segfault, where no signal handler of this program runs.
// The kernel watches the creating thread, which Start keeps locked and alive.
func setPdeathsig(attr *syscall.SysProcAttr) {
	attr.Pdeathsig = syscall.SIGKILL
}
