package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type responseWriter struct{}

func main() {

	resp, err := http.Get("http://www.sfeir.com")

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rw := responseWriter{}
	io.Copy(rw, resp.Body)
}

func (responseWriter) Write(bs []byte) (int, error) {
	fmt.Printf(string(bs))

	return len(bs), nil
}
