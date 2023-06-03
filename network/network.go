package network

import "net"

var Global_Network map[string]*Network //保存全局的网络信息
func Global_Network_Init() {
	Global_Network = make(map[string]*Network)
}

type Network struct { //单个新创建的网络信息包含网络名、子网划分、网络驱动
	Name   string     `json:"name"`   //网络名
	Subnet *net.IPNet `json:"subnet"` //子网
	Ips    string     `json:"ips"`    //网络划分
	Driver Driver     `json:"driver"` //网络驱动
}

func (n *Network) CreateNetwork() error {
	//1.创建网络驱动
	return n.Driver.Create(n)
}
