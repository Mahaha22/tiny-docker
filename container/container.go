package container

import (
	"fmt"
	"os/exec"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/namespace"
	"tiny-docker/utils"
)

var Global_ContainerMap map[string]*Container //索引所有容器
func Global_ContainerMap_Init() {
	Global_ContainerMap = make(map[string]*Container)
}

// 容器的数据结构
type Container struct {
	ContainerId string
	Name        string
	CpuLimit    string
	MemLimit    string
	RealPid     int //容器在主机上真实的进程id
	//NetWork
}

// 实例化一个容器
func CreateContainer(req *cmdline.Request) *Container {
	return &Container{
		ContainerId: utils.GetUniqueId(),
		CpuLimit:    req.Args.Cpu,
		MemLimit:    req.Args.Mem,
	}
}

// 容器相关配置初始化，例如namesapce和cgroup
func (c *Container) Init() error {
	//加载容器镜像
	/*
		待补充
	*/
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = namespace.Cloneflag                //独立命名空间配置
	stdinPipe, _ := cmd.StdinPipe()                      //创建一个连接子进程输入流的管道用于父进程向内传递信息
	mountinfo := "mount -t proc proc /proc -o private\n" //挂载proc private标志，表示这个挂载点是与其他挂载点相互独立的，不会影响其他挂载点，也不会被其他挂载点所影响。

	//加载容器资源配置信息
	//pid := os.Getpid()
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start() fail %v", err)
	}
	c.RealPid = cmd.Process.Pid
	fmt.Println("pid = ", c.RealPid)
	stdinPipe.Write([]byte(mountinfo))
	//cmd.Wait()
	return nil
}
