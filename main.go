package main

import (
	"flag"
	"fmt"
	"os"
	"osint-master/handleIP"
	handleusername "osint-master/handleUsername"
	"strings"
)

// ---------- Handlers ----------

func handleName(fullName string) string {
	parts := strings.Split(fullName, " ")
	first := parts[0]
	last := ""
	if len(parts) > 1 {
		last = parts[1]
	}

	return fmt.Sprintf(`First name: %s
Last name: %s
Phone Number: +1234567890
Address: Address123, CITY, COUNTRY-CODE
LinkedIn: linkedin.com/in/XX.XX
Facebook: facebook.com/XX.XX
`, first, last)
}

/*
func handleIP(ip string) string {
	return fmt.Sprintf(`IP Address: %s
ISP: Google LLC
City: Mountain View
Country: COUNTRY
ASN: 15169
Known Issues: No reported abuse
`, ip)
}
*/
/*
func handleUsername(username string) string {
	return fmt.Sprintf(`Username: %s
Facebook: Found
Twitter: Found
LinkedIn: Found
Instagram: Not Found
GitHub: Found
Recent Activity: Active on GitHub, last post 1 days ago
`, username)
}*/

func handleDomain(domain string) string {
	return fmt.Sprintf(`Main Domain: %s

Subdomains found: 3
  - www.%s (IP: 123.123.123.123)
  - mail.%s (IP: 123.123.123.123)
  - test.%s (IP: 123.123.123.123)

Potential Subdomain Takeover Risks:
  - Subdomain: test.%s
    CNAME record points to a non-existent AWS S3 bucket
    Recommended Action: Fix DNS record
`, domain, domain, domain, domain, domain)
}

// ---------- Utility ----------

func saveToFile(filename, data string) {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	fmt.Println("Data saved in", filename)
}

// ---------- Main ----------

func main() {
	// Custom help
	flag.Usage = func() {
		fmt.Println("Welcome to osintmaster multi-function Tool\n")
		fmt.Println("OPTIONS:")
		fmt.Println(`  -n  "Full Name"        Search information by full name`)
		fmt.Println(`  -i  "IP Address"       Search information by IP address`)
		fmt.Println(`  -u  "Username"         Search information by username`)
		fmt.Println(`  -d  "Domain"           Enumerate subdomains`)
		fmt.Println(`  -o  "FileName"         File name to save output`)
	}

	// Flags
	name := flag.String("n", "", "")
	ip := flag.String("i", "", "")
	username := flag.String("u", "", "")
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
	if *domain != "" {
		options++
	}

	if options == 0 {
		fmt.Println("Error: You must provide one of -n, -i, -u, or -d")
		return
	}

	if options > 1 {
		fmt.Println("Error: Use only ONE of -n, -i, -u, or -d")
		return
	}

	// Dispatch
	var result string

	switch {
	case *name != "":
		result = handleName(*name)

	case *ip != "":
		//result =HandleIP(*ip)
		handleIP.HandleIP(*ip)

	case *username != "":
		//result = handleUsername(*username)
		result = handleusername.HandleUsername(*username)

	case *domain != "":
		result = handleDomain(*domain)
	}

	// Output
	fmt.Println(result)

	if *output != "" {
		saveToFile(*output, result)
	}
}
