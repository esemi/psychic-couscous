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

const movieFilenameConfigPath = "files.movies"

type movieManager struct {
	logger   *zap.Logger
	repo     domain.MovieRepository
	filepath string
	config   *viper.Viper
}

//compile time check.
var _ domain.LoaderRelations = (*movieManager)(nil)

//DefMovieManagerName is definition name.
const DefMovieManagerName = "movieManager"

//DefMovieManager is definition getter.
func DefMovieManager() di.Def {
	return di.Def{
		Name: DefMovieManagerName,
		Tags: []di.Tag{{
			Name: domain.DefTagLoaderRelations,
		}, {
			Name: domain.DefTagLoaderEntities,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var (
				logger  = ctn.Get(gzzap.BundleName).(*zap.Logger).Named(DefMovieManagerName)
				config  = ctn.Get(gzViper.BundleName).(*gzViper.Viper).Sub("app.imdb")
				dirData = ctn.Get(pwd.BundleName).(*pwd.PWD)
				repo    = ctn.Get(database.DefMoviesRepositoryName).(domain.MovieRepository)
			)

			return &movieManager{
				logger:   logger,
				repo:     repo,
				filepath: fmt.Sprintf("%s/%s/%s", dirData.CurrentDir(), config.GetString("download_path"), config.GetString(movieFilenameConfigPath)),
				config:   config,
			}, nil
		},
	}
}

func (m *movieManager) LoadRelations(ctx context.Context) error {
	return nil
}

func (m *movieManager) Name() string {
	return DefMovieManagerName
}

func (m *movieManager) LoadEntities(ctx context.Context) (err error) {
	var complete = make(chan struct{})
	go func() {
		err = m.repo.LoadFromCSV(m.config.GetString(movieFilenameConfigPath))
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
