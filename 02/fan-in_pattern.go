package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	ch1, err := read("file1.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file 1 %v", err))
	}

	ch2, err := read("file2.csv")
	if err != nil {
		panic(fmt.Errorf("Could not read file 2 %v", err))
	}

	exit := make(chan struct{})
	chM := merge1(ch1, ch2)

	go func() {
		for v := range chM {
			fmt.Println(v)
		}

		close(exit)
	}()

	<-exit
	fmt.Println("Process completed")
}

func read(file string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file %v", err)
	}

	ch := make(chan []string)
	cr := csv.NewReader(f)

	go func() {
		for {
			record, err := cr.Read()
			if err == io.EOF {
				close(ch)
				return
			}

			ch <- record
		}
	}()

	return ch, nil
}

func merge1(cs ...<-chan []string) <-chan []string {
	var wg sync.WaitGroup
	out := make(chan []string)
	send := func(c <-chan []string) {
		for n := range c {
			out <- n
		}

		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()

		close(out)
	}()

	return out
}
