package utils

import (
	"errors"
	"net"
)

func GetLocalIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return net.IP{}, err
	}

	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if inter.Flags&net.FlagLoopback != 0 {
			continue // Loopback Interface
		}

		addresses, err := inter.Addrs()
		if err != nil {
			return net.IP{}, err
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			return ip, nil
		}
	}
	return net.IP{}, errors.New("Are you connected to the internet?")
}
