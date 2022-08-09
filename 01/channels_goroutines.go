package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		fmt.Println(time.Now(), "sleeping")

		time.Sleep(2 * time.Second)

		ch <- "Hello"
	}()

	fmt.Println(time.Now(), "waiting for message")

	v := <-ch

	fmt.Println(time.Now(), "message received ", v)
}
