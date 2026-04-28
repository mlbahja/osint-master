package handlers

import (
	"fmt"
	"os"
)

func SaveToFile(filename, data string) {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
	fmt.Println("Data saved in", filename)
}
