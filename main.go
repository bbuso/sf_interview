package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type BufferedResponse interface {
	io.Reader
}

func readCSVFiles(paths []string) {
	var wg sync.WaitGroup
	for _, filepath := range paths {
		wg.Add(1)
		go fmt.Println(readCSVFile(filepath, &wg))
	}
	wg.Wait()
}

func readCSVFile(filePath string, wg *sync.WaitGroup) bool {

	defer wg.Done()
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("unable to read file: "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = 0
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	set := make(map[string]bool)

	for _, character := range records {
		if set[character[1]] {
			return false
		}
		set[character[1]] = true
	}

	return true

}

func ByteCompressor(reader io.ReadCloser) BufferedResponse {
	read_bytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	var zipped bytes.Buffer
	writer := gzip.NewWriter(&zipped)
	writer.Write(read_bytes)
	writer.Close()
	r := bytes.NewReader(zipped.Bytes())

	return r
}

func main() {

}
