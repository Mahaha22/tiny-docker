package main

import (
	"log"
	"net"
	"tiny-docker/container"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/term"
	"tiny-docker/server/service"
	TermSvc "tiny-docker/server/service/term"

	"google.golang.org/grpc"
)

func init() {
	container.Global_ContainerMap_Init()
}
func main() {
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
