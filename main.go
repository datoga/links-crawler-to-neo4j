package main

import (
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Web site url is missing")
		os.Exit(1)
	}

	driver, session, connErr := connectToNeo4j()

	if connErr != nil {
		fmt.Println("Error connecting to Database:", connErr)
		os.Exit(1)
	}

	defer driver.Close()

	defer session.Close()

	url := os.Args[1]
	ev := make(chan link)
	r := retriever{visited: make(map[string]bool)}
	r.addEvent("newLink", ev)

	go func(session *neo4j.Session) {
		for {
			l := <-ev
			fmt.Println(l.source + " -> " + l.target)
			_, err := createNode(session, &l)

			if err != nil {
				fmt.Println("Failed to create node:", err)
			}

		}
	}(&session)

	r.crawl(url)

	fmt.Println("Creation of relationship between nodes.. ")
	_, qErr := createNodesRelationship(&session)

	if qErr == nil {
		fmt.Println("Nodes updated")
	} else {
		fmt.Println("Error while updating nodes:", qErr)
	}

}
