package cgroup

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"tiny-docker/conf"
)

const cpu_name = "cpu"

type CpuSubsystem struct {
}

func (c *CpuSubsystem) Init(container_name string, res *conf.Cgroupflag) error {
	//subPath = /sys/fs/cgroup/cpu/tiny-docker/container_name
	SubPath, err := GetSubPath(container_name, cpu_name)
	if err != nil {
		return fmt.Errorf("GetSubPath errors : %v", err)
	}

	//往里面写具体的配置
	if res.Cpulimit != "" {
		if err := os.WriteFile(path.Join(SubPath, "cpu.cfs_quota_us"), []byte(res.Cpulimit), 0644); err != nil {
			return fmt.Errorf("limit cpu fail : %v", err)
		}
	}
	return nil
}

func (c *CpuSubsystem) Apply(container_name string, pid int) error {
	SubPath, err := GetSubPath(container_name, cpu_name)
	if err != nil {
		return fmt.Errorf("GetSubPath errors : %v", err)
	}
	if err := os.WriteFile(path.Join(SubPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("pid = %v join cgroup fail : %v", pid, err)
	}
	return nil
}
func (c *CpuSubsystem) Delete(container_name string) error {
	SubPath, err := GetSubPath(container_name, cpu_name)
	if err != nil {
		return fmt.Errorf("GetSubPath errors : %v", err)
	}
	cmd := exec.Command("rmdir", SubPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v cgroup delete fail : %v", container_name, err)
	}
	return nil
}
