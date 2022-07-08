package main

import (
	"regexp"
	"fmt"
	"os"
)

func checkForValidDomain(domain string) (DomainCheck bool) {
	var DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9\.-]+\.[a-zA-Z]{2,}$`)

	if DomainRegex.MatchString(domain) {
		DomainCheck = true
	} else {
		DomainCheck = false
	}

	return
}

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
