package service

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *ContainerService) RunContainer(ctx context.Context, req *cmdline.Request) (*cmdline.RunResponse, error) {
	//实现具体的业务逻辑
	newContainer := container.CreateContainer(req) //实例化一个容器
	if err := newContainer.Init(); err != nil {    //初始化容器
		return &cmdline.RunResponse{
			ContainerId: "", //容器创建失败返回空
		}, err
	}
	container.Global_ContainerMap[newContainer.ContainerId] = newContainer //放入Global_ContainerMap
	return &cmdline.RunResponse{
		ContainerId: newContainer.ContainerId, //容器创建成功返回真实的容器id
	}, nil
}

func (r *ContainerService) PsContainer(context.Context, *emptypb.Empty) (*cmdline.ContainerInfo, error) {
	info := &cmdline.ContainerInfo{}
	for _, c := range container.Global_ContainerMap {
		tmp := &cmdline.Container{
			ContainerId: c.ContainerId,
			Image:       c.Image,
			CreateTime:  c.CreateTime,
			Status:      fmt.Sprintf("%v", c.Status),
			//Ports:      ,
			Name:    c.Name,
			Command: c.Command,
		}
		for vol1, vol2 := range c.Volmnt {
			tmp.VolumeMount = append(tmp.VolumeMount, vol1, vol2)
		}
		info.Containers = append(info.Containers, tmp)
	}
	return info, nil
}

func (r *ContainerService) ExecContainer(ctx context.Context, req *cmdline.Request) (*cmdline.ContainerStdout, error) {
	containerid := req.Cmd[0]
	cmd_agrs := ""
	for _, v := range req.Cmd[1:] {
		cmd_agrs += v + " "
	}
	cmd_agrs += "\n"

	//判断容器是否存在
	if container, ok := container.Global_ContainerMap[containerid]; ok {
		pid := container.RealPid
		cmd := exec.Command("/bin/bash", "-c", "nsenter --all -t "+strconv.Itoa(pid))
		stdinPipe, _ := cmd.StdinPipe()
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrPipe, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			return nil, err
		}

		//chrootinfo := "chroot /root/busybox /bin/sh\n" //设置容器根目录
		chrootinfo := "chroot /mnt/tiny-docker/" + container.ContainerId + "/merge /bin/sh\n"
		pathinfo := "export PATH=:/bin\n"   //设置环境变量
		stdinPipe.Write([]byte(chrootinfo)) //根路径与容器保持一致
		stdinPipe.Write([]byte(pathinfo))   //环境变量根容器保持一致

		ch := make(chan struct{})
		buf := make([]byte, 1024)
		ret := &cmdline.ContainerStdout{}
		//接收stdout和stderr其一
		go func() {
			n, _ := stdoutPipe.Read(buf)
			ret.Outinfo = string(buf[:n])
			//fmt.Println("outinfo = ", ret.Outinfo)
			ch <- struct{}{}
		}()
		go func() {
			n, _ := stderrPipe.Read(buf)
			ret.Errinfo = string(buf[:n])
			//fmt.Println("errinfo = ", ret.Errinfo)
			ch <- struct{}{}
		}()
		//也有可能没有任何输出,1s后自动退出
		// go func() {
		// 	time.Sleep(time.Second)
		// 	ch <- struct{}{}
		// }()
		stdinPipe.Write([]byte(cmd_agrs)) //执行命令
		<-ch                              //接收到stdout或stderr的输出即可解除阻塞
		fmt.Println("ret = ", ret)
		return ret, nil
	} else {
		return nil, fmt.Errorf("container %v is not existed", containerid)
	}
}

func (r *ContainerService) KillContainer(ctx context.Context, req *cmdline.Request) (*cmdline.RunResponse, error) {
	for _, v := range req.Cmd {
		if c, ok := container.Global_ContainerMap[v]; ok { //容器存在
			container.KillVolume(c)
			container.KillContainer(c)
			return &cmdline.RunResponse{}, nil
		}
	}
	return &cmdline.RunResponse{}, nil
}
