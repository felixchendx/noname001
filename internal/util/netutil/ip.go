package netutil

import (
	"net"
	"strings"
)

// credits: // https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
func CollectIPs() ([]string, error) {
	// TODO: nic type ?
	netIfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	collectedIPs := make([]string, 0)
	for _, ni := range netIfaces {
		addrs, err := ni.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet : ip = v.IP
			case *net.IPAddr: ip = v.IP
			}

			if ip != nil {
				if ip.IsLoopback() { continue }

				ipString := ip.String()
				if IsIPv4(ipString) {
					collectedIPs = append(collectedIPs, ipString)
				}
			}
		}
	}

	return collectedIPs, nil
}

// credits: // https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6
func IsIPv4(ipAddr string) (bool) {
	return strings.Count(ipAddr, ":") < 2
}
func IsIPv6(ipAddr string) (bool) {
	return strings.Count(ipAddr, ":") >= 2
}
