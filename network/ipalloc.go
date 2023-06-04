package network

import "net"

type IPAlloc struct { //保存子网中以分配或者未分配的ip 例如192.168.0.0/24这个子网中以分配或者未分配ip
	Ipbits    []bool     `json:"ipbits"` //使用位图法保存ip的分配
	Subnet    *net.IPNet `json:"subnet"`
	Allocated int        `json:"allocated"` //已分配总数
}

func NewIPAllocator(subnet *net.IPNet) *IPAlloc {
	ones, _ := subnet.Mask.Size()
	bitMapSize := 1 << uint(32-ones)
	return &IPAlloc{
		Ipbits:    make([]bool, bitMapSize),
		Subnet:    subnet,
		Allocated: 0,
	}
}

// 分配ip
func (a *IPAlloc) Allocate() (net.IP, bool) {
	for i, used := range a.Ipbits {
		if i == 0 { //网络地址默认不分配
			continue
		}
		if !used {
			ip := make(net.IP, 4)
			copy(ip, a.Subnet.IP)
			//这里可以写成循环，但是为了便于理解写成这样
			//举例[192.168.0.0]需要分配的网络的下表是577则对应[0,0,2,65]
			//则分配的ip就是[192.168.2.65]
			ip[3] += byte(i % 256)
			ip[2] += byte((i / 256) % 256)
			ip[1] += byte((i / 256 / 256) % 256)
			ip[0] += byte((i / 256 / 256 / 256) % 256)
			a.Ipbits[i] = true
			a.Allocated++
			return ip, true
		}
	}
	return nil, false
}

// 释放ip
func (a *IPAlloc) Realese(ip net.IP) {
	ip0 := a.Subnet.IP.To4() //子网ip
	ip = ip.To4()            //需要的ipv4的格式
	i := 0
	i += int(ip[3] - ip0[3])
	i += int(ip[2]-ip0[2]) * 256
	i += int(ip[1]-ip0[1]) * 256 * 256
	i += int(ip[0]-ip0[0]) * 256 * 256 * 256
	if i >= len(a.Ipbits) { //如果越界
		return
	}
	a.Ipbits[i] = false
	a.Allocated--
}
