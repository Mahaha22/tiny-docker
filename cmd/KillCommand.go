package cmd

import (
	"context"
	"fmt"
	"tiny-docker/grpc/cmdline"
	"tiny-docker/grpc/conn"
)

func KillCommand(containerIds []string) error {
	req := &cmdline.Request{
		Cmd: containerIds,
	}
	client, err := conn.GrpcClient_Single()
	if err != nil {
		return fmt.Errorf("grpc client create err = %v", err)
	}
	//rpc调用杀死容器
	_, err = client.KillContainer(context.Background(), req)
	if err != nil {
		return fmt.Errorf("kill container err = %v", err)
	}
	return nil
}
