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
Autor nikolai.furmanczak@gmail.com`)
}
