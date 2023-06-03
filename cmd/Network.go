package cmd

import (
	"context"
	"fmt"
	"net"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"

	"github.com/urfave/cli"
)

func CreateNetwork(ctx *cli.Context) error {
	subnet := ctx.String("subnet")                      //解析出cidr地址
	if _, _, err := net.ParseCIDR(subnet); err != nil { //判断cidr是否有效
		return fmt.Errorf("subnet format is err = %v", err)
	}
	driver := ctx.String("d")
	name := ctx.Args()[0]
	req := cmdline.Network{
		Subnet: subnet,
		Driver: driver,
		Name:   name,
	}
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return err
	}
	_, err = client.CreateNetwork(context.Background(), &req)
	if err != nil {
		return err
	}
	fmt.Printf("\033[32mCreate Network %v Sucessfully\033[0m\n", name)
	return nil
}
