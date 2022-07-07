package main

import (
	"regexp"
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
