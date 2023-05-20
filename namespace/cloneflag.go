/*
这里是容器初始化时对于namespace的隔离
默认隔离IPC、UTS、NET、NEWNS、UTS
并且做了属组映射，让容器在自己的内部拥有root权限

我暂时还没有想要精细化的调整这个地方，所以先放着。
如果你感兴趣可以做如 tiny-docker run --uts --ipc yourdocker 指定隔离某些命名空间
*/
package namespace

import "syscall"

var Cloneflag = &syscall.SysProcAttr{
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
}
