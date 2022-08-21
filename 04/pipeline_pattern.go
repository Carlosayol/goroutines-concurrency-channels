package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	filePath = "../files/"
)

func main() {
	recordsC, err := readCSV(filePath + "filePipeline.csv")
	if err != nil {
		fmt.Printf("fatal error %e", err)
	}

	for val := range sanitize(titleize(recordsC)) {
		fmt.Printf("%v\n", val)
	}
}

func readCSV(file string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening file %v", err)
	}

	ch := make(chan []string)

	go func() {
		cr := csv.NewReader(f)
		cr.FieldsPerRecord = 3

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

func sanitize(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			if len(val[0]) > 3 {
				fmt.Println("skipped ", val)
				continue
			}

			ch <- val
		}

		close(ch)
	}()

	return ch
}

func titleize(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			val[0] = strings.Title(val[0])
			val[1], val[2] = val[2], val[1]

			ch <- val
		}

		close(ch)
	}()

	return ch
}
