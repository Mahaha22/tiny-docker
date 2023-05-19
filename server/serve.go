package main

import (
	"log"
	"net"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/server/service"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("listener fail :", err)
	}

	rpcServer := grpc.NewServer()

	//注册服务
	cmdline.RegisterServiceServer(rpcServer, &service.ContainerService{})

	//启动服务
	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Server fail :", err)
	}

}
