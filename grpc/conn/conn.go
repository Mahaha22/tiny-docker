package conn

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

// 返回客户端一个服务器连接conn
func getConn() (*grpc.ClientConn, error) {
	var err error
	if conn == nil {
		conn, err = grpc.Dial(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
