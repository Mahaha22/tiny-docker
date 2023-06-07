package network

type DriverType int

const (
	Bridge  DriverType = iota //默认的网络类型
	Host                      //容器与主机共享网络网络命名空间
	Overlay                   //集群中的通信方式
	None                      //无网络
)

type Driver interface {
	Create(*Network) error
	Remove() error
}

func NewDriver(Dtype string) Driver {
	switch Dtype {
	case "bridge":
		return &BridgeDriver{}
	// case "overlay":
	// 	nw.Driver.Dtype = network.Overlay
	case "none":
		return &NoneDriver{}
	case "host":
		return &HostDriver{}
	default:
		return &BridgeDriver{} //bridge作为默认驱动
	}
}
