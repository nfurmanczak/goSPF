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

	// Check if the user started the application with a
	if len(os.Args) > 1 {
		domain = strings.ToLower(os.Args[1])

		// The function "checkForValidDomain" checks is the var "domain" contains a valid domain
		// This is done via a simple RegEx which expects a dot and then a TLD. This method will not detect all invalid domains.
		if checkForValidDomain(domain) == false {
			fmt.Println("Error:", domain, "is not a valid domain.")
			os.Exit(3)
		}
		
	} else {
		// Exit the application with exit code 2 when a domain as transfer parameter is missing
		fmt.Println("Error: Domain missing.")
		fmt.Println("Usage: ./spf_check example-domain.org")
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

	fmt.Println(spfMap)

	spfRecord = findRedirect(spfRecord)
	
	findAllQualifier(spfRecord)

	var includes = []string{}
	includes = findIncludeInSPFRecord(spfRecord)

	// We need to collect all SPF records in one string slice. This slice need to contains the SPF record
	// form the var spfRecord and from all SPF Records from the include slice
	var spfRecords []string

	fmt.Println(len(includes))

	if len(includes) != 0 {
		findAllIncludes(&includes)

		for _, x := range includes {
			txtrecords, dns_error := net.LookupTXT(x)
			if dns_error != nil {
				fmt.Println("Warning: No TXT Record found for Inclue Domain", x)
			}

			spfRecords = append(spfRecords, findSPFRecord(txtrecords))
		}
	} else {
		spfRecords = append(spfRecords, spfRecord)
	}

	tstList := findIP4(spfRecords)

	ip4addr := findIPv4Addresses(tstList)
	ip4nets := findIPv4Networks(tstList)

	fmt.Println("IP4 Adresses: ")
	for _, x := range ip4addr {
		fmt.Println("=>", x)
	}

	fmt.Println("IP4 CIDS: ")
	for _, x := range ip4nets {
		fmt.Println("=>", x)
	}

}
