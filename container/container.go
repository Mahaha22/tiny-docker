package container

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	os.Chdir("/root/go/tiny-docker/container/lib") //钩子程序的位置
	cmd := exec.Command("./task")                  //启动容器的钩子
	r, w, _ := os.Pipe()                           //用于跟这个钩子程序通信
	cmd.ExtraFiles = []*os.File{w}                 //将管道的一端传递给钩子
	if err := cmd.Start(); err != nil {
		fmt.Println("cmd err = ", err)
		return err
	}

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	r.Close()
	pid, _ := strconv.Atoi(string(buf[:n-1]))
	c.RealPid = pid

	//2.容器内的第一个程序已经启动，针对这个容器以他的唯一标识container_id在/sys/fs/cgroup/<subsystem>/tiny-docker/<container_id>创建新的subsystem
	if err := cgroup.SetCgroup(c.ContainerId, &c.CgroupRes, c.RealPid); err != nil {
		return fmt.Errorf("cgroup资源配置出错 :%v", err)
	}
	//cmd.Wait()
	return nil
}

func (c *Container) Remove() {
	process, _ := os.FindProcess(c.RealPid)
	process.Signal(syscall.Signal(9))
	fmt.Println(c.RealPid, " is killed")
}
