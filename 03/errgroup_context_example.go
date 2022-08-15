package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	filePath = "../files/"
)

func main() {
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	wait := errGroupContext(ctx)

	<-wait
}

func errGroupContext(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{}, 1)
	g, ctx := errgroup.WithContext(ctx)

	for _, file := range []string{filePath + "file1.csv", filePath + "file2.csv", filePath + "file3.csv"} {
		file := file

		g.Go(func() error {
			ch, err := readTimeout(file)
			if err != nil {
				return fmt.Errorf("error reading %v", err)
			}

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Context completed %v\n", ctx.Err())

					return ctx.Err()
				case line, ok := <-ch:
					if !ok {
						return nil
					}

					fmt.Println(line)
				}
			}
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			fmt.Printf("Error reading files %v", err)
		}

		close(ch)
	}()

	return ch
}

func readTimeout(file string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening file %w", err)
	}

	ch := make(chan []string)

	go func() {
		cr := csv.NewReader(f)
		time.Sleep(time.Millisecond)
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
