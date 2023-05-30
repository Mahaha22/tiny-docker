package container

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"tiny-docker/cgroup"
	"tiny-docker/conf"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/overlayfs"
	"tiny-docker/utils"
)

var Global_ContainerMap map[string]*Container //索引所有容器
func Global_ContainerMap_Init() {
	Global_ContainerMap = make(map[string]*Container)
}

type state int //容器的运行状态
const (
	RUNNING state = iota
	STOPPED
	UNKNOWN
	EXITDED
)

// 容器的数据结构
type Container struct {
	ContainerId  string               //容器的唯一标识
	Name         string               //容器的自定义名称
	CgroupRes    conf.Cgroupflag      //cgroup资源
	NameSpaceRes *syscall.SysProcAttr //Namespace隔离标准当作克隆标志
	RealPid      int                  //容器在主机上真实的进程id
	Volmnt       map[string]string    //挂载卷映射 [host]->container
	CreateTime   string               //创建时间
	Status       state                //容器的运行状态
	Image        string               //容器镜像
	//Ports
	//NetWork
}

// 实例化一个容器
func CreateContainer(req *cmdline.Request) *Container {
	c := &Container{
		Volmnt: make(map[string]string),
	}
	//解析出卷挂载的参数
	for _, mntPair := range req.Args.Volmnt {
		//mntPair = hostdir:containerdir
		paths := strings.Split(mntPair, ":")
		c.Volmnt[paths[0]] = paths[1]
	}

	c.ContainerId = utils.GetUniqueId()       //生成容器id
	c.Name = req.Args.Name                    //容器名
	c.Image = req.Args.ImageId                //容器镜像
	c.Status = RUNNING                        //容器状态
	c.CreateTime = time.Now().Format("15:04") //创建时间
	c.CgroupRes = conf.Cgroupflag{            //容器资源限制
		Cpulimit: req.Args.Cpu,
		Memlimit: req.Args.Mem,
	}
	c.NameSpaceRes = conf.Cloneflag //容器命名空间隔离
	return c
}

// 容器相关配置初始化，例如namesapce和cgroup
func (c *Container) Init() error {
	//1.加载容器镜像挂载overlay存储
	image_path := "/root/busybox"
	err := overlayfs.MountOverlay(image_path, c.ContainerId)
	if err != nil {
		fmt.Println("overlayfs mount err = ", err)
		return err
	}
	//2.挂载容器卷
	err = overlayfs.MountFS(c.Volmnt, c.ContainerId)
	if err != nil {
		fmt.Println("volume mount err = ", err)
		return err
	}

	//3.从挂载好的overlay存储中启动容器
	os.Chdir("/root/go/tiny-docker/container/lib") //钩子程序的位置
	cmd := exec.Command("./task", c.ContainerId)   //启动容器的钩子
	r, w, _ := os.Pipe()                           //用于跟这个钩子程序通信
	cmd.ExtraFiles = []*os.File{w}                 //将管道的一端传递给钩子
	if err := cmd.Start(); err != nil {
		fmt.Println("cmd err = ", err)
		return err
	}

	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	r.Close()
	pid, _ := strconv.Atoi(string(buf[:n-1])) //钩子程序用于容器的启动，并从管道的另一端返回容器启动的真实pid
	c.RealPid = pid

	//4.创建cgroup资源
	//容器内的第一个程序已经启动，针对这个容器以他的唯一标识container_id在/sys/fs/cgroup/<subsystem>/tiny-docker/<container_id>创建新的subsystem
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
