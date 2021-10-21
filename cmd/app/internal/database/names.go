package database

import (
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

// DefNamesRepositoryName is definition name.
const DefNamesRepositoryName = "personRepository"

// DefNamesRepository is definition getter.
func DefNamesRepository() di.Def {
	return di.Def{
		Name: DefNamesRepositoryName,
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

func (p *personRepository) Save(name domain.Person) (err error) {
	session := p.neo.NewSession(neo.SessionConfig{})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo.Transaction) (interface{}, error) {
		_, err := tx.Run("CREATE (p:Person { id: $id, name: $name, birthY: $birthY, deathY: $deathY, professions: $professions, knownForTitles: $knownForTitles })", map[string]interface{}{
			"id":             name.ID,
			"name":           name.PrimaryName,
			"birthY":         name.BirthYear,
			"deathY":         name.DeathYear,
			"professions":    name.PrimaryProfessions,
			"knownForTitles": name.KnownForTitles,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return
}
