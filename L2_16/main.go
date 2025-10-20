package main

import (
	parsingflags "L2_16/parsingFlags"
	"fmt"
)

func main() {
	config, err := parsingflags.ParseFlags()
	if err != nil {
		panic(err)
	}
	fmt.Println(config)
}
