package main

import (
	"net"
	"regexp"
)

// Do i really need this function?
func validateIPv4Addr(input string) (returnValue bool) {
	var ipv4Regex = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	if ipv4Regex.MatchString(input) {
		returnValue = true
	} else {
		returnValue = false
	}

	return
}

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
