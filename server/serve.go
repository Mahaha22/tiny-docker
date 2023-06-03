package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/term"
	"tiny-docker/network"
	"tiny-docker/server/service"
	TermSvc "tiny-docker/server/service/term"

	"google.golang.org/grpc"
)

func init() {
	container.Global_ContainerMap_Init() //容器map表初始化
	network.Global_Network_Init()        //网络map表初始化
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
		container.KillallVolume()    //2.清除所有容器卷和挂载
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
