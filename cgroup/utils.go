package cgroup

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

const ContainerCgroupName = "tiny-docker"

func findCgroupMountPoint(Subsystem string) (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", fmt.Errorf("findCgroupMountPoint open err:%v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f) //创建一个读取文件的扫描器，逐行读取信息
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, " ")
		for _, v := range strings.Split(fields[len(fields)-1], ",") {
			if v == Subsystem {
				return fields[4], nil //可以使用命令 cat /proc/self/mountinfo | grep cpu 分析这个函数的作业
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("findCgroupMountPoint scanner err:%v", err)
	}
	return "", nil
}

func GetSubPath(container_name, subsystem string) (string, error) {
	//也可以简单直接的认为
	//CgropRootPoint = /sys/fs/cgroup/<subsystem>
	CgropRootPoint, err := findCgroupMountPoint(subsystem)
	if err != nil {
		return "", fmt.Errorf("findCgroupMountPoint errors : %v", err)
	} else if CgropRootPoint == "" {
		return "", fmt.Errorf("findCgroupMountPoint is empty")
	}
	folderPath := path.Join(CgropRootPoint, ContainerCgroupName, container_name)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		//判断 例如/sys/fs/cgroup/cpu/tiny-docker/cpu 文件夹是否存在
		//文件夹不存在
		err = os.Mkdir(folderPath, 0755)
		if err != nil {
			return "", err
		}
	}
	return folderPath, nil
}
