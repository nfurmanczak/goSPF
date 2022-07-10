package main

import (
	"regexp"
	"strings"
)

func finxMX(spfRecords []string) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).
	// var MXwithDomainRegex = regexp.MustCompile(`mx:\S+`)
	// var MXwithoutDomainRegex = regexp.MustCompile(`mx\s`)

	//mxWithoutDomain := regexp.MustCompile(`ip4:\S+`)
	//mxWithDomain := regexp.MustCompile(`ip4:\S+`)

	// Find mx record for the own domain

	// Finx mx records that are refereding to a other domain
}

func findIP6(spfRecord []string) (ip6list []string) {
	// DUMMY FUNCTION
	return
}

func findIP4(spfRecords []string) (ip4list []string) {

	ip4Regex := regexp.MustCompile(`ip4:\S+`)

	for _, spfRecord := range spfRecords {
		for _, x := range ip4Regex.FindAllString(spfRecord, -1) {
			ip4list = append(ip4list, strings.Replace(x, "ip4:", "", 1))
		}
	}

	return
}

func findIPv4Networks(ip4list []string) (ip4networks []string) {

	for _, line := range ip4list {
		if strings.Contains(line, "/") {
			ip4networks = append(ip4networks, line)
		}
	}

	return
}

func findIPv4Addresses(ip4list []string) (ip4addresses []string) {

	for _, line := range ip4list {
		if !strings.Contains(line, "/") {
			ip4addresses = append(ip4addresses, line)
		}
	}

	return
}
