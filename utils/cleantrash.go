package utils

import (
	"fmt"
	"os/exec"
	"syscall"
)

func Cleantrash() {
	//1.清理挂载残余
	for i := 0; i <= 20; i++ {
		syscall.Unmount("tiny-docker", 0)
	}
	//2.清理容器根目录残余
	args := "rm -rf /mnt/tiny-docker/*"
	cmd := exec.Command("/bin/bash", "-c", args)
	if err := cmd.Run(); err != nil {
		fmt.Println("清理容器根目录残余失败，请手动清理")
	}
}
