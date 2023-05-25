package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

var registeredInitializers = make(map[string]func())

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	}
	registeredInitializers[name] = initializer
}
func Init() bool {
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()
		return true
	}
	return false
}

func init() {
	fmt.Println("****:", os.Args[0])
	log.Printf("init start, os.Args = %+v\n", os.Args)
	Register("childProcess", childProcess)
	if Init() {
		os.Exit(0)
	}
}
func childProcess() {
	pipe := os.NewFile(uintptr(3), "pipe")
	pid := os.Getpid()
	pipe.WriteString(strconv.Itoa(pid))
	setUpMount()
	_, err := exec.LookPath("/bin/sh")
	fmt.Println("son = ", os.Getpid())
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	if err := syscall.Exec("/bin/sh", []string{"sleep", "300"}, os.Environ()); err != nil {
		fmt.Println("err = ", err)
		return
	}
}
func Self() string {
	return "/proc/self/exe"
}
func Command(args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: Self(),
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS |
				syscall.CLONE_NEWNET |
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
		},
	}
}
func setUpMount() {
	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	fmt.Println("mount 1 = ", err)
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	fmt.Println("mount 2 = ", err)

	if err := syscall.Chroot("/root/busybox"); err != nil {
		fmt.Println("Chroot 失败：", err)
		return
	}
	fmt.Println("成功设置容器根目录为 /root/busybox")
	if err := syscall.Chdir("/"); err != nil {
		fmt.Println("Chdir 失败：", err)
		return
	}

}
func main() {
	log.Printf("main start, os.Args = %+v\n", os.Args)
	cmd := Command("childProcess")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("fatherProcess,pid = ", os.Getpid())
	if err := cmd.Start(); err != nil {
		log.Panicf("failed to run command: %s", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Panicf("failed to wait command: %s", err)
	}
	log.Println("main exit")
}
