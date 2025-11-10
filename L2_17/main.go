package main

import (
	"log"
	"mytelnet/config"
	"mytelnet/telnet"
)

func main() {

	config := config.MustLoad()
	cl := telnet.New(config)

	if err := cl.Connect(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer cl.Cleanup()

	cl.Start()

}
