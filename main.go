package main

import (
	"flag"
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

func getStdinHandler(headers []string, std io.Reader) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		for name, value := range getHeaders(headers) {
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

func getHeaders(raw []string) map[string]string {
	headers := getDefaultHeaders()
	for name, value := range getParsedHeaders(raw) {
		headers[name] = value
	}
	return headers
}

type headerFlags []string

func (h headerFlags) String() string {
	return strings.Join(h, "|")
}
func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func getParsedFlags(cliFlags []string) (int, []string) {
	var port int
	var showUsage bool
	var headers headerFlags

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.BoolVar(&showUsage, "help", false, "show this help")

	flags.IntVar(&port, "port", 8080, "serve on this port")
	flags.IntVar(&port, "p", 8080, "serve on this port")

	flags.Var(&headers, "header", "additional header(s)")
	flags.Var(&headers, "h", "additional header(s)")

	err := flags.Parse(cliFlags)
	if err != nil {
		panic(err)
	}

	if showUsage {
		flags.Usage()
		os.Exit(0)
	}

	return port, headers
}

func serveStdin(port int, headers []string, std io.Reader) {
	http.HandleFunc("/", getStdinHandler(headers, std))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	port, headers := getParsedFlags(os.Args[1:])
	serveStdin(port, headers, os.Stdin)
}
