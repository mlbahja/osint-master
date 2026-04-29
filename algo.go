package main

import (
	"fmt"
	"strings"
)

func UsernameGenrate(s string) []string {
	var result []string
	parts := strings.Split(s, " ")
	if len(parts) > 1 {
		for _, c := range parts {
			fmt.Println(c)
		}
			
	}


	return result

}
func main() {

}