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
		cli.StringFlag{ //指定容器镜像
			Name:  "i",
			Usage: `containe image's id`,
		},
		cli.StringFlag{ //指定网络
			Name:  "net",
			Usage: `set container name`,
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: `HostPort:ContainerPort`,
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
			return fmt.Errorf("查询失败%v", err)
		}
		return err
	},
}

// 给容器发送命令
var exec = cli.Command{
	Name:  "exec",
	Usage: "send a cmd to container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: `enter a container and create new tty`,
		},
		// cli.StringFlag{
		// 	Name:  "id",
		// 	Usage: `container's id`,
		// },
	},
	Action: func(context *cli.Context) error {
		//判断启动一个容器需要的最少参数
		if len(context.Args()) < 2 {
			return fmt.Errorf("missing container command %v", len(context.Args()))
		}
		//rpc调用
		err := cmd.ExecCommand(context)
		if err != nil {
			fmt.Println("exec fail err = ", err)
			return err
		}
		return err
	},
}

// 终结容器
var kill = cli.Command{
	Name:  "kill",
	Usage: "kill a container",
	Action: func(context *cli.Context) error {
		//判断终结一个容器需要的最少参数
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command %v", len(context.Args()))
		}
		//rpc调用
		err := cmd.KillCommand(context.Args())
		if err != nil {
			return fmt.Errorf("kill container err = %v", err)
		}
		return nil
	},
}

// 管理容器网络
var network = cli.Command{
	Name:  "network",
	Usage: "manage container network",
	//network下包含多个子命令
	//create 创建网络
	//ls     显示所有网络
	//delete 删除网络
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: `new a container`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subnet",
					Usage: `Subnet in CIDR / 192.168.0.0/24`,
				},
				cli.StringFlag{
					Name:  "d",
					Usage: `network driver`,
				},
			},
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 { //至少有一个参数指定网络的名字
					return fmt.Errorf("missing network name")
				}
				//grpc远程调用创建新的网络
				err := cmd.CreateNetwork(context)
				if err != nil {
					return fmt.Errorf("create new network err = %v", err)
				}
				return nil
			},
		},
		{
			Name:  "ls",
			Usage: `show all network`,
			Action: func(context *cli.Context) error {
				//触发grpc远程调用
				return cmd.ListNetwork()
			},
		},
		{
			Name:  "delete",
			Usage: `delete network config`,
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("need network's name")
				}
				return cmd.DeleteNetwork(context.Args())
			},
		},
	},
}
