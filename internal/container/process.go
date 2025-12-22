package container

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// NewContainerProcess 创建容器进程
func NewContainerProcess() (*exec.Cmd, *os.File) {
	// 在父进程创建匿名管道
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Errorf("pipe error: %v", err)
		return nil, nil
	}
	// 准备 cmd
	// /proc/self/exe 相当于在执行自身 即在 doki 进程内执行 doki 命令 创建新进程
	cmd := exec.Command("/proc/self/exe", "init")
	// 做隔离
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // 主机隔离
			syscall.CLONE_NEWPID | // PID 隔离
			syscall.CLONE_NEWNS | // namespace 隔离
			syscall.CLONE_NEWNET | // 网络隔离
			syscall.CLONE_NEWIPC, // 进程间通信隔离
		// TODO 用户隔离 暂时还没做 目标是实现 rootless

	}
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}
