package main

import "fmt"

func main() {

	func() {
		defer fmt.Println("defer 1")
		defer fmt.Println("defer 2")
		fmt.Println("Hello World")
		panic("panic here")

	}()

}
