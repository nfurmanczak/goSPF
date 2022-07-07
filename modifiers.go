package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func findRedirect(spfRecord string) (redirectSPF string) {
	var redirect string

	if strings.Contains(spfRecord, "redirect=") {
		redirectRegex := regexp.MustCompile(`redirect=\S+`)
		redirect = redirectRegex.FindAllString(spfRecord, 1)[0]
		redirect = strings.Replace(redirect, "redirect=", "", 1)
		foundtxtrecord, dns_error := net.LookupTXT(redirect)
		redirectSPF = findSPFRecord(foundtxtrecord)

		if dns_error != nil {
			fmt.Println("Error: No TXT DNS-Reord found")
			os.Exit(3)
		}

	} else {
		redirectSPF = spfRecord
	}

	return
}
