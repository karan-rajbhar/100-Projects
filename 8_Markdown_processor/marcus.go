package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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

	// Write the opening HTML tags
	fmt.Fprintln(out, "<html>\n<body>")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			level := strings.Count(strings.SplitN(line, " ", 2)[0], "#")
			text := strings.TrimSpace(line[level:])
			fmt.Fprintf(out, "<h%d>%s</h%d>\n", level, text, level)
		} else {
			fmt.Fprintf(out, "<p>%s</p>\n", line)
		}
	}

	// Write the closing HTML tags
	fmt.Fprintln(out, "</body>\n</html>")

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
