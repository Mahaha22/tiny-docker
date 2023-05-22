package container

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func TestContainerIsAlive(t *testing.T) {

	ret := ContainerIsAlive(&Container{RealPid: 36597})
	fmt.Println("ret", ret)
}

func TestMonitor(t *testing.T) {
	Global_ContainerMap = make(map[string]*Container)
	Global_ContainerMap["cc896bce"] = &Container{
		RealPid: 1122,
	}
	Monitor()
}

func TestRemove(t *testing.T) {
	process, _ := os.FindProcess(29816)
	process.Signal(syscall.Signal(9))
}
