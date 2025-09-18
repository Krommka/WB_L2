package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"os"
)

func main() {
	ntpTime, err := ntp.Time("pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "NTP error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(ntpTime.Format("02.01.2006 Mon 15:04:05 MST"))

}
