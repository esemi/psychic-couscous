package database

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

type db struct {
	neo neo4j.Driver
}
