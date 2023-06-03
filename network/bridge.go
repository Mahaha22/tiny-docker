// // 操作bridge网络驱动，也是默认的网络驱动
package network

import (
	"net"
	"os/exec"
	"tiny-docker/utils"
)

type BridgeDriver struct {
	//Dtype DriverType `json:"dtype"` //驱动类型
	Name string    `json:"name"` //名字
	Ip   net.IPNet `json:"ip"`   //ip地址
	Brd  net.IPNet `json:"brd"`  //广播地址
}

func (b *BridgeDriver) Create(nw *Network) error {
	b.Name = "td-" + utils.GetUniqueId()
	//1.添加网桥驱动
	cmd := exec.Command("brctl", "addbr", b.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	//2.配置网络
	//3.启动网桥
	return nil
}

func (b *BridgeDriver) Remove() error {
	//1.停止网桥
	//2.删除网桥
	cmd := exec.Command("brctl", "delbr", b.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
