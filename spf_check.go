package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func findSPFRecord(txtrecords []string) (foundSPFrecord string) {

	var spfCounter int = 0

	for _, txtrecord := range txtrecords {
		if match, _ := regexp.MatchString("^v=spf1", txtrecord); match {
			foundSPFrecord = txtrecord
			spfCounter = spfCounter + 1
		}
	}

	if spfCounter > 1 {
		fmt.Println("Error: More then one SPF RR found.")
		os.Exit(2)
	}

	return foundSPFrecord
}

func main() {

	var domain string

	if len(os.Args) > 1 {
		domain = strings.ToLower(os.Args[1])
	} else {
		fmt.Println("Error: Domain missing.")
		fmt.Println("Usage: ./spf_check example-domain.org")
		os.Exit(2)
	}

	if checkForValidDomain(domain) == false {
		fmt.Println("Error:", domain, "is not a valid domain.")
		os.Exit(3)
	}

	txtrecords, dns_error := net.LookupTXT(domain)

	if dns_error != nil {
		fmt.Println("Error: No TXT DNS-Record found")
		os.Exit(3)
	}

	var spfRecord string = findSPFRecord(txtrecords)

	fmt.Println(spfRecord)

	if spfRecord == "null" {
		fmt.Println("Error: No SPF record found for Domain: ", domain)
		os.Exit(2)
	}

	spfRecord = findRedirect(spfRecord)
	findAllQualifier(spfRecord)

	findIncludeInSPFRecord(spfRecord)

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

	for _, x := range includes {
		fmt.Println("=> include:", x)
	}

	for _, x := range spfRecords {
		fmt.Println(x)
	}

}
