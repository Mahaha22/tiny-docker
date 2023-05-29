package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/term"
	"tiny-docker/overlayfs"
	"tiny-docker/server/service"
	TermSvc "tiny-docker/server/service/term"

	"google.golang.org/grpc"
)

func init() {
	container.Global_ContainerMap_Init()
	go container.Monitor()
}
func main() {
	//1.前置处理
	//注册一些信号，可以让程序优雅的停止
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//服务器退出后的清理工作
	go func() {
		<-sigChan
		for _, v := range container.Global_ContainerMap { //1.清除并卸载所有挂载点
			if err := overlayfs.DeleteOverlayMnt(v.ContainerId); err != nil {
				fmt.Println("DeleteOverlayMnt err = ", err)
			}
		}
		container.KillallContainer() //2.服务器退出信号时，清除所有容器
		os.Exit(0)
	}()

	//2.服务器启动
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("listener fail :", err)
	}

	rpcServer := grpc.NewServer()
	//注册服务
	cmdline.RegisterServiceServer(rpcServer, &service.ContainerService{})
	term.RegisterTermServer(rpcServer, &TermSvc.TermService{})

	//启动服务
	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Server fail :", err)
	}

}
