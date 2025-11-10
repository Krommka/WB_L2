package main

import (
	"fmt"
	"main/parsingFlags"
	readinput "main/readInput"
	"main/sorting"
	"os"
)

func main() {
	config, files, err := parsingflags.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	lines, err := readinput.ReadInput(files)
	if err != nil {
		if err != readinput.ErrLargeFile {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		external.Sorting(config)
	}

	sorting.InternalSort(lines, config)

}
