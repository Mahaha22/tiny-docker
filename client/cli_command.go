package main

import (
	"fmt"
	"tiny-docker/cmd"

	"github.com/urfave/cli"
)

// 用于启动一个容器
var run = cli.Command{
	Name:  "run",
	Usage: `run a new container`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: `enable tty`,
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: `limit the use of cpu`,
		},
		cli.StringFlag{
			Name:  "mem",
			Usage: `limit the use of mem`,
		},
	},
	Action: func(context *cli.Context) error {
		//判断启动一个容器需要的最少参数
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command %v", len(context.Args()))
		}
		//发送run容器指令
		err := cmd.RunCommand(context)
		if err != nil {
			return fmt.Errorf("\nrun 容器启动失败: %v", err)
		}
		return err
	},
}
