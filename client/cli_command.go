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
		cli.BoolFlag{ //进入容器终端
			Name:  "it",
			Usage: `enable tty`,
		},
		cli.StringSliceFlag{ //容器卷挂载
			Name:  "v",
			Usage: `mount volume -- vol1 : vol2`,
		},
		cli.StringFlag{
			Name:  "i",
			Usage: `containe image's id`,
		},
		cli.StringFlag{ //容器名
			Name:  "name",
			Usage: `set container name`,
		},
		cli.StringFlag{ //cpu用量限制
			Name:  "cpu",
			Usage: `limit the use of cpu`,
		},
		cli.StringFlag{ //内存用量限制
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

// 查看运行中的容器
var ps = cli.Command{
	Name:  "ps",
	Usage: "show all container",
	Action: func(context *cli.Context) error {

		//发送ps指令
		err := cmd.PsCommand()
		if err != nil {
			return fmt.Errorf("\nrun 容器启动失败: %v", err)
		}
		return err
	},
}
