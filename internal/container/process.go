package container

import (
	"os/exec"

	"golang.org/x/sys/unix"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags:
	}
}
