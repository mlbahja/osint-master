package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type IPInfo struct {
	// ip-api.com fields
	Status      string  `json:"status"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"regionName"`
	AS          string  `json:"as"`
	ASN         string  `json:"asn"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Query       string  `json:"query"`
	Timezone    string  `json:"timezone"`
	Zip         string  `json:"zip"`

	// AbuseIPDB fields (for known issues)
	AbuseConfidence int    `json:"abuseConfidenceScore,omitempty"`
	TotalReports    int    `json:"totalReports,omitempty"`
	LastReported    string `json:"lastReportedAt,omitempty"`
	UsageType       string `json:"usageType,omitempty"`
	Domain          string `json:"domain,omitempty"`
}

// GetIPInfo fetches geolocation and ISP information for an IP address
func GetIPInfo(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,countryCode,region,regionName,city,zip,lat,lon,timezone,isp,org,as,asname,query", ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IP data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var info IPInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if info.Status != "success" {
		return nil, fmt.Errorf("API error: %s", info.Status)
	}

	return &info, nil
}

// CheckAbuseDB checks if IP has been reported for abuse (using free AbuseIPDB API)
func CheckAbuseDB(ip string) (string, error) {
	// Note: AbuseIPDB requires an API key for free tier
	// This is a placeholder - you can sign up for free key at https://www.abuseipdb.com/
	apiKey := "" // Add your API key here or read from env

	if apiKey == "" {
		return "API key required for abuse database (get free key from abuseipdb.com)", nil
	}

	url := fmt.Sprintf("https://api.abuseipdb.com/api/v2/check?ipAddress=%s&maxAgeInDays=90", ip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Key", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse response (simplified)
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if data, ok := result["data"].(map[string]interface{}); ok {
		if confidence, ok := data["abuseConfidenceScore"].(float64); ok && confidence > 0 {
			return fmt.Sprintf("IP has abuse reports (confidence: %.0f%%)", confidence), nil
		}
	}

	return "No reported abuse in last 90 days", nil
}

// HandleIP is the main function for IP lookup
func HandleIP(ip string) string {
	// Basic IP validation
	if !isValidIP(ip) {
		return fmt.Sprintf("Error: '%s' is not a valid IP address\n", ip)
	}

	output := fmt.Sprintf("IP Address: %s\n", ip)
	output += strings.Repeat("=", 50) + "\n\n"

	// Get IP information
	output += "[INFO] Fetching geolocation and ISP data...\n\n"

	info, err := GetIPInfo(ip)
	if err != nil {
		return fmt.Sprintf("Error: %v\n", err)
	}

	// Display basic information
	output += "📍 GEOGRAPHIC INFORMATION:\n"
	output += fmt.Sprintf("  - Country: %s (%s)\n", info.Country, info.CountryCode)
	if info.Region != "" {
		output += fmt.Sprintf("  - Region/State: %s\n", info.Region)
	}
	if info.City != "" {
		output += fmt.Sprintf("  - City: %s\n", info.City)
	}
	if info.Zip != "" {
		output += fmt.Sprintf("  - Postal Code: %s\n", info.Zip)
	}
	if info.Timezone != "" {
		output += fmt.Sprintf("  - Timezone: %s\n", info.Timezone)
	}
	output += fmt.Sprintf("  - Coordinates: %.4f, %.4f\n", info.Lat, info.Lon)

	output += "\n🌐 NETWORK INFORMATION:\n"
	if info.ISP != "" {
		output += fmt.Sprintf("  - ISP: %s\n", info.ISP)
	}
	if info.Org != "" && info.Org != info.ISP {
		output += fmt.Sprintf("  - Organization: %s\n", info.Org)
	}
	if info.AS != "" {
		output += fmt.Sprintf("  - ASN: %s\n", info.AS)
	}

	output += "\n🔍 ADDITIONAL INFORMATION:\n"
	output += fmt.Sprintf("  - Google Maps: https://www.google.com/maps?q=%.4f,%.4f\n", info.Lat, info.Lon)
	output += fmt.Sprintf("  - IP Info: https://ipinfo.io/%s\n", ip)
	output += fmt.Sprintf("  - VirusTotal: https://www.virustotal.com/gui/ip-address/%s\n", ip)

	// Check for abuse reports (optional - requires API key)
	output += "\n⚠️  ABUSE REPORT STATUS:\n"
	abuseStatus, err := CheckAbuseDB(ip)
	if err != nil {
		output += fmt.Sprintf("  - Abuse check failed: %v\n", err)
	} else {
		output += fmt.Sprintf("  - %s\n", abuseStatus)
	}

	// Additional OSINT tips
	output += "\n" + strings.Repeat("=", 50) + "\n"
	output += "OSINT TIPS:\n"
	output += "  - Check if this IP appears in data breach logs\n"
	output += "  - Use Shodan to find open ports/services (shodan.io)\n"
	output += "  - Check reverse DNS for domain associations\n"
	output += "  - Search IP in threat intelligence feeds\n"

	return output
}

// isValidIP performs basic IP address validation
func isValidIP(ip string) bool {
	// Simple validation - checks for IPv4 format
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false

			
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}
