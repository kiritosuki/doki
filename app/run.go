package app

import (
	"errors"
	"fmt"
	"os"

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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 处理选项
	if runFlags.interactiveFlag {
		cmd.Stdin = os.Stdin
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
