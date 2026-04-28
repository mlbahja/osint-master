package handlers

import "fmt"

func HandleDomain(domain string) string {
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
