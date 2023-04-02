package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	// fmt.Println("Command-line arguments:", args[0])

	if len(args) > 0 {

		// Open the file for reading
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// Create a scanner to read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Print each line of the file
			fmt.Println(scanner.Text())
		}

		// Check for any errors during scanning
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("File Not provided")
	}
}
