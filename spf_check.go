// This software was written and tested with golang version 1.18 and 1.19.
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// This map contains the domain as key and SPF record as value
	spfMap := make(map[string]string)

	// These two slices contain IPv4 and IPv6 addresses that the user specified when invoking the program.
	// It is checked whether these IP addresses are part of the SPF record. The user is not forced to specify IP addresses.
	// These two slices can therefore also be empty.
	var UserIP4Check []string
	var UserIP6Check []string

	// Var to enable or disable the verbose mode
	var verbose_mode bool = false
	var domain string = ""

	// Check if the user started the application with a valid values
	if len(os.Args) > 1 {
		domain = strings.ToLower(os.Args[1])

		// Check if other values (e.g. IP addresses, debug mode, ...) are present
		for _, arg := range os.Args {
			if strings.ToLower(arg) == "help" {
				help()
				os.Exit(0)
			}

			if strings.ToLower(arg) == "version" {
				version()
				os.Exit(0)
			}

			if strings.ToLower(arg) == "verbose" {
				verbose_mode = true
			}

			// Try to parse the string from the args to a IP address
			ipaddr := net.ParseIP(arg)

			//... check if this IP address is a IPv4 or IPv6 address.
			// We will add this addresses to the slice
			if ipaddr != nil {
				if strings.Contains(ipaddr.String(), ":") {
					// Found a IPv6 address ...
					UserIP6Check = append(UserIP6Check, ipaddr.String())
				} else {
					// Found a IPv4 address
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
		fmt.Println("Usage: ./spf_check example.org [1.2.3.4] [2001:12::1b12:0:0:1a1] [verbose] [version] [help]")
		os.Exit(3)
	}

	// Get TXT records from the domain via DNS lookup ...
	txtrecords, dns_error := net.LookupTXT(domain)

	// ... exit application with exit code 2 (critical) when the DNS lookup is not possible (e.g. timeout, ... )
	if dns_error != nil {
		fmt.Println("Error: No TXT DNS-Record found for", domain)
		os.Exit(2)
	}

	// The slice txtrecords can contain multiple txt records. The function will search for SPF records and return the
	// found records as a string. The function will exit the application with error code 3 if zero or more then one SPF record is found
	var spfRecord string = findSPFRecord(txtrecords)

	spfMap[domain] = findSPFRecord(txtrecords)

	// A SPF record can contain a redirect which points to an SPF record from a different domain
	spfRecord = findRedirect(spfRecord)

	fmt.Println("SPF-Record:", spfRecord)

	findAllQualifier(spfRecord, verbose_mode)

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

	if verbose_mode == true {
		for y, x := range spfMap {
			fmt.Println(y, x)
		}
	}

	ipv4slice := findIP4(spfMap)
	ip6slice := findIP6(spfMap)

	ip4addr := findIP4Addresses(ipv4slice)
	ip4nets := findIP4Networks(ipv4slice)
	ip6addr := findIP6Addresses(ip6slice)
	ip6nets := findIP6Networks(ip6slice)
	aIPs := findARecord(spfMap)
	mxIPs := findMXRecord(spfMap)

	//mergedIPs := append([]string{}, append(aIPs, mxIPs...) ...)

	if (len(aIPs) != 0) || (len(mxIPs) != 0) {
		mergeSlices(append(aIPs, mxIPs...), &ip4addr, &ip6addr)
	}

	if verbose_mode == true {
		verbosePrintIPs(domain, ip4addr, ip4nets, ip6addr, ip6nets)
	}

	var exit_status bool = true

	if len(UserIP4Check) != 0 {
		UserIP4Check = compareIPAddr(ip4addr, UserIP4Check)

		if len(ip4nets) != 0 {
			UserIP4Check = checkNetworks(ip4nets, UserIP4Check)
		}

		if len(UserIP4Check) != 0 {
			fmt.Printf("%d IPv4 addresse(s) are not part of the SPF-record:\n", len(UserIP4Check))
			exit_status = false
			for _, i := range UserIP4Check {
				fmt.Println("-", i)
			}
		} else {
			fmt.Println("All IPv4 addresses are covered by the SPF record.")
		}
	}

	if len(UserIP6Check) != 0 {
		UserIP6Check = compareIPAddr(ip6addr, UserIP6Check)

		if len(ip6nets) != 0 {
			UserIP6Check = checkNetworks(ip6nets, UserIP6Check)
		}

		if len(UserIP6Check) != 0 {
			fmt.Printf("%d IPv6 addresse(s) are not part of the SPF record\n:", len(UserIP4Check))
			exit_status = false
			for _, i := range UserIP6Check {
				fmt.Println("-", i)
			}

		} else {
			fmt.Println("All IPv6 addresses are covered by the SPF record.")
		}
	}

	if len(UserIP4Check) != 0 || len(UserIP6Check) != 0 {
		if exit_status {
			os.Exit(0)
		} else {
			os.Exit(2)
		}
	}

}
