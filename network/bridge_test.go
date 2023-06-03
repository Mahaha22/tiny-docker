package network

import (
	"net"
	"testing"
	"tiny-docker/grpc/cmdline"
)

func TestBridge(t *testing.T) {
	req := cmdline.Network{
		Name:   "haohao",
		Subnet: "192.168.0.1/28",
		Driver: "bridge",
	}
	_, subnet, _ := net.ParseCIDR(req.Subnet)
	nw := &Network{
		Name:   req.Name,
		Subnet: subnet,
		Driver: NewDriver(req.Driver),
	}
	//Global_Network[nw.Name] = nw //加入全局map表中
	nw.CreateNetwork() //根据配置新建网络
}
