package IpUtils

import (
	"net"
)

/**
获取本地IP
*/
func GetLocalIp() (string, error) {
	ip, err := GetOutBoundIP()
	if len(ip) < 1 {
		return GetLocalIpByInterfaceAddrs()
	}
	return ip, err
}

func GetLocalIpByInterfaceAddrs() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok {
			if !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", nil
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = localAddr.IP.String()
	return
}
