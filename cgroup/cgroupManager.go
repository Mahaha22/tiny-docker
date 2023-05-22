package cgroup

import (
	"tiny-docker/conf"
)

type Subsystem interface {
	Init(container_name string, res *conf.Cgroupflag) error //初始化一个资源子系统
	Apply(container_name string, pid int) error
	Delete(container_name string) error
}

var (
	subsystems = []Subsystem{
		&CpuSubsystem{},
		&MemSubsystem{},
	}
)

// 为每一个资源创建以自己容器名字命名的cgroup子系统
func SetCgroup(name string, res *conf.Cgroupflag, pid int) error {
	for _, subsystem := range subsystems {
		err := subsystem.Init(name, res)
		if err != nil {
			return err
		}
		err = subsystem.Apply(name, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// 当需要进入容器内部的时候，这个新建的终端也需要加入容器的cgroup
func ApplyCgroup(name string, pid int) error {
	for _, subsystem := range subsystems {
		err := subsystem.Apply(name, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// 销毁容器的cgroup子系统
func DestroyCgroup(name string) error {
	for _, subsystem := range subsystems {
		err := subsystem.Delete(name)
		if err != nil {
			return err
		}
	}
	return nil
}
