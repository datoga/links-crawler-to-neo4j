package main

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func connectToNeo4j() (neo4j.Driver, neo4j.Session, error) {

	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	//driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth(
	//	"neo4j", "alice!in!wonderland", ""), configForNeo4j40)

	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.NoAuth(), configForNeo4j40)

	if err != nil {
		return nil, nil, err
	}

	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}
	session, err := driver.NewSession(sessionConfig)
	if err != nil {
		return nil, nil, err
	}

	return driver, session, nil
}

func createNode(session *neo4j.Session, l *Link) (neo4j.Result, error) {
	r, err := (*session).Run("CREATE (:WebLink{source: $source, target: $target}) ", map[string]interface{}{
		"source": l.source,
		"target": l.target,
	})

	if err != nil {
		return nil, err
	}

	return r, err
}

func createNodesRelationship(session *neo4j.Session) (neo4j.Result, error) {
	r, err := (*session).Run("MATCH (a:WebLink),(b:WebLink) WHERE a.target = b.source CREATE (a)-[r:point_to]->(b)", map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	return r, err
}
