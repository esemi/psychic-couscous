package database

import (
	"fmt"

	neo "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sarulabs/di/v2"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/neo4j"
)

type movieRepository struct {
	db
}

// compile time check.
var _ domain.MovieRepository = (*movieRepository)(nil)

// DefMoviesRepositoryName is definition name.
const DefMoviesRepositoryName = "movieRepository"

// DefMoviesRepository is definition getter.
func DefMoviesRepository() di.Def {
	return di.Def{
		Name: DefMoviesRepositoryName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var neo = ctn.Get(neo4j.BundleName).(neo4j.Driver)

			return &movieRepository{db{neo: neo}}, nil
		},
	}
}

func (p *movieRepository) GetByID(id domain.MovieID) (domain.Movie, error) {
	panic("implement me")
}

func (p *movieRepository) GetByName(name string) (domain.Movie, error) {
	panic("implement me")
}

func (p *movieRepository) Save(movie domain.Movie) (err error) {
	session := p.neo.NewSession(neo.SessionConfig{})
	defer session.Close()
	_, err = session.WriteTransaction(func(tx neo.Transaction) (interface{}, error) {
		_, err := tx.Run(`CREATE 
(p:Movie { 
	id: $id, 
	titleType: $titleType,
	primaryTitle: $primaryTitle,
	originalTitle: $originalTitle,
	isAdult: $isAdult,
	startYear: $startYear, 
	endYear: $endYear, 
	runtimeMinutes: $runtimeMinutes, 
	genres: $genres 
})`,
			map[string]interface{}{
				"id":             movie.ID,
				"titleType":      movie.TitleType,
				"primaryTitle":   movie.PrimaryTitle,
				"originalTitle":  movie.OriginalTitle,
				"isAdult":        movie.Adult,
				"startYear":      movie.StartYear,
				"endYear":        movie.EndYear,
				"runtimeMinutes": movie.RuntimeMinutes,
				"genres":         movie.Genres,
			},
		)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return
}

func (p *movieRepository) LoadFromCSV(filename string) (err error) {
	var session = p.neo.NewSession(neo.SessionConfig{})
	_, err = session.Run(fmt.Sprintf(
		`USING PERIODIC COMMIT 500 LOAD CSV WITH HEADERS FROM 'file:///%s' AS line FIELDTERMINATOR '\t'
	CREATE (p:Movie { 
		id: line.tconst, 
		titleType: line.titleType,
		primaryTitle: line.primaryTitle,
		originalTitle: line.originalTitle,
		isAdult: toBoolean(line.isAdult),
		startYear: toInteger(line.startYear), 
		endYear: toIntegerOrNull(line.endYear), 
		runtimeMinutes: toInteger(line.runtimeMinutes), 
		genres: split(line.genres, ',')
	})`, filename), nil)

	return err
}
