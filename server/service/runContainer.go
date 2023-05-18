package service

import (
	"context"
	"tiny-docker/grpc/run"
)

type ContainerService struct {
	run.UnimplementedRunServiceServer
}

func (r *ContainerService) RunContainer(ctx context.Context, req *run.RunRequest) (*run.RunResponse, error) {
	//实现具体的业务逻辑
	//启动docker
	//解析请求

	// conf := flag{

	// }
	// id, err := InitContainer()

	return &run.RunResponse{
		ContainerId: 100,
	}, nil
}
