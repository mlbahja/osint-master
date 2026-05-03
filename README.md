# 🔍 OSINT Master - Multi-Function OSINT Tool

## 📋 Overview

OSINT Master is a powerful command-line tool for passive reconnaissance using publicly available data. It performs comprehensive information gathering based on user inputs such as full names, IP addresses, usernames, and domains. This tool is designed for educational purposes to help understand Open-Source Intelligence (OSINT) techniques and their applications in cybersecurity.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 👤 **Full Name Search** | Extract name components and provide OSINT search links |
| 🌐 **IP Address Lookup** | Geolocation, ISP info, coordinates, and threat intelligence links |
| 👥 **Username Search** | Check username existence on 8+ social platforms |
| 🔄 **Username Generation** | Generate and test 80+ username variations from a full name |
| 🏢 **Subdomain Enumeration** | Discover subdomains and detect takeover risks |
| 💾 **File Output** | Save all results to a text file |

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher
- Internet connection for API calls

### Steps

```bash
# Clone the repository
git clone https://github.com/yourusername/osint-master.git
cd osint-master

# Initialize Go module
go mod init osint-master

# Install dependencies
go get github.com/PuerkitoBio/goquery

# Build the tool
go build -o osintmaster main.go

# Optional: Install globally
sudo cp osintmaster /usr/local/bin/