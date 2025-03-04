package server

import (
	"net"
)

type NetworkInfo struct {
	LocalIp string
}

func (ni *NetworkInfo) GetLocalIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ni.LocalIp = ipNet.IP.String()
		}
	}
}
