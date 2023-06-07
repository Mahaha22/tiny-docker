// 与主机共享网络
package network

type HostDriver struct {
	Name string
}

func (h *HostDriver) Create(*Network) error {
	//Host网络什么都不需要创建，本质上是容器在启动的时候不需要新建网络命名空间
	return nil
}
func (h *HostDriver) Remove() error {
	return nil
}
