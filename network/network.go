package network

import (
	"fmt"
	"net"
)

var Global_Network map[string]*Network //保存全局的网络信息
func Global_Network_Init() {
	Global_Network = make(map[string]*Network)
}
func RemoveAllNetwork() {
	for name, nw := range Global_Network {
		err := nw.RemoveNetwork()
		if err != nil {
			fmt.Printf("network %s delete error:%v\n", name, err)
		}
	}
}

type Network struct { //单个新创建的网络信息包含网络名、子网划分、网络驱动
	Name    string     `json:"name"`    //网络名
	Subnet  *net.IPNet `json:"subnet"`  //子网
	Ipalloc *IPAlloc   `json:"ipalloc"` //网络划分
	Driver  Driver     `json:"driver"`  //网络驱动
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
