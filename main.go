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
	ev := make(chan link)
	r := retriever{visited: make(map[string]bool)}
	r.addEvent("newLink", ev)

	go func() {
		for {
			l := <-ev
			fmt.Println(l.source + " -> " + l.target)
		}
	}()

	r.crawl(url)

}
