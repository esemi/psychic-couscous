package database

import (
	"fmt"

	neo "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sarulabs/di/v2"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/neo4j"
)

type personRepository struct {
	db
}

// compile time check.
var _ domain.PersonRepository = (*personRepository)(nil)

// DefPersonsRepositoryName is definition name.
const DefPersonsRepositoryName = "personRepository"

// DefPersonsRepository is definition getter.
func DefPersonsRepository() di.Def {
	return di.Def{
		Name: DefPersonsRepositoryName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var neo = ctn.Get(neo4j.BundleName).(neo4j.Driver)

			return &personRepository{db{neo: neo}}, nil
		},
	}
}

func (p *personRepository) GetByID(id domain.PersonID) (domain.Person, error) {
	panic("implement me")
}

func (p *personRepository) GetByName(name string) (domain.Person, error) {
	panic("implement me")
}

func (p *personRepository) Save(person domain.Person) (err error) {
	session := p.neo.NewSession(neo.SessionConfig{})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo.Transaction) (interface{}, error) {
		_, err := tx.Run(`CREATE 
(p:Person { 
	id: $id, 
	name: $name, 
	birthY: $birthY, 
	deathY: $deathY, 
	professions: $professions, 
	knownForTitles: $knownForTitles 
})`,
			map[string]interface{}{
				"id":             person.ID,
				"name":           person.PrimaryName,
				"birthY":         person.BirthYear,
				"deathY":         person.DeathYear,
				"professions":    person.PrimaryProfessions,
				"knownForTitles": person.KnownForTitles,
			},
		)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return
}

func (p *personRepository) LoadFromCSV(filename string) (err error) {
	var session = p.neo.NewSession(neo.SessionConfig{})
	_, err = session.Run(fmt.Sprintf(
		`USING PERIODIC COMMIT 500 LOAD CSV WITH HEADERS FROM 'file:///%s' AS line FIELDTERMINATOR '\t'
	CREATE (p:Person { 
		id: line.nconst, 
		name: line.primaryName, 
		birthY: toInteger(line.birthYear), 
		deathY: toIntegerOrNull(line.deathYear), 
		professions: split(line.primaryProfession, ','), 
		knownForTitles: split(line.knownForTitles, ',') 
	})`, filename), nil)

	return err
}

func (p *personRepository) LoadRelationsFromCSV(filename string) (err error) {
	var session = p.neo.NewSession(neo.SessionConfig{})
	_, err = session.Run(fmt.Sprintf(
		`USING PERIODIC COMMIT 500 LOAD CSV WITH HEADERS FROM 'file:///%s' AS line FIELDTERMINATOR '\t'
	UNWIND split(line.knownForTitles, ',') AS kft
	MATCH 
		(person:Person {id: line.nconst}),
		(movie:Movie {id: kft})
	CREATE (person)-[:KNOWN_FOR]->(movie)
`, filename), nil)

	return err
}
