package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func findAllIncludes(includes *[]string) {
	// This function finds all possible includes for the SPF record. A include from a SPF record can also contain one or more other includes.
	// It makes only sense to execute this function, if the var spfRecord contains minimum one include.
	// The SPF RFC limits the number to maximum 10 includes (include_counter). This function will print a warning message if it will find more then
	// 10 includes and exit the program at 15 includes. The number 15 has no particular meaning and is only used to avoid infinity loops.

	var include_counter int = 0 // count how many includes are found
	var tmp_spfRecord string
	var tmp_includes []string // A copy from includes
	var new_includes []string
	var trans_includes []string

	// The var includes shoud contain all found includes when this function was executed. We need to copy all elements from
	// includes to an other slice because we need to flush the var tmp_include
	tmp_includes = *includes

	// Endlos loop which will only be exit when no more includes are found.
	for true {

		// Loop all elements in tmp_includes ...
		for _, x := range tmp_includes {

			// Get all TXT records from include domain
			txtrecords, dns_error := net.LookupTXT(x)

			// Print a warning when we can not get the TXT record(s) from the domain in the var x
			// This can have different reasons: DNS timeout, no TXT record available, etc. ...
			// It is better to print a warning, skip all other steps and continue with the next domain in tmp_includes
			if dns_error != nil {
				fmt.Println("Warning: Can not get TXT-Record from Domain:", x, "for include search.")
				continue
			}

			// txtrecords can contain multiple TXT record. We need to find a valid SPF record
			tmp_spfRecord = findSPFRecord(txtrecords)
			new_includes = findIncludeInSPFRecord(tmp_spfRecord)

			// Check if tmp_spfRecord (new_incudes) contains also more include tags ...
			if new_includes != nil {

				include_counter += 1

				// Print a warning if we found more then 10 nested includes. IP and domains in this SPF records will be ignored
				if include_counter >= 10 {
					fmt.Println("Warning: More than 10 nested includes detected.")

					// It is possible to get stuck in a endles loop when includes refer to each other. We need to exit the program at in this case
				} else if include_counter >= 15 {
					fmt.Println("Error: Possible loop with include domains.")
					os.Exit(2)
				}

				// Add all new found domains from include mechanism to the slie
				for _, x := range new_includes {
					*includes = append(*includes, x)
					trans_includes = append(trans_includes, x)
				}

				new_includes = nil

			} else {
				break
			}
		}

		if len(trans_includes) != 0 {

			// Remove all elemente in tmp_includes and prepare this slice for the next for loop run
			tmp_includes = nil

			// Copy all new found include domains into the
			for _, x := range trans_includes {
				tmp_includes = append(tmp_includes, x)
			}
			trans_includes = nil
		} else {
			// No more include domains found, we need to exit this for loop
			break
		}
	}

	//fmt.Println("Gefundene Includes:")
	//for _, x := range *includes {
	//	fmt.Println("-", x)
	//}

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

func findIncludeInSPFRecord(spfRecord string) (foundInclude []string) {

	foundInclude = []string{}

	if strings.Contains(spfRecord, "include:") {
		includeRegex := regexp.MustCompile(`include:(\S+)`)
		for _, x := range includeRegex.FindAllString(spfRecord, -1) {
			foundInclude = append(foundInclude, strings.Replace(x, "include:", "", 1))
		}
	} else {
		foundInclude = nil
	}

	return
}

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

func findSingleIP4Networks(record string) (ipv4Networks []string) {
	var validIPv4Network = regexp.MustCompile(`ip4:([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\/[0-9]+)`)

	ipv4Networks = validIPv4Network.FindAllString(record, -1)

	return
}

func checkForValidDomain(domain string) (DomainCheck bool) {
	var DomainRegex = regexp.MustCompile(`^[a-zA-Z0-9\.-]+\.[a-zA-Z]{2,}$`)

	if DomainRegex.MatchString(domain) {
		DomainCheck = true
	} else {
		DomainCheck = false
	}

	return
}

func findAllMechanism(spfRecord string) {
	var allRegex = regexp.MustCompile(`[-?+~]{1}all$`)
	var all = allRegex.FindString(spfRecord)

	if all != "" {
		if all == "-all" {
			fmt.Println("Hardfail found (-all)")
		} else if all == "~all" {
			fmt.Println("Softfail found (~all)")
		} else if all == "?all" {
			fmt.Println("Neutral found (?all)")
		} else if all == "+all" {
			fmt.Println("Error: +all is invalid.")
			os.Exit(2)
		} else {
			fmt.Println("Error: all found but is invalid.")
			os.Exit(2)
		}
	} else {
		fmt.Println("Error: Invalid SPF RR. No all mechanism found.")
		os.Exit(2)
	}

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
		fmt.Println("Error: No TXT DNS-Reord found")
		os.Exit(3)
	}

	var spfRecord string = findSPFRecord(txtrecords)

	fmt.Println(spfRecord)

	if spfRecord == "null" {
		fmt.Println("Error: No SPF record found for Domain: ", domain)
		os.Exit(2)
	}

	spfRecord = findRedirect(spfRecord)
	findAllMechanism(spfRecord)

	findIncludeInSPFRecord(spfRecord)

	var includes = []string{}
	includes = findIncludeInSPFRecord(spfRecord)

	// If
	if len(includes) != 0 {
		findAllIncludes(&includes)

		for _, x := range includes {
			fmt.Println("-", x)
		}
	}

}
