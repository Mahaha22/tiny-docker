package service

import (
	"context"
	"fmt"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
)

func (r *ContainerService) PsContainer(context.Context, *cmdline.Request) (*cmdline.ContainerInfo, error) {
	info := &cmdline.ContainerInfo{}
	for _, c := range container.Global_ContainerMap {
		tmp := &cmdline.Container{
			ContainerId: c.ContainerId,
			Image:       c.Image,
			CreateTime:  c.CreateTime,
			Status:      fmt.Sprintf("%v", c.Status),
			//Ports:      ,
			Name: c.Name,
		}
		for vol1, vol2 := range c.Volmnt {
			tmp.VolumeMount = append(tmp.VolumeMount, vol1, vol2)
		}
		info.Containers = append(info.Containers, tmp)
	}
	return info, nil
}