package manager

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	gzViper "github.com/gozix/viper/v2"
	gzzap "github.com/gozix/zap/v2"
	"github.com/sarulabs/di/v2"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/database"
	"gitlab.backend.keenetic.link/imdb-graph/app/cmd/app/internal/domain"
	"gitlab.backend.keenetic.link/imdb-graph/app/gozix/pwd"
	"go.uber.org/zap"
)

const nameFilenameConfigPath = "files.names"

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
				repo    = ctn.Get(database.DefNamesRepositoryName).(domain.PersonRepository)
			)

			return &nameManager{
				logger:   logger,
				repo:     repo,
				filepath: fmt.Sprintf("%s/%s/%s", dirData.CurrentDir(), config.GetString("download_path"), config.GetString(nameFilenameConfigPath)),
				config:   config,
			}, nil
		},
	}
}

func (m *nameManager) Save() error {

	return nil
}

func (m *nameManager) LoadRelations(ctx context.Context) error {
	return nil

	f, err := os.Open(m.filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	cr := csv.NewReader(gr)
	cr.Comma = '\t'
	var i int64
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rec, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		i++
		if i == 1 {
			continue
		}

		var birthY, deathY int64
		birthY, _ = strconv.ParseInt(rec[2], 10, 8)
		deathY, _ = strconv.ParseInt(rec[3], 10, 8)

		var profs = make([]string, 0)
		if rec[4] != "" {
			profs = strings.Split(rec[4], ",")
		}
		var kft = make([]domain.TitleID, 0)
		if rec[5] != "" {
			for _, str := range strings.Split(rec[5], ",") {
				kft = append(kft, domain.TitleID(str))
			}
		}

		item := domain.Person{
			ID:                 domain.PersonID(rec[0]),
			PrimaryName:        rec[1],
			BirthYear:          int16(birthY),
			DeathYear:          int16(deathY),
			PrimaryProfessions: profs,
			KnownForTitles:     kft,
		}
		if err = m.repo.Save(item); err != nil {
			return err
		}

		if i%10 == 0 {
			m.logger.Info("Processed lines", zap.Int64("lines", i))
		}
	}
}

func (m *nameManager) Name() string {
	return DefPersonManagerName
}

func (m *nameManager) LoadEntities(ctx context.Context) (err error) {
	var complete = make(chan struct{})
	go func() {
		err = m.repo.LoadFromCSV(m.config.GetString(nameFilenameConfigPath))
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
