package service

import (
	"context"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
)

func (r *ContainerService) RunContainer(ctx context.Context, req *cmdline.Request) (*cmdline.RunResponse, error) {
	//实现具体的业务逻辑
	newContainer := container.CreateContainer(req)                         //实例化一个容器
	container.Global_ContainerMap[newContainer.ContainerId] = newContainer //放入Global_ContainerMap
	if err := newContainer.Init(); err != nil {                            //初始化容器
		return &cmdline.RunResponse{
			ContainerId: "", //容器创建失败返回空
		}, err
	}
	return &cmdline.RunResponse{
		ContainerId: newContainer.ContainerId, //容器创建成功返回真实的容器id
	}, nil
}
