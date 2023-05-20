package service

import "tiny-docker/grpc/cmdline"

type ContainerService struct {
	cmdline.UnimplementedServiceServer
}
