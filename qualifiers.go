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
			fmt.Println("Hardfail found (-all).")
		} else if all == "~all" {
			fmt.Println("Softfail found (~all).")
		} else if all == "?all" {
			fmt.Println("Neutral found (?all). This is not recommended. Try to use -all.")
		} else if all == "+all" {
			fmt.Println("Error: +all is invalid. Please check SPF-Record!")
			os.Exit(2)
		} else {
			fmt.Println("Error: all found but is invalid. Please check SPF-Record!")
			os.Exit(2)
		}
	} else {
		fmt.Println("Error: Invalid SPF-Record. No all mechanism found. Please check SPF-Record!")
		os.Exit(2)
	}
}
