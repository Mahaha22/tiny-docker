package service

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
)

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
