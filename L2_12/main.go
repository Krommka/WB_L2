package main

import (
	"fmt"
	"os"
)

func main() {
	config, files, err := ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	lines, err := ReadInput(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(lines) == 0 {
		return
	}

	output, matchCount, err := Grep(lines, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, line := range output {
		fmt.Println(line)
	}

	if matchCount == 0 && !config.Count {
		os.Exit(1)
	}
}
