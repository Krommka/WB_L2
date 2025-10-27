package main

import (
	"fmt"
	"time"
)

func main() {
	or := orFunc
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func orFunc(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return channels[0]
	}

	res := make(chan interface{})

	go func() {
		defer close(res)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-orFunc(channels[2:]...):
		}
	}()
	return res
}
