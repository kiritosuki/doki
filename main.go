package main

import (
	"os"

	"github.com/kiritosuki/doki/internal"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// log -> logrus logrus用log别名替代了 因为logrus的所有操作都是适应原生log的 无需担心

const usage = "doki is a simple container runtime implementation like docker."

func main() {
	app := cli.NewApp()
	app.Name = "doki"
	app.Usage = usage

	// 设置改app的command指令集 doki command
	// 注意：initCommand等 command是一个结构体 一个变量 不是一个函数
	app.Commands = []*cli.Command{
		internal.InitCommand,
		internal.RunCommand,
	}

	// 设置日志输出的形式json和标准输出
	app.Before = func(context *cli.Context) error {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		// 如果程序运行错误 打印错误日志 并立即终止程序(Fatal)
		log.Fatal(err)
	}
}
