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
		//发起启动一个容器请求
		conn, err := GetConn()
		if err != nil {
			return err
		}
		err = cmd.RunCommand(context, conn)
		return err
	},
}
