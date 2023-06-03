package service

import (
	"context"
	"net"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/network"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ContainerService) CreateNetwork(ctx context.Context, req *cmdline.Network) (*emptypb.Empty, error) {
	_, subnet, _ := net.ParseCIDR(req.Subnet)
	nw := &network.Network{
		Name:   req.Name,
		Subnet: subnet,
		Driver: network.NewDriver(req.Driver),
	}
	network.Global_Network[nw.Name] = nw //加入全局map表中
	err := nw.CreateNetwork()            //根据配置新建网络
	if err != nil {
		return nil, err
	}
	return nil, nil
}
