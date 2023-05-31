package service

import (
	"context"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
)

func (r *ContainerService) KillContainer(ctx context.Context, req *cmdline.Request) (*cmdline.RunResponse, error) {
	for _, v := range req.Cmd {
		if c, ok := container.Global_ContainerMap[v]; ok { //容器存在
			container.KillVolume(c)
			container.KillContainer(c)
			return &cmdline.RunResponse{}, nil
		}
	}
	return &cmdline.RunResponse{}, nil
}
