package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

func child() {
	fmt.Println("[exe]", "pid:", os.Getpid())
	//err := syscall.Sethostname([]byte("mycontainer"))
	// if err != nil {
	// 	fmt.Println("sethostname err = ", err)
	// 	return
	// }
	//os.Chdir("/root/busybox")
	overlayfs_path := path.Join("/mnt/tiny-docker", os.Args[1], "merge")
	fmt.Println("overlaypath = ", overlayfs_path)
	os.Chdir(overlayfs_path)
	err := syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		fmt.Println("mount err = ", err)
		return
	}
	err = syscall.Chroot(".")
	if err != nil {
		fmt.Println("chroot err = ", err)
		return
	}
	//cmd := exec.Command("/bin/sh", "-c", "/bin/sleep 3000") //模拟容器中启动一个服务
	cmd := exec.Command("/bin/sh", "-c", os.Args[2])
	cmd.Env = append(cmd.Env, ":/bin") //添加环境变量
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return
	}
	cmd.Wait()
	syscall.Unmount("proc", 0)
}

func main() {
	if os.Args[0] == "/proc/self/exe" {
		//子进程用于真正的启动一个容器
		child()
	} else {
		//父进程用于创建隔离的命名空间
		fmt.Println("[main]", "pid:", os.Getpid())
		cmd := exec.Command("/proc/self/exe", os.Args[1], os.Args[2])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWIPC |
				syscall.CLONE_NEWUSER |
				syscall.CLONE_NEWNS |
				syscall.CLONE_NEWPID,
			UidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: 0, Size: 1},
			},
			GidMappings: []syscall.SysProcIDMap{
				{ContainerID: 0, HostID: 0, Size: 1},
			},
			Unshareflags: syscall.CLONE_NEWNS,
		}
		if os.Args[3] != "host" { //如是果非host模式
			cmd.SysProcAttr.Cloneflags |= syscall.CLONE_NEWNET
		}
		if err := cmd.Start(); err != nil {
			return
		}
		w := os.NewFile(uintptr(3), "pipe")
		fmt.Fprintln(w, cmd.Process.Pid)
		w.Close()
		cmd.Wait()
	}
}
