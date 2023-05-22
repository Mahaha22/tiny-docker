package cgroup

import (
	"testing"
	"tiny-docker/conf"
)

func TestCpu(t *testing.T) {
	Cpu := CpuSubsystem{}
	res := conf.Cgroupflag{
		Cpulimit: "10000",
		Memlimit: "100m",
	}
	Cpu.Init("container", res)
	Cpu.Apply("container", 38406)
	Cpu.Delete("container")
}

func TestMem(t *testing.T) {
	Mem := MemSubsystem{}
	res := conf.Cgroupflag{
		Cpulimit: "10000",
		Memlimit: "100m",
	}
	Mem.Init("container", res)
	Mem.Apply("container", 40068)
	Mem.Delete("container")
}
