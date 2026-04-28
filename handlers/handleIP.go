package handlers

import "fmt"

// La majuscule est obligatoire pour l'exportation
func HandleIP(ip string) string {
	fmt.Println("Traitement de l'IP...")
	return fmt.Sprintf(`IP Address: %s
ISP: Google LLC
City: Mountain View
Country: COUNTRY
ASN: 15169
Known Issues: No reported abuse
`, ip)
}
