package main 

import ( 
	"fmt"
)

func help() { 
	fmt.Println(
	`
	Help     
	The tool goSPF can be used to analyse SPF record or for the monitoring of a SPF record. 
	The program must be started with at least one domain as transfer parameter. The tool will 
	search for a SPF record and will display all ressaurces. 
	
	Example: 
	./gospf example.org

	The tool can also check if a IPv4 or IPv6 address is part of SPF record. The tool will oif cause aks
	
	Example: 
	./gospf exmaple.org 127.0.0.1 192.168.2.1 
	 
	`)
}