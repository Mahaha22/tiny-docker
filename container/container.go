package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"tiny-docker/cgroup"
	"tiny-docker/conf"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/utils"
)

var Global_ContainerMap map[string]*Container //索引所有容器
func Global_ContainerMap_Init() {
	Global_ContainerMap = make(map[string]*Container)
}

// 容器的数据结构
type Container struct {
	ContainerId  string               //容器的唯一标识
	Name         string               //容器的自定义名称
	CgroupRes    conf.Cgroupflag      //cgroup资源
	NameSpaceRes *syscall.SysProcAttr //Namespace隔离标准当作克隆标志
	RealPid      int                  //容器在主机上真实的进程id
	//NetWork
}

// 实例化一个容器
func CreateContainer(req *cmdline.Request) *Container {
	return &Container{
		ContainerId: utils.GetUniqueId(),
		Name:        req.Args.Name,
		CgroupRes: conf.Cgroupflag{
			Cpulimit: req.Args.Cpu,
			Memlimit: req.Args.Mem,
		},
		NameSpaceRes: conf.Cloneflag,
	}
}

// 容器相关配置初始化，例如namesapce和cgroup
func (c *Container) Init() error {
	//1.加载容器镜像
	/*
		待补充
	*/
	//cmd := exec.Command("sleep", "30000")
	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = c.NameSpaceRes                     //独立命名空间配置
	stdinPipe, _ := cmd.StdinPipe()                      //创建一个连接子进程输入流的管道用于父进程向内传递信息
	mountinfo := "mount -t proc proc /proc -o private\n" //挂载proc private标志，表示这个挂载点是与其他挂载点相互独立的，不会影响其他挂载点，也不会被其他挂载点所影响。
	chrootinfo := "chroot /root/busybox /bin/sh\n"       //设置根目录
	pathinfo := "export PATH=:/bin\n"                    //设置环境变量
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start() fail %v", err)
	}
	c.RealPid = cmd.Process.Pid
	//2.容器内的第一个程序已经启动，针对这个容器以他的唯一标识container_id在/sys/fs/cgroup/<subsystem>/tiny-docker/<container_id>创建新的subsystem
	if err := cgroup.SetCgroup(c.ContainerId, &c.CgroupRes, c.RealPid); err != nil {
		return fmt.Errorf("cgroup资源配置出错 :%v", err)
	}

	stdinPipe.Write([]byte(mountinfo))
	stdinPipe.Write([]byte(chrootinfo))
	stdinPipe.Write([]byte(pathinfo))
	//cmd.Wait()
	return nil
}
func (c *Container) Remove() {
	process, _ := os.FindProcess(c.RealPid)
	process.Signal(syscall.Signal(9))
	fmt.Println(c.RealPid, " is killed")
}
