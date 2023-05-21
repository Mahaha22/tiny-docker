// cgroup资源限制的配置信息 暂定2个关键的指标
package conf

type Cgroupflag struct {
	Cpulimit string
	Memlimit string
}
