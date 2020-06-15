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

	links, err := retrieve(url)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, link := range links {
		fmt.Println(link)
	}
}
