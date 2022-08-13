package main

import (
	"net/netip"
)

func compareIPAddr(mergesIPS []string, UserIPs2Check []string) []string {

	for _, IP := range mergesIPS {
		for index, checkIP := range UserIPs2Check {
			if IP == checkIP {
				//fmt.Println("IP match!!!!!")
				//cIP = append(s[:index], s[index+1:]
				UserIPs2Check = append(UserIPs2Check[:index], UserIPs2Check[index+1:]...)
			}
		}
	}

	return UserIPs2Check
}

func checkNetworks(ip4nets []string, UserIPs2Check []string) []string {

	for _, IP := range ip4nets {
		network, _ := netip.ParsePrefix(IP)

		for index, checkIP := range UserIPs2Check {
			testIP, _ := netip.ParseAddr(checkIP)

			if network.Contains(testIP) {
				UserIPs2Check = append(UserIPs2Check[:index], UserIPs2Check[index+1:]...)
			}
		}
	}

	return UserIPs2Check
}

//exampleIP, _ := netip.ParseAddr("52.82.175.255")

/*
	fmt.Println("Check if IP is in CIDR Network")
	for _, x := range ip4nets {
		network, _ := netip.ParsePrefix(x)

		if network.Contains(exampleIP) {
			fmt.Println(exampleIP, "is part from the network", x)
		}
	}
*/
