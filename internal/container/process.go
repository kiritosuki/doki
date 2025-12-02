package container

import (
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 准备一个新进程 作为容器运行的进程
func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	// 在Linux系统中 当前正在运行的程序会被放在/proc/self/exe
	// 这里是准备一个新的进程 执行当前正在运行的程序
	// 例如：
	// doki run /bin/bash
	// 这里当前的程序就是doki
	// 这里就是准备一个新的进程 运行 doki init /bin/bash
	cmd := exec.Command("/proc/self/exe", args...)
	// 设置 namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | // 隔离 hostname
			syscall.CLONE_NEWPID | // 隔离 PID
			syscall.CLONE_NEWNS | // 隔离挂载点 mount points
			syscall.CLONE_NEWNET | // 隔离网络
			syscall.CLONE_NEWIPC, // 隔离 IPC (进程间通信：内存 信号等)
		//syscall.CLONE_NEWUSER, 补充 隔离 user docker源码中没有做用户隔离 导致docker命令都需要sudo
		//TODO 后续可能会加上用户隔离 来免除sudo 有点复杂

	}

	if tty {
		// 如果允许 tty 就把 cmd 的标准输入输出和 操作系统终端绑定
		// 这样传给终端的命令会传给 tty
		// tty 打印的结果也会显示在终端上
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
