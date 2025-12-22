package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

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
	// 获取容器 init 进程对象 cmd
	cmd, writePipe := container.NewContainerProcess()

	// 处理选项
	// 这里父进程向匿名管道写入数据 虽然匿名管道有缓冲区 但缓冲区一旦被占满就会进程阻塞
	// 如果此时子进程没有启动 就无法从匿名管道里读取数据 就会一直阻塞
	// 所以需要保证父进程向匿名管道写入数据的时机为 子进程启动之后
	runner, err := processRunFlags(cmd, writePipe, command, commandArgs)
	if err != nil {
		return err
	}
	return runner()
}

// processRunFlags 用于处理选项 返回最终执行对象
func processRunFlags(cmd *exec.Cmd, writePipe *os.File, command string, commandArgs []string) (func() error, error) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// -i / -t 处理 将终端输入与 init 进程输入 IO 连通
	if runFlags.interactiveFlag || runFlags.ttyFlag {
		cmd.Stdin = os.Stdin
	}

	// -t 处理 外层 tty 忽略信号 传给内层 pty 处理
	if runFlags.ttyFlag {
		signal.Ignore(syscall.SIGINT, syscall.SIGQUIT)
	}

	// -t 处理 tty 的细节处理通过管道发送到 init 进程处理
	return func() error {
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = sendPipeArgs(writePipe, command, commandArgs, runFlags.ttyFlag)
		if err != nil {
			return err
		}
		return cmd.Wait()
	}, nil
}

// sendPipeArgs 用于父进程向匿名管道发送参数
func sendPipeArgs(writePipe *os.File, command string, commandArgs []string, ttyFlag bool) error {
	defer writePipe.Close()
	initArgs := container.InitArgs{
		Command:   command,
		Args:      commandArgs,
		EnableTty: ttyFlag,
	}
	initArgsBytes, err := json.Marshal(initArgs)
	if err != nil {
		return err
	}
	// 循环写入
	for len(initArgsBytes) > 0 {
		length, err := writePipe.Write(initArgsBytes)
		if err != nil {
			return err
		}
		initArgsBytes = initArgsBytes[length:]
	}
	return nil
}
