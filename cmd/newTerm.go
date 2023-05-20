package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"tiny-docker/grpc/conn"
	"tiny-docker/grpc/term"
)

func newTerm(containerId string) {
	fmt.Println("new term")
	client, err := conn.GrpcClient_Double()
	if err != nil {
		log.Fatal("newTerm client error : ", err)
	}
	stream, err := client.Newterm(context.Background())
	if err != nil {
		log.Fatal("newTerm stream error : ", err)
	}
	//接收服务器的信息
	go recv(stream)
	//首先把需要开启交互的容器的id传过去建立交互连接
	if err := stream.Send(&term.Request{Input: containerId}); err != nil {
		log.Fatalf("failed to send containerId: %v", err)
	}
	//从终端读取用户输入，并将其发送到Shell服务中执行
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\033[32m[tiny-docker]# \033[0m") //此处将会显示为绿色
	for {
		if !scanner.Scan() {
			break
		}
		cmd := scanner.Text()
		if err := stream.Send(&term.Request{Input: cmd}); err != nil {
			log.Fatalf("failed to send command: %v", err)
		}
		if cmd == "exit" {
			break
		}
	}
	//time.Sleep(time.Second * 1)
}

func recv(stream term.Term_NewtermClient) {
	for {
		out, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}
		fmt.Print(out.GetOutput())
		fmt.Print("\033[32m[tiny-docker]# \033[0m")
	}
}
