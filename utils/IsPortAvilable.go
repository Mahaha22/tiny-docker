// 判断某个端口是否被占用
package utils

import (
	"fmt"
	"net"
)

func IsPortAvilable(port string) error {
	host := "localhost"
	// 尝试在指定的主机和端口上监听连接
	_, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		//fmt.Printf("Port %d is already in use\n", port)
		return err
	}
	return nil
}
