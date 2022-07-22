package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func validateIPv4Addr(input string) (returnValue bool) {
	var ipv4Regex = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	if ipv4Regex.MatchString(input) {
		returnValue = true
	} else {
		returnValue = false
	}

	return
}

func findMXRecord(spfRecords map[string](string)) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).
	var MXRecordWithDomainRegex = regexp.MustCompile(`mx:\S+`)
	var MXRecordWithoutDomainRegex = regexp.MustCompile(`mx\s`)

	mxMap := make(map[string]([]string))

	for domain, spfrr := range spfRecords {
		for _, mxTag := range MXRecordWithDomainRegex.FindAllString(spfrr, -1) {
			MXRecordDomain := strings.Replace(mxTag, "mx:", "", 1)
			records, dns_error := net.LookupMX(MXRecordDomain)

			if dns_error == nil {
				for _, mxhost := range records {
					ips, dns_error := net.LookupIP(mxhost.Host)

					if dns_error == nil {
						for _, ip := range ips {
							mxMap[MXRecordDomain] = append(mxMap[MXRecordDomain], ip.String())
						}
					}
				}
			}
		}

		for range MXRecordWithoutDomainRegex.FindAllString(spfrr, -1) {
			records, dns_error := net.LookupMX(domain)

			if dns_error == nil {
				for _, mxhost := range records {

					ips, _ := net.LookupIP(mxhost.Host)

					for _, ip := range ips {
						mxMap[domain] = append(mxMap[domain], ip.String())
					}
				}
			}
		}

	}

	for domain, ips := range mxMap {
		fmt.Println("Domain:", domain)
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}

}

func findARecord(spfRecords map[string](string)) {
	var ARecordWithDomainRegex = regexp.MustCompile(`a:\S+`)
	var ARecordWithoutDomainRegex = regexp.MustCompile(`a\s`)

	aMap := make(map[string]([]string))

	for domain, spfrr := range spfRecords {

		for _, x := range ARecordWithDomainRegex.FindAllString(spfrr, -1) {
			ARecordDomain := strings.Replace(x, "a:", "", 1)
			records, dns_error := net.LookupIP(ARecordDomain)

			if dns_error == nil {
				for _, ip := range records {
					aMap[ARecordDomain] = append(aMap[ARecordDomain], ip.String())
				}
			}
		}

		for range ARecordWithoutDomainRegex.FindAllString(spfrr, -1) {
			records, dns_error := net.LookupIP(domain)

			if dns_error == nil {
				for _, ip := range records {
					aMap[domain] = append(aMap[domain], ip.String())
				}
			}
		}

	}

	//fmt.Println("---------------")
	//fmt.Println(aMap)
	//fmt.Println("---------------")

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
