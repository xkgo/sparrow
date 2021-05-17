package IpUtils

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestGetLocalIp(t *testing.T) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok {
			fmt.Printf("%+v\n", ipnet)
			if !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Println(ipnet.IP.String())
				}
			}
		}
	}
}

func TestGetOutBoundIP(t *testing.T) {
	ip, err := GetOutBoundIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)
}
