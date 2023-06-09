package container

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
	"tiny-docker/cgroup"
	"tiny-docker/conf"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/network"
	"tiny-docker/overlayfs"
	"tiny-docker/utils"
)

var Global_ContainerMap map[string]*Container //索引所有容器
func init() {
	Global_ContainerMap_Init() //容器map表初始化
	go Monitor()
}

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

type ContainerNet struct {
	Ip                 net.IP //容器的ip
	Port               map[string]string
	Network_Membership *network.Network //容器所属网络
}

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
	Command      string               //启动命令
	Net          ContainerNet         //容器网络
}

// 实例化一个容器
func CreateContainer(req *cmdline.Request) (*Container, error) {
	c := &Container{
		Volmnt: make(map[string]string),
	}
	//解析出卷挂载的参数
	for _, mntPair := range req.Args.Volmnt {
		//mntPair = hostdir:containerdir
		paths := strings.Split(mntPair, ":")
		c.Volmnt[paths[0]] = paths[1]
	}

	c.ContainerId = utils.GetUniqueId() //生成容器id
	if req.Args.Name == "" {            //容器名
		c.Name = c.ContainerId
	} else {
		c.Name = req.Args.Name
	}
	//要保证容器的名字唯一
	for _, v := range Global_ContainerMap {
		if c.Name == v.Name {
			return nil, fmt.Errorf("container %v existed", c.Name)
		}
	}
	c.Image = req.Args.ImageId                //容器镜像
	c.Status = RUNNING                        //容器状态
	c.CreateTime = time.Now().Format("15:04") //创建时间
	c.CgroupRes = conf.Cgroupflag{            //容器资源限制
		Cpulimit: req.Args.Cpu,
		Memlimit: req.Args.Mem,
	}
	c.NameSpaceRes = conf.Cloneflag //容器命名空间隔离
	for _, cmd := range req.Cmd {
		c.Command += cmd + " "
	}
	if req.Args.Net != "" { //创建容器时指定了网络
		nw, ok := network.Global_Network[req.Args.Net]
		if !ok {
			return nil, fmt.Errorf("network occour wrong")
		}
		c.Net.Network_Membership = nw
	} else { //如果不指定网络，那么就设置为默认Tiny-docker网络
		nw, ok := network.Global_Network["Tiny-docker"]
		if !ok {
			return nil, fmt.Errorf("network occour wrong")
		}
		c.Net.Network_Membership = nw
	}
	if req.Args.Ports != nil { //如果指定了端口映射
		c.Net.Port = make(map[string]string)
		for _, item := range req.Args.Ports { //[80:8080 1234:90]
			p := strings.Split(item, ":")
			c.Net.Port[p[0]] = p[1] //[80->8080],[1234->90]
		}
	}
	//c.Command += "\n"
	return c, nil
}

// 容器相关配置初始化，例如namesapce和cgroup
func (c *Container) Init() error {
	//1.加载容器镜像挂载overlay存储
	image_path := c.Image
	err := overlayfs.MountOverlay(image_path, c.ContainerId)
	if err != nil {
		fmt.Println("overlayfs mount err = ", err)
		return err
	}
	//2.挂载容器卷
	if c.Volmnt != nil {
		err = overlayfs.MountFS(c.Volmnt, c.ContainerId)
		if err != nil {
			fmt.Println("volume mount err = ", err)
			return err
		}
	}

	//3.从挂载好的overlay存储中启动容器
	os.Chdir("/root/go/tiny-docker/container/lib") //钩子程序的位置
	//如果容器指定的网络是host，那么就没有必要隔离网络命名空间,这里给容器启动钩子传递一个标记
	var netflag string
	if _, ok := c.Net.Network_Membership.Driver.(*network.HostDriver); ok {
		//如果是host驱动
		netflag = "host"
	}
	cmd := exec.Command("./task", c.ContainerId, c.Command, netflag) //启动容器的钩子+容器名+第一条启动命令
	r, w, _ := os.Pipe()                                             //用于跟这个钩子程序通信
	cmd.ExtraFiles = []*os.File{w}                                   //将管道的一端传递给钩子
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
		KillContainer(c)
		return fmt.Errorf("cgroup资源配置出错 :%v", err)
	}

	//5.配置网络资源
	c.Net.Ip, err = network.ApplyNetwork(c.RealPid, c.Net.Network_Membership, c.Net.Port)
	if err != nil {
		KillContainer(c)
		return fmt.Errorf("网络配置出错 :%v", err)
		//清除网络配置
	}
	//cmd.Wait()
	return nil
}

func (c *Container) Remove() {
	process, _ := os.FindProcess(c.RealPid)
	process.Signal(syscall.Signal(9))
	fmt.Println(c.RealPid, " is killed")
}
