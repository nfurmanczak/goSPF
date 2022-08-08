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

func findMXRecord(spfRecords map[string](string)) (mxIPs []string) {
	// The mx mechanism can point to the original domain (mx) or to another domain (mx:example.org).
	var MXRecordWithDomainRegex = regexp.MustCompile(`mx:\S+`)
	var MXRecordWithoutDomainRegex = regexp.MustCompile(`mx\s`)

	for domain, spfrr := range spfRecords {
		for _, mxTag := range MXRecordWithDomainRegex.FindAllString(spfrr, -1) {
			MXRecordDomain := strings.Replace(mxTag, "mx:", "", 1)
			records, dns_error := net.LookupMX(MXRecordDomain)

			if dns_error == nil {
				for _, mxhost := range records {
					ips, dns_error := net.LookupIP(mxhost.Host)

					if dns_error == nil {
						for _, ip := range ips {
							mxIPs = append(mxIPs, ip.String())
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
						mxIPs = append(mxIPs, ip.String())
					}
				}
			}
		}

	}

	return
}

func findARecord(spfRecords map[string](string)) (aIPs []string) {

	var ARecordWithDomainRegex = regexp.MustCompile(`a:\S+`)
	var ARecordWithoutDomainRegex = regexp.MustCompile(`a\s`)

	//aMap := make(map[string]([]string))

	for domain, spfrr := range spfRecords {

		for _, aTag := range ARecordWithDomainRegex.FindAllString(spfrr, -1) {
			ARecordDomain := strings.Replace(aTag, "a:", "", 1)
			records, dns_error := net.LookupIP(ARecordDomain)

			if dns_error == nil {
				for _, ip := range records {
					aIPs = append(aIPs, ip.String())
				}
			} else {
				fmt.Println("DNS error")
			}
		}

		for range ARecordWithoutDomainRegex.FindAllString(spfrr, -1) {
			records, dns_error := net.LookupIP(domain)

			if dns_error == nil {
				for _, ip := range records {
					aIPs = append(aIPs, ip.String())
				}
			} else {
				fmt.Println("DNS error")
			}
		}

	}

	return
}

func findIP6(spfRecords map[string](string)) (ip6list []string) {
	ip6Regex := regexp.MustCompile(`ip6:\S+`)

	for _, spfRecord := range spfRecords {
		for _, ipTag := range ip6Regex.FindAllString(spfRecord, -1) {
			ip6 := strings.Replace(ipTag, "ip6:", "", 1)
			ip6list = (append(ip6list, ip6))
		}
	}

	return
}

func findIP6Networks(ip6list []string) (ip6networks []string) {
	for _, line := range ip6list {
		_, _, err := net.ParseCIDR(line)
		if err == nil {
			ip6networks = append(ip6networks, line)
		}
	}

	return
}

func findIP6Addresses(ip6list []string) (ip6addresses []string) {
	for _, line := range ip6list {
		ip := net.ParseIP(line)
		if ip != nil {
			ip6addresses = append(ip6addresses, line)
		}
	}

	return
}

func findIP4(spfRecords map[string](string)) (ip4list []string) {
	ip4Regex := regexp.MustCompile(`ip4:\S+`)

	for _, spfRecord := range spfRecords {
		for _, ipTag := range ip4Regex.FindAllString(spfRecord, -1) {

			ip := strings.Split(ipTag, ":")
			ip4list = append(ip4list, ip[1])
		}
	}

	return
}

func findIP4Networks(ip4list []string) (ip4networks []string) {

	for _, line := range ip4list {
		if strings.Contains(line, "/") {
			ip4networks = append(ip4networks, line)
		}
	}

	return
}

func findIP4Addresses(ip4list []string) (ip4addresses []string) {

	for _, line := range ip4list {
		if !strings.Contains(line, "/") {
			ip4addresses = append(ip4addresses, line)
		}
	}

	return
}
