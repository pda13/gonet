package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

// gonet dns-lookup --domain google.com

var dnsLookupCmd = &cobra.Command{
	Use:   "dns-lookup",
	Short: "Get domain information using DNS-lookup: A/AAAA, MX, CNAME",
	Run:   GetDnsLookupCommand(),
}

func init() {
	rootCmd.AddCommand(dnsLookupCmd)

	dnsLookupCmd.Flags().String("domain", "ya.ru", "domain")
}

func GetDnsLookupCommand() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		domain, err := cmd.Flags().GetString("domain")
		if err != nil {
			log.Fatal(err)
		}

		// (IPv4)
		ipv4Addrs, err := net.LookupIP(domain)
		if err != nil {
			fmt.Printf("error fetching IPv4: %v\n", err)
		} else {
			fmt.Println("IPv4:")
			for _, ip := range ipv4Addrs {
				if ip.To4() != nil {
					fmt.Println(ip)
				}
			}
		}

		// (IPv6)
		fmt.Println("\nIPv6:")
		for _, ip := range ipv4Addrs {
			if ip.To4() == nil {
				fmt.Println(ip)
			}
		}

		// MX (mail servers)
		mxRecords, err := net.LookupMX(domain)
		if err != nil {
			fmt.Printf("error fetching MX: %v\n", err)
		} else {
			fmt.Println("\nMX (mail-servers):")
			for _, mx := range mxRecords {
				fmt.Printf("Host: %s, Pref: %d\n", mx.Host, mx.Pref)
			}
		}

		// CNAME
		cname, err := net.LookupCNAME(domain)
		if err != nil {
			fmt.Printf("error fetching CNAME: %v\n", err)
		} else {
			fmt.Println("\nCNAME:")
			fmt.Println(cname)
		}

		// TXT
		txtRecords, err := net.LookupTXT(domain)
		if err != nil {
			fmt.Printf("error fetching TXT: %v\n", err)
		} else {
			fmt.Println("\nTXT:")
			for _, txt := range txtRecords {
				fmt.Println(txt)
			}
		}

		// NS (DNS-servers)
		nsRecords, err := net.LookupNS(domain)
		if err != nil {
			fmt.Printf("error fetching NS: %v\n", err)
		} else {
			fmt.Println("\nNS (DNS-servers):")
			for _, ns := range nsRecords {
				fmt.Println(ns.Host)
			}
		}
	}
}
