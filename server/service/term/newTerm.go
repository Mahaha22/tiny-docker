package service

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"tiny-docker/container"
	"tiny-docker/grpc/term"
)

type TermService struct {
	term.UnimplementedTermServer
}

func (t *TermService) Newterm(stream term.Term_NewtermServer) error {
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
		go recvProcess(stdinPipe, stream)  //输入流
		go sendProcess(stdoutPipe, stream) //输出流
		go errProcess(stderrPipe, stream)  //错误流
		//wg.Wait()
		cmd.Wait()
		return nil

	} else {
		return fmt.Errorf("\033[31mError:\033[0m No such container: %v", containerId.Input)
	}
}

func recvProcess(stdin io.WriteCloser, stream term.Term_NewtermServer) {
	defer func() {
		if err := recover(); err != nil {
			//异常处理,保证服务器不会因为panic而退出
			fmt.Println("客户端退出")
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
			fmt.Println("客户端退出")
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
			fmt.Println("客户端退出")
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
