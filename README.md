# goSPF
SPF validator written in golang. This tool can read and validatei SPF records (Sender Policy Framework). You can also check if IPv4 oder IPv6 addresses are covered by the SPF record. A SPF record is a special TXT record in the DNS zone. This TXT record contains IP addresses or networks. This addresses are alllows to send emails for the domain. 

More information about SPF (Sender Policy Framework): 
- [https://en.wikipedia.org/wiki/Sender_Policy_Framework](https://en.wikipedia.org/wiki/Sender_Policy_Framework)
- [https://datatracker.ietf.org/doc/rfc7208/](https://datatracker.ietf.org/doc/rfc7208/)

Example: 
goSPF needs to be always started with a domain. 

	goSPF google.com
	SPF-Record: v=spf1 include:_spf.google.com ~all

goSPF will just print the found SPF record. Depending on whether the program finds an SPF record, it exits with an adequate exit code. The exit code is '0' for a valid SPF record and '2' if no SPF was found or the SPF record is invalid. 

A SPF record is valid when: 
- It is a TXT record 
- Is the only SPF record for this FQDN 
- Starts with v=spf1
- Ends with -all, ~all or ?all




 
  
