package network

type IPSubnet struct { //保存子网中以分配或者未分配的ip 例如192.168.0.0/24这个子网中以分配或者未分配ip
	Ip string `json:"ip"`
}
