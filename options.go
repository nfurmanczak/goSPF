package main

import (
	"fmt"
)

func help() {
	fmt.Println(`The tool goSPF can be used to analyse SPF record or for the monitoring of a SPF record. 
The program must be started with at least one domain as transfer parameter. The tool will 
search for a SPF record and will display all ressaurces. 
	
Example: 
./gospf example.org

The tool can also check if a IPv4 or IPv6 address is part of the SPF record. It can handle 
one single IP or a list of multiple IPs.  
	
Example: 
./gospf exmaple.org 127.0.0.1 192.168.2.1 2001:db8:1::ab9:C0A8:102`)
}

func version() {

	fmt.Println(`Version 0.1
goSPF - Check and validate SPF record in DNS zones 
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
