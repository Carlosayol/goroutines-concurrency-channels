package main

import "fmt"

func main() {
	var ch chan int
	ch = make(chan int)
	// ch := make(chan int)

	ch <- 10

	v := <-ch
	fmt.Println("received: ", v)
}
