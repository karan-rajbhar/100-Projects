package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Spec struct {
	Value string `json:"value"`
	End   bool   `json:"end"`
}

func main() {
	specJSON := `{
		"!" : {
            "value" : "<strong>",
            "end" : true
        },
        "~" : {
            "value" : "<em>",
            "end" : true
        },
        "*" : {
            "value" : "<ul><li>",
            "end" : false
        },
        ")": {
            "value" : "<ol><li>",
            "end" : false
        },
        "^": {
            "value" : "<a href=\"",
            "end" : false
        },
        "{": {
            "value" : "\">",
            "end" : true
        }
	}`
	var spec map[string]Spec
	err := json.Unmarshal([]byte(specJSON), &spec)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a markdown file to convert.")
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	out, err := os.Create("output.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()
	// Write the start of the HTML body
	fmt.Fprintln(out, "<html><body>")
	scanner := bufio.NewScanner(file)
	reStart := regexp.MustCompile(`^\d+\)`)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		var builder strings.Builder
		// Check if the first character of the line is in the spec
		_, ok := spec[string(line[0])]
		// If it's not in the spec and the line does not start with an ordered list, start a <p> tag
		if !ok && !reStart.MatchString(line) {
			builder.WriteString("<p>")
		}
		insideTag := false
		tagValue := ""
		insideList := false
		// Create a stack to keep track of the opened tags
		var stack []string

		for i, runeValue := range line {
			char := string(runeValue)
			if value, ok := spec[char]; ok {
				if !insideTag {
					// Start a new tag
					builder.WriteString(value.Value)
					tagValue = value.Value
					insideTag = true

					// Push the tag to the stack
					stack = append(stack, tagValue)
				} else {
					// Pop the last tag from the stack
					lastTag := stack[len(stack)-1]
					stack = stack[:len(stack)-1]

					// End the last tag
					builder.WriteString("</" + strings.Trim(lastTag, "<>") + ">")
					insideTag = false
				}
			} else if char == ")" {
				// If we're not inside a list, start a new list
				if !insideList {
					builder.WriteString("<ol><li>")
					insideList = true
				} else {
					// If we're inside a list, start a new list item
					builder.WriteString("</li><li>")
				}
			} else {
				// If we're not inside a tag, append the character to the output
				builder.WriteString(char)
			}

			// If we're at the end of a line, end the paragraph
			if i == len(line)-1 {
				builder.WriteString("</p>")
			}
		}

		// Close any remaining tags
		for len(stack) > 0 {
			lastTag := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			builder.WriteString("</" + strings.Trim(lastTag, "<>") + ">")
		}

		// If we're inside a list, close it
		if insideList {
			builder.WriteString("</li></ol>")
		}

		fmt.Fprintln(out, builder.String())
	}

	fmt.Fprintln(out, "</body></html>")
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
