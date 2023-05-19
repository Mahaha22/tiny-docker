package conn

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 返回客户端一个服务器连接conn
func getConn() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("116.62.227.94:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
