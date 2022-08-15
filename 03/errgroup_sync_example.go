package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	filePath = "../files/"
)

func main() {
	wait := waitGroups()

	<-wait
}

func waitGroups() <-chan struct{} {
	ch := make(chan struct{}, 1)
	var wg sync.WaitGroup

	for _, file := range []string{filePath + "file1.csv", filePath + "file2.csv", filePath + "file3.csv"} {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ch, err := read(file)
			if err != nil {
				fmt.Printf("error reading %v", err)
			}

			for line := range ch {
				fmt.Println(line)
			}
		}()
	}

	go func() {
		wg.Wait()

		close(ch)
	}()

	return ch
}

func read(file string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening file %w", err)
	}

	ch := make(chan []string)

	go func() {
		cr := csv.NewReader(f)

		for {
			record, err := cr.Read()
			if errors.Is(err, io.EOF) {
				close(ch)

				return
			}

			ch <- record
		}
	}()

	return ch, nil
}
