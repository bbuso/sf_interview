package main

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://www.google.com/", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	response := ByteCompressor(resp.Body)
	fmt.Println(reflect.TypeOf(response))
	// test_err := os.WriteFile("./test.html.gz", response, 0644)
	// if test_err != nil {
	// 	log.Fatal(err)
	// }

	// readCSVFile("./testdata/TestProcessEligibleChannel2_0_TestProcessEligibleChannel2_0_CODES.csv")
	files, _ := ioutil.ReadDir("./testdata")
	paths := []string{}
	for _, file := range files {
		filepath := "./testdata/" + file.Name()
		paths = append(paths, filepath)
	}
	readCSVFiles(paths)

}
