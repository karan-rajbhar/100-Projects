package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 4 {
		fmt.Println("Error: Please provide three arguments: old string, new string, and file path.")
		os.Exit(1)
	}

	oldStr := args[1]
	newStr := args[2]
	filePath := args[3]

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	newContent := strings.Replace(string(content), oldStr, newStr, -1)

	err = ioutil.WriteFile(filePath, []byte(newContent), 0)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	fmt.Println("File updated successfully.")
}
