package container

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

// InitCmd 是 init 命令
var InitCmd = &cobra.Command{
	Use:    "init",
	Hidden: true,
	Short:  "initialize the container, don't call it outside",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

// runInit 是 init 的核心实现
func runInit() error {
	// 将 mount 挂载设置为 private 避免子进程传播到父进程
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}
	mountFlags := syscall.MS_NOEXEC | // 禁止在挂载点执行可执行文件
		syscall.MS_NOSUID | // 禁止在挂载点用 su 提升权限
		syscall.MS_NODEV // 禁止在挂载点使用 /dev 设备文件

	// 挂载 /proc
	err = syscall.Mount("proc", "/proc", "proc", uintptr(mountFlags), "")
	if err != nil {
		return err
	}
	readPipe := os.NewFile(uintptr(3), "readPipe")
	defer readPipe.Close()
	// 这里用 io 包的 readAll 函数 不用自己写循环 可以确保读完
	// 但是 io 包没有对应的 writeALl 函数
	bytes, err := io.ReadAll(readPipe)
	if err != nil {
		// 报错关闭 fd 资源
		return err
	}
	var initArgs InitArgs
	// 反序列化要传递指针
	err = json.Unmarshal(bytes, &initArgs)
	if err != nil {
		return err
	}
	command := initArgs.Command
	args := initArgs.Args
	// 用户进程对象 cmd
	cmd := exec.Command(command, args...)
	enableTty := initArgs.EnableTty

	if enableTty {
		log.Info("TODO: 处理 tty 逻辑")
	}

	//// -t 单独处理
	//if enableTty {
	//	// 给用户进程如 bash 进程单独分配 pty
	//	// pty.start 内部会实现:
	//	// 把用户进程作为新的 session 的 leader
	//	// 把该 pty 作为控制终端
	//	// 把用户进程设置为前台进程组的 leader
	//	ptmx, err := pty.Start(cmd)
	//	if err != nil {
	//		return err
	//	}
	//	defer ptmx.Close()
	//	// 将 pty 的 IO 通路与 init 进程的 IO通路连通
	//	// pty master 的辅助输入输出是阻塞操作 必须用协程
	//	go func() { io.Copy(ptmx, os.Stdin) }()
	//	go func() { io.Copy(os.Stdout, ptmx) }()
	//
	//	return cmd.Wait()
	//}

	// 其余情况 将用户进程的 IO 通路与 init 进程连通
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
