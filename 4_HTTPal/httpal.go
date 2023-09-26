package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: Please provide URL.")
		os.Exit(1)
	}

	URL := args[1]

	// Make the HTTP GET request
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	err = os.WriteFile("response.txt", body, 0644)
	if err != nil {
		fmt.Println("Error saving response to file:", err)
		return
	}

}
