package container

import (
	"encoding/json"
	"io"
	"os"
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
	// 这里用 io 包的 readAll 函数 不用自己写循环 可以确保读完
	// 但是 io 包没有对应的 writeALl 函数
	bytes, err := io.ReadAll(readPipe)
	if err != nil {
		// 报错关闭 fd 资源
		readPipe.Close()
		return err
	}
	// 不用 defer 关闭 原因是 defer 理论上在 syscall.Exec 之后执行
	// 但 syscall.Exec 执行成功后 go runtime 会被整个替换掉
	// fd.close 也不复存在 也没有机会执行 写 defer 无意义
	err = readPipe.Close()
	if err != nil {
		return err
	}
	var initArgs InitArgs
	// 反序列化要传递指针
	err = json.Unmarshal(bytes, &initArgs)
	if err != nil {
		return err
	}
	command := initArgs.Command
	// syscall.Exec 是最底层的执行函数
	// 会把要执行的命令替换当前进程
	// 要求 args[0] == command
	args := append([]string{command}, initArgs.Args...)
	return syscall.Exec(command, args, os.Environ())
}
