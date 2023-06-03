package cmd

import (
	"context"
	"fmt"
	"tiny-docker/grpc/conn"
)

func PsCommand() error {
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return fmt.Errorf("\nclient创建失败 : %v", err)
	}
	//ps远程调用
	containerinfo, err := client.PsContainer(context.Background(), nil)
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
