package main

import (
	"context"
	"fmt"
	"log"
	"tiny-docker/grpc/cmdline"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := cmdline.NewServiceClient(conn)
	fmt.Printf("%v", client)
	res, err := client.RunContainer(context.Background(), &cmdline.Request{
		Args: &cmdline.Flag{
			Cpu: "100",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
