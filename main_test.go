package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestCSVReader(t *testing.T) {
	files, _ := ioutil.ReadDir("./testdata")
	paths := []string{}
	for _, file := range files {
		filepath := "./testdata/" + file.Name()
		paths = append(paths, filepath)
	}
	readCSVFiles(paths)
}

func TestByteReader(t *testing.T) {
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
}
