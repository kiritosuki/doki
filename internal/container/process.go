package container

import (
	"os/exec"
	"syscall"
)

// NewContainerProcess 创建容器进程
func NewContainerProcess(command string, commandArgs ...string) *exec.Cmd {
	args := append([]string{"init", command}, commandArgs...)
	// 准备 cmd
	// /proc/self/exe 相当于在执行自身 即在 doki 进程内执行 doki 命令
	cmd := exec.Command("/proc/self/exe", args...)
	// 做隔离
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // 主机隔离
			syscall.CLONE_NEWPID | // PID 隔离
			syscall.CLONE_NEWNS | // namespace 隔离
			syscall.CLONE_NEWNET | // 网络隔离
			syscall.CLONE_NEWIPC, // 进程间通信隔离
		// TODO 用户隔离 暂时还没做 目标是实现 rootless
	}
	return cmd
}
