package container

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
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
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
