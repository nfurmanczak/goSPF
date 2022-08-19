// GO version 1.19
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	var domain string

	// This map contains the domain as key and SPF record as value
	spfMap := make(map[string]string)

	var UserIP4Check []string
	var UserIP6Check []string

	// Check if the user started the application with a valid value
	if len(os.Args) > 1 {
		domain = strings.ToLower(os.Args[1])

		for _, arg := range os.Args {
			if strings.ToLower(arg) == "help" {
				help()
				os.Exit(0)
			}

			if strings.ToLower(arg) == "version" {
				version()
				os.Exit(0)
			}

			ipaddr := net.ParseIP(arg)

			if ipaddr != nil {
				if strings.Contains(ipaddr.String(), ":") {
					// Found a IPv6 address ...
					UserIP6Check = append(UserIP6Check, ipaddr.String())
				} else {
					// Not a IPv6 address
					UserIP4Check = append(UserIP4Check, ipaddr.String())
				}
			}
		}

		// The function "checkForValidDomain" checks is the var "domain" contains a valid domain
		// This is done via a simple RegEx which expects a dot and then a TLD. This method will not detect all invalid domains.
		if !checkForValidDomain(domain) {
			fmt.Println("Error:", domain, "is not a valid domain.")
			os.Exit(3)
		}
	} else {
		// Exit the application with exit code 2 when a domain as transfer parameter is missing
		fmt.Println("Error: Domain missing.")
		fmt.Println("Usage: ./spf_check example-domain.org [1.2.3.4] [2001:db8:1::ab9:C0A8:102] [debug] [version] [help] [monitor]")
		fmt.Println("")
		os.Exit(2)
	}

	// Get TXT records from the domain via DNS lookup ...
	txtrecords, dns_error := net.LookupTXT(domain)

	// ... exit application with exit code 3 when the DNS lookup is not possible (e.g. timeout, .. )
	if dns_error != nil {
		fmt.Println("Error: No TXT DNS-Record found")
		os.Exit(3)
	}

	// The slice txtrecords can contain multiple txt records. The function will search for SPF records and return the
	// found records as a string. The function will exit the application with error code 3 if zero or more then one SPF record is found
	var spfRecord string = findSPFRecord(txtrecords)

	spfMap[domain] = findSPFRecord(txtrecords)

	// A SPF record can contain a redirect which points to an SPF record from a different domain
	spfRecord = findRedirect(spfRecord)

	fmt.Println("SPF after redirect:", spfRecord)

	findAllQualifier(spfRecord)

	var includes = []string{}
	includes = findIncludeInSPFRecord(spfRecord)

	// We need to collect all SPF records in one string slice. This slice need to contains the SPF record
	// form the var spfRecord and from all SPF Records from the include slice
	var spfRecords []string

	if len(includes) != 0 {
		findAllIncludes(&includes)

		for _, domain := range includes {
			txtrecords, dns_error := net.LookupTXT(domain)
			if dns_error != nil {
				fmt.Println("Warning: No TXT Record found for Inclue Domain", domain)
			} else {
				spfRecords = append(spfRecords, findSPFRecord(txtrecords))
				spfMap[domain] = findSPFRecord(txtrecords)
			}
		}
	} else {
		spfRecords = append(spfRecords, spfRecord)
		spfMap[domain] = spfRecord
	}

	ipv4slice := findIP4(spfMap)
	ip6slice := findIP6(spfMap)

	ip4addr := findIP4Addresses(ipv4slice)
	ip4nets := findIP4Networks(ipv4slice)
	ip6addr := findIP6Addresses(ip6slice)
	ip6nets := findIP6Networks(ip6slice)

	if len(ip4addr) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("IPv4 Addresses:")
		fmt.Println("---------------------------")
		for _, x := range ip4addr {
			fmt.Println("-", x)
		}
	}

	if len(ip4nets) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("IPv4 Networks:")
		fmt.Println("---------------------------")
		for _, x := range ip4nets {
			fmt.Println("-", x)
		}
	}

	if len(ip6addr) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("IPv6 Addresses:")
		fmt.Println("---------------------------")
		for _, x := range ip6addr {
			fmt.Println("-", x)
		}
	}

	if len(ip6nets) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("IPv6 Networks:")
		fmt.Println("---------------------------")
		for _, x := range ip6nets {
			fmt.Println("-", x)
		}
	}

	//exampleIP, _ := netip.ParseAddr("52.82.175.255")

	/*
		fmt.Println("Check if IP is in CIDR Network")
		for _, x := range ip4nets {
			network, _ := netip.ParsePrefix(x)

			if network.Contains(exampleIP) {
				fmt.Println(exampleIP, "is part from the network", x)
			}
		}
	*/
	aIPs := findARecord(spfMap)
	mxIPs := findMXRecord(spfMap)

	if len(aIPs) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("A Includes:")
		fmt.Println("---------------------------")

		for _, ip := range aIPs {
			fmt.Println("-", ip)
		}
	}

	if len(mxIPs) > 0 {
		fmt.Println("---------------------------")
		fmt.Println("MX Includes:")
		fmt.Println("---------------------------")

		for _, ip := range mxIPs {
			fmt.Println("-", ip)
		}
	}

	if len(UserIP4Check) > 0 || len(UserIP6Check) > 0 {
		fmt.Println("")
		fmt.Println("////////////////////////////")
		fmt.Println("IPs to check")
		fmt.Println("////////////////////////////")

		fmt.Println("IPv4:")
		for _, x := range UserIP4Check {
			fmt.Println("=>", x)
		}

		fmt.Println("IPv6:")
		for _, x := range UserIP6Check {
			fmt.Println("=>", x)
		}
	}

	//mergedIPs := append([]string{}, append(aIPs, mxIPs...) ...)

	if (len(aIPs) != 0) || (len(mxIPs) != 0) {
		mergeSlices(append(aIPs, mxIPs...), &ip4addr, &ip6addr)
	}

	if len(UserIP4Check) != 0 {
		UserIP4Check = compareIPAddr(ip4addr, UserIP4Check)

		if len(ip4nets) != 0 {
			UserIP4Check = checkNetworks(ip4nets, UserIP4Check)
		}

		if len(UserIP4Check) != 0 {
			fmt.Println(len(UserIP4Check), "IPv4 addresse(s) are not part of the SPF-record:")

			for _, i := range UserIP4Check {
				fmt.Println("-", i)
			}
		} else {
			fmt.Println("All IPv4 addr are coverd with the SPF record")
		}
	}

	if len(UserIP6Check) != 0 {
		UserIP6Check = compareIPAddr(ip6addr, UserIP6Check)

		if len(ip6nets) != 0 {
			UserIP6Check = checkNetworks(ip6nets, UserIP6Check)
		}

		if len(UserIP6Check) != 0 {
			fmt.Println(len(UserIP6Check), "IPv6 addresse(s) are not part of the SPF record:")

			for _, i := range UserIP6Check {
				fmt.Println("-", i)
			}
		} else {
			fmt.Println("All IPv6 addr are covers with the SPF record")
		}
	}

}
