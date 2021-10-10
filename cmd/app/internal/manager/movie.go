package manager

import (
	neo "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type movieManager struct {
	neoClient neo.Session
}

func (m *movieManager) Save() error {

	return nil
}
