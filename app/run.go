package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/kiritosuki/doki/internal/container"
	"github.com/spf13/cobra"
)

type RunOpts struct {
	interactiveFlag bool
	ttyFlag         bool
}

var runFlags = &RunOpts{}

// runCmd 是 run 命令
var runCmd = &cobra.Command{
	Use:   "run [args...]",
	Short: "run a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRun(args)
	},
}

// init 是 runCmd 的初始化配置函数
func init() {
	runCmd.Flags().BoolVarP(&runFlags.interactiveFlag, "interactive", "i", false, "open stdin")
	runCmd.Flags().BoolVarP(&runFlags.ttyFlag, "tty", "t", false, "allocate a tty")
}

// runRun 是 run 命令的核心实现函数
func runRun(args []string) error {
	if len(args) == 0 {
		return errors.New("image required, usage: doki run IMAGE [COMMAND] [ARGS...]")
	}
	// image 是镜像名 也是doki run 后紧跟的第一个参数
	image := args[0]
	command := "/bin/sh"
	var commandArgs []string
	if len(args) >= 2 {
		// 获取命令
		command = args[1]
	}
	if len(args) >= 3 {
		// 获取命令参数
		commandArgs = args[2:]
	}
	// TODO 处理镜像
	fmt.Printf("处理镜像...image: %s\n", image)
	// 获取容器进程对象 cmd
	cmd := container.NewContainerProcess(command, commandArgs...)

	// 处理选项
	runner, err := processRunFlags(cmd)
	if err != nil {
		return err
	}
	return runner()
}

func processRunFlags(cmd *exec.Cmd) (func() error, error) {
	// -t 处理
	if runFlags.ttyFlag {
		// 子进程作为新的会话
		cmd.SysProcAttr.Setsid = true
		// 给子进程分配 tty
		cmd.SysProcAttr.Setctty = true

		return func() error {
			// 新分配一对 pty 把新的 pty slave 作为 cmd 的控制终端
			// 返回 pty master 即 ptmx
			// 通信方式可以这么理解:
			// doki - pty master - pty slave - bash
			ptmx, err := pty.Start(cmd)
			if err != nil {
				return err
			}
			defer ptmx.Close()
			// 这里 pty master 复制输入输出 都是阻塞操作 必须要用协程
			go func() { io.Copy(ptmx, os.Stdin) }()
			go func() { io.Copy(os.Stdout, ptmx) }()
			return cmd.Wait()
		}, nil
	}

	// -i 处理
	if runFlags.interactiveFlag {
		cmd.Stdin = os.Stdin
	}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return func() error {
		err := cmd.Start()
		if err != nil {
			return err
		}
		return cmd.Wait()
	}, nil
}
