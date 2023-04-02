package main

import (
	"fmt"
	"os"
)

func main() {
	// Get the current working directory

	args := os.Args[1:]
	if len(args) > 0 {

		// Open the directory
		dir, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer dir.Close()

		// Read the directory entries
		entries, err := dir.Readdir(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the directory entries
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Println(entry.Name())
			} else {
				fmt.Println(entry.Name())
			}
		}
	} else {

		cwd, err := os.Getwd()
		// Open the directory
		dir, err := os.Open(cwd)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer dir.Close()

		// Read the directory entries
		entries, err := dir.Readdir(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the directory entries
		for _, entry := range entries {
			if entry.IsDir() {
				fmt.Println(entry.Name())
			} else {
				fmt.Println(entry.Name())
			}
		}
	}
}
