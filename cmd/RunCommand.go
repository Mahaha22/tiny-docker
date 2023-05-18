package cmd

import (
	"context"
	"fmt"
	"log"
	"tiny-docker/grpc/run"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func RunCommand(ctx *cli.Context, conn *grpc.ClientConn) error {
	defer conn.Close()
	//解析出设定的配置
	req := &run.RunRequest{
		It:  ctx.Bool("it"),
		Cpu: ctx.String("cpu"),
		Mem: ctx.String("mem"),
		Cmd: ctx.Args(),
	}
	/*
		$ ./tiny-docker run -it -mem 100m cpu 10000 /bin/bash
		cmd = /bin/bash
		conf = {
			It = true;
			Cpu = 10000;
			mem = 100m;
		}
	*/

	//创建一个grpc客户端
	grpc_client := run.NewRunServiceClient(conn)
	//远程调用
	response, err := grpc_client.RunContainer(context.Background(), req)
	if err != nil {
		log.Fatal("grpc调用失败 : ", err)
	}
	fmt.Println("id = ", response.ContainerId)
	return nil
}
