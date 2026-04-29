package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type Platform struct {
	Name     string
	URL      string
	Category string
}

func GetPlatforms() []Platform {
	return []Platform{
		// Social Media
		{Name: "GitHub", URL: "https://github.com/%s", Category: "Development"},
		{Name: "Twitter", URL: "https://twitter.com/%s", Category: "Social"},
		{Name: "Instagram", URL: "https://instagram.com/%s", Category: "Social"},
		{Name: "Reddit", URL: "https://reddit.com/user/%s", Category: "Social"},
		{Name: "LinkedIn", URL: "https://linkedin.com/in/%s", Category: "Professional"},
		{Name: "Medium", URL: "https://medium.com/@%s", Category: "Blogging"},
		{Name: "GitLab", URL: "https://gitlab.com/%s", Category: "Development"},
		{Name: "Pinterest", URL: "https://pinterest.com/%s", Category: "Social"},
		{Name: "Tumblr", URL: "https://%s.tumblr.com", Category: "Blogging"},
		{Name: "YouTube", URL: "https://youtube.com/@%s", Category: "Social"},
	}
}

type CheckResult struct {
	Platform string
	URL      string
	Found    bool
	Category string
}

func checkPlatform(username string, platform Platform, client *http.Client) CheckResult {
	url := fmt.Sprintf(platform.URL, username)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return CheckResult{
			Platform: platform.Name,
			URL:      url,
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
			Found:    false,
			Category: platform.Category,
		}
	}
	defer resp.Body.Close()

	// Profile exists if status is 200 OK or 302 Found (redirect)
	found := resp.StatusCode == 200 || resp.StatusCode == 302

	return CheckResult{
		Platform: platform.Name,
		URL:      url,
		Found:    found,
		Category: platform.Category,
	}
}

func HandleUsername(username string) string {
	// Remove @ if present
	username = strings.TrimPrefix(username, "@")

	if username == "" {
		return "Error: Please provide a valid username\n"
	}

	// Basic username validation
	for _, r := range username {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.' && r != '_' && r != '-' {
			return fmt.Sprintf("Error: Username '%s' contains invalid characters\n", username)
		}
	}

	output := fmt.Sprintf("Username: %s\n", username)
	output += strings.Repeat("=", 50) + "\n\n"

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	platforms := GetPlatforms()
	var foundResults []CheckResult

	// Check each platform
	output += "Checking platforms...\n\n"

	for _, platform := range platforms {
		result := checkPlatform(username, platform, client)
		if result.Found {
			foundResults = append(foundResults, result)
		}
	}

	// Display results by category
	if len(foundResults) > 0 {
		output += fmt.Sprintf("FOUND %d profile(s):\n", len(foundResults))
		output += strings.Repeat("-", 40) + "\n\n"

		// Group by category
		categories := map[string][]CheckResult{}
		for _, result := range foundResults {
			categories[result.Category] = append(categories[result.Category], result)
		}

		for category, results := range categories {
			output += fmt.Sprintf("[%s]\n", category)
			for _, result := range results {
				output += fmt.Sprintf("  - %s: %s\n", result.Platform, result.URL)
			}
			output += "\n"
		}
	} else {
		output += "No profiles found for this username\n\n"
	}

	// Add summary
	output += strings.Repeat("=", 50) + "\n"
	output += fmt.Sprintf("Summary: %d/%d platforms checked\n", len(foundResults), len(platforms))

	// Add OSINT tips
	output += "\nOSINT Tips:\n"
	output += "  - Try username variations (add numbers, underscores)\n"
	output += "  - Search Google: site:github.com \"%s\"\n"
	output += "  - Check data breach databases for this username\n"

	// Legal disclaimer
	output += "\n" + strings.Repeat("=", 50) + "\n"
	output += "Legal: Only use for authorized security testing or your own accounts.\n"

	return output
}
