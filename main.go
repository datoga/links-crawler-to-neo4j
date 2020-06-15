package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Web site url is missing")
		os.Exit(1)
	}

	url := os.Args[1]

	retreive(url)
}
