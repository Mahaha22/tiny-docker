package cmd

import (
	"context"
	"fmt"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"

	"github.com/urfave/cli"
)

func RunCommand(ctx *cli.Context) error {
	//解析出设定的配置
	req := &cmdline.Request{
		Args: &cmdline.Flag{
			It:      ctx.Bool("it"),
			ImageId: ctx.String("i"),
			Name:    ctx.String("name"),
			Cpu:     ctx.String("cpu"),
			Mem:     ctx.String("mem"),
			Volmnt:  ctx.StringSlice("v"),
		},
		Cmd: ctx.Args(),
	}
	/*
		$ ./tiny-docker run -it -mem 100m cpu 10000 /bin/bash
		cmd = /bin/bash
		flag = {
			It = true;
			Cpu = 10000;
			mem = 100m;
		}
	*/

	client, err := conn.GrpcClient_Single()
	if err != nil {
		return fmt.Errorf("\nclient创建失败 : %v", err)
	}
	response, err := client.RunContainer(context.Background(), req)
	if err != nil {
		return fmt.Errorf("\ngrpc RunContainer()调用失败 : %v", err)
	}
	//fmt.Printf("Create ContainerId \033[32m%v\033[0m Sucessfully\n", response.ContainerId)
	fmt.Printf("\033[32mCreate ContainerId %v Sucessfully\033[0m\n", response.ContainerId)
	// 此处判断是否需要交互
	if req.Args.It {
		//拿着服务器返回的id，建立一个双向流grpc，去连指定容器的 nsenter --target 1000 --mount --uts --ipc --net --pid bash
		err := newTerm(response.ContainerId)
		if err != nil {
			return fmt.Errorf("term err = ", err)
		}
	}
	return nil
}
