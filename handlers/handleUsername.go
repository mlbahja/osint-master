package handlers

import (
	"fmt"
	"net/http"
	"osint-master/help"
	"strings"
	"time"
)

type Platform struct {
	Name     string
	URL      string
	Category string
}

func GetPlatforms() []Platform {
	return []Platform{
		{Name: "GitHub", URL: "https://github.com/%s", Category: "Development"},
		{Name: "Twitter", URL: "https://twitter.com/%s", Category: "Social"},
		{Name: "Instagram", URL: "https://instagram.com/%s", Category: "Social"},
		{Name: "Reddit", URL: "https://reddit.com/user/%s", Category: "Social"},
		{Name: "LinkedIn", URL: "https://linkedin.com/in/%s", Category: "Professional"},
		{Name: "Medium", URL: "https://medium.com/@%s", Category: "Blogging"},
		{Name: "GitLab", URL: "https://gitlab.com/%s", Category: "Development"},
		{Name: "Pinterest", URL: "https://pinterest.com/%s", Category: "Social"},
	}
}

type CheckResult struct {
	Platform string
	URL      string
	Username string
	Found    bool
	Category string
}

func checkUsernameOnPlatform(username string, platform Platform, client *http.Client) CheckResult {
	url := fmt.Sprintf(platform.URL, username)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return CheckResult{
			Platform: platform.Name,
			URL:      url,
			Username: username,
			Found:    false,
			Category: platform.Category,
		}
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{
			Platform: platform.Name,
			URL:      url,
			Username: username,
			Found:    false,
			Category: platform.Category,
		}
	}
	defer resp.Body.Close()

	// Profile exists if status is 200 OK or 302 Found
	found := resp.StatusCode == 200 || resp.StatusCode == 302

	return CheckResult{
		Platform: platform.Name,
		URL:      url,
		Username: username,
		Found:    found,
		Category: platform.Category,
	}
}

// HandleUsername checks if a username exists on platforms
func HandleUsername(username string) string {
	username = strings.TrimPrefix(username, "@")

	if username == "" {
		return "Error: Please provide a valid username\n"
	}

	output := fmt.Sprintf("Checking username: %s\n", username)
	output += strings.Repeat("=", 50) + "\n\n"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	platforms := GetPlatforms()
	var foundResults []CheckResult

	output += "Searching platforms...\n\n"

	for _, platform := range platforms {
		result := checkUsernameOnPlatform(username, platform, client)
		if result.Found {
			foundResults = append(foundResults, result)
			output += fmt.Sprintf("[FOUND] %s: %s\n", result.Platform, result.URL)
		}
	}

	if len(foundResults) == 0 {
		output += "\nNo profiles found for this username\n"
	}

	output += "\n" + strings.Repeat("=", 50) + "\n"
	output += fmt.Sprintf("Total profiles found: %d/%d\n", len(foundResults), len(platforms))

	return output
}

// HandleUsernameFromName generates usernames from a name and checks them
func HandleUsernameFromName(fullName string) string {
	// First, generate possible usernames
	usernames := help.GenerateFromFullName(fullName)

	if len(usernames) == 0 {
		return "Error: Could not generate usernames from provided name\n"
	}

	output := fmt.Sprintf("Full Name: %s\n", fullName)
	output += fmt.Sprintf("Generated %d possible username variations\n", len(usernames))
	output += strings.Repeat("=", 50) + "\n\n"

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	platforms := GetPlatforms()
	allResults := make(map[string][]CheckResult) // username -> results

	// Check first 20 usernames to avoid rate limiting
	maxUsernames := 20
	if len(usernames) < maxUsernames {
		maxUsernames = len(usernames)
	}

	output += "Checking username variations on platforms...\n\n"

	for i := 0; i < maxUsernames; i++ {
		username := usernames[i]
		var foundForThisUser []CheckResult

		for _, platform := range platforms {
			result := checkUsernameOnPlatform(username, platform, client)
			if result.Found {
				foundForThisUser = append(foundForThisUser, result)
			}
		}

		if len(foundForThisUser) > 0 {
			allResults[username] = foundForThisUser
			output += fmt.Sprintf("\n[VALID USERNAME] %s\n", username)
			for _, result := range foundForThisUser {
				output += fmt.Sprintf("    - %s: %s\n", result.Platform, result.URL)
			}
		}
	}

	if len(allResults) == 0 {
		output += "\nNo valid usernames found from generated variations\n"
		output += "Try a different name or check common patterns manually\n"
	} else {
		output += "\n" + strings.Repeat("=", 50) + "\n"
		output += fmt.Sprintf("Found %d valid username(s) out of %d checked\n", len(allResults), maxUsernames)
	}

	output += "\nOSINT TIPS:\n"
	output += "  - Use the -u flag to check specific usernames\n"
	output += "  - Try common years (1985, 1990) or birth years\n"
	output += "  - Check email addresses derived from these usernames\n"

	return output
}

// Add this function to your existing handlers/handleUsername.go file
// This function will GENERATE usernames from help package AND TEST them

// HandleGenerateAndTestUsernames generates usernames from a full name and tests them
func HandleGenerateAndTestUsernames(fullName string) string {
	// Step 1: Generate usernames using your help package
	usernames := help.GenerateFromFullName(fullName)

	if len(usernames) == 0 {
		return fmt.Sprintf("Error: Could not generate usernames from '%s'\n", fullName)
	}

	output := fmt.Sprintf("Full Name: %s\n", fullName)
	output += fmt.Sprintf("Generated %d username variations\n", len(usernames))
	output += strings.Repeat("=", 60) + "\n\n"

	// Step 2: Test only the most relevant usernames (limit to 30 to avoid rate limiting)
	maxToCheck := 30
	if len(usernames) < maxToCheck {
		maxToCheck = len(usernames)
	}

	output += fmt.Sprintf("Testing %d most likely username variations on social platforms...\n", maxToCheck)
	output += strings.Repeat("-", 60) + "\n\n"

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	platforms := GetPlatforms()
	var validUsernames []string
	validCount := 0

	for i := 0; i < maxToCheck; i++ {
		username := usernames[i]
		foundOn := []string{}

		// Test this username on all platforms
		for _, platform := range platforms {
			url := fmt.Sprintf(platform.URL, username)
			req, err := http.NewRequest("HEAD", url, nil)
			if err != nil {
				continue
			}

			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

			resp, err := client.Do(req)
			if err != nil {
				continue
			}
			resp.Body.Close()

			if resp.StatusCode == 200 || resp.StatusCode == 302 {
				foundOn = append(foundOn, platform.Name)
			}
		}

		if len(foundOn) > 0 {
			validCount++
			validUsernames = append(validUsernames, username)
			output += fmt.Sprintf("\n[VALID] %s - Found on: %s\n", username, strings.Join(foundOn, ", "))
		} else {
			// Optional: show which ones are invalid (comment out to reduce spam)
			// output += fmt.Sprintf("[NOT FOUND] %s\n", username)
		}
	}

	// Summary
	output += "\n" + strings.Repeat("=", 60) + "\n"
	output += fmt.Sprintf("\nSUMMARY:\n")
	output += fmt.Sprintf("  - Total username variations generated: %d\n", len(usernames))
	output += fmt.Sprintf("  - Usernames tested: %d\n", maxToCheck)
	output += fmt.Sprintf("  - Valid usernames found: %d\n", validCount)

	if validCount > 0 {
		output += "\nVALID USERNAMES TO INVESTIGATE:\n"
		for i, username := range validUsernames {
			if i < 10 {
				output += fmt.Sprintf("  %d. %s\n", i+1, username)
			}
		}
	}

	output += "\n" + strings.Repeat("=", 60) + "\n"
	output += "OSINT TIPS:\n"
	output += "  - Use -u \"username\" to check a specific username in detail\n"
	output += "  - The most valuable usernames are those found on multiple platforms\n"
	output += "  - Check email patterns: username@gmail.com, username@outlook.com\n"

	return output
}
