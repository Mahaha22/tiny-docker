# Tiny-Docker
   Tiny Docker是一个使用Golang语言实现的精简版Docker项目，旨在模仿runC实现容器管理的基本功能。该项目采用了CS架构，客户端和服务器使用GRPC框架进行交互。可以实现高效的容器远程管理。  
   ![结构图](./assets/tiny-docker.png)

# 核心功能实现简要讲解
## 容器
- [启动容器](#启动容器) 
    - [xxx](#xxx)
    - [xxx](#xxx)
- [进入容器](#进入容器) 
    - [xxx](#xxx)
    - [xxx](#xxx)
- [执行容器命令](#执行容器命令)
    - [xxx](#xxx)
    - [xxx](#xxx)
- [查看运行中的容器](#查看运行中的容器)
    - [xxx](#xxx)
    - [xxx](#xxx)
- [删除容器](#删除容器)
    - [xxx](#xxx)
## 资源限制
- [cpu资源限制](#cpu资源限制)
- [memory资源限制](#mem资源限制)
## 容器卷挂载
- [aufs文件系统](#aufs)
- [overlay文件系统](#overlay)
## 网络配置
- [网络配置基本原理](#网络配置基本原理)
  - [bridge](#bridge)
- [网络管理](#创建网络)
  - [create](#nw_create)
  - [list](#nw_list)
  - [delete](#nw_delete)

## 网络配置 <a id="网络配置基本原理"></a>
### bridge网络<a id="bridge"></a> 
```shell
$ #新建网络命名空间n1、n2 和网桥 br0
$ ip netns add n1 
$ ip netns add n2 
$ brctl addbr br0
$ #新建两对虚拟网卡
$ ip link add veth1_br type veth peer name veth1_c
$ ip link add veth2_br type veth peer name veth2_c
$ #将网卡的一端放入网络命名空间n1、n2
$ ip link set veth1_c netns n1     
$ ip link set veth2_c netns n2  
$ #将网卡的一端放入网桥br0
$ ip link set veth1_br master br0
$ ip link set veth2_br master br0
$ #给网桥设置ip地址和广播地址，并启动up
$ ip addr add 192.168.0.1/16 brd + dev br0 
$ ip link set br0 up
$ ip netns exec n1 ip addr add 192.168.0.2/16 dev veth1_c
$ #给放入命名空间n1、n2的网卡设置ip地址，并启动up
$ ip netns exec n1 ip addr add 192.168.0.2/16 brd + dev veth1_c
$ ip netns exec n1 ip link set veth1_c up
$ ip netns exec n2 ip addr add 192.168.0.3/16 brd + dev veth2_c
$ ip netns exec n2 ip link set veth2_c up
$ #启动虚拟网卡的另一端
$ ip link set veth1_br up
$ ip link set veth2_br up

$ #以上设置完成以后，在192.168.0.0/16这个网段中的 n1 n2 br0都是可以互相通信的
$ ip netns exec n1 ping -c 1 192.168.0.3
PING 192.168.0.3 (192.168.0.3) 56(84) bytes of data.
64 bytes from 192.168.0.3: icmp_seq=1 ttl=64 time=0.185 ms

--- 192.168.0.3 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.185/0.185/0.185/0.000 ms
$ ip netns exec n2 ping -c 1 192.168.0.2
PING 192.168.0.2 (192.168.0.2) 56(84) bytes of data.
64 bytes from 192.168.0.2: icmp_seq=1 ttl=64 time=0.081 ms

--- 192.168.0.2 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 0ms
rtt min/avg/max/mdev = 0.081/0.081/0.081/0.000 ms
```
> 注意，此时n1、n2网络命名空间是不可以上外网，外网也无法访问n1、n2;  
> 需要访问外网和提供服务需要配置路由表和iptables流量转发  

**实现访问外网**  
原理是iptables中nat表中SNAT功能，SNAT 可以在出站数据包经过网络设备时修改其源IP地址，使其看起来像来自于网络设备的另一个IP地址。这对于将内部私有 IP地址转换为公共IP地址非常有用，以便与外部网络通信。  
`iptables -t nat -A POSTROUTING -s <源 IP 地址范围> -j SNAT --to-source <目标 IP 地址>`  
如果我们不能确定需要转换的公网地址可以简单一点让iptables动态的去检测
`iptables -t nat -A POSTROUTING -s <源IP地址> -j MASQUERADE`  

```shell
$ #进入网络命名空间内部设置默认网关
$ ip netns exec n1 ip route add default via 192.168.0.1 dev veth1_c
$ ip netns exec n2 ip route add default via 192.168.0.1 dev veth2_c
$ ip netns exec n1 route -n
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         192.168.0.1     0.0.0.0         UG    0      0        0 veth1_c
192.168.0.0     0.0.0.0         255.255.0.0     U     0      0        0 veth1_c

$ #现在在n1、n2网络中是无法访问www.baidu.com
$ #需要设置SNAT转发
$ iptables -t nat -A POSTROUTING -s 192.168.0.0/16 -j MASQUERADE
$ ip netns exec n1 ping www.baidu.com
PING www.a.shifen.com (180.101.50.188) 56(84) bytes of data.
64 bytes from 180.101.50.188 (180.101.50.188): icmp_seq=1 ttl=49 time=9.60 ms
64 bytes from 180.101.50.188 (180.101.50.188): icmp_seq=2 ttl=49 time=9.57 ms
64 bytes from 180.101.50.188 (180.101.50.188): icmp_seq=3 ttl=49 time=9.62 ms
^C
--- www.a.shifen.com ping statistics ---
3 packets transmitted, 3 received, 0% packet loss, time 5ms
rtt min/avg/max/mdev = 9.566/9.593/9.616/0.020 ms
```
> 注意，当新启动一个容器镜像的时候虽然所有的网络配置都正确，但是也不一定可以访问www.baidu.com,这是因为容器内部的域名解析配置不正确，如果出现这种情况请检查容器的/etc/resolve.conf。最简单的做法是直接复制主机的/etc/resolve.conf













## **1.容器根文件系统挂载实现**
<a id="section1"></a>
以下对功能实现的原理做举例说明

1、获得一个redis的完整工作目录
```shell
[root@haohao ~]# docker run -d redis:latest sleep 300
9c664a921a275e5d1b31fc395ef3d4f9734b06a75cef331aeda4cca073052cb6
[root@haohao ~]# docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED         STATUS         PORTS      NAMES
9c664a921a27   redis:latest   "docker-entrypoint.s…"   4 seconds ago   Up 2 seconds   6379/tcp   friendly_mcclintock
[root@haohao ~]# docker export -o redis.tar 9c664a921a27
[root@haohao ~]# mkdir redis
[root@haohao ~]# tar -xvf redis.tar -C ~/redis
[root@haohao ~]# cd redis/
[root@haohao redis]# ls
bin  boot  data  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
[root@haohao redis]# pwd
/root/redis
```
2、使用overlay挂载方式挂载redis工作目录
```shell
[root@haohao redis]# mkdir -p /mnt/redis/{merge,work,upper}
[root@haohao redis]# cd /mnt/redis/
[root@haohao redis]# tree .
        .
        ├── merge //合并层
        ├── upper //读写层
        └── work  //工作层
[root@haohao redis]# mount -t overlay redis-overlay -o lowerdir=/root/redis,upperdir=/mnt/redis/upper,workdir=/mnt/redis/work merge
[root@haohao redis]# mount
redis-overlay on /mnt/redis/merge type overlay (rw,relatime,lowerdir=/root/redis,upperdir=/mnt/redis/upper,workdir=/mnt/redis/work)
#说明已经将redis工作目录作为只读层以overlay的方式挂载到了/mnt/redis/中
[root@haohao redis]# tree . -L 2
        .
        ├── merge
        │   ├── bin
        │   ├── boot
        │   ├── data
        │   ├── dev
        │   ├── etc
        │   ├── home
        │   ├── lib
        │   ├── lib64
        │   ├── media
        │   ├── mnt
        │   ├── opt
        │   ├── proc
        │   ├── root
        │   ├── run
        │   ├── sbin
        │   ├── srv
        │   ├── sys
        │   ├── tmp
        │   ├── usr
        │   └── var
        ├── upper
        └── work
            └── work
24 directories, 0 files
```
将`/mnt/redis/merge`作为容器的根文件系统目录`/`，从`/`这里启动容器，既可以做到不影响容器镜像的前提下进行读写
具体实现见`container`、`overlayfs`包

## **2.容器卷挂载的实现**
容器的根文件系统挂载在`/mnt/redis/merge`,那么在容器启动之前，根文件系统挂载完成以后，我们可以先把需要挂载的卷挂载到指定目录
例如我们需要把`/root/vol`挂载到容器中的`/root/vol`
核心实现是直接挂载
`mount --bind src target`

`mount --bind` 命令是 Linux 系统中的一个命令，它的作用是将一个目录挂载到另一个目录上，从而让这两个目录共享同一个文件系统。这个过程也被称为“绑定挂载”（bind mount）。

具体来说，`mount --bind` 命令可以将一个目录（源目录）挂载到另一个目录（目标目录）上，使得目标目录中的文件和源目录中的文件指向同一块磁盘空间。这样，如果你在目标目录中创建、修改或删除文件，相应的操作也会反映在源目录中。

绑定挂载通常被用来实现以下功能：

- 将一个目录的内容复制到另一个目录中；
- 将一个目录挂载到另一个目录上，并以只读方式访问；
- 将一个目录挂载到另一个目录上，并以可写方式访问。

需要注意的是，`mount --bind` 命令只能在 Linux 系统中使用，并且需要在 root 权限下执行。

`举例`
```shell
[root@haohao vol]# pwd
/root/vol
[root@haohao vol]# tree .
        .
        ├── a.txt
        ├── b.txt
        ├── test1
        │   ├── a.txt
        │   └── ok.txt
        └── test2
[root@haohao vol]# mount --bind /root/vol /mnt/redis/merge/root/vol
mount: /mnt/redis/merge/root/vol: mount point does not exist.
[root@haohao vol]# mkdir /mnt/redis/merge/root/vol
[root@haohao vol]# mount --bind /root/vol /mnt/redis/merge/root/vol
[root@haohao root]# pwd
/mnt/redis/merge/root
[root@haohao root]# tree .
        .
        └── vol
            ├── a.txt
            ├── b.txt
            ├── test1
            │   ├── a.txt
            │   └── ok.txt
            └── test2
```
## **3.exec执行容器命令实现**
一个核心的命令就是`nsenter --all -t <pid> <cmd>`  
举例`nsenter --all -t 12233 ls`即可进入pid为12233的程序的所有命名空间执行ls命令  

***主要流程***  
1.客户端执行命令`tiny-docker exec <containerId> <cmd>`[cmd/ExecCommand.go]  
2.向服务器发起`ExecContainer()远程调用`[server/execContainer.go]  
3.ExecContainer()会根据请求的容器id和需要执行的命令调用`nsenter --all -t <pid> <cmd>`并重定向此程序的`STDOUT、STDERR`并返回给客户端。