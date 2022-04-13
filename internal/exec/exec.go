package exec

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type Cmd struct {
	cmd        *exec.Cmd
	waitStream chan error
}

func New(version string) *Cmd {
	cmd := exec.Command("java", "-Dnogui=true", "-jar", fmt.Sprintf("JMusicBot-%s.jar", version))
	setNewProcessGroup(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Cmd{
		cmd: cmd,
	}
}

func setNewProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

func (e *Cmd) Start() error {
	if err := e.cmd.Start(); err != nil {
		return err
	}
	e.waitStream = make(chan error)
	go func() {
		defer close(e.waitStream)
		err := e.cmd.Wait()
		if errExit, ok := err.(*exec.ExitError); ok {
			e.waitStream <- fmt.Errorf("JMusicBot exited with code error %d", errExit.ExitCode())
			return
		}
		e.waitStream <- err
	}()
	return nil
}

func (e *Cmd) Wait() <-chan error {
	return e.waitStream
}

func (e *Cmd) Kill() error {
	if e.cmd.Process != nil {
		if err := e.cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}
	}
	return <-e.waitStream
}
