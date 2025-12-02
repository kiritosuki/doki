package internal

import (
	"fmt"

	"github.com/kiritosuki/doki/internal/app"
	"github.com/kiritosuki/doki/internal/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var RunCommand = &cli.Command{
	Name:  "run",
	Usage: `create a container from image`,
	Flags: []cli.Flag{
		// 这里把 -i -t 合二为一了 允许tty
		&cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	// run 真正执行的函数是这个
	// 1. 判断参数是否包含 command
	// 2. 获取用户指定的 command
	// 3. 执行该命令
	Action: func(context *cli.Context) error {
		if context.Args().Len() < 1 {
			return fmt.Errorf("missing command")
		}
		command := context.Args().Get(0)
		tty := context.Bool("it")
		app.Run(tty, command)
		return nil
	},
}

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "init container process, don't call it outside",

	// 获取传递过来的 command 参数
	// 执行初始化操作
	Action: func(context *cli.Context) error {
		log.Info("init container process...")
		command := context.Args().Get(0)
		log.Infof("command: %s", command)
		err := container.RunContainerInitProcess(command, nil)
		return err
	},
}
