package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type responseWriter struct{}

func (responseWriter) Write(bs []byte) (int, error) {
	fmt.Printf(string(bs))

	return len(bs), nil
}

func retreive(uri string) {
	resp, err := http.Get(uri)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	rw := responseWriter{}
	io.Copy(rw, resp.Body)
}
