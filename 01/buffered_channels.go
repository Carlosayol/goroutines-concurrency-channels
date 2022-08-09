package main

import (
	"fmt"
	"time"
)

// a buffered channel is the same as a channel but with a given lenght of capacity

func main() {
	ch := make(chan int, 2)

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(time.Now(), i, "sending value")
			ch <- i
			fmt.Println(time.Now(), i, "sent")
		}

		fmt.Println(time.Now(), "all values sent")
	}()

	time.Sleep(2 * time.Second)

	fmt.Println(time.Now(), "waiting for messages")

	fmt.Println(time.Now(), "received", <-ch)
	fmt.Println(time.Now(), "received", <-ch)
	fmt.Println(time.Now(), "received", <-ch)

	fmt.Println(time.Now(), "end")
}
