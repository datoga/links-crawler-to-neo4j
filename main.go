package main

import (
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {

	if *uriFlag == "" {
		fmt.Println("Web site url is missing")
		os.Exit(1)
	}

	if *levelsFlag <= 0 {
		fmt.Println("Levels flags must be greater than 0")
		os.Exit(1)
	}

	driver, session, connErr := connectToNeo4j()

	if connErr != nil {
		fmt.Println("Error connecting to Database:", connErr)
		os.Exit(1)
	}

	defer driver.Close()

	defer session.Close()

	r := NewRetriever(*levelsFlag, *asyncFlag)

	endCh := make(chan bool, 1)

	nodeIndex := 0

	go func(session *neo4j.Session) {
		for l := range r.Links {
			nodeIndex++

			fmt.Printf("%d: %s -> %s\n", nodeIndex, l.source, l.target)
			_, err := createNode(session, &l)

			if err != nil {
				fmt.Println("Failed to create node:", err)
			}
		}

		endCh <- true

	}(&session)

	r.Crawl(*uriFlag)

	<-endCh

	fmt.Println("Creation of relationship between nodes.. ")
	_, qErr := createNodesRelationship(&session)

	if qErr == nil {
		fmt.Printf("%d Nodes updated\n", nodeIndex)
	} else {
		fmt.Println("Error while updating nodes:", qErr)
	}
}
