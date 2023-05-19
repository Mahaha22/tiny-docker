package conn

import (
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/term"
)

// 创建grpc客户端,一元grpc
func GrpcClient_Single() (cmdline.ServiceClient, error) {
	conn, err := getConn()
	//fmt.Println(conn)
	if err != nil {
		return nil, err
		//log.Fatal("RunCommand 连接服务器失败 : ", err)
	}
	return cmdline.NewServiceClient(conn), nil
}

// 双向grpc
func GrpcClient_Double() (term.TermClient, error) {
	conn, err := getConn()
	if err != nil {
		return nil, err
		//log.Fatal("RunCommand 连接服务器失败 : ", err)
	}

	return term.NewTermClient(conn), nil
}
