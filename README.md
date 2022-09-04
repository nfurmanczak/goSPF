# goSPF
SPF validator written in golang. This tool can read and validatei SPF records (Sender Policy Framework). You can also check if IPv4 oder IPv6 addresses are covered by the SPF record. A SPF record is a special TXT record in the DNS zone. This TXT record contains IP addresses or networks. This addresses are alllows to send emails for the domain. 

More information about SPF (Sender Policy Framework): 
- [https://en.wikipedia.org/wiki/Sender_Policy_Framework](https://en.wikipedia.org/wiki/Sender_Policy_Framework)
- [https://datatracker.ietf.org/doc/rfc7208/](https://datatracker.ietf.org/doc/rfc7208/)


## Compile source code 

I would recommend to install the golang package on your system. The golang package is available for Linux, Windows and macOS (Intel and ARM). You can find more information about the installation or packages for Windows on the golang webpage: [https://go.dev/doc/install](https://go.dev/doc/install)

Don't forget to add the dir /usr/local/go/bin to the PATH var on your macOS or Linux system: 

	export PATH=$PATH:/usr/local/go/bin

Check the version to verify if golang is installed correctly: 

	go version 
	go version go1.18.5 linux/amd64


Clone the git repo: 

	git clone https://github.com/nfurmanczak/goSPF.git

Change the dir and run go build to compile the code: 

	cd goSPF
	go build 

The command go build will create the executable file goSPF. You can check the version from goSPF to check if the tool works: 

	./goSPF version 


## Usage

goSPF needs to be always started with a domain. The tool will search for all TXT record in the DNS zone and checks which TXT record is a valid SPF record. 

	goSPF google.com
	SPF-Record: v=spf1 include:_spf.google.com ~all

goSPF will print the found SPF record. Depending on whether the program finds an SPF record, it exits with an adequate exit code. The exit code is '0' for a valid SPF record and '2' if no SPF was found or the SPF record is invalid. 

A SPF record is valid when: 
- It is a TXT record 
- Is the only SPF record for this FQDN 
- Starts with v=spf1
- Ends with -all, ~all or ?all

goSPF can also print all IP addresses and networks which are covered by the SPF record. Just start the tool with the verbose mode. 

	goSPF google.com verbose 
	SPF-Record: v=spf1 include:_spf.google.com ~all
	Softfail found (~all).
	The SPF-record for google.com contains the following IP addresses and networks:
	IPv4 networks:
	- 172.217.0.0/19
	- 172.217.32.0/20
	...

Another key function is to check if given IP address is covered by the SPF record. The tool can also check if a IP address if part of a network which is defined in the SPF record. 
 
	goSPF google.com 35.190.247.1             
	SPF-Record: v=spf1 include:_spf.google.com ~all
	All IPv4 addresses are covered by the SPF record.

	goSPF google.com 192.168.0.1 127.0.0.1         
	SPF-Record: v=spf1 include:_spf.google.com ~all
	2 IPv4 addresse(s) are not part of the SPF-record:
	- 192.168.0.1
	- 127.0.0.1

goSPF can work with IPv4 and IPv6 addresses. The addresses must be separated only by whitespaces. IPv4 and IPv6 addresses do not have to be specified separately. The addresses can also be mixed. If all IP addresses are present in the SPF record, goSPF exits with the exit code '0'. If even only one address is not present in the SPF record, goSPF exits with the exit code '2'. 
