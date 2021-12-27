package manager

import (
	"context"
	"fmt"
	"github.com/spf13/viper"

	gzViper "github.com/gozix/viper/v2"
	gzzap "github.com/gozix/zap/v2"
	"github.com/sarulabs/di/v2"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/database"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
	"go.uber.org/zap"
)

const personFilenameConfigPath = "files.persons"

type nameManager struct {
	logger   *zap.Logger
	repo     domain.PersonRepository
	filepath string
	config   *viper.Viper
}

//compile time check.
var _ domain.LoaderRelations = (*nameManager)(nil)

//DefPersonManagerName is definition name.
const DefPersonManagerName = "personManager"

//DefPersonManager is definition getter.
func DefPersonManager() di.Def {
	return di.Def{
		Name: DefPersonManagerName,
		Tags: []di.Tag{{
			Name: domain.DefTagLoaderRelations,
		}, {
			Name: domain.DefTagLoaderEntities,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var (
				logger  = ctn.Get(gzzap.BundleName).(*zap.Logger).Named(DefPersonManagerName)
				config  = ctn.Get(gzViper.BundleName).(*gzViper.Viper).Sub("app.imdb")
				dirData = ctn.Get(pwd.BundleName).(*pwd.PWD)
				repo    = ctn.Get(database.DefPersonsRepositoryName).(domain.PersonRepository)
			)

			return &nameManager{
				logger:   logger,
				repo:     repo,
				filepath: fmt.Sprintf("%s/%s/%s", dirData.CurrentDir(), config.GetString("download_path"), config.GetString(personFilenameConfigPath)),
				config:   config,
			}, nil
		},
	}
}

func (m *nameManager) Save() error {

	return nil
}

func (m *nameManager) LoadRelations(ctx context.Context) (err error) {
	var complete = make(chan struct{})
	go func() {
		err = m.repo.LoadRelationsFromCSV(m.config.GetString(personFilenameConfigPath))
		complete <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-complete:
			return
		}
	}
}

func (m *nameManager) Name() string {
	return DefPersonManagerName
}

func (m *nameManager) LoadEntities(ctx context.Context) (err error) {
	var complete = make(chan struct{})
	go func() {
		err = m.repo.LoadFromCSV(m.config.GetString(personFilenameConfigPath))
		complete <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-complete:
			return
		}
	}
}
