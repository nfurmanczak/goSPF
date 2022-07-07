package main

import (
	"fmt"
	"regexp"
)

func finxMXMechanism(spfRecord []string) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).
	// var MXwithDomainRegex = regexp.MustCompile(`mx:\S+`)
	// var MXwithoutDomainRegex = regexp.MustCompile(`mx\s`)

	for _, x := range spfRecord {
		fmt.Println(x)
	}
}

func findSingleIP4Networks(record string) (ipv4Networks []string) {
	var validIPv4Network = regexp.MustCompile(`ip4:([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\/[0-9]+)`)

	ipv4Networks = validIPv4Network.FindAllString(record, -1)

	return
}

func findIP4Addresses(spfRecord string) {
	var IPv4Addresses []string

	var validIPv4Addresses = regexp.MustCompile(`ip4:([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\/[0-9]+)`)

	IPv4Addresses = validIPv4Addresses.FindAllString(spfRecord, -1)

	fmt.Println(IPv4Addresses)
}
