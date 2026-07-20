package handlers

import (
	"fmt"
	"os"
	"path/filepath"
)

func SaveToFile(filename, data string) {
	dir := filepath.Dir(filename)

	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	}

	err := os.WriteFile(filename, []byte(data), 0o644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("Data saved in", filename)
}




func SaveAsPdf(file, content string) {
	//pdf := fpdf.new("p")
}
