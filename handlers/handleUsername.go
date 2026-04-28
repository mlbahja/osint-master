package handlers

import "fmt"

func HandleUsername(username string) string {
	fmt.Println("this is the username : ")
	return fmt.Sprintf(`Username: %s
Facebook: Found
Twitter: Found
LinkedIn: Found
Instagram: Not Found
GitHub: Found
Recent Activity: Active on GitHub, last post 1 days ago
`, username)
}
