package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	typeFlag := flag.String("type", "", "commit type")
	scope := flag.String("scope", "", "commit scope (optional)")
	description := flag.String("description", "", "commit description")
	body := flag.String("body", "", "commit body (optional)")
	footer := flag.String("footer", "", "commit footer (optional)")

	flag.Parse()

	for *typeFlag == "" {
		*typeFlag = readInput("Enter commit type: ")
	}
	for *description == "" {
		*description = readInput("Enter commit description: ")
	}

	// Optional flags
	if *scope == "" {
		*scope = readInput("Enter commit scope (optional): ")
	}
	if *body == "" {
		*body = readInput("Enter commit body (optional): ")
	}
	if *footer == "" {
		*footer = readInput("Enter commit footer (optional): ")
	}

	var commitMessage string
	if *scope != "" {
		commitMessage = fmt.Sprintf("%s(%s): %s\n\n%s\n\n%s", *typeFlag, *scope, *description, *body, *footer)
	} else {
		commitMessage = fmt.Sprintf("%s: %s\n\n%s\n\n%s", *typeFlag, *description, *body, *footer)
	}
	lines := strings.Split(commitMessage, "\n")
	nonBlankLines := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonBlankLines = append(nonBlankLines, line)
		}
	}
	commitMessage = strings.Join(nonBlankLines, "\n")
	fmt.Println("Commit message: ", commitMessage)

	addCmd := exec.Command("git", "add", ".")
	addCmd.Stderr = os.Stderr
	addCmd.Stdout = os.Stdout
	err := addCmd.Run()
	if err != nil {
		fmt.Println("Error adding changes: ", err)
		return
	}

	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error committing: ", err)
	}
}

// Specs
// Command for creating a commit, command could even be called gcommit or whatever
// Have flags for or ask for:
// 			type
// 			scope (optional)
// 			description — concise
// 			body (optional)
// 			footer (optional)
// If the required type and description aren't present, error
// If everything's there, create the Git commit
