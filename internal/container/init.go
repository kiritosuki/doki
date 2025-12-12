package container

import (
	"errors"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

// InitCmd 是 init 命令
var InitCmd = &cobra.Command{
	Use:    "init [args...]",
	Hidden: true,
	Short:  "initialize the container, don't call it outside",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(args)
	},
}

// runInit 是 init 的核心实现
func runInit(args []string) error {
	// 将 mount 挂载设置为 private 避免传播
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
	if len(args) < 1 {
		return errors.New("COMMAND required by init")
	}
	command := args[0]
	return syscall.Exec(command, args, os.Environ())
}
