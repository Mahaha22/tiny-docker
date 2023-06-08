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
	"tiny-docker/network"
	"tiny-docker/server/service"
	TermSvc "tiny-docker/server/service/term"
	"tiny-docker/utils"

	"google.golang.org/grpc"
)

func init() {
	//创建服务启动需要的路径
	var filedir []string
	filedir = append(filedir, "/sys/fs/cgroup/cpu/tiny-docker")    //cpu subsystem
	filedir = append(filedir, "/sys/fs/cgroup/memory/tiny-docker") //mem subsystem
	filedir = append(filedir, "/mnt/tiny-docker/")                 //容器目录保存路径
	for _, v := range filedir {
		if err := utils.CreateDirectoryIfNotExists(v); err != nil {
			fmt.Println("路径初始化失败,tiny-docker退出")
			os.Exit(0)
		}
	}
}
func main() {
	//1.前置处理
	//注册一些信号，可以让程序优雅的停止
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	//服务器退出后的清理工作
	go func() {
		<-sigChan
		container.KillallVolume()    //1.清除所有容器卷和挂载
		container.KillallContainer() //2.服务器退出信号时，清除所有容器
		network.RemoveAllNetwork()   //3.删除所有网络配置
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
