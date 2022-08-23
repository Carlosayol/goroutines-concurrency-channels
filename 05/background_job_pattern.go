package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Process ID", os.Getegid())

	listenForWork()
	<-waitToExit()

	fmt.Println("Process completed")
}

func listenForWork() {
	const workersN int = 5

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM)

	workersC := make(chan struct{}, workersN)

	go func() {
		for {
			<-sc

			workersC <- struct{}{}
		}
	}()

	go func() {
		var workers int

		for range workersC {
			workerID := (workers % workersN) + 1
			workers++

			fmt.Printf("%d<-\n", workerID)
			go func() {
				doWork(workerID)
			}()
		}
	}()
}

func doWork(id int) {
	fmt.Printf("<-%d starting\n", id)
	time.Sleep(3 * time.Second)
	fmt.Printf("<-%d completed\n", id)
}

func waitToExit() <-chan struct{} {
	runC := make(chan struct{}, 1)
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, os.Interrupt)
	go func() {
		defer close(runC)

		<-sc
	}()

	return runC
}
