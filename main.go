package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func getStdin(std io.Reader) string {
	stdin, err := ioutil.ReadAll(std)
	if err != nil {
		panic("Unable to read from stdin")
	}
	return string(stdin)
}

func main() {
	stdin := getStdin(os.Stdin)
	fmt.Println("got", stdin)
}
