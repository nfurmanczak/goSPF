package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func findMXRecord(spfRecords map[string](string)) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).

	var MXRecordWithDomainRegex = regexp.MustCompile(`mx:\S+`)
	var MXRecordWithoutDomainRegex = regexp.MustCompile(`mx\s`)

	// Map is not a good data struct. I need to change this
	mxMap := make(map[string](string))

	for domain, spfrr := range spfRecords {

		for _, x := range MXRecordWithDomainRegex.FindAllString(spfrr, -1) {
			MXRecordDomain := strings.Replace(x, "mx:", "", 1)
			records, dns_error := net.LookupMX(MXRecordDomain)
			fmt.Println(records)

			if dns_error == nil {
				for _, ip := range records {
					fmt.Println(ip.Host)
					mxMap[MXRecordDomain] = ip.Host
				}
			}
		}

		for range MXRecordWithoutDomainRegex.FindAllString(spfrr, -1) {

			records, dns_error := net.LookupMX(domain)
			fmt.Println(records)

			if dns_error == nil {

				for _, ip := range records {
					fmt.Println(ip.Host)
					mxMap[domain] = ip.Host
				}
			}
		}

	}

	for domain, ip := range mxMap {
		fmt.Println(domain, " ", ip)
	}

}

func findARecord(spfRecords map[string](string)) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).
	var ARecordWithDomainRegex = regexp.MustCompile(`a:\S+`)
	var ARecordWithoutDomainRegex = regexp.MustCompile(`a\s`)

	//mxWithoutDomain := regexp.MustCompile(`ip4:\S+`)
	//mxWithDomain := regexp.MustCompile(`ip4:\S+`)

	aMap := make(map[string](string))

	for domain, spfrr := range spfRecords {

		for _, x := range ARecordWithDomainRegex.FindAllString(spfrr, -1) {
			ARecordDomain := strings.Replace(x, "a:", "", 1)
			records, dns_error := net.LookupIP(ARecordDomain)

			if dns_error == nil {
				for _, ip := range records {
					aMap[ARecordDomain] = ip.String()
				}
			}
		}

		for range ARecordWithoutDomainRegex.FindAllString(spfrr, -1) {

			records, dns_error := net.LookupIP(domain)

			if dns_error == nil {
				for _, ip := range records {
					aMap[domain] = ip.String()
				}
			}
		}

	}

	for domain, ip := range aMap {
		fmt.Println(domain, " ", ip)
	}

}

func findIP6(spfRecord []string) (ip6list []string) {
	// DUMMY FUNCTION
	return
}

func findIP4(spfRecords map[string](string)) (ip4list []string) {

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
