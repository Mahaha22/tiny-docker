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
		Name:    req.Name,                       //网络名
		Subnet:  subnet,                         //子网划分
		Driver:  network.NewDriver(req.Driver),  //网络驱动
		Ipalloc: network.NewIPAllocator(subnet), //初始化一个ip分配器，这里面保存所以已分配的ip
	}

	network.Global_Network[nw.Name] = nw //加入全局map表中
	err := nw.CreateNetwork()            //根据配置新建网络
	if err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (r *ContainerService) ListNetwork(context.Context, *emptypb.Empty) (*cmdline.Networks, error) {
	res := &cmdline.Networks{}
	for _, nw := range network.Global_Network {
		res.Nws = append(res.Nws, &cmdline.Network{
			Name:   nw.Name,
			Subnet: nw.Subnet.String(),
			Driver: GetDriverStr(nw.Driver),
		})
	}
	return res, nil
}

func GetDriverStr(driver network.Driver) string {
	//类型推断
	switch driver.(type) {
	case *network.BridgeDriver: //bridge
		{
			return "bridge"
		}
		//host
		//none
		//overlay
	}
	return ""
}
