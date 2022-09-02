package main

import (
	"fmt"
	"os"
	"regexp"
)

func findAllQualifier(spfRecord string, verbose_mode bool) {
	var allRegex = regexp.MustCompile(`[-?+~]{1}all$`)
	var all = allRegex.FindString(spfRecord)
	var verbose_message string 

	if all != "" {
		if all == "-all" {
			verbose_message="Hardfail found (-all)."
		} else if all == "~all" {
			verbose_message="Softfail found (~all)."
		} else if all == "?all" {
			verbose_message="Neutral found (?all). This is not recommended. Try to use -all."
		} else if all == "+all" {
			verbose_message="Error: +all is invalid. Please check SPF-Record!"
			os.Exit(2)
		} else {
			fmt.Println("Error: all found but is invalid. Please check SPF-Record!")
			os.Exit(2)
		}
	} else {
		fmt.Println("Error: Invalid SPF-Record. No all mechanism found. Please check SPF-Record!")
		os.Exit(2)
	}
	
	if verbose_mode { 
		fmt.Println(verbose_message)
	}
}
