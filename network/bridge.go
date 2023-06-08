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
	Name     string `json:"name"`     //名字
	Ip       net.IP `json:"ip"`       //ip地址
	Brd      net.IP `json:"brd"`      //广播地址
	Iptables string `json:"iptables"` //保存iptables信息，用于退出时清理
}

func (b *BridgeDriver) Create(nw *Network) error {
	if nw.Name != "" {
		b.Name = nw.Name
	} else {
		b.Name = "td-" + utils.GetUniqueId()
	}
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
	//4.设置iptables nat表POSTROUTING
	args = fmt.Sprintf("sudo iptables -t nat -A POSTROUTING -s %v/%v -j MASQUERADE", nw.Subnet.IP.String(), size)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return err
	}
	b.Iptables = fmt.Sprintf("sudo iptables -t nat -D POSTROUTING -s %v/%v -j MASQUERADE", nw.Subnet.IP.String(), size)
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
	//3.清理iptables
	fmt.Println(b.Iptables)
	cmd = exec.Command("/bin/bash", "-c", b.Iptables)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func SetBridgeNetwork(pid int, nw *Network, ports map[string]string) (net.IP, error) {
	//1.创建虚拟网卡对
	ip, ok := nw.Ipalloc.Allocate() //容器的ip地址
	if !ok {
		return nil, fmt.Errorf("ip allocate fail")
	}
	veth_br := "br" + utils.GetUniqueId() //放在网桥的一端
	veth_c := "c" + utils.GetUniqueId()   //放在容器内部的一端

	//1.新建网卡对
	args := fmt.Sprintf("ip link add %s type veth peer name %s", veth_br, veth_c)
	fmt.Println(args)
	cmd := exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	//2.将veth_c放入容器
	args = fmt.Sprintf("ip link set %s netns %v", veth_c, pid)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	//3.将veth_br放入网桥
	br, _ := nw.Driver.(*BridgeDriver)
	args = fmt.Sprintf("ip link set %s master %s", veth_br, br.Name)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	//4.给放入容器的网卡设置ip
	size, _ := nw.Subnet.Mask.Size()
	args = fmt.Sprintf("nsenter -n -t %v ip addr add %s/%d brd + dev %s", pid, ip, size, veth_c)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		fmt.Println("4 err  = ", err)
		return nil, err
	}

	//5.将veth_br启动
	args = fmt.Sprintf("ip link set dev %s up", veth_br)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	//6.将veth_c启动
	args = fmt.Sprintf("nsenter -n -t %v ip link set dev %s up", pid, veth_c)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		fmt.Println("6 err = ", err)
		return nil, err
	}
	//7.进入容器内部设置默认路由
	b, _ := nw.Driver.(*BridgeDriver)
	gateway := b.Ip
	args = fmt.Sprintf("nsenter -n -t %v ip route add default via %v dev %v", pid, gateway, veth_c)
	fmt.Println(args)
	cmd = exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		fmt.Println("6 err = ", err)
		return nil, err
	}
	//8.设置iptables nat表POSTROUTING
	//一般来说这一步在驱动创建的时候就已经完成
	//9.设置iptables nat表PREROUTING
	//举例 iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 192.168.0.2:80
	for HostPort, ContainerPort := range ports {
		//判断主机端口是否被占用
		err := utils.IsPortAvilable(HostPort)
		if err != nil {
			HostPortTable[HostPort] = true
			return nil, fmt.Errorf("port %v is in use", HostPort)
		}
		_, ok := HostPortTable[HostPort]
		if ok {
			return nil, fmt.Errorf("port %v is in use", HostPort)
		}
		args = fmt.Sprintf("iptables -t nat -A PREROUTING -p tcp -m tcp --dport %v -j DNAT --to-destination %v:%v", HostPort, ip, ContainerPort)
		fmt.Println(args)
		cmd = exec.Command("/bin/bash", "-c", args)
		if err := cmd.Run(); err != nil {
			fmt.Println("6 err = ", err)
			return nil, err
		}
	}
	return ip, nil
}
