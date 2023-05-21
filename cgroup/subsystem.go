package cgroup

import (
	"tiny-docker/conf"
)

type Subsystem interface {
	Init(res conf.Cgroupflag) error //初始化一个资源子系统
	Apply(path string, pid int) error
	Delete(path string) error
}
