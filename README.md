# Tiny-Docker
   Tiny Docker是一个使用Golang语言实现的精简版Docker项目，旨在模仿runC实现容器管理的基本功能。该项目采用了CS架构，客户端和服务器使用GRPC框架进行交互。可以实现高效的容器远程管理。

# 工程项目详解
```shell
├── client                        //客户端
│   ├── cli_command.go            //客户端调用的功能
│   ├── client.go                 //客户端入口
│   └── test
├── cli.go
├── cmd                           //命令 例如./tiny-docker create / run / ps 等各种功能的具体实现
│   └── RunCommand.go             //run的实现
├── go.mod
├── go.sum
├── grpc                          //protoc生成的client-server之间信息交互的文件
│   ├── cmdline                   //普通命令行的格式 例如 ./tiny-docker run -it -cpu 10% -mem 100m bash
│   │   ├── cmdline_grpc.pb.go
│   │   └── cmdline.pb.go
│   ├── conn                      //用于grpc连接的建立
│   │   ├── conn.go
│   │   └── grpc_client.go
│   └── term                      //用于与容器终端交互的数据格式
│       ├── term_grpc.pb.go
│       └── term.pb.go
├── protof-src                    //protobuf文件，生成go-grpc格式
│   ├── cmdline.proto             //命令行数据的定义
│   └── term.proto                //终端数据交互的定义
└── server                        //服务端
    ├── serve.go                  //服务端入口
    └── service                   //可供使用的服务
        └── runContainer.go       //见名知意
```
