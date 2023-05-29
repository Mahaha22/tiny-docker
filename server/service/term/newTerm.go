package service

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"syscall"
	"tiny-docker/cgroup"
	"tiny-docker/container"
	"tiny-docker/grpc/term"
)

type TermService struct {
	term.UnimplementedTermServer
}

var termExitch chan struct{}

func (t *TermService) Newterm(stream term.Term_NewtermServer) error {
	termExitch = make(chan struct{}, 1)
	//fmt.Println("this")
	defer func() {
		if err := recover(); err != nil {
			//异常处理,保证服务器不会因为panic而退出
			fmt.Println("panic err: ", err)
		}
	}()

	containerId, err := stream.Recv() //第一条发过来的是新建终端需要执行的命令
	if err != nil {
		return err
	}
	//fmt.Println("cli  = ", containerId)
	if container, ok := container.Global_ContainerMap[containerId.Input]; ok { //如果存在此容器
		pid := container.RealPid
		cmd := exec.Command("/bin/bash", "-c", "nsenter --all -t "+strconv.Itoa(pid))
		stdinPipe, _ := cmd.StdinPipe()
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrPipe, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("\033[31mError:\033[0m New term fails %v", err)
		}

		// var wg sync.WaitGroup
		// wg.Add(3)
		//chrootinfo := "chroot /root/busybox /bin/sh\n" //设置容器根目录
		chrootinfo := "chroot /mnt/tiny-docker/" + container.ContainerId + "/merge /bin/sh\n"
		pathinfo := "export PATH=:/bin\n" //设置环境变量
		stdinPipe.Write([]byte(chrootinfo))
		stdinPipe.Write([]byte(pathinfo))

		go recvProcess(stdinPipe, stream)  //输入流
		go sendProcess(stdoutPipe, stream) //输出流
		go errProcess(stderrPipe, stream)  //错误流
		//wg.Wait()
		//新创建的终端需要加入容器的cgroup以受到资源的限制
		cgroup.ApplyCgroup(container.ContainerId, cmd.Process.Pid)

		//在终端内部可以通过exit主动退出，但是当客户端意外退出时，也应具备关闭远程终端的能力，避免浪费资源
		go Killcmd(cmd.Process.Pid, stdinPipe)
		cmd.Wait()
		return nil

	} else {
		return fmt.Errorf("\033[31mError:\033[0m No such container: %v", containerId.Input)
	}
}

func Killcmd(pid int, stdin io.WriteCloser) {
	<-termExitch
	syscall.Kill(pid, 9)
}

func recvProcess(stdin io.WriteCloser, stream term.Term_NewtermServer) {
	defer func() {
		if err := recover(); err != nil {
			//异常处理,保证服务器不会因为panic而退出
			//fmt.Println("客户端退出")
			//stdin.Write([]byte("exit\n"))
			termExitch <- struct{}{}
		}
	}()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if _, err := stdin.Write([]byte(req.Input + "\n")); err != nil {
			fmt.Println("recvProvess fail")
		}
	}
	stdin.Close()
}
func sendProcess(stdout io.ReadCloser, stream term.Term_NewtermServer) {
	defer func() {
		if err := recover(); err != nil {
			//异常处理,保证服务器不会因为panic而退出
			//fmt.Println("客户端退出")
			termExitch <- struct{}{}
		}
	}()
	for {
		buf := make([]byte, 1024)
		n, err := stdout.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		if n > 0 {
			if err := stream.Send(&term.Response{Output: string(buf[:n])}); err != nil {
				fmt.Println("sendProcess fail")
			}
		}
	}
	stdout.Close()
}
func errProcess(stderr io.ReadCloser, stream term.Term_NewtermServer) {
	defer func() {
		if err := recover(); err != nil {
			//异常处理,保证服务器不会因为panic而退出
			//fmt.Println("客户端退出")
			termExitch <- struct{}{}
		}
	}()
	for {
		buf := make([]byte, 1024)
		n, err := stderr.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		if n > 0 {
			if err := stream.Send(&term.Response{Output: string(buf[:n])}); err != nil {
				fmt.Println("errProcess fail")
			}
		}
	}
	stderr.Close()
}
