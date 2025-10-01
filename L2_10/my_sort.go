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
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(lines) == 0 {
		return
	}

	if config.CheckSorted {
		if sorting.IsSorted(lines, config) {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "sort: disorder in input\n")
			os.Exit(1)
		}
	}

	sorted := sorting.SortLines(lines, config)

	for _, line := range sorted {
		fmt.Println(line)
	}

}
