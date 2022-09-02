package main

import (
	"net"
)

func mergeSlices(ips []string, ip4addr *[]string, ip6addr *[]string) {
	for _, ip := range ips {
		ipaddr := net.ParseIP(ip)

		if ipaddr.To4() == nil {
			*ip6addr = append(*ip6addr, ip)
		} else {
			*ip4addr = append(*ip4addr, ip)
		}
	}
}
