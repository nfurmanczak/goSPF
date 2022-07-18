package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

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

func findAllIncludes(includes *[]string) {
	// This function finds all possible includes for the SPF record. A include from a SPF record can also contain one or more other includes.
	// It makes only sense to execute this function, if the var spfRecord contains minimum one include.
	// The SPF RFC limits the number to maximum 10 includes (include_counter). This function will print a warning message if it will find more then
	// 10 includes and exit the program at 15 includes. The number 15 has no particular meaning and is only used to avoid infinity loops.
	//

	var include_counter int = 0 // count how many includes are found
	var tmp_spfRecord string
	var tmp_includes []string // A copy from includes
	var new_includes []string
	var trans_includes []string

	// The var includes should contain all found includes when this function was executed. We need to copy all elements from
	// includes to an other slice because we need to flush the var tmp_include
	tmp_includes = *includes

	// Infinity loop which will only be exit when no more includes are found.
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

				// Print a warning if we found more than 10 nested includes. IP and domains in this SPF records will be ignored
				if include_counter >= 10 {
					fmt.Println("Warning: More than 10 nested includes detected.")

					// It is possible to get stuck in an endless loop when includes refer to each other. We need to exit the program at in this case
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
				txtrecords = nil
			}
		}
		//
		// END LOOP
		//

		if len(trans_includes) != 0 {

			// Remove all elements in tmp_includes and prepare this slice for the next for loop run
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

}

func findAllIncludesWithMaps(includes *[]string) {
	// This function finds all possible includes for the SPF record. A include from a SPF record can also contain one or more other includes.
	// It makes only sense to execute this function, if the var spfRecord contains minimum one include.
	// The SPF RFC limits the number to maximum 10 includes (include_counter). This function will print a warning message if it will find more then
	// 10 includes and exit the program at 15 includes. The number 15 has no particular meaning and is only used to avoid infinity loops.

	var include_counter int = 0 // count how many includes are found
	var tmp_spfRecord string
	var tmp_includes []string // A copy from includes
	var new_includes []string
	var trans_includes []string

	// The var includes should contain all found includes when this function was executed. We need to copy all elements from
	// includes to an other slice because we need to flush the var tmp_include
	tmp_includes = *includes

	// Infinity loop which will only be exit when no more includes are found.
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

				// Print a warning if we found more than 10 nested includes. IP and domains in this SPF records will be ignored
				if include_counter >= 10 {
					fmt.Println("Warning: More than 10 nested includes detected.")

					// It is possible to get stuck in a endless loop when includes refer to each other. We need to exit the program at in this case
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
				txtrecords = nil
			}
		}
		//
		// END LOOP
		//

		if len(trans_includes) != 0 {

			// Remove all elements in tmp_includes and prepare this slice for the next for loop run
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

}
