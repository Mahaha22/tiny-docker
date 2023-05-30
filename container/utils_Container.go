// 当服务器退出的时候，清除所有的容器
package container

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"tiny-docker/cgroup"
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
	for k, v := range Global_ContainerMap {
		v.Remove()
		delete(Global_ContainerMap, k)
		cgroup.DestroyCgroup(v.ContainerId)
	}
}
