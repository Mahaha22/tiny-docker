package network

import (
	"fmt"
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
		Name:    req.Name,
		Subnet:  subnet,
		Driver:  NewDriver(req.Driver),
		Ipalloc: NewIPAllocator(subnet),
	}
	//Global_Network[nw.Name] = nw //加入全局map表中
	nw.CreateNetwork() //根据配置新建网络
	nw.Driver.Remove() //清除网络
}

func TestAllocateIp(t *testing.T) {
	_, subnet, _ := net.ParseCIDR("192.168.0.55/16")
	allocator := NewIPAllocator(subnet)

	for i := 0; i < 578; i++ {
		ip, success := allocator.Allocate()
		if success {
			fmt.Println("Allocated IP:", ip)
		} else {
			fmt.Println("No available IP addresses.")
		}
	}
	iptmp := net.ParseIP("192.168.2.33")
	allocator.Realese(iptmp)
	for i := 0; i <= 3; i++ {
		ip, success := allocator.Allocate()
		if success {
			fmt.Println("Allocated IP:", ip)
		} else {
			fmt.Println("No available IP addresses.")
		}
	}

	// 获取广播地址
	ip := make(net.IP, len(subnet.IP))
	for i := 0; i < len(subnet.IP); i++ {
		ip[i] = subnet.IP[i] | ^subnet.Mask[i] //掩码取反与ip地址取或
	}
	fmt.Println(ip)
}
