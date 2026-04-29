package main

import (
	"flag"
	"fmt"
	"os"
	"osint-master/handlers"
)

func main() {
	// Custom help
	flag.Usage = func() {
		fmt.Println("Welcome to osintmaster multi-function Tool\n")
		fmt.Println("OPTIONS:")
		fmt.Println(`  -n  "Full Name"        Search information by full name`)
		fmt.Println(`  -i  "IP Address"       Search information by IP address`)
		fmt.Println(`  -u  "Username"         Search information by username`)
		fmt.Println(`  -g  "Full Name"        Generate and test usernames from full name`)
		fmt.Println(`  -d  "Domain"           Enumerate subdomains`)
		fmt.Println(`  -o  "FileName"         File name to save output`)
	}

	name := flag.String("n", "", "")
	ip := flag.String("i", "", "")
	username := flag.String("u", "", "")
	generate := flag.String("g", "", "") // This is the flag for generate and test
	domain := flag.String("d", "", "")
	output := flag.String("o", "", "")

	flag.Parse()

	// Show help if no args
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	// Ensure only one main option is used
	options := 0
	if *name != "" {
		options++
	}
	if *ip != "" {
		options++
	}
	if *username != "" {
		options++
	}
	if *generate != "" {
		options++
	}
	if *domain != "" {
		options++
	}

	if options == 0 {
		fmt.Println("Error: You must provide one of -n, -i, -u, -g, or -d")
		return
	}

	if options > 1 {
		fmt.Println("Error: Use only ONE of -n, -i, -u, -g, or -d")
		return
	}

	// Dispatch
	var result string

	switch {
	case *name != "":
		result = handlers.HandleNameScraper(*name)

	case *ip != "":
		result = handlers.HandleIP(*ip)

	case *username != "":
		result = handlers.HandleUsername(*username)

	case *generate != "":
		// This will generate AND test usernames
		result = handlers.HandleGenerateAndTestUsernames(*generate)

	case *domain != "":
		result = handlers.HandleDomain(*domain)
	}

	// Output
	fmt.Println(result)

	if *output != "" {
		handlers.SaveToFile(*output, result)
	}
}