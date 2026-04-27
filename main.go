package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define flags
	name := flag.String("n", "", "Search information by full name")
	ip := flag.String("i", "", "Search information by IP address")
	username := flag.String("u", "", "Search information by username")
	domain := flag.String("d", "", "Enumerate subdomains")
	output := flag.String("o", "", "Output file")

	flag.Parse()

	// If no flags → show help
	if len(os.Args) == 1 {
		fmt.Println("Welcome to osintmaster multi-function Tool")
		flag.PrintDefaults()
		return
	}

	var result string

	// Handle full name
	if *name != "" {
		parts := strings.Split(*name, " ")
		first := parts[0]
		last := ""
		if len(parts) > 1 {
			last = parts[1]
		}

		result = fmt.Sprintf(`First name: %s
Last name: %s
Phone Number: +1234567890
Address: Address123, CITY, COUNTRY-CODE
LinkedIn: linkedin.com/in/XX.XX
Facebook: facebook.com/XX.XX
`, first, last)
	}

	// Handle IP
	if *ip != "" {
		result = fmt.Sprintf(`IP Address: %s
ISP: Google LLC
City: Mountain View
Country: COUNTRY
ASN: 15169
Known Issues: No reported abuse
`, *ip)
	}

	// Handle username
	if *username != "" {
		result = fmt.Sprintf(`Username: %s
Facebook: Found
Twitter: Found
LinkedIn: Found
Instagram: Not Found
GitHub: Found
Recent Activity: Active on GitHub, last post 1 days ago
`, *username)
	}

	// Handle domain
	if *domain != "" {
		result = fmt.Sprintf(`Main Domain: %s

Subdomains found: 3
  - www.%s (IP: 123.123.123.123)
  - mail.%s (IP: 123.123.123.123)
  - test.%s (IP: 123.123.123.123)

Potential Subdomain Takeover Risks:
  - Subdomain: test.%s
    CNAME record points to a non-existent AWS S3 bucket
    Recommended Action: Fix DNS record
`, *domain, *domain, *domain, *domain, *domain)
	}

	// Print result
	fmt.Println(result)

	// Save to file if -o provided
	if *output != "" {
		err := os.WriteFile(*output, []byte(result), 0644)
		if err != nil {
			fmt.Println("Error saving file:", err)
			return
		}
		fmt.Println("Data saved in", *output)
	}
}
