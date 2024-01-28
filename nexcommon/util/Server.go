package util

import (
	"net"
	"os"
)

func GenerateIPAddress(name string) (ip string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var first string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
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
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			if first == "" {
				first = ip.String()
			}
			if name != "" {
				if iface.Name != name {
					continue
				}
			}
			return ip.String(), nil
		}
	}
	if ip == "" {
		return first, nil
	}
	return "", err
}

func GenerateHostname() (string, error) {
	return os.Hostname()
}
