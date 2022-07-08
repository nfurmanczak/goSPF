/*
   This function search for a "redirect" in a SPF record. A redirect refers to an other domain which includes a valid SPF record. 
   It does not include another SPF record like an include. In opposite to an include a redirect is also a MODIFIER. Modifiers are using 
   always the equals sign (=) and not a colon (:). A redirect is relatively simple and does not require an all at the end. 
   However, the SPF record being referenced must have an all. Otherwise an SPF record would not be valid. 
   
  Examples for a valid SPF record with a redirect
  
  		"v=spf1 redirect=spfrr_.domain.org" 
		"v=spf1 redirect=example.net"
		
  A common mistake is to use ':' instead of '='. The function can detect this mistake and will exit the application with error code 2. 
  
  If this function will find a valid redirect, it will lookup for a valid TXT and SPF record in this domains and returns a new SPF record. 
  If no redirect is found, the old SPF record is returned. 
*/ 


package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

/* 
	findRediretct expects a valid string with a SPF record. Return value is a new string with a valid SPF record if we have find a redirect or the 
*/ 

func findRedirect(spfRecord string) (redirectSPF string) {
	
	var redirect string

	// Exit if invalid redirect found  
	if strings.Contains(spfRecord, "redirect:") { 
		fmt.Println("Error: Invalid redirect with ':' instead of '='.")
		os.Exit(2)
	}

	if strings.Contains(spfRecord, "redirect=") {
		
		// Redirect found, get TXT and SPF record from redirect domain
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
		// No redirect found, return SPF record from the transfer parameters.
		redirectSPF = spfRecord
	}

	return
}
