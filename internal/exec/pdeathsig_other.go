//go:build unix && !linux

package exec

import "syscall"

func setPdeathsig(attr *syscall.SysProcAttr) {}
