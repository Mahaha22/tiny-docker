package network

type NoneDriver struct {
	Name string `json:"name"` //名字
}

func (n *NoneDriver) Create(*Network) error {
	//none网络什么都不需要创建，本质上是容器在启动的时候新建一个网络命名空间但是却不配置任何网络
	return nil
}
func (n *NoneDriver) Remove() error {
	return nil
}
