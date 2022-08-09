package main

import "fmt"

func main() {
	go hello()
}

func hello() {
	fmt.Println("This message wont appear")
}
