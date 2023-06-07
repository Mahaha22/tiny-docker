package network

import (
	"fmt"
	"net"
	"os"
)

func init() {
	Global_Network = make(map[string]*Network)
	//初始化一个默认网络
	_, subnet, _ := net.ParseCIDR("192.168.0.0/20")
	nw := &Network{
		Name:    "Tiny-docker",          //网络名
		Subnet:  subnet,                 //子网划分
		Driver:  NewDriver("bridge"),    //网络驱动
		Ipalloc: NewIPAllocator(subnet), //初始化一个ip分配器，这里面保存所以已分配的ip
	}
	Global_Network[nw.Name] = nw //加入全局map表中
	err := nw.CreateNetwork()    //根据配置新建网络
	if err != nil {
		fmt.Println("容器网络初始化失败")
		RemoveAllNetwork()
		os.Exit(-1)
	}
}

var Global_Network map[string]*Network //保存全局的网络信息
func RemoveAllNetwork() {
	for name, nw := range Global_Network {
		err := nw.RemoveNetwork()
		if err != nil {
			fmt.Printf("network %s delete error:%v\n", name, err)
		}
	}
}

type Network struct { //单个新创建的网络信息包含网络名、子网划分、网络驱动
	Name    string            `json:"name"`    //网络名
	Subnet  *net.IPNet        `json:"subnet"`  //子网
	Ipalloc *IPAlloc          `json:"ipalloc"` //网络划分
	Driver  Driver            `json:"driver"`  //网络驱动
	Port    map[string]string `json:"driver"`  //port映射
}

func (n *Network) CreateNetwork() error {
	//这里要做断言

	//1.创建网络驱动
	return n.Driver.Create(n)
}

func (n *Network) RemoveNetwork() error {
	err := n.Driver.Remove() //删除驱动
	if err != nil {
		return err
	}
	delete(Global_Network, n.Name) //从network网络map表中删除
	return nil
}

func ApplyNetwork(pid int, nw *Network) (net.IP, error) {
	switch nw.Driver.(type) {
	case *BridgeDriver:
		{
			ip, err := SetBridgeNetwork(pid, nw)
			return ip, err
		}
	case *HostDriver:
		{
			hostip, err := GetHostIp()
			return hostip, err
		}
	case *NoneDriver:
		{
			return net.IP{}, nil
		}
	}
	return net.IP{}, nil
}

func GetHostIp() (net.IP, error) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get hostname:", err)
		return net.IP{}, err
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Failed to lookup IP:", err)
		return net.IP{}, err
	}

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4, nil
		}
	}
	return net.IP{}, err
}
