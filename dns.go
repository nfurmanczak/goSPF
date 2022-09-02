package main

import (
	"regexp"
	"fmt"
	"os"
)

// checkForValidDomain
// This functions checks for valid domains via a simple RegEx. Function expect one string and return TRUE if the domain 
// is valid and return FALSE if the domain invalid. 
func checkForValidDomain(domain string) (DomainCheck bool) {
	
	// This RegEx was designed to be simple and basic. It's checks for valid characters and a dot. 
	// Not all invalid domains may be detected. However, at the latest with the DNS request for TXT records, the program will notice that there is no valid DNS record and terminate. 
	// In addition, new TLDs can always be added by ICANN.  
	var DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9\.-]+\.[a-zA-Z]{2,}$`)

	if DomainRegex.MatchString(domain) {
		DomainCheck = true
	} else {
		DomainCheck = false
	}

	return
}

// findSPFRecord
// A domain can contain multiple TXT records. This function looks for records which contains the Regex
// "^v=spf1". This could be a valid SPF record. We also count how much TXT record we find with "^v=spf1"
// More than one SPF record is invalid and we will exit the application with an error message and error code 3. 
func findSPFRecord(txtrecords []string) (foundSPFrecord string) {

	var spfCounter int = 0

	for _, txtrecord := range txtrecords {
		if match, _ := regexp.MatchString("^v=spf1", txtrecord); match {
			foundSPFrecord = txtrecord
			spfCounter = spfCounter + 1
		}
	}

	if spfCounter > 1 {
		fmt.Println("Error: More then one SPF record found. SPF setup is faulty!")
		os.Exit(2)
	} else if spfCounter == 0 {
		fmt.Println("Error: No SPF record found!")
		os.Exit(2)
	}
	
	return foundSPFrecord	
}
