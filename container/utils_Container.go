// 当服务器退出的时候，清除所有的容器
package container

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"tiny-docker/cgroup"
	"tiny-docker/network"
	"tiny-docker/overlayfs"
)

// 用于检测容器是否存活
func ContainerIsAlive(c *Container) state {
	statfile := fmt.Sprintf("/proc/%d/stat", c.RealPid)
	file, _ := os.Open(statfile)
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, " ")
		STAT := str[2]
		fmt.Println("pid = ", c.RealPid, " stat is ", STAT)
		if STAT == "Z" {
			return EXITDED
		} else if STAT == "T" {
			return STOPPED
		}
	}
	return RUNNING
}

// 监视容器的存活状态
func Monitor() {
	for {
		for k, v := range Global_ContainerMap {
			//fmt.Println("v = ", v, "pid = ", v.RealPid, "alive = ", ContainerIsAlive(v))
			S := ContainerIsAlive(v)
			if S == EXITDED { //如果容器死亡，直接清除
				delete(Global_ContainerMap, k)
				//当容器退出时需要对资源进行销毁
				err := cgroup.DestroyCgroup(k)
				fmt.Println("monitor err:", err)
			} else if S == STOPPED {
				Global_ContainerMap[k].Status = STOPPED //更改一下容器的状态
			}
		}
		//fmt.Println("***", Global_ContainerMap)
		time.Sleep(time.Second)
	}
}

// 当服务器退出时需清除所有容器信息
func KillallContainer() {
	for _, v := range Global_ContainerMap {
		KillContainer(v)
	}
}

// 清除容器
func KillContainer(c *Container) {
	c.Remove()
	delete(Global_ContainerMap, c.ContainerId)
	cgroup.DestroyCgroup(c.ContainerId)
	ReleaseNetWork(c)
}

func KillallVolume() {
	for _, c := range Global_ContainerMap { //1.清除并卸载所有挂载点
		KillVolume(c)
	}
}

func KillVolume(c *Container) {
	// 1.1卸载容器卷
	if err := overlayfs.RemoveMountFs(c.Volmnt, c.ContainerId); err != nil {
		fmt.Println("RemoveMountFs err = ", err)
	}
	// 1.2卸载容器根文件系统
	if err := overlayfs.DeleteOverlayMnt(c.ContainerId); err != nil {
		fmt.Println("DeleteOverlayMnt err = ", err)
	}
}

func ReleaseNetWork(c *Container) {
	//释放占用的ip
	if c.Net.Ip != nil {
		c.Net.Network_Membership.Ipalloc.Realese(c.Net.Ip)
	}
	//释放端口映射
	if c.Net.Port != nil {
		for HostPort, ContainerPort := range c.Net.Port {
			delete(network.HostPortTable, HostPort)
			args := fmt.Sprintf("iptables -t nat -D PREROUTING -p tcp -m tcp --dport %v -j DNAT --to-destination %v:%v", HostPort, c.Net.Ip, ContainerPort)
			fmt.Println(args)
			cmd := exec.Command("/bin/bash", "-c", args)
			if err := cmd.Run(); err != nil {
				fmt.Println("6 err = ", err)
			}
		}
	}
}
