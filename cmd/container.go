package cmd

import (
	"context"
	"fmt"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"

	"github.com/urfave/cli"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RunCommand(ctx *cli.Context) error {
	//解析出设定的配置
	req := &cmdline.Request{
		Args: &cmdline.Flag{
			It:      ctx.Bool("it"),
			ImageId: ctx.String("i"),
			Net:     ctx.String("net"),
			Name:    ctx.String("name"),
			Cpu:     ctx.String("cpu"),
			Mem:     ctx.String("mem"),
			Volmnt:  ctx.StringSlice("v"),
			Ports:   ctx.StringSlice("p"),
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
			return fmt.Errorf("term err = %v", err)
		}
	}
	return nil
}

func PsCommand() error {
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return fmt.Errorf("\nclient创建失败 : %v", err)
	}
	//ps远程调用
	containerinfo, err := client.PsContainer(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("ps container err = ", err)
		return err
	}
	// v := cmdline.Container{
	// 	ContainerId: "ab62e06d9d7f",
	// 	Command:     "redis:latest",
	// 	Image:       "docker-entrypoint.s…,",
	// 	CreateTime:  "2 seconds ago",
	// 	Status:      "RUNNING",
	// 	Ports:       "6379/tcp",
	// 	Name:        "redis",
	// 	VolumeMount: "/root/vol:/root/vol",
	// }
	fmt.Printf("%-15s %-10s %-15s %-15s %-10s %-10s %-20s %s\n", "CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "VOLUME", "NAMES")
	ports := "  -   "
	mnt := "  -   "
	for _, v := range containerinfo.Containers {
		if v.Ports != nil {
			ports = ""
			len := len(v.Ports)
			for i := 0; i < len; i = i + 2 {
				ports += v.Ports[i] + ":" + v.Ports[i+1]
				ports += ";"
			}
		}
		if v.VolumeMount != nil {
			mnt = ""
			len := len(v.VolumeMount)
			for i := 0; i < len; i = i + 2 {
				mnt += v.VolumeMount[i] + ":" + v.VolumeMount[i+1]
				mnt += ";"
			}
		}
		state := "  -   "
		if v.Status == "0" {
			state = ""
			state += "RUNNING"
		} else if v.Status == "1" {
			state = ""
			state += "STOPPED"
		}
		fmt.Printf("%-15s %-10s %-15s %-15s %-10s %-10s %-20s %s\n", v.ContainerId, v.Image, v.Command, v.CreateTime, state, ports, mnt, v.Name)
	}
	return nil
}

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
			return fmt.Errorf("term err = %v", err)
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

func KillCommand(containerIds []string) error {
	req := &cmdline.Request{
		Cmd: containerIds,
	}
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return fmt.Errorf("grpc client create err = %v", err)
	}
	//rpc调用杀死容器
	_, err = client.KillContainer(context.Background(), req)
	if err != nil {
		return fmt.Errorf("kill container err = %v", err)
	}
	return nil
}
