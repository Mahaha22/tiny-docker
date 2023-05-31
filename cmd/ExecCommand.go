package cmd

import (
	"context"
	"fmt"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"

	"github.com/urfave/cli"
)

func ExecCommand(ctx *cli.Context) error {
	it := ctx.Bool("it")

	req := &cmdline.Request{
		Cmd: ctx.Args(),
	}
	//分两种情况
	//1.需要建立新终端
	if it {
		err := newTerm(ctx.Args()[0])
		if err != nil {
			return fmt.Errorf("term err = ", err)
		}
	} else { //2.不需要建立新终端
		client, err := conn.GrpcClient_Single()
		if err != nil {
			return fmt.Errorf("\nclient创建失败 : %v", err)
		}
		res, err := client.ExecContainer(context.Background(), req)
		if err != nil {
			return err
		}
		//fmt.Println(res)
		if res.Errinfo != "" {
			fmt.Printf("\033[31m[%v]:%v\033[0m\n", ctx.Args()[0], res.Errinfo)
		} else {
			fmt.Printf("\033[32m[%v]\033[0m:\n%v\n", ctx.Args()[0], res.Outinfo)
		}
	}
	return nil
}
