package main

import (
	"fmt"
	"os"
)

func main() {

	config, files, err := ParseFlags()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v", err)
		os.Exit(1)
	}

	lines, err := ReadInput(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
	}

	if len(lines) == 0 {
		return
	}

	output := Cut(lines, config)

	for _, line := range output {
		fmt.Println(line)
	}

}
