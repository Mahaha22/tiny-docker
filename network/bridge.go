// // 操作bridge网络驱动，也是默认的网络驱动
package network

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"tiny-docker/utils"
)

type BridgeDriver struct {
	//Dtype DriverType `json:"dtype"` //驱动类型
	Name string `json:"name"` //名字
	Ip   net.IP `json:"ip"`   //ip地址
	Brd  net.IP `json:"brd"`  //广播地址
}

func (b *BridgeDriver) Create(nw *Network) error {
	b.Name = "td-" + utils.GetUniqueId()
	//1.添加网桥驱动
	cmd := exec.Command("brctl", "addbr", b.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	//2.配置网络
	ip, ret := nw.Ipalloc.Allocate()
	if !ret { //分配ip失败s
		return fmt.Errorf("分配ip地址失败")
	}
	b.Ip = ip //网桥的ip地址
	brd := make(net.IP, 4)
	for i := 0; i < len(nw.Subnet.IP); i++ {
		brd[i] = nw.Subnet.IP[i] | ^nw.Subnet.Mask[i] //掩码取反与ip地址取或
	}
	b.Brd = brd
	size, _ := nw.Subnet.Mask.Size()
	ipstr := b.Ip.String() + "/" + strconv.Itoa(size)
	args := fmt.Sprintf("ip addr add %s brd %s dev %s", ipstr, b.Brd, b.Name)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		fmt.Println("err = ", err)
		return err
	}
	//3.启动网桥
	args = fmt.Sprintf("ip link set %s up", b.Name)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (b *BridgeDriver) Remove() error {
	//1.停止网桥
	args := fmt.Sprintf("ip link set %s down", b.Name)
	cmd := exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return err
	}
	//2.删除网桥
	cmd = exec.Command("brctl", "delbr", b.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
