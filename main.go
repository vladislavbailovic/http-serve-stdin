package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getStdin(std io.Reader) string {
	stdin, err := ioutil.ReadAll(std)
	if err != nil {
		panic("Unable to read from stdin")
	}
	return string(stdin)
}

func getStdinHandler(std io.Reader) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, getStdin(std))
	}
}

func main() {
	http.HandleFunc("/", getStdinHandler(os.Stdin))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
