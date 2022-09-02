package main

import (
	"fmt"
)

func help() {
	fmt.Println(`This tool can be used to analyse or monitor SPF record.  
The program must be started with at least the domain as transfer parameter. It will 
automatically search for an SPF record. 

Example: 
./goSPF example.org

You can get more information about the SPF record with the [verbose] mode. This includes 
information about the all mechanism and all IP addresses and networks which are covered by 
the SPF record. Depending on whether an SPF record was found, the program exits with 
a different exit code. If no SPF record or no valid SPF record is found, the program 
terminates with the exit code 2. With a valid SPF record the program terminates with 
the exit code 0. 

Example:
./goSPF example.org verbose 


The tool can also check if a IPv4 or IPv6 address is part of the SPF record. It can handle 
one single IP or a list of multiple IPv4 or IPv6 addresses. The order of the addresses 
does not matter. 

Example: 
./goSPF exmaple.org 127.0.0.1 192.168.2.1 2001:db8:1::ab9:C0A8:102

If all specified IP addresses are part of the SPF record, the program exits with 
exit code 0, otherwise with exit code 2. `)
}

func version() {

	fmt.Println(`Version 0.1
goSPF - Check and validate SPF record in DNS zones. 
Autor nikolai@furmanczak.de`)
}

func verbosePrintIPs(domain string, ip4addr []string, ip4nets []string, ip6addr []string, ip6nets []string) {
	fmt.Println("The SPF-record for", domain, "contains the following IP addresses and networks:")

	if len(ip4addr) != 0 {
		fmt.Println("IPv4 addresses:")
		for _, ip := range ip4addr {
			fmt.Println("-", ip)
		}
	}

	if len(ip4nets) != 0 {
		fmt.Println("IPv4 networks:")
		for _, network := range ip4nets {
			fmt.Println("-", network)
		}
	}

	if len(ip6addr) != 0 {
		fmt.Println("IPv6 addresses:")
		for _, ip := range ip6addr {
			fmt.Println("-", ip)
		}
	}

	if len(ip6nets) != 0 {
		fmt.Println("IPv6 networks:")
		for _, network := range ip6nets {
			fmt.Println("-", network)
		}
	}

}
