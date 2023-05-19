package service

import (
	"context"
	"tiny-docker/grpc/cmdline"
)

type ContainerService struct {
	cmdline.UnimplementedServiceServer
}

func (r *ContainerService) RunContainer(ctx context.Context, req *cmdline.Request) (*cmdline.RunResponse, error) {
	//实现具体的业务逻辑
	//启动docker
	//解析请求

	// conf := flag{

	// }
	// id, err := InitContainer()

	return &cmdline.RunResponse{
		ContainerId: 100,
	}, nil
}
