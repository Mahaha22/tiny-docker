package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 返回客户端一个服务器连接conn
func GetConn() (*grpc.ClientConn, error) {
	//1.连接服务器,带证书访问
	conn, err := grpc.Dial(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	return conn, err
	// defer conn.Close()
	// //2.调用Product.pb.go中的NewProdService方法
	// productServiceClient := service.NewProdServiceClient(conn)

	// //3.直接像调用本地方法一样调用GetProductStock方法
	// resp, err := productServiceClient.GetProductStock(context.Background(), &service.ProductRequest{ProdId: 233})
	// if err != nil {
	// 	log.Fatal("调用失败:", err)
	// }
	// fmt.Println("调用grpc方法成功,ProdStock=", resp.ProdStock)
}
