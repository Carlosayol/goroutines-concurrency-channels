package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 2)
	exit := make(chan struct{})

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(time.Now(), i, "sending")
			ch <- i
			fmt.Println(time.Now(), i, "sent")

			time.Sleep(1 * time.Second)
		}

		fmt.Println(time.Now(), "all completed")
		close(ch)
	}()

	go func() {
		// we use select when there are more than 1 channel
		// for {
		// 	select {
		// 	case v, open := <-ch:
		// 		if !open {
		// 			close(exit)
		// 			return
		// 		}

		// 		fmt.Println(time.Now(), v, "received")
		// 	}
		// }

		// we use range when there is only 1 channel
		for v := range ch {
			fmt.Println(time.Now(), v, "received")
		}

		close(exit)
	}()

	fmt.Println(time.Now(), "waiting for completion")
	<-exit
	fmt.Println(time.Now(), "exiting")
}
