# Tiny-Docker
   Tiny Docker是一个使用Golang语言实现的精简版Docker项目，旨在模仿runC实现容器管理的基本功能。该项目采用了CS架构，客户端和服务器使用GRPC框架进行交互。可以实现高效的容器远程管理。

# 工程项目详解
```shell
.
├── cgroup                          //Cgroup管理子模块
│   ├── cgroupManager.go           
│   ├── cgroup_test.go
│   ├── cpu.go
│   ├── memory.go
│   └── utils.go
├── client                          //客户端入口
│   ├── cli_command.go
│   ├── client.go
│   └── test
├── cmd                             //具体命令实现
│   ├── newTerm.go
│   └── RunCommand.go
├── conf                            //配置数据结构
│   ├── cgroupflag.go
│   └── cloneflag.go
├── container                       //容器管理子模块
│   ├── container.go
│   ├── container_test.go
│   └── utils_Container.go
├── go.mod
├── go.sum
├── grpc                            //grpc通信数据结构
│   ├── cmdline                     //命令行数据结构
│   │   ├── cmdline_grpc.pb.go      
│   │   └── cmdline.pb.go   
│   ├── conn                        //cli-sever连接
│   │   ├── conn.go
│   │   └── grpc_client.go
│   └── term                        //远程终端交互数据结构
│       ├── term_grpc.pb.go
│       └── term.pb.go 
├── protof-src                      //protobuuf源文件
│   ├── cmdline.proto
│   └── term.proto
├── README.md
├── server                          //服务器入口
│   ├── serve.go
│   └── service                     //具体服务
│       ├── containerService.go
│       ├── runContainer.go
│       └── term
│           └── newTerm.go
└── utils                           //项目用到的工具如随机哈希值
    ├── hash.go
    └── hash_test.go
```
