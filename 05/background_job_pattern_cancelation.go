package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Scheduler struct {
	workers   int
	msgC      chan struct{}
	signalC   chan os.Signal
	waitGroup sync.WaitGroup
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

func NewScheduler(workers, buffer int) *Scheduler {
	return &Scheduler{
		workers: workers,
		msgC:    make(chan struct{}, buffer),
		signalC: make(chan os.Signal, 1),
	}
}

func (s *Scheduler) ListenForWork() {
	go func() {
		signal.Notify(s.signalC, syscall.SIGTERM)

		for {
			<-s.signalC

			s.msgC <- struct{}{}
		}
	}()

	s.waitGroup.Add(s.workers)

	for i := 0; i < s.workers; i++ {
		i := i

		go func() {
			for {
				select {
				case _, open := <-s.msgC:
					if !open {
						fmt.Printf("%d closing\n", i+1)

						s.waitGroup.Done()
						return
					}

					fmt.Printf("%d<- Processing\n", i)
				}
			}
		}()
	}
}

func (s *Scheduler) Exit() {
	close(s.msgC)
	s.waitGroup.Wait()
}

func main() {
	fmt.Println("Process ID", os.Getpid())

	s := NewScheduler(5, 10)
	s.ListenForWork()

	fmt.Println("Starting")

	<-waitToExit()
	s.Exit()

	fmt.Println("Exiting")
}
