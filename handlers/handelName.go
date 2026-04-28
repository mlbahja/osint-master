package handlers

import (
	"fmt"
	"strings"
)

func HandleName(fullName string) string {
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
