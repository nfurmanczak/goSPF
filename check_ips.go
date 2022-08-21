package main

import (
	"net/netip"
)

func compareIPAddr(mergesIPS []string, UserIPs2Check []string) []string {

	for _, IP := range mergesIPS {
		for index, checkIP := range UserIPs2Check {
			if IP == checkIP {
				UserIPs2Check[index] = "null"
			}
		}
	}

	notCovertIP := []string{}

	for _, IP := range UserIPs2Check {
		if IP != "null" {
			notCovertIP = append(notCovertIP, IP)
		}
	}

	return notCovertIP
}

func checkNetworks(ip4nets []string, UserIPs2Check []string) []string {

	for _, IP := range ip4nets {
		network, _ := netip.ParsePrefix(IP)

		for index, checkIP := range UserIPs2Check {
			testIP, _ := netip.ParseAddr(checkIP)

			if network.Contains(testIP) {
				UserIPs2Check[index] = "null"
			}
		}
	}

	notCovertIP := []string{}

	for _, IP := range UserIPs2Check {
		if IP != "null" {
			notCovertIP = append(notCovertIP, IP)
		}
	}

	return notCovertIP
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
