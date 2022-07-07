package main

import (
	"fmt"
	"os"
	"regexp"
)

func findAllQualifier(spfRecord string) {
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
