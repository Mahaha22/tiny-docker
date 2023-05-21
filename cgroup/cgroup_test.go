package cgroup

import (
	"testing"
	"tiny-docker/conf"
)

func TestCpu(t *testing.T) {
	Cpu := CpuSubsystem{}
	res := conf.Cgroupflag{
		Cpulimit: "10000",
	}
	Cpu.Init("container", res)
	Cpu.Apply("container", 1000)
	Cpu.Delete("container")
}
