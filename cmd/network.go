package cmd

import (
	"context"
	"fmt"
	"net"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"

	"github.com/urfave/cli"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 创建网络
func CreateNetwork(ctx *cli.Context) error {
	driver := ctx.String("d")
	var subnet string
	if driver == "" || driver == "bridge" { //如果使用bridge驱动需要给出subnet地址
		subnet = ctx.String("subnet")                       //解析出cidr地址
		if _, _, err := net.ParseCIDR(subnet); err != nil { //判断cidr是否有效
			return fmt.Errorf("subnet format is err = %v", err)
		}
	}
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

// 显示网络
func ListNetwork() error {
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return err //创建grpc客户端出现问题，直接返回
	}
	networks, err := client.ListNetwork(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}
	fmt.Printf("%-15s %-15s %s\n", "Subnet", "Name", "Driver")
	for _, nw := range networks.Nws {
		fmt.Printf("%-15s %-15s %s\n", nw.Subnet, nw.Name, nw.Driver)
	}
	return nil
}

// 删除网络
func DeleteNetwork(args []string) error {
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return err //创建grpc客户端出现问题，直接返回
	}
	for _, nw_name := range args {
		_, err = client.DelNetwork(context.Background(), &cmdline.Network{Name: nw_name})
		if err != nil {
			return err
		}
	}
	return nil
}
