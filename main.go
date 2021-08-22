package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
		for name, value := range getDefaultHeaders() {
			resp.Header().Set(name, value)
		}
		resp.WriteHeader(http.StatusOK)

		fmt.Fprintf(resp, getStdin(std))
	}
}

func getDefaultHeaders() map[string]string {
	return map[string]string{
		"content-type": "text/plain; charset=utf-8",
	}
}

func getParsedHeaders(raw []string) map[string]string {
	headers := make(map[string]string)
	for _, rawHeader := range raw {
		splits := strings.Split(rawHeader, ":")
		if len(splits) < 2 {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(splits[0]))
		value := strings.TrimSpace(strings.Join(splits[1:], ":"))
		headers[name] = value
	}
	return headers
}

func serveStdin(port int, std io.Reader) {
	http.HandleFunc("/", getStdinHandler(std))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	serveStdin(8080, os.Stdin)
}
