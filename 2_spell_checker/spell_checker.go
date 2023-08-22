package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/makifdb/spellcheck"
)

func main() {
	args := os.Args[1:]
	// Init spellchecker
	sc, err := spellcheck.New()
	if err != nil {
		fmt.Println(err)
	}
	var typoList []string
	if len(args) > 0 {
		file, err := os.Open(args[0])

		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		// Define a regular expression to match special characters.
		regex := regexp.MustCompile("[^a-zA-Z0-9]+")
		linePosition := 1
		for scanner.Scan() {
			line := scanner.Text()
			words := strings.Fields(line)
			wordPosition := 1
			for _, word := range words {

				word = strings.ToLower(word)
				word = regex.ReplaceAllString(word, "")
				ok := sc.SearchDirect(word)
				if !ok {
					fmt.Println("-Line ", linePosition, ", Col ", wordPosition, ":", word, "appears to be a typo")

					typoList = append(typoList, fmt.Sprintf("-Lines %d , Col %d: \"%s\" appears to be a typo", linePosition, wordPosition, word))
				}
				wordPosition++
			}
			linePosition++

		}
	}

	if len(typoList) == 0 {
		return
	} else {
		fmt.Println("Typos Found:")
		for _, typo := range typoList {
			fmt.Println(typo)
		}
	}

}
